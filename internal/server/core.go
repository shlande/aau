package server

import (
	"context"
	"errors"
	port2 "github.com/shlande/dmhy-rss/internal/server/port"
	http2 "github.com/shlande/dmhy-rss/internal/server/port/http"
	"github.com/shlande/dmhy-rss/pkg/classify"
	store2 "github.com/shlande/dmhy-rss/pkg/controller/store"
	subscriber2 "github.com/shlande/dmhy-rss/pkg/controller/subscriber"
	state2 "github.com/shlande/dmhy-rss/pkg/controller/subscriber/state"
	worker2 "github.com/shlande/dmhy-rss/pkg/controller/worker"
	source2 "github.com/shlande/dmhy-rss/pkg/data/source"
	"github.com/shlande/dmhy-rss/pkg/log"
	"github.com/shlande/dmhy-rss/pkg/parser"
	"github.com/sirupsen/logrus"
	"time"
)

func NewServer(parser2 parser.Parser,
	provider2 source2.Provider,
	subscriber2 *subscriber2.Multi,
	perm store2.Store, temp store2.Store) *Server {
	server := &Server{
		Entry:    log.NewEntry("server"),
		ctx:      context.Background(),
		Parser:   parser2,
		Provider: provider2,
		Multi:    subscriber2,
		perm:     perm,
		temp:     temp,
		tasks:    make(map[string]*worker2.Worker),
	}
	server.AddSubscriber(state2.New(perm, server.GetWorker))
	server.Load()
	return server
}

type Server struct {
	*logrus.Entry
	ctx context.Context
	parser.Parser
	source2.Provider
	*subscriber2.Multi
	perm  store2.Store
	temp  store2.Store
	tasks map[string]*worker2.Worker
}

func (c *Server) GetWorker(id string) (*worker2.Worker, error) {
	wk, has := c.tasks[id]
	if !has {
		return nil, errors.New("没有找到worker")
	}
	return wk, nil
}

func (c *Server) Load() {
	workers, err := c.perm.ListWorker()
	if err != nil {
		panic(err)
	}
	for _, v := range workers {
		w := v.Recover(c.ctx, c.Provider, c.Parser, c.Multi)
		c.Infoln("加载worker id:", w.Id)
		c.tasks[w.Id] = w
	}
}

func (c *Server) StartHttp(addr string) {
	http2.Start(addr, c)
}

func (c *Server) AddSubscriber(subscriber2 subscriber2.Subscriber) {
	c.Multi.Combine(subscriber2)
}

func (c *Server) WatchStatus(workerId string) (*port2.WorkerInfo, error) {
	worker := c.tasks[workerId]
	if worker == nil {
		return nil, errors.New("没有找到任务")
	}
	return port2.NewWorkerInfo(worker), nil
}

func (c *Server) Search(ctx context.Context, words string) ([]*port2.CollectionSummary, error) {
	infos, err := c.Provider.Keywords(ctx, words)
	if err != nil {
		c.Warnln("无法获取到数据: " + err.Error())
		return nil, err
	}
	detail, err := c.Parse(infos...)
	if err != nil {
		c.Warnln("解析数据出现错误: " + err.Error())
		return nil, err
	}
	cls := classify.Classify(detail)
	c.temp.Save(cls...)
	summary := make([]*port2.CollectionSummary, 0, len(cls))
	for _, v := range cls {
		summary = append(summary, port2.NewCollectionSummary(v))
	}
	return summary, err
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
	task := worker2.NewWorker(cl, updateTime, c.Provider, c.Parser, c.Multi)
	go task.Run(c.ctx)
	c.tasks[task.Id] = task
	return nil
}

func (c *Server) GetCollection(collectinoId string) (*classify.Collection, error) {
	task, err := c.temp.Get(collectinoId)
	if errors.Is(err, store2.ErrNotFound) {
		task, err = c.perm.Get(collectinoId)
		if err != nil {
			return nil, err
		}
	}
	return task, nil
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

func (c *Server) WatchList() []*port2.WorkerInfo {
	panic("implement me")
}
