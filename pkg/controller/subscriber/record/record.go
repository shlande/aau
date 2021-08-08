package record

import (
	"context"
	"log"
	"strconv"
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

func (r *Recorder) Created(_ context.Context, collection *classify.Collection) {
	for _, v := range collection.Items {
		r.Added(nil, v)
	}
}

func (r *Recorder) Added(_ context.Context, detail *parser.Detail) {
	err := r.Set(detail.Name+"-"+strconv.Itoa(detail.Episode), newRecord(detail))
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
