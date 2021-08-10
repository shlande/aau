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

func NewManager(provider *tools.CollectionProvider,
	sb subscriber.Subscriber,
	pm store.MissionInterface,
	pc store.CollectionInterface,
	pl store.LogInterface) *Manager {
	m := &Manager{
		msq:                 workqueue.NewDelayingQueue(),
		tracker:             map[string]struct{}{},
		CollectionProvider:  provider,
		MissionInterface:    pm,
		CollectionInterface: pc,
		LogInterface:        pl,
		Subscriber:          sb,
		addMsChan:           make(chan *mission.Mission),
		shutdown:            make(chan struct{}),
		ttlTimer:            time.NewTimer(0),
	}

	// 从数据库中加载出来持久的mission
	mss, err := pm.GetAll(true)
	if err != nil && err != store.ErrNotFound {
		panic(err)
	}
	// 添加mission到manager中
	for _, v := range mss {
		m.addMission(v)
		logrus.Infoln("正在恢复任务：", v.Id(), v.Name)
	}
	return m
}

type Manager struct {
	// 这个只是用来追踪任务的
	tracker map[string]struct{}
	// addMsChan 传递需要固定的任务
	addMsChan chan *mission.Mission

	msq workqueue.DelayingInterface

	ctx context.Context

	*tools.CollectionProvider

	subscriber.Subscriber

	store.MissionInterface
	store.CollectionInterface
	store.LogInterface

	shutdown chan struct{}

	ttlTimer *time.Timer
}

func (m *Manager) ListActiveMissionId() (ids []string) {
	for k, _ := range m.tracker {
		ids = append(ids, k)
	}
	return
}

func (m *Manager) AddChan() chan<- *mission.Mission {
	return m.addMsChan
}

func (m *Manager) Run(ctx context.Context) {
	m.ctx = ctx
	go func() {
		<-ctx.Done()
		logrus.Print("Manger正在退出")
		m.msq.ShutDown()
		<-m.shutdown
		logrus.Print("Manger退出完成")
	}()

	go m.addLoop()

	for m.run() {
		m.waitTTL()
	}
}

func (m *Manager) waitTTL() {
	// 尝试清空上次的timer
	select {
	case <-m.ttlTimer.C:
	}

	m.ttlTimer.Reset(time.Minute)

	// 等待时间
	select {
	case <-m.ttlTimer.C:
	case <-m.ctx.Done():
	}
}

func (m *Manager) addLoop() {
	for {
		select {
		case <-m.ctx.Done():
			return
		case ms := <-m.addMsChan:
			// 从外部读取到到mission，一定要更新
			err := m.MissionInterface.Save(ms)
			if err != nil {
				logrus.Errorln("保存mission失败：", err)
			}
			err = m.CollectionInterface.Save(ms.Collection)
			if err != nil {
				logrus.Errorln("保存mission失败：", err)
			}
			m.addMission(ms)
		}
	}
}

func (m *Manager) addMission(ms *mission.Mission) {
	if _, has := m.tracker[ms.Id()]; has {
		return
	}
	m.tracker[ms.Id()] = struct{}{}
	m.msq.AddAfter(ms.Id(), ms.GetNextUpdateDelay())
}

func (m *Manager) run() bool {
	val, sd := m.msq.Get()
	if sd {
		close(m.shutdown)
		return false
	}
	// 从store中取，保证内容是最新的
	ms, err := m.MissionInterface.Get(val.(string))
	if err != nil {
		logrus.Errorln("无法更新mission:", err)
		return true
	}

	m.update(ms)
	// 保存mission
	err = m.MissionInterface.Save(ms)
	if err != nil {
		logrus.Errorln(err)
	}

	switch ms.Status {
	case mission.Finish:
		delete(m.tracker, ms.Id())
		m.msq.Done(ms.Id())
	default:
		m.msq.AddAfter(ms.Id(), ms.GetNextUpdateDelay())
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
	if cl == nil || (cl.Latest < ms.GetExpectedLatest() && len(cl.Items) <= len(ms.Items)) {
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
			if m.Subscriber != nil {
				m.Subscriber.Added(v)
			}
		}
		// 保存
		err := m.CollectionInterface.Save(ms.Collection)
		if err != nil {
			logrus.Errorln("无法保存collection：", err)
		}
	}

	// 记录日志
	if m.LogInterface != nil {
		err := m.LogInterface.Save(ms.Id(), log)
		if err != nil {
			logrus.Errorln("无法保存日志：", err)
		}
	}
}

func (m *Manager) getUpdateTimeout() (context.Context, func()) {
	ctx := m.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithTimeout(ctx, time.Second*10)
}
