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
	Do(ctx context.Context) Machine
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

func (w *waiting) Do(ctx context.Context) Machine {
	select {
	case <-ctx.Done():
		return nil
	case <-w.Timer.C:
		return w.next()
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

func (w *update) Do(ctx context.Context) Machine {
	w.sleep(ctx)
	infos, err := w.provider.Keywords(ctx, w.Name)
	if err != nil {
		w.addLog(newLog(UpdateFail, err.Error()))
		return w.retry()
	}
	details, err := w.parser.Parse(infos...)
	if err != nil {
		w.addLog(newLog(UpdateFail, err.Error()))
		return w.retry()
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
		w.addLog(newLog(UpdateFail, "没有找到符合要求的数据"))
		return w.retry()
	}
	for _, v := range details {
		w.Collection.Add(v)
		if w.worker.subscriber != nil {
			w.worker.subscriber.Added(ctx, v)
		}
	}
	// 更新完成，开始下载
	w.addLog(newLog(UpdateFinish, ""))
	return w.next()
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
