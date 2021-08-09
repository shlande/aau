package main

import (
	"github.com/shlande/dmhy-rss/pkg/controller/subscriber"
	"github.com/shlande/dmhy-rss/pkg/controller/subscriber/log"
	"github.com/shlande/dmhy-rss/pkg/controller/subscriber/record"
)

type SubscribeConfig struct {
	RecordPath string
}

func buildSubscribe(config SubscribeConfig) subscriber.Subscriber {
	multi := &subscriber.Multi{}
	multi.Combine(log.NewLog())
	if len(config.RecordPath) != 0 {
		multi.Combine(record.NewRecord(record.NewJsonKVFromFile(config.RecordPath)))
	}
	return multi
}
