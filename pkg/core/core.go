package core

import (
	"context"
	"errors"
	"github.com/shlande/dmhy-rss/pkg/classify"
	"github.com/shlande/dmhy-rss/pkg/downloader"
	"github.com/shlande/dmhy-rss/pkg/parser"
	"github.com/shlande/dmhy-rss/pkg/provider"
	"github.com/shlande/dmhy-rss/pkg/store"
	"github.com/shlande/dmhy-rss/pkg/task"
	"log"
	"time"
)

type Core struct {
	ctx context.Context
	parser.Parser
	provider.Provider
	downloader.Downloader
	perm  store.Store
	temp  store.Store
	tasks map[string]task.Worker
}

func (c *Core) Keywords(ctx context.Context, words string) []*classify.Collection {
	detail, err := c.Parse(c.Provider.Keywords(ctx, words)...)
	if err != nil {
		log.Println(err)
		return nil
	}
	cls := classify.Classify(detail)
	c.temp.Save(cls...)
	return cls
}

func (c *Core) Watch(collectionId string, updateTime time.Weekday) error {
	cl, err := c.temp.Get(collectionId)
	if err != nil {
		return err
	}
	// 如果已经监控了
	if _, has := c.tasks[collectionId]; has {
		return errors.New("已经监控过了")
	}
	task := task.NewWorker(cl, updateTime, c.Provider, c.Parser)
	go task.Run(c.ctx)
	c.tasks[task.Id] = task
	return nil
}

func (c *Core) UnWatch(collectionId string) error {
	// 如果已经监控了
	task, has := c.tasks[collectionId]
	if !has {
		return errors.New("没有找到task")
	}
	task.Stop()
	return nil
}

//func (c *Core) WatchList() [] {
//	panic("implement me")
//}

func (c *Core) Log(collectionId string) []*task.Log {
	panic("implement me")
}
