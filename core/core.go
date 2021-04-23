package core

import (
	"context"
	"fmt"
	"time"
)

// Core 核心逻辑处理未知
type Core struct {
	cls []*Collection
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
	for _, cls := range c.cls {
		if cls.CheckUpdate() {
			err := cls.Update(ctx)
			if err != nil {
				fmt.Println(err)
			}
			return
		}
	}
	// 开始更新collection

}
