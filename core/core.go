package core

import (
	"context"
	"errors"
	"fmt"
	iface "github.com/shlande/dmhy-rss"
	"github.com/shlande/dmhy-rss/pkg/parser/common"
	"time"
)

// Core 核心逻辑处理未知
type Core struct {
	// 用于存放已经确定需要更新到信息
	pined map[string]*Collection
	// 用于存放临时信息
	temp map[string]*Collection
	dw   iface.Downloader
}

func (c *Core) Search(ctx context.Context, keyword string) []*Collection {
	cls, err := common.FindCollectionsByKeywords(ctx, keyword)
	if err != nil {
		return nil
	}
	var clts []*Collection
	for _, cl := range cls {
		clt := NewCollection(cl, time.Monday, 0)
		// 暂时存放信息
		c.temp[clt.Hash] = clt
		clts = append(clts, clt)
	}
	return clts
}

func (c *Core) Add(hash string, updateTime time.Weekday, episode int) error {
	if _, ok := c.pined[hash]; ok {
		return errors.New("该collection已经加入到收录列表中了")
	}
	if clt, ok := c.temp[hash]; ok {
		clt = NewCollection(clt.Info, updateTime, episode)
		c.pined[clt.Hash] = clt
		return nil
	}
	return errors.New("未知的collection")
}

func (c *Core) Delete(hash string) error {
	if clt, ok := c.temp[hash]; ok {
		clt.Terminate(context.TODO())
		return nil
	}
	return errors.New("未知的collection")
}

func (c *Core) List() []*Collection {
	panic("implement me")
}

func (c *Core) Run(ctx context.Context) {
	// 每五分钟进行一次更新状态
	ticker := time.NewTicker(time.Minute * 5)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.updateOnce(ctx)
		}
	}
}

func (c *Core) updateOnce(ctx context.Context) {
	// 查找出一个需要更新的任务
	for _, cls := range c.pined {
		if cls.CheckUpdate(ctx) {
			err := cls.Update(ctx)
			if err != nil {
				fmt.Println(err)
			}
			return
		}
	}
}
