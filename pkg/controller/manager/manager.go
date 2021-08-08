package manager

import (
	"context"
	"errors"
	"fmt"
	"github.com/shlande/dmhy-rss/pkg/controller/mission"
	"github.com/shlande/dmhy-rss/pkg/controller/store"
	"github.com/shlande/dmhy-rss/pkg/controller/subscriber"
	"github.com/shlande/dmhy-rss/pkg/data"
	"github.com/shlande/dmhy-rss/pkg/data/tools"
	"github.com/shlande/dmhy-rss/third_part/workqueue"
	"github.com/sirupsen/logrus"
	"time"
)

func NewManager(provider tools.CollectionProvider, p store.Interface) *Manager {
	return &Manager{
		msq:                workqueue.NewDelayingQueue(),
		CollectionProvider: provider,
		Store:              p,
		shutdown:           make(chan struct{}),
	}
}

type Manager struct {
	msq workqueue.DelayingInterface

	ctx context.Context

	tools.CollectionProvider

	subscriber.Subscriber

	Store store.Interface

	shutdown chan struct{}
}

func (m *Manager) Run(ctx context.Context) {
	go func() {
		<-ctx.Done()
		logrus.Print("Manger正在退出")
		m.msq.ShutDown()
		<-m.shutdown
		logrus.Print("Manger退出完成")
	}()
	for m.run() {
	}
}

func (m *Manager) run() bool {
	val, sd := m.msq.Get()
	if sd {
		close(m.shutdown)
		return false
	}

	ms := val.(*mission.Mission)
	m.update(ms)

	switch ms.Status {
	case mission.Finish:
		m.msq.Done(ms)
	default:
		m.msq.AddAfter(ms, ms.GetNextUpdateDelay())
	}
	return true
}

func (m *Manager) update(ms *mission.Mission) {
	// 标记状态为更新
	if ms.Status == mission.Waiting {
		m.done(ms, nil)
	}

	ctx, cf := m.getUpdateTimeout()
	defer cf()
	// 查找资源信息
	cls, err := m.CollectionProvider.Search(ctx, ms.Collection.Animation)
	if err != nil {
		m.done(ms, errors.New("无法获取资源信息："+err.Error()))
	}
	// 查找于当前资源相关的信息
	var cl *data.Collection
	for _, v := range cls {
		if v.Id() == ms.Collection.Id() {
			cl = v
			break
		}
	}
	// 更新失败
	if cl == nil || (cl.Latest < ms.GetNextExpectedResource() && len(cl.Items) <= len(ms.Items)) {
		m.done(ms, fmt.Errorf("最新的内容还未更新：local: %v, remote:%v", ms.Latest, cl.Latest))
		return
	}
	m.done(ms, ms.Merge(cl))
}

func (m *Manager) done(ms *mission.Mission, val interface{}) {
	// 获取日志，更新mission状态
	log := ms.Next(val)
	logrus.Print(log)
	if log.Action == mission.UpdateSuccess {
		// 发布事件
		for _, v := range val.([]*data.Source) {
			m.Added(v)
		}
	}

	// 记录日志
	err := m.Store.Log().Save(ms.Id(), log)
	if err != nil {
		logrus.Errorln("无法保存日志：", err)
	}
}

func (m *Manager) getUpdateTimeout() (context.Context, func()) {
	return context.WithTimeout(m.ctx, time.Second*10)
}
