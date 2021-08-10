package subscriber

import (
	"github.com/shlande/dmhy-rss/pkg/data"
)

func Combine(subs ...Subscriber) *Multi {
	return &Multi{subs: subs}
}

type Multi struct {
	subs []Subscriber
}

func (w *Multi) Created(collection *data.Collection) {
	for _, v := range w.subs {
		v.Created(collection)
	}
}

func (w *Multi) Added(detail *data.Source, collection *data.Collection) {
	for _, v := range w.subs {
		v.Added(detail, collection)
	}
}

func (w *Multi) Combine(subs ...Subscriber) *Multi {
	w.subs = append(w.subs, subs...)
	return w
}
