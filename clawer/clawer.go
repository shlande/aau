package clawer

import (
	"context"
	"github.com/mmcdole/gofeed"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
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
	resp, err := http.Get("https://share.dmhy.org/topics/rss/rss.xml?sort_id=2&keyword=" + url.QueryEscape(keywords))
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
func parseTitle(title string) *Detail {
	// 分割所有的书名号
	eps := regexp.MustCompile(`[ |\[]([0-9]{1,3})[ |\]]`).FindString(title)
	episode, _ := strconv.ParseInt(eps[1:3], 10, 8)

	// 这里不直接使用单个词，因为可能会出现很高记录的误匹配
	var lan Language
	gb := regexp.MustCompile(`[\[|【]GB[\]|】]|简体|简中|简繁|簡繁|簡日|简日|[_|\[]CHS[_|\]]`).MatchString(title)
	big := regexp.MustCompile(`[\[|【]BIG5[\]|】]|繁體|繁中|简繁|簡繁|繁日|[_|\[]CHT[_|\]]`).MatchString(title)
	jp := regexp.MustCompile(`簡日|繁日|简日`).MatchString(title)
	if gb {
		lan = lan | GB
	}
	if big {
		lan = lan | BIG5
	}
	if jp {
		lan = lan | JP
	}

	var quality Quality
	p720 := regexp.MustCompile(`720[p|P]`).MatchString(title)
	if p720 {
		quality = P720
	}
	// TODO：还有1080x60fps
	p1080 := regexp.MustCompile(`1080[p|P]`).MatchString(title)
	if p1080 {
		quality = P1080
	}
	// 1440p匹配
	//k2 := regexp.MustCompile(``)

	var sub = Internal
	external := regexp.MustCompile(`外挂|外掛`).MatchString(title)
	if external {
		sub = External
	}
	return &Detail{
		Name:     "",
		Language: lan,
		Quality:  quality,
		Episode:  int(episode),
		SubType:  sub,
	}
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
