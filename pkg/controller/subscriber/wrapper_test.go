package subscriber

import (
	"context"
	log2 "github.com/shlande/dmhy-rss/pkg/controller/subscriber/log"
	record2 "github.com/shlande/dmhy-rss/pkg/controller/subscriber/record"
	source2 "github.com/shlande/dmhy-rss/pkg/data/source"
	"testing"
)

func TestCombine(t *testing.T) {
	wp := Combine(record2.NewRecord(record2.NewJsonKVFromFile("./record/test/test.log")), log2.NewLog())
	wp.Added(context.Background(), &parser.Detail{
		TitleInfo: &parser.TitleInfo{Name: "test", Episode: 10},
		Info: &source2.Info{
			Title:      "你好👌",
			MagnetUrl:  "Sooageae",
			TorrentUrl: "dfag",
		},
	})
}
