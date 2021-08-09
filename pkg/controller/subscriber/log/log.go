package log

import (
	"github.com/shlande/dmhy-rss/pkg/data"
	"github.com/sirupsen/logrus"
)

func NewLog() *Log {
	return &Log{}
}

type Log struct{}

func (l *Log) Created(collection *data.Collection) {
	logrus.Infoln("新的监控添加 id:", collection.Id(), "name:", collection.Name)
}

func (l *Log) Added(detail *data.Source) {
	logrus.Infoln("新内容添加", detail.Name, detail.Episode)
}
