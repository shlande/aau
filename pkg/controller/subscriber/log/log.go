package log

import (
	"context"
	"github.com/sirupsen/logrus"
)

func NewLog() *Log {
	return &Log{Entry: log.NewEntry("record")}
}

type Log struct {
	*logrus.Entry
}

func (l *Log) Created(_ context.Context, collection *classify.Collection) {
	l.Logger.Infoln("新的监控添加 id:", collection.Id(), "name:", collection.Name)
}

func (l *Log) Added(_ context.Context, detail *parser.Detail) {
	l.Logger.Infoln("新内容添加", detail.Name, detail.Episode)
}
