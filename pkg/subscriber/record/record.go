package record

import (
	"context"
	"github.com/shlande/dmhy-rss/pkg/parser"
	"log"
	"time"
)

type KVSetter interface {
	Set(key string, value interface{}) error
}

func NewRecord(setter KVSetter) *Recorder {
	return &Recorder{KVSetter: setter}
}

type Recorder struct {
	KVSetter
}

func (r *Recorder) Added(ctx context.Context, detail *parser.Detail) {
	err := r.Set(detail.Name, newRecord(detail))
	if err != nil {
		log.Println("无法记录数据：" + err.Error())
	}
}

func newRecord(detail *parser.Detail) *record {
	url := detail.MagnetUrl
	if len(url) == 0 {
		url = detail.TorrentUrl
	}
	return &record{
		Name:    detail.Name,
		Episode: detail.Episode,
		Url:     url,
		Time:    time.Now(),
	}
}

type record struct {
	Name    string
	Episode int
	Url     string
	Time    time.Time
}
