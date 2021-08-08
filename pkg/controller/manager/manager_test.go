package manager

import (
	"github.com/shlande/dmhy-rss/pkg/controller/mission"
	"github.com/shlande/dmhy-rss/pkg/data"
	"github.com/shlande/dmhy-rss/pkg/data/parser"
	"github.com/shlande/dmhy-rss/pkg/data/source/dmhy"
	"github.com/shlande/dmhy-rss/pkg/data/tools"
	"testing"
	"time"
)

func TestTracking(t *testing.T) {
	worker := NewManager(tools.CollectionProvider{
		Parser:   parser.New(),
		Provider: dmhy.NewProvider(),
	}, nil)

	ms := mission.NewMission(&data.Animation{
		Name:          "かげきしょうじょ!!",
		Translated:    "歌剧少女!!",
		AirDate:       time.Now().Truncate(time.Hour * 24 * 30 * 2),
		AirWeekday:    time.Sunday,
		AirTime:       time.Hour * 3,
		TotalEpisodes: 11,
		Category:      "tv",
	}, data.Metadata{
		Fansub:   []string{"喵萌奶茶屋"},
		Quality:  data.P1080,
		Type:     data.Episode,
		Language: data.GB | data.BIG5,
		SubType:  data.Internal,
	})

	// 一次更新，应该是等待状态
	worker.update(ms)
	if ms.Status != mission.Waiting {
		panic("should waiting")
	}

	// 删除部分的内容
	ms.Collection.Latest--
	ms.Collection.Items = ms.Collection.Items[:len(ms.Collection.Items)-1]
	// 此时应该更新依然更新成功

	worker.update(ms)
	if ms.Status != mission.Waiting {
		panic("should updating")
	}

	// 此时应该更新失败，因为没有新的内容
	worker.update(ms)
	if ms.Status != mission.Updating {
		panic("should updating")
	}
}

func TestOnceCollect(t *testing.T) {
	worker := NewManager(tools.CollectionProvider{
		Parser:   parser.New(),
		Provider: dmhy.NewProvider(),
	}, nil)

	ms := mission.NewMission(&data.Animation{
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
	})
	// 一次更新，应该就已经是完成状态
	worker.update(ms)
	if ms.Status != mission.Finish {
		panic("should finish")
	}
}
