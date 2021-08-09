package record

import (
	"github.com/shlande/dmhy-rss/pkg/data"
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

func (r *Recorder) Created(collection *data.Collection) {
	for _, v := range collection.Items {
		r.Added(v)
	}
}

func (r *Recorder) Added(detail *data.Source) {
	err := r.Set(detail.Name+"-"+strconv.Itoa(detail.Episode), newRecord(detail))
	if err != nil {
		log.Println("无法记录数据：" + err.Error())
	}
}

func newRecord(source *data.Source) *record {
	url := source.MagnetUrl
	if len(url) == 0 {
		url = source.TorrentUrl
	}
	return &record{
		Name:    source.Name,
		Episode: source.Episode,
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
