package dmhy

import (
	"context"
	"github.com/mmcdole/gofeed"
	pvd "github.com/shlande/dmhy-rss/pkg/provider"
	"net/http"
	"net/url"
)

func NewProvider() *provider {
	return &provider{fp: gofeed.NewParser()}
}

type provider struct {
	fp *gofeed.Parser
}

func (p provider) Keywords(ctx context.Context, keywords string) (infos []*pvd.Info, err error) {
	resp, err := http.Get("https://www.dmhy.org/topics/rss/rss.xml?sort_id=2&keyword=" + url.QueryEscape(keywords))
	if err != nil {
		return nil, err
	}
	feed, err := p.fp.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	for _, i := range feed.Items {
		info := &pvd.Info{
			Title: i.Title,
		}
		// TODO:多翻译组情况
		for _, i := range i.Authors {
			info.Fansub = append(info.Fansub, i.Name)
		}
		if len(i.Categories) > 0 {
			info.RawCategory = i.Categories[0]
		}
		if len(i.Enclosures) > 0 {
			info.MagnetUrl = i.Enclosures[0].URL
		}
		info.CreateTime = i.PublishedParsed
		infos = append(infos, info)
	}
	return
}
