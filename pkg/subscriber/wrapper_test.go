package subscriber

import (
	"context"
	"github.com/shlande/dmhy-rss/pkg/parser"
	"github.com/shlande/dmhy-rss/pkg/provider"
	"github.com/shlande/dmhy-rss/pkg/subscriber/log"
	"github.com/shlande/dmhy-rss/pkg/subscriber/record"
	"testing"
)

func TestCombine(t *testing.T) {
	wp := Combine(record.NewRecord(record.NewJsonKVFromFile("./record/test/test.log")), log.NewLog())
	wp.Added(context.Background(), &parser.Detail{
		TitleInfo: &parser.TitleInfo{Name: "test", Episode: 10},
		Info: &provider.Info{
			Title:      "你好👌",
			MagnetUrl:  "Sooageae",
			TorrentUrl: "dfag",
		},
	})
}
