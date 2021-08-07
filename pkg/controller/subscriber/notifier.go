package subscriber

import (
	"context"
	"github.com/shlande/dmhy-rss/pkg/classify"
	"github.com/shlande/dmhy-rss/pkg/parser"
)

// Subscriber 外部下载接口
type Subscriber interface {
	// Created 存放新创建的collection
	Created(ctx context.Context, collection *classify.Collection)
	Added(ctx context.Context, detail *parser.Detail)
}
