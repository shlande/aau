package worker

import (
	"context"
	subscriber "github.com/shlande/dmhy-rss/pkg/controller/subscriber"
	"github.com/shlande/dmhy-rss/pkg/data"
	"github.com/shlande/dmhy-rss/pkg/data/tools"
	"log"
	"time"
)

func Recover(status Status, collection *data.Collection, logs []*Log) *RecoverHelper {
	return &RecoverHelper{
		Status:     status,
		Id:         collection.Id(),
		Collection: collection,
		logs:       logs,
	}
}

type RecoverHelper Worker

func (r *RecoverHelper) Recover(ctx context.Context, pvd tools.CollectionProvider, sub subscriber.Subscriber) *Worker {
	r.provider = pvd
	r.subscriber = sub
	r.end = make(chan struct{}, 1)
	wk := (*Worker)(r)
	go wk.run(ctx)
	return wk
}

func NewWorker(collection *data.Collection, pvd tools.CollectionProvider, sub subscriber.Subscriber) *Worker {
	return &Worker{
		Id:         collection.Id(),
		Collection: collection,
		provider:   pvd,
		subscriber: sub,
		end:        make(chan struct{}, 1),
	}
}

// 基础资源
type Worker struct {
	SkipTime int

	Id string
	Status

	cf  func()
	end chan struct{}

	*data.Collection

	provider tools.CollectionProvider

	subscriber subscriber.Subscriber
	logs       []*Log
}

func (w *Worker) Run(ctx context.Context) {
	ctx, w.cf = context.WithCancel(ctx)
	// 首先尝试把内容添加到store中
	w.subscriber.Created(ctx, w.Collection)
}

func (w *Worker) run(ctx context.Context) {
	var m Machine = &waiting{Worker: w, Timer: time.NewTimer(getNextUpdateTime(w.AirWeekday, w.Collection.LastUpdate))}
	for {
		m = m.Do(ctx)
		if m == nil {
			break
		}
		w.Status = m.Status()
		if w.Status == Finish {
			log.Println("更新完成")
			break
		}
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

// TODO: 跳过后需要重新计算
func (w *Worker) getExpectedEpisode() int {
	return len(w.Items) + 1
}
