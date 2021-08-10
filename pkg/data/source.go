package data

import "github.com/shlande/dmhy-rss/pkg/data/source"

type Source struct {
	Name string
	*source.Info
	Episode int
	Metadata
}

func (s Source) GetDownloadUrl() string {
	if len(s.MagnetUrl) != 0 {
		return s.MagnetUrl
	}
	return s.TorrentUrl
}
