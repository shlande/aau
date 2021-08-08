package subscriber

import (
	"github.com/shlande/dmhy-rss/pkg/data"
)

// Subscriber 外部下载接口
type Subscriber interface {
	// Created 存放新创建的collection
	Created(collection *data.Collection)
	Added(detail *data.Source)
}
