package task

import (
	"context"
	"github.com/shlande/dmhy-rss/pkg/classify"
	"github.com/shlande/dmhy-rss/pkg/downloader"
	"github.com/shlande/dmhy-rss/pkg/parser"
	"github.com/shlande/dmhy-rss/pkg/provider"
	"time"
)

func NewWorker(collection *classify.Collection, updateTime time.Weekday, pvd provider.Provider, ps parser.Parser) *worker {
	return &worker{
		parser:     ps,
		Id:         collection.Id(),
		Latest:     nil,
		Collection: collection,
		provider:   pvd,
		UpdateTime: updateTime,
		dl:         nil,
	}
}

// 基础资源
type worker struct {
	Id     string
	parser parser.Parser
	Latest *parser.Detail
	*classify.Collection
	UpdateTime time.Weekday
	provider   provider.Provider
	dl         downloader.Downloader
}

func (w *worker) Run(ctx context.Context) {
	panic("implement me")
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
