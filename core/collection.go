package core

import (
	"context"
	"errors"
	"github.com/shlande/dmhy-rss/clawer"
	"log"
)

type Collection struct {
	// 将名称，类型，翻译组都进行hash计算得到结果
	Hash string
	Status
	WantEpisode int
	Info        *clawer.Collection
	Policy
	// 补充信息
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
		return errors.New("没有找到相关信息")
	}
	if len(items) > 1 {
		log.Println("找到多条信息，使用第一条")
		c.Info.AddItem(items[0])
	}
	c.Info.AddItem(items[0])
	return err
}
