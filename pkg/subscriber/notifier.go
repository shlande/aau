package subscriber

import (
	"context"
	"github.com/shlande/dmhy-rss/pkg/parser"
)

// Subscriber 外部下载接口
type Subscriber interface {
	Added(ctx context.Context, detail *parser.Detail)
}
