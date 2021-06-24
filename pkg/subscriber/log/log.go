package log

import (
	"context"
	"github.com/shlande/dmhy-rss/pkg/log"
	"github.com/shlande/dmhy-rss/pkg/parser"
	"github.com/sirupsen/logrus"
)

func NewLog() *Log {
	return &Log{Entry: log.NewEntry("record")}
}

type Log struct {
	*logrus.Entry
}

func (l *Log) Added(_ context.Context, detail *parser.Detail) {
	l.Logger.Infoln("新内容添加", detail.Name, detail.Episode)
}
