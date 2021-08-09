package manager

import (
	"github.com/shlande/dmhy-rss/pkg/controller/mission"
	"github.com/shlande/dmhy-rss/pkg/controller/store/memory"
	"github.com/shlande/dmhy-rss/pkg/data"
	"github.com/shlande/dmhy-rss/pkg/data/parser"
	"github.com/shlande/dmhy-rss/pkg/data/source/dmhy"
	"github.com/shlande/dmhy-rss/pkg/data/tools"
	"testing"
	"time"
)

var ms = mission.NewMission(
	&data.Animation{
		Id:            "bgm317680",
		Name:          "かげきしょうじょ!!",
		Translated:    "歌剧少女!!",
		AirDate:       time.Now().Truncate(time.Hour * 24 * 30 * 2),
		AirBreak:      time.Hour * 24 * 7,
		TotalEpisodes: 11,
		Category:      "tv",
	},
	data.Metadata{
		Fansub:   []string{"喵萌奶茶屋"},
		Quality:  data.P1080,
		Type:     data.Episode,
		Language: data.GB | data.BIG5,
		SubType:  data.Internal,
	},
)

func TestTracking(t *testing.T) {
	p := memory.New()
	worker := NewManager(tools.CollectionProvider{
		Parser:   parser.New(),
		Provider: dmhy.NewProvider(),
	}, p.Mission(), p.Collection(), p.Log())

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
	}, nil, nil, nil)

	// 一次更新，应该就已经是完成状态
	worker.update(ms)
	if ms.Status != mission.Finish {
		panic("should finish")
	}
}
