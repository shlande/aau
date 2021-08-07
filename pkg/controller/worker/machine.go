package worker

import (
	"context"
	"github.com/shlande/dmhy-rss/pkg/data"
	"log"
	"time"
)

// Machine 状态机，用于控制流程
type Machine interface {
	Status() Status
	Skip() (Machine, *Log)
	Do(ctx context.Context) Machine
}

type finish struct {
	*Worker
}

func (f finish) Status() Status {
	return Finish
}

func (f finish) Skip() (Machine, *Log) {
	panic("implement me")
}

func (f finish) Do(ctx context.Context) Machine {
	panic("implement me")
}

// waiting 等待下一次更新
type waiting struct {
	*Worker
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
	return &update{Worker: w.Worker}
}

type update struct {
	*Worker
	*time.Timer
}

func (w *update) Skip() (Machine, *Log) {
	panic("implement me")
}

func (w *update) Status() Status {
	return Update
}

func (w *update) retry() Machine {
	return &update{Worker: w.Worker, Timer: w.getTimer()}
}

func (w *update) Do(ctx context.Context) Machine {
	if w.sleep(ctx) {
		return nil
	}
	cls, err := w.provider.Search(ctx, w.Animation)
	if err != nil {
		w.addLog(newLog(UpdateFail, err.Error()))
		return w.retry()
	}
	var cl *data.Collection
	for _, v := range cls {
		if v.Id() == w.Collection.Id() {
			cl = v
			break
		}
	}

	// 更新失败
	if cl == nil || cl.Latest < w.getExpectedEpisode() {
		w.addLog(newLog(UpdateFail, "没有找到符合要求的数据"))
		return w.retry()
	}

	for _, v := range cl.Items {
		err := w.Collection.Add(v)
		if err != nil && err != data.ErrEpisodeExist {
			log.Println(err)
		} else {
			if w.Worker.subscriber != nil {
				w.Worker.subscriber.Added(ctx, v)
			}
		}
	}
	// 更新完成，开始下载
	w.addLog(newLog(UpdateFinish, ""))
	// 判断是否完成
	if w.Collection.IsFull() {
		return finish{Worker: w.Worker}
	}
	return w.next()
}

func (w *update) next() Machine {
	return &waiting{Worker: w.Worker, Timer: w.getTimer()}
}

func (w *update) getTimer() *time.Timer {
	timer := w.Timer
	if timer == nil {
		timer = time.NewTimer(w.getNextUpdateTime())
	} else {
		timer.Reset(w.getNextUpdateTime())
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

func (w *update) getNextUpdateTime() time.Duration {
	return getNextUpdateTime(w.AirWeekday, w.LastUpdate)
}

func getNextUpdateTime(weekday time.Weekday, lastUpdate time.Time) time.Duration {
	// 获取当前是第几周
	var day int
	if lastUpdate.Weekday() >= weekday {
		day = int(weekday + 7 - lastUpdate.Weekday())
	} else {
		day = int(weekday - lastUpdate.Weekday())
	}
	return time.Date(lastUpdate.Year(), lastUpdate.Month(), lastUpdate.Day()+day, 0, 0, 0, 0, time.Local).Sub(lastUpdate)
}
