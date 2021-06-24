package worker

import (
	"context"
	"github.com/shlande/dmhy-rss/pkg/classify"
	"github.com/shlande/dmhy-rss/pkg/parser/common"
	"github.com/shlande/dmhy-rss/pkg/provider/dmhy"
	"github.com/shlande/dmhy-rss/pkg/subscriber"
	"testing"
	"time"
)

func TestWorker(t *testing.T) {
	ctx := context.Background()
	// 首先尝试获取一个detail
	parser := common.Parse{}
	provider := dmhy.NewProvider()
	var dl subscriber.Subscriber = nil
	info, err := provider.Keywords(ctx, "无职转生")
	if err != nil {
		panic(err)
	}
	details, err := parser.Parse(info...)
	if err != nil {
		panic(err)
	}
	collection := classify.Classify(details)
	if err != nil {
		panic(err)
	}
	cl := collection[0]
	cl.Latest = cl.Latest - 2
	worker := NewWorker(cl, time.Sunday, provider, parser, dl)

	var m Machine = &waiting{Worker: worker, Timer: time.NewTimer(0)}
	// 第一次，应该跳转到waiting状态
	m = m.Do(ctx)
	if m.Status() != Update {
		panic("0")
	}
	// 准备更新
	m = m.Do(ctx)
	if m.Status() != Wait {
		panic("1")
	}
	// 这里更新应该完成，手动把之前设置的latest调整回去
	worker.Collection.Latest += 2
	// 准备下载
	m = m.Do(ctx)
	m = m.Do(ctx)
	if m.Status() != Update {
		panic("2")
	}
}
