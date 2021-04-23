package clawer

import (
	"context"
	"github.com/mmcdole/gofeed"
	"io"
	"net/http"
	"net/url"
)

type Option struct {
	Name    string
	Episode int
	Fansub  []string
	Category
	Quality
	SubType
	Language
}

func FindCollectionsByKeywords(ctx context.Context, keywords string) ([]*Collection, error) {
	items, err := find(ctx, keywords)
	if err != nil {
		return nil, err
	}
	return classify(items), nil
}

func FindItems(ctx context.Context, option *Option) ([]*Item, error) {
	return nil, nil
}

// 单纯的获取rss并解析成item
func find(ctx context.Context, keywords string) ([]*Item, error) {
	resp, err := http.Get("https://share.dmhy.org/rss/rss.html?keyword=" + url.QueryEscape(keywords))
	if err != nil {
		return nil, err
	}
	return parse(resp.Body)
}

// 把item归类成collection
func classify(items []*Item) []*Collection {
	return nil
}

// 根据option筛选出指定的item
func sort(items []*Item, option *Option) []*Item {
	return nil
}

// 尝试解析头部，获取到一些信息
func parseTitle(title string) Detail {
	return Detail{}
}

func parseCategory(category string) Category {
	if category == "動畫" {
		return Animate
	}
	if category == "季度全集" {
		return FullSession
	}
	return Unknown
}

func parse(reader io.Reader) (rcds []*Item, err error) {
	fp := gofeed.NewParser()
	feed, err := fp.Parse(reader)
	if err != nil {
		return nil, err
	}
	for _, i := range feed.Items {
		var fansub []string
		var category, magnet string
		for _, i := range i.Authors {
			fansub = append(fansub, i.Name)
		}
		if len(i.Categories) > 0 {
			category = i.Categories[0]
		}
		if len(i.Enclosures) > 0 {
			magnet = i.Enclosures[0].URL
		}
		rcds = append(rcds, &Item{
			Fansub:    fansub,
			Title:     i.Title,
			Category:  parseCategory(category),
			MagnetUrl: magnet,
			Link:      i.Link,
			PubDate:   i.PublishedParsed,
		})
	}
	return rcds, nil
}
