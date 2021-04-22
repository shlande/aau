package main

import (
	"context"
	"fmt"
	"github.com/mmcdole/gofeed"
	"io"
	"time"
)

type (
	Category string
	Quality  string
	Language string
	SubType  string
)

const (
	Unknown     Category = "未知"
	FullSession Category = "季度全集"
	Animate     Category = "动画"
	Movie       Category = "电影"

	P720  Quality = "720p"
	P1080 Quality = "1080p"
	K2    Quality = "2k"

	GB  Language = "简体"
	BIG Language = "繁体"
	JP  Language = "日语"

	Internal SubType = "内嵌"
	External SubType = "外挂"
)

func ParseCategory(category string) Category {
	fmt.Println(category)
	if category == "動畫" {
		return Animate
	}
	if category == "季度全集" {
		return FullSession
	}
	return Unknown
}

// Record 爬到的单个记录
type Record struct {
	// 字幕组
	Fansub []string
	// 原始标题
	Title string
	// 类型
	Category
	// 详细信息
	Detail
	// 磁力链接地址
	MagnetUrl string
	// 原文链接
	Link string
	// 发布时间
	PubDate *time.Time
}

type Detail struct {
	// 简体还是繁体
	Languages []Language
	// 画质
	Quality
	// 集数
	Episode int
	// 字幕类型
	SubType
}

func FindByKeywords(ctx context.Context, keyword string) {

	//feed,err := fp.Parse()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(feed)
}

func Parse(reader io.Reader) (rcds []*Record, err error) {
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
		rcds = append(rcds, &Record{
			Fansub:    fansub,
			Title:     i.Title,
			Category:  ParseCategory(category),
			MagnetUrl: magnet,
			Link:      i.Link,
			PubDate:   i.PublishedParsed,
		})
	}
	return rcds, nil
}

// 尝试解析头部，获取到一些信息
func parseDetail(title string) Detail {
	return Detail{}
}
