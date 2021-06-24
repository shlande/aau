package log

import "github.com/sirupsen/logrus"

func SetLogger(lg *logrus.Logger) {
	logger = lg
}

var logger *logrus.Logger

func NewEntry(name string) *logrus.Entry {
	if logger == nil {
		logger = logrus.New()
	}
	return logrus.WithField("component", name)
}
