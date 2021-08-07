package dmhy

import (
	"context"
	"github.com/mmcdole/gofeed"
	"github.com/shlande/dmhy-rss/pkg/data/source"
	"net/http"
	"net/url"
)

func NewProvider() *provider {
	return &provider{fp: gofeed.NewParser()}
}

type provider struct {
	fp *gofeed.Parser
}

func (p provider) Keywords(ctx context.Context, keywords string) (infos []*source.Info, err error) {
	resp, err := http.Get("https://www.dmhy.org/topics/rss/rss.xml?sort_id=2&keyword=" + url.QueryEscape(keywords))
	if err != nil {
		return nil, err
	}
	feed, err := p.fp.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	for _, i := range feed.Items {
		info := &source.Info{
			Name: i.Title,
		}
		if len(i.Enclosures) > 0 {
			info.MagnetUrl = i.Enclosures[0].URL
		}
		info.CreateTime = *i.PublishedParsed
		infos = append(infos, info)
	}
	return
}
