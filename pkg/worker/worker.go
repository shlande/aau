package worker

import (
	"context"
	"fmt"
	"github.com/shlande/dmhy-rss/pkg/classify"
	"github.com/shlande/dmhy-rss/pkg/downloader"
	"github.com/shlande/dmhy-rss/pkg/parser"
	"github.com/shlande/dmhy-rss/pkg/provider"
	"time"
)

// Worker 负责监控更新collection
// 当一个collection需要更新时发送消息
type Worker interface {
	Id() string
	Run(ctx context.Context)
	// Log 输出日志
	Log() []*Log
	Stop()
	// 结束并删除任务
	Terminate()
}

func NewWorker(collection *classify.Collection, updateTime time.Weekday, pvd provider.Provider, ps parser.Parser, dl downloader.Downloader) *worker {
	return &worker{
		parser:     ps,
		Id:         collection.Id(),
		Collection: collection,
		provider:   pvd,
		UpdateTime: updateTime,
		dl:         dl,
	}
}

// 基础资源
type worker struct {
	Id     string
	parser parser.Parser
	*classify.Collection
	UpdateTime time.Weekday
	provider   provider.Provider
	dl         downloader.Downloader
}

func (w *worker) Run(ctx context.Context) {
	var log *Log
	var m Machine = &waiting{worker: w, timer: time.NewTimer(0)}
	for {
		m, log = m.Do(ctx)
		fmt.Println(log)
	}
}

func (w *worker) Terminate() {
	panic("implement me")
}

func (w *worker) Log() []*Log {
	panic("implement me")
}

func (w *worker) Stop() {
	panic("implement me")
}
