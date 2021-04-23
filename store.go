package main

import (
	"context"
	"github.com/shlande/dmhy-rss/clawer"
)

// Store 用来保存爬到的动漫
type Store interface {
	// Store 保存数据
	Store(ctx context.Context, record *clawer.Item) error
	// FindByName 通过名称搜索
	FindByName(ctx context.Context, name string) ([]*clawer.Item, error)
}
