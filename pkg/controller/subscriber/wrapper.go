package subscriber

import (
	"context"
	"github.com/shlande/dmhy-rss/pkg/classify"
	"github.com/shlande/dmhy-rss/pkg/parser"
)

func Combine(subs ...Subscriber) *Multi {
	return &Multi{subs: subs}
}

type Multi struct {
	subs []Subscriber
}

func (w *Multi) Created(ctx context.Context, collection *classify.Collection) {
	for _, v := range w.subs {
		v.Created(ctx, collection)
	}
}

func (w *Multi) Added(ctx context.Context, detail *parser.Detail) {
	for _, v := range w.subs {
		v.Added(ctx, detail)
	}
}

func (w *Multi) Combine(subs ...Subscriber) *Multi {
	w.subs = append(w.subs, subs...)
	return w
}
