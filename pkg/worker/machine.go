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
	return &update{w.worker}
}

type update struct {
	*worker
}

func (u *update) Skip() (Machine, *Log) {
	panic("implement me")
}

func (u *update) Status() Status {
	return Update
}

func (u *update) retry() Machine {
	return &waiting{worker: u.worker, timer: time.NewTimer(time.Hour)}
}

func (w *update) Do(ctx context.Context) (Machine, *Log) {
	// 查找合适的内容
	infos, err := w.provider.Keywords(ctx, w.Name)
	if err != nil {
		return w.retry(), &Log{}
	}
	details, err := w.parser.Parse(infos...)
	if err != nil {
		return w.retry(), &Log{
			Action:   UpdateFail,
			EmitTime: time.Now(),
			Message:  err.Error(),
		}
	}
	episode := 1
	if w.Latest != nil {
		episode = w.Latest.Episode + 1
	}
	details = classify.Find(details, &classify.Option{
		Name:     w.Name,
		Episode:  episode,
		Fansub:   w.Fansub,
		Category: w.Category,
		Quality:  w.Quality,
		SubType:  w.SubType,
		Language: w.Language,
	})
	// 更新失败
	if len(details) < 1 {
		return w.retry(), &Log{
			Action:   UpdateFail,
			EmitTime: time.Now(),
			Message:  "没有找到符合条件的内容",
		}
	}
	w.Latest = details[0]
	w.Collection.Add(w.Latest)
	// 更新完成，开始下载
	log := &Log{
		Action:   UpdateFinish,
		EmitTime: time.Now(),
	}
	if len(details) > 1 {
		log.Message = "找到多条记录，使用第一条"
	}
	return &download{w.worker}, log
}

type download struct {
	*worker
}

func (d *download) Skip() (Machine, *Log) {
	panic("implement me")
}

func (d *download) Status() Status {
	return Download
}

func (d *download) getUrl() string {
	url := d.Latest.MagnetUrl
	if len(url) == 0 {
		url = d.Latest.TorrentUrl
	}
	return url
}

func (d *download) Do(ctx context.Context) (Machine, *Log) {
	// 获取下载地址
	url := d.getUrl()
	if len(url) == 0 {
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
	err := d.dl.Add(ctx, url)
	if err != nil {
		return d.next(), &Log{
			Action:   DownloadCancel,
			EmitTime: time.Now(),
			Message:  "下载失败：" + err.Error(),
		}
	}
	// 如果成功，则进行观察
	ticker := time.NewTimer(time.Minute)
	defer ticker.Stop()
	var (
		faild   int
		process float64
	)
	for {
		select {
		case <-ctx.Done():
			return nil, nil
		case <-ticker.C:
			process, err = d.dl.Check(ctx, url)
			if err != nil {
				faild += 1
			}
			// 最大重试次数
			if faild >= 5 {
				return d.next(), &Log{
					Action:   DownloadCancel,
					EmitTime: time.Now(),
					Message:  "下载检查错误次数达到最大限制：" + err.Error(),
				}
			}
			if process == 1 {
				return d.next(), &Log{
					Action:   DownloadFinish,
					EmitTime: time.Now(),
				}
			}
		}
	}
}

func (d *download) next() Machine {
	return &waiting{d.worker, d.getNextUpdateTime()}
}

func (d *download) getNextUpdateTime() *time.Timer {
	return time.NewTimer(getNextUpdateTime(d.UpdateTime))
}

func getNextUpdateTime(weekday time.Weekday) time.Duration {
	return time.Second
}
