package core

import (
	"context"
	"errors"
	"fmt"
	iface "github.com/shlande/dmhy-rss"
	"github.com/shlande/dmhy-rss/clawer"
	"log"
	"time"
)

func NewCollection(collection *clawer.Collection, updateTime time.Weekday, episodes int) *Collection {
	return &Collection{
		// 进行hash计算
		Hash:        "hash-here",
		WantEpisode: episodes,
		Info:        collection,
		Policy:      NewPolicy(updateTime),
	}

}

type Collection struct {
	// 将名称，类型，翻译组都进行hash计算得到结果
	Hash string
	Status
	WantEpisode int
	Info        *clawer.Collection
	*Policy
	dl iface.Downloader
	// 补充信息
	latest *clawer.Item
}

func (c *Collection) Update(ctx context.Context) error {
	items, err := clawer.FindItems(ctx, &clawer.Option{
		Name:     c.Info.Name,
		Episode:  c.Info.Latest + 1,
		Fansub:   c.Info.Fansub,
		Category: c.Info.Category,
		Quality:  c.Info.Quality,
		SubType:  c.Info.SubType,
		Language: c.Info.Language,
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
	c.Info.AddItem(c.latest)
	c.Updated(msg)
	// 开始下载
	return c.dl.Add(ctx, c.latest.MagnetUrl)
}

func (c *Collection) CheckUpdate(ctx context.Context) bool {
	if c.Policy.Status == Download {
		// 下载完成了
		if c.checkDownload(ctx) {
			c.Policy.FinishDownload(fmt.Sprintf("完成下载%v", c.latest.Name))
		}
		// 下载超时
		if c.Policy.CheckTimeout() {
			c.dl.Delete(ctx, c.latest.MagnetUrl)
			c.Policy.DownloadTimeout(fmt.Sprintf("下载超时，文件:%v", c.latest.Name))
		}
		return false
	}
	return c.Policy.CheckUpdate()
}

func (c *Collection) checkDownload(ctx context.Context) bool {
	process, err := c.dl.Check(ctx, c.latest.MagnetUrl)
	if err != nil {
		log.Println("检测下载状态失败：", err.Error())
		return false
	}
	return process == 1
}

func (c *Collection) Terminate(ctx context.Context) error {
	if c.Policy.Status == Download {
		err := c.dl.Delete(ctx, c.latest.MagnetUrl)
		if err != nil {
			return err
		}
	}
	c.Policy.Terminate("用户删除collection")
	return nil
}
