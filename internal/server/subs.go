package server

import (
	"github.com/shlande/dmhy-rss/pkg/controller/subscriber"
	"github.com/shlande/dmhy-rss/pkg/controller/subscriber/download"
	"github.com/shlande/dmhy-rss/pkg/controller/subscriber/log"
	"github.com/shlande/dmhy-rss/pkg/controller/subscriber/record"
)

type SubscribeConfig struct {
	RecordPath    string
	Aria2BasePath string
	Aria2Path     string
	Aria2Secret   string
}

func BuildSubscriber(config SubscribeConfig) subscriber.Subscriber {
	multi := &subscriber.Multi{}
	multi.Combine(log.NewLog())
	if len(config.RecordPath) != 0 {
		multi.Combine(record.NewRecord(record.NewJsonKVFromFile(config.RecordPath)))
	}
	if len(config.Aria2Path) != 0 {
		multi.Combine(download.New(
			download.NewAria2(config.Aria2Path, config.Aria2Secret),
			download.NewPath("/downloads"),
		))
	}
	return multi
}
