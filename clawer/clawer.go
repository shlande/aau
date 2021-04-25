package clawer

import (
	"context"
	"github.com/mmcdole/gofeed"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
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
	var matches []*Item
	// var fansub,episode,quality,subtype,lan bool
	for _, i := range items {
		// 暂时不进行name匹配
		for _, fs := range option.Fansub {
			for _, ifs := range i.Fansub {
				if fs == ifs {
					goto EPISODE
				}
			}
		}
		continue
	EPISODE:
		if option.Episode != 0 && i.Episode != option.Episode {
			continue
		}
		if option.Quality != UnknownQuality && option.Quality != i.Quality {
			continue
		}
		if option.SubType != UnknownSubType && i.SubType != option.SubType {
			continue
		}
		if option.Category != UnknownCategory && i.Category != option.Category {
			continue
		}
		// 暂时只支持And匹配
		if option.Language != UnknownLanguage && i.Language&option.Language != option.Language {
			continue
		}
		matches = append(matches, i)
	}
	return matches
}

// 尝试解析头部，获取到一些信息
func parseTitle(title string) *Detail {
	name := parseName(title, false)
	// 分割所有的书名号
	eps := regexp.MustCompile(`[ |\[]([0-9]{1,3})[ |\]]`).FindString(title)
	if len(eps) == 4 {
		eps = eps[1:3]
	}
	episode, _ := strconv.ParseInt(eps, 10, 8)

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
		Name:     name,
		Language: lan,
		Quality:  quality,
		Episode:  int(episode),
		SubType:  sub,
	}
}

func parseName(title string, keep bool) string {
	var cname, left string
	if !keep {
		// 查找第一个右括号书名号位置,这个位置一般是字幕组名称截止的位置
		begin := regexp.MustCompile(`\]|】`).FindStringIndex(title)
		// 删除掉字幕组
		title = title[begin[1]:]
	}

	bw := regexp.MustCompile(`[\[|【]`).FindStringIndex(title)
	// 如果内容是以括号开始的，那么name应该就是这个括号内的内容
	if bw[0] == 0 {
		ew := regexp.MustCompile(`\]|】`).FindStringIndex(title)
		cname = title[bw[1]:ew[0]]
		left = title[ew[1]:]
	} else {
		// 获取没被书名号括起来的内容
		cname = title[:bw[0]]
		left = title[bw[1]:]
	}
	// 判断是否是有效的名称,有些字幕组会加一月新番等多余的名字
	if regexp.MustCompile(`新番`).MatchString(cname) {
		// 这类型的暂时无法匹配
		return parseName(left, true)
	}
	// 对有效名称进行最后的处理
	// 通过反斜杠来分割名称中中英文
	{
		if raw := strings.Split(cname, "\\"); len(raw) != 1 {
			cname = raw[0]
		}
		if raw := strings.Split(cname, "/"); len(raw) != 1 {
			cname = raw[0]
		}
		if raw := strings.Split(cname, "_"); len(raw) != 1 {
			cname = raw[0]
		}
		// 删除仅限港澳台地区的字样
		{
			cname = regexp.MustCompile(`[\[|【|(|（](仅|僅)限港澳台(地区|地區)*[\)|\]|】|）]`).ReplaceAllString(cname, "")
		}
		// 删除两侧的空白符号
		cname = regexp.MustCompile(`^ *`).ReplaceAllString(cname, "")
		cname = regexp.MustCompile(` *$`).ReplaceAllString(cname, "")
	}
	return cname
}

func parseCategory(category string) Category {
	if category == "動畫" {
		return Animate
	}
	if category == "季度全集" {
		return FullSession
	}
	return UnknownCategory
}

func parse(reader io.Reader) (rcds []*Item, err error) {
	fp := gofeed.NewParser()
	feed, err := fp.Parse(reader)
	if err != nil {
		return nil, err
	}
	for _, i := range feed.Items {
		var fansub []string
		var magnet string
		var category Category
		for _, i := range i.Authors {
			fansub = append(fansub, i.Name)
		}
		if len(i.Categories) > 0 {
			category = parseCategory(i.Categories[0])
		}
		if len(i.Enclosures) > 0 {
			magnet = i.Enclosures[0].URL
		}
		detail := parseTitle(i.Title)
		if category == Animate && detail.Episode == 0 {
			log.Println("解析错误，无法获取集数：" + i.Title)
		}
		rcds = append(rcds, &Item{
			Fansub:    fansub,
			Title:     i.Title,
			Category:  category,
			MagnetUrl: magnet,
			Detail:    detail,
			Link:      i.Link,
			PubDate:   i.PublishedParsed,
		})
	}
	return rcds, nil
}
