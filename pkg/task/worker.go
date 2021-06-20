package task

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/shlande/dmhy-rss/pkg/classify"
	"github.com/shlande/dmhy-rss/pkg/downloader"
	"github.com/shlande/dmhy-rss/pkg/parser"
	"github.com/shlande/dmhy-rss/pkg/provider"
	"log"
	"time"
)

type
struct {

}

func NewWorker(collection *classify.Collection, updateTime time.Weekday, pvd provider.Provider) *worker {
	return &worker{
		provider:    pvd,
		Id:          string(collection.Id()),
		WantEpisode: episodes,
		Collection:  collection,
		Policy:      NewPolicy(updateTime),
	}
}

// worker 包含了所有需要执行的任务
type worker struct {
	// 将名称，类型，翻译组都进行hash计算得到结果
	Id          string `json:"hash"`
	// 补充信息
	Latest *parser.Detail
	*classify.Collection

	provider provider.Provider `json:"_"`
	policy *Policy
	dl downloader.Downloader
}

// Update 更新内容
func (c *worker) Update() error {
	c.CheckUpdate()
	items, err := classify.Find(ctx, &common.Option{
		Name:     c.Collection.Name,
		Episode:  c.Collection.Latest + 1,
		Fansub:   c.Collection.Fansub,
		Category: c.Collection.Category,
		Quality:  c.Collection.Quality,
		SubType:  c.Collection.SubType,
		Language: c.Collection.Language,
	})
	if err != nil {
		return err
	}
	if len(items) < 1 {
		err = errors.New("没有找到相关信息")
		c.UpdateFail(err)
		return err
	}
	var msg string
	c.latest = items[0]
	if len(items) > 1 {
		msg = "找到多条信息，使用第一条"
		log.Println(msg)
	}
	c.Collection.AddItem(c.latest)
	c.Updated(msg)
	// 开始下载
	return c.dl.Add(ctx, c.latest.MagnetUrl)
}

func (c *worker) CheckUpdate(ctx context.Context) bool {
	if c.policy.Status == Download {
		// 下载完成了
		if c.checkDownload(ctx) {
			c.policy.FinishDownload(fmt.Sprintf("完成下载%v", c.Latest.Name))
		}
		// 下载超时
		if c.policy.CheckTimeout() {
			c.dl.Delete(ctx, c.Latest.MagnetUrl)
			c.policy.DownloadTimeout(fmt.Sprintf("下载超时，文件:%v", c.Latest.Name))
		}
		return false
	}
	return c.policy.CheckUpdate()
}

func (c *worker) checkDownload(ctx context.Context) bool {
	process, err := c.dl.Check(ctx, c.latest.MagnetUrl)
	if err != nil {
		log.Println("检测下载状态失败：", err.Error())
		return false
	}
	return process == 1
}

func (c *worker) Terminate(ctx context.Context) error {
	if c.Policy.Status == Download {
		err := c.dl.Delete(ctx, c.latest.MagnetUrl)
		if err != nil {
			return err
		}
	}
	c.Policy.Terminate("用户删除collection")
	return nil
}
