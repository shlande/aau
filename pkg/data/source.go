package data

import "github.com/shlande/dmhy-rss/pkg/data/source"

type Source struct {
	*source.Info
	Episode int
	Metadata
}
