package subscriber

import (
	"context"
	"github.com/shlande/dmhy-rss/pkg/parser"
)

func Combine(subs ...Subscriber) *wrapper {
	return &wrapper{subs: subs}
}

type wrapper struct {
	subs []Subscriber
}

func (w *wrapper) Added(ctx context.Context, detail *parser.Detail) {
	for _, v := range w.subs {
		v.Added(ctx, detail)
	}
}

func (w *wrapper) Combine(subs ...Subscriber) *wrapper {
	w.subs = append(w.subs, subs...)
	return w
}
