package worker

import (
	"context"
	"github.com/shlande/dmhy-rss/pkg/classify"
	"time"
)

// Machine 状态机，用于控制流程
type Machine interface {
	Status() Status
	Skip() (Machine, *Log)
	Do(ctx context.Context) (Machine, *Log)
}

// waiting 等待下一次更新
type waiting struct {
	*worker
	*time.Timer
}

func (w *waiting) Skip() (Machine, *Log) {
	panic("implement me")
}

func (w *waiting) Status() Status {
	return Wait
}

func (w *waiting) Do(ctx context.Context) (Machine, *Log) {
	select {
	case <-ctx.Done():
		return nil, &Log{
			Action:   Terminate,
			EmitTime: time.Now(),
			Message:  "ctx down",
		}
	case <-w.Timer.C:
		return w.next(), nil
	}
}

func (w *waiting) next() Machine {
	return &update{worker: w.worker}
}

type update struct {
	*worker
	*time.Timer
}

func (w *update) Skip() (Machine, *Log) {
	panic("implement me")
}

func (w *update) Status() Status {
	return Update
}

func (w *update) retry() Machine {
	return &update{worker: w.worker, Timer: w.getTimer()}
}

func (w *update) Do(ctx context.Context) (Machine, *Log) {
	w.sleep(ctx)
	infos, err := w.provider.Keywords(ctx, w.Name)
	if err != nil {
		return w.retry(), newLog(UpdateFail, err.Error())
	}
	details, err := w.parser.Parse(infos...)
	if err != nil {
		return w.retry(), newLog(UpdateFail, err.Error())
	}
	details = classify.Find(details, &classify.Condition{
		Name:     w.Name,
		Fansub:   w.Fansub,
		Category: w.Category,
		Quality:  w.Quality,
		SubType:  w.SubType,
		Language: w.Language,
	}, classify.After(w.Latest))
	// 更新失败
	if len(details) < 1 {
		return w.retry(), newLog(UpdateFail, "没有找到符合要求的数据")
	}
	for _, v := range details {
		w.Collection.Add(v)
		if w.worker.subscriber != nil {
			w.worker.subscriber.Added(ctx, v)
		}
	}
	// 更新完成，开始下载
	return w.next(), newLog(UpdateFinish, "")
}

func (w *update) next() Machine {
	return &waiting{worker: w.worker, Timer: w.getTimer()}
}

func (w *update) getTimer() *time.Timer {
	timer := w.Timer
	if timer == nil {
		timer = time.NewTimer(getNextUpdateTime(w.UpdateTime))
	} else {
		timer.Reset(getNextUpdateTime(w.UpdateTime))
	}
	return timer
}

func (w *update) sleep(ctx context.Context) (ctxDone bool) {
	if w.Timer == nil {
		return false
	}
	for {
		select {
		case <-ctx.Done():
			return true
		case <-w.Timer.C:
			return false
		}
	}
}

func getNextUpdateTime(weekday time.Weekday) time.Duration {
	return time.Second
}
