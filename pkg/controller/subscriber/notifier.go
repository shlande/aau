package subscriber

import (
	"context"
	"github.com/shlande/dmhy-rss/pkg/data"
)

// Subscriber 外部下载接口
type Subscriber interface {
	// Created 存放新创建的collection
	Created(ctx context.Context, collection *data.Collection)
	Added(ctx context.Context, detail *data.Source)
}
