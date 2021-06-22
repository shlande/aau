package worker

import (
	"context"
	"github.com/shlande/dmhy-rss/pkg/classify"
	"log"
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
	timer *time.Timer
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
	case <-w.timer.C:
		return w.next(), nil
	}
}

func (w *waiting) next() Machine {
	w.timer.Stop()
	return &update{worker: w.worker}
}

type update struct {
	*worker
	*time.Timer
}

func (u *update) Skip() (Machine, *Log) {
	panic("implement me")
}

func (u *update) Status() Status {
	return Update
}

func (u *update) retry() Machine {
	return &update{worker: u.worker, Timer: time.NewTimer(time.Hour)}
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
	details = classify.Find(details, &classify.Option{
		Name:     w.Name,
		Episode:  w.Latest + 1,
		Fansub:   w.Fansub,
		Category: w.Category,
		Quality:  w.Quality,
		SubType:  w.SubType,
		Language: w.Language,
	})
	// 更新失败
	if len(details) < 1 {
		return w.retry(), newLog(UpdateFail, "没有找到符合要求的数据")
	}
	urls := make([]string, 0, len(details))
	for _, v := range details {
		if len(v.MagnetUrl) != 0 {
			urls = append(urls, v.MagnetUrl)
		} else if len(v.TorrentUrl) != 0 {
			urls = append(urls, v.TorrentUrl)
		}
		w.Collection.Add(v)
	}
	// 更新完成，开始下载
	return &download{w.worker, urls}, newLog(UpdateFinish, "")
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

type download struct {
	*worker
	urls []string
}

func (d *download) Skip() (Machine, *Log) {
	panic("implement me")
}

func (d *download) Status() Status {
	return Download
}

func (d *download) Do(ctx context.Context) (Machine, *Log) {
	// 获取下载地址
	if len(d.urls) == 0 {
		return d.next(), &Log{
			Action:   DownloadFinish,
			EmitTime: time.Now(),
			Message:  "没有下载地址，跳过下载",
		}
	}
	if d.dl == nil {
		return d.next(), &Log{
			Action:   DownloadCancel,
			EmitTime: time.Now(),
			Message:  "没有指定下载器，跳过下载",
		}
	}
	for _, v := range d.urls {
		err := d.dl.Add(ctx, v)
		if err != nil {
			log.Println("无法添加下载任务:" + err.Error())
		}
	}
	// FIXME: 启动协程管理下载？
	return d.next(), newLog(DownloadFinish, "")
}

func (d *download) next() Machine {
	return &waiting{d.worker, time.NewTimer(getNextUpdateTime(d.UpdateTime))}
}

func getNextUpdateTime(weekday time.Weekday) time.Duration {
	return time.Second
}
