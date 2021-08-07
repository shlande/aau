package worker

import (
	"context"
	"github.com/shlande/dmhy-rss/pkg/data"
	"github.com/shlande/dmhy-rss/pkg/data/parser"
	"github.com/shlande/dmhy-rss/pkg/data/source/dmhy"
	"github.com/shlande/dmhy-rss/pkg/data/tools"
	"testing"
	"time"
)

func TestWorker(t *testing.T) {
	ctx := context.Background()

	worker := NewWorker(data.NewCollection(&data.Animation{
		Name:          "無職転生～異世界行ったら本気だす～",
		Translated:    "无职转生～到了异世界就拿出真本事～",
		AirDate:       time.Now().Truncate(time.Hour * 24 * 30 * 2),
		AirWeekday:    time.Sunday,
		AirTime:       time.Hour * 3,
		TotalEpisodes: 11,
		Category:      "tv",
	}, data.Metadata{
		Fansub:   []string{"桜都字幕組"},
		Quality:  data.P1080,
		Type:     data.Episode,
		Language: 4,
		SubType:  data.Internal,
	}), tools.CollectionProvider{
		Parser:   parser.New(),
		Provider: dmhy.NewProvider(),
	}, nil)

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

func TestOnceCollect(t *testing.T) {
	ctx := context.Background()

	worker := NewWorker(data.NewCollection(&data.Animation{
		Name:          "無職転生～異世界行ったら本気だす～",
		Translated:    "无职转生～到了异世界就拿出真本事～",
		AirDate:       time.Now().Truncate(time.Hour * 24 * 30 * 2),
		AirWeekday:    time.Sunday,
		AirTime:       time.Hour * 3,
		TotalEpisodes: 11,
		Category:      "tv",
	}, data.Metadata{
		Fansub:   []string{"桜都字幕組"},
		Quality:  data.P1080,
		Type:     data.Episode,
		Language: 4,
		SubType:  data.Internal,
	}), tools.CollectionProvider{
		Parser:   parser.New(),
		Provider: dmhy.NewProvider(),
	}, nil)

	var m Machine = &waiting{Worker: worker, Timer: time.NewTimer(0)}
	// 第一次，应该跳转到waiting状态
	m = m.Do(ctx)
	if m.Status() != Update {
		panic("0")
	}
	// 准备更新
	m = m.Do(ctx)
	if m.Status() != Finish {
		panic("1")
	}
}
