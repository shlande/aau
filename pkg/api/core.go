package api

import (
	"context"
	"errors"
	"github.com/shlande/dmhy-rss/pkg/classify"
	"github.com/shlande/dmhy-rss/pkg/parser"
	"github.com/shlande/dmhy-rss/pkg/provider"
	"github.com/shlande/dmhy-rss/pkg/store"
	"github.com/shlande/dmhy-rss/pkg/subscriber"
	"github.com/shlande/dmhy-rss/pkg/worker"
	"time"
)

type Server struct {
	ctx context.Context
	parser.Parser
	provider.Provider
	subscriber.Subscriber
	perm  store.Store
	temp  store.Store
	tasks map[string]worker.Worker
}

func (c *Server) Keywords(ctx context.Context, words string) ([]*classify.Collection, error) {
	infos, err := c.Provider.Keywords(ctx, words)
	if err != nil {
		return nil, err
	}
	detail, err := c.Parse(infos...)
	if err != nil {
		return nil, err
	}
	cls := classify.Classify(detail)
	c.temp.Save(cls...)
	return cls, nil
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

func (c *Server) UnWatch(collectionId string) error {
	// 如果已经监控了
	task, has := c.tasks[collectionId]
	if !has {
		return errors.New("没有找到task")
	}
	task.Terminate()
	delete(c.tasks, task.Id())
	return nil
}

//func (c *Server) WatchList() [] {
//	panic("implement me")
//}

func (c *Server) Log(collectionId string) []*worker.Log {
	panic("implement me")
}
