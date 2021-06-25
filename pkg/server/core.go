package server

import (
	"context"
	"errors"
	"github.com/shlande/dmhy-rss/pkg/classify"
	"github.com/shlande/dmhy-rss/pkg/log"
	"github.com/shlande/dmhy-rss/pkg/parser"
	"github.com/shlande/dmhy-rss/pkg/provider"
	"github.com/shlande/dmhy-rss/pkg/server/port"
	"github.com/shlande/dmhy-rss/pkg/server/port/http"
	"github.com/shlande/dmhy-rss/pkg/store"
	"github.com/shlande/dmhy-rss/pkg/subscriber"
	"github.com/shlande/dmhy-rss/pkg/worker"
	"github.com/sirupsen/logrus"
	"time"
)

func NewServer(parser2 parser.Parser,
	provider2 provider.Provider,
	subscriber2 subscriber.Subscriber,
	perm store.Store, temp store.Store) *Server {
	return &Server{
		Entry:      log.NewEntry("server"),
		ctx:        context.Background(),
		Parser:     parser2,
		Provider:   provider2,
		Subscriber: subscriber2,
		perm:       perm,
		temp:       temp,
		tasks:      make(map[string]*worker.Worker),
	}
}

type Server struct {
	*logrus.Entry
	ctx context.Context
	parser.Parser
	provider.Provider
	subscriber.Subscriber
	perm  store.Store
	temp  store.Store
	tasks map[string]*worker.Worker
}

func (c *Server) StartHttp(addr string) {
	http.Start(addr, c)
}

func (c *Server) GetWorker(workerId string) *port.WorkerInfo {
	worker := c.tasks[workerId]
	return port.NewWorkerInfo(worker)
}

func (c *Server) Search(ctx context.Context, words string) []*classify.Collection {
	infos, err := c.Provider.Keywords(ctx, words)
	if err != nil {
		c.Warnln("无法获取到数据: " + err.Error())
		return nil
	}
	detail, err := c.Parse(infos...)
	if err != nil {
		c.Warnln("解析数据出现错误: " + err.Error())
		return nil
	}
	cls := classify.Classify(detail)
	c.temp.Save(cls...)
	return cls
}

func (c *Server) Watch(collectionId string, updateTime time.Weekday) error {
	cl, err := c.temp.Get(collectionId)
	if err != nil {
		return err
	}
	// 如果已经监控了
	if _, has := c.tasks[collectionId]; has {
		return errors.New("已经监控过了")
	}
	task := worker.NewWorker(cl, updateTime, c.Provider, c.Parser, c.Subscriber)
	go task.Run(c.ctx)
	c.tasks[task.Id] = task
	return nil
}

func (c *Server) GetCollection(collectinoId string) *classify.Collection {
	task, err := c.temp.Get(collectinoId)
	if errors.Is(err, store.ErrNotFound) {
		task, _ = c.perm.Get(collectinoId)
	}
	return task
}

func (c *Server) UnWatch(collectionId string) error {
	// 如果已经监控了
	task, has := c.tasks[collectionId]
	if !has {
		return errors.New("没有找到task")
	}
	task.Terminate()
	c.Info("结束监控任务：" + collectionId)
	delete(c.tasks, collectionId)
	return nil
}

func (c *Server) WatchList() []*port.WorkerInfo {
	panic("implement me")
}
