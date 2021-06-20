package task

import (
	"context"
	"errors"
	"fmt"
	"github.com/shlande/dmhy-rss/pkg/classify"
	"github.com/shlande/dmhy-rss/pkg/downloader"
	"github.com/shlande/dmhy-rss/pkg/parser"
	"github.com/shlande/dmhy-rss/pkg/provider"
	"log"
	"time"
)

func NewWorker(collection *classify.Collection, updateTime time.Weekday, pvd provider.Provider, ps parser.Parser) *worker {
	return &worker{
		parser:     ps,
		Id:         collection.Id(),
		Latest:     nil,
		Collection: collection,
		provider:   pvd,
		policy:     NewPolicy(updateTime),
		dl:         nil,
	}
}

// worker 包含了所有需要执行的任务
type worker struct {
	parser parser.Parser
	// 将名称，类型，翻译组都进行hash计算得到结果
	Id  string `json:"hash"`
	cxt context.Context
	// 补充信息
	time.Timer
	Latest *parser.Detail
	*classify.Collection
	provider provider.Provider `json:"_"`
	policy   *Policy
	dl       downloader.Downloader
}

func (w *worker) Log() []*Log {
	panic("implement me")
}

func (w *worker) Stop() {
	panic("implement me")
}

func (w *worker) Run(ctx context.Context) {
	w.cxt = ctx
	// 首先尝试更新一次
	ctx, cf := w.timeLimitedCtx()
	w.CheckUpdate(ctx)
	cf()
	for {
		select {
		case <-ctx.Done():
			return
		case <-w.Timer.C:
			err := w.Update()
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func (w *worker) timeLimitedCtx() (context.Context, func()) {
	return context.WithTimeout(w.cxt, time.Second*5)
}

// Update 更新内容
func (w *worker) Update() error {
	ctx, cf := w.timeLimitedCtx()
	defer cf()
	if !w.CheckUpdate(ctx) {
		// 调整时间后移，一个小时候再检查一次
		w.Timer.Reset(time.Hour)
		return nil
	}
	// 更新
	items, err := w.find(ctx)
	if err != nil {
		return err
	}
	if len(items) < 1 {
		err = errors.New("没有找到相关信息")
		w.policy.UpdateFail(err)
		return err
	}
	var msg string
	w.Latest = items[0]
	if len(items) > 1 {
		msg = "找到多条信息，使用第一条"
		log.Println(msg)
	}
	w.Collection.Add(w.Latest)
	w.policy.Updated(msg)
	// 开始下载
	if len(w.getUrl()) == 0 {
		w.policy.FinishDownload("没有下载连接，跳过下载环节")
		return nil
	}
	return w.dl.Add(ctx, w.getUrl())
}

func (w *worker) find(ctx context.Context) ([]*parser.Detail, error) {
	infos := w.provider.Keywords(ctx, w.Name)
	details, err := w.parser.Parse(infos...)
	if err != nil {
		return nil, err
	}
	details = classify.Find(details, &classify.Option{
		Name:     w.Name,
		Episode:  w.Latest.Episode + 1,
		Fansub:   w.Fansub,
		Category: w.Category,
		Quality:  w.Quality,
		SubType:  w.SubType,
		Language: w.Language,
	})
	return details, nil
}

func (w *worker) CheckUpdate(ctx context.Context) bool {
	if w.policy.Status == Download {
		// 下载完成了
		if w.checkDownload(ctx) {
			w.policy.FinishDownload(fmt.Sprintf("完成下载%v", w.Latest.Name))
		}
		// 下载超时
		if w.policy.CheckTimeout() {
			w.dl.Delete(ctx, w.Latest.MagnetUrl)
			w.policy.DownloadTimeout(fmt.Sprintf("下载超时，文件:%v", w.Latest.Name))
		}
		return false
	}
	return w.policy.CheckUpdate()
}

func (w *worker) getUrl() string {
	url := w.Latest.MagnetUrl
	if len(url) == 0 {
		url = w.Latest.TorrentUrl
	}
	return url
}

func (w *worker) checkDownload(ctx context.Context) bool {
	process, err := w.dl.Check(ctx, w.getUrl())
	if err != nil {
		log.Println("检测下载状态失败：", err.Error())
		return false
	}
	return process == 1
}

func (w *worker) Terminate() error {
	if w.policy.Status == Download {
		ctx, cf := w.timeLimitedCtx()
		defer cf()
		err := w.dl.Delete(ctx, w.getUrl())
		if err != nil {
			return err
		}
	}
	w.policy.Terminate("用户删除collection")
	return nil
}
