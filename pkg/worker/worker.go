package worker

import (
	"context"
	"github.com/shlande/dmhy-rss/pkg/classify"
	"github.com/shlande/dmhy-rss/pkg/parser"
	"github.com/shlande/dmhy-rss/pkg/provider"
	"github.com/shlande/dmhy-rss/pkg/subscriber"
	"time"
)

func Recover(status Status, collection *classify.Collection, updateTime time.Weekday, logs []*Log) *RecoverHelper {
	return &RecoverHelper{
		Status:     status,
		Id:         collection.Id(),
		Collection: collection,
		UpdateTime: updateTime,
		logs:       logs,
	}
}

type RecoverHelper Worker

func (r *RecoverHelper) Recover(pvd provider.Provider, ps parser.Parser, sub subscriber.Subscriber) *Worker {
	r.provider = pvd
	r.parser = ps
	r.subscriber = sub
	r.end = make(chan struct{}, 1)
	return (*Worker)(r)
}

func NewWorker(collection *classify.Collection, updateTime time.Weekday, pvd provider.Provider, ps parser.Parser, sub subscriber.Subscriber) *Worker {
	return &Worker{
		parser:     ps,
		Id:         collection.Id(),
		Collection: collection,
		provider:   pvd,
		UpdateTime: updateTime,
		subscriber: sub,
		end:        make(chan struct{}, 1),
	}
}

// 基础资源
type Worker struct {
	Id string
	Status
	cf     func()
	end    chan struct{}
	parser parser.Parser
	*classify.Collection
	UpdateTime time.Weekday
	provider   provider.Provider
	subscriber subscriber.Subscriber
	logs       []*Log
}

func (w *Worker) Run(ctx context.Context) {
	ctx, w.cf = context.WithCancel(ctx)
	// 首先尝试把内容添加到store中
	w.subscriber.Created(ctx, w.Collection)
	var m Machine = &waiting{Worker: w, Timer: time.NewTimer(getNextUpdateTime(w.UpdateTime, w.Collection.LastUpdate))}
	for {
		m = m.Do(ctx)
		if m == nil {
			break
		}
		w.Status = m.Status()
	}
	w.end <- struct{}{}
}

func (w *Worker) Terminate() {
	w.cf()
	<-w.end
}

func (w *Worker) addLog(log *Log) {
	w.logs = append(w.logs, log)
}

func (w *Worker) Log() []*Log {
	return w.logs
}

func (w *Worker) Stop() {
	panic("implement me")
}
