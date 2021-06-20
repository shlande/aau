package parser

import "github.com/shlande/dmhy-rss/pkg/provider"

// Parser is a tools used to parse useful info from common
type Parser interface {
	ParseTitle(title string) (*TitleInfo, error)
	Parse(info *provider.Info) (*Detail, error)
}

func ParseCategory(category string) Category {
	if category == "動畫" {
		return Animate
	}
	if category == "季度全集" {
		return FullSession
	}
	return UnknownCategory
}

type TitleInfo struct {
	Name string
	// 分类
	Category
	// 简体还是繁体
	Language
	// 画质
	Quality
	// 集数
	Episode int
	// 字幕类型
	SubType
}

// Detail 一些详细的信息，需要从标题中提取出来的信息
type Detail struct {
	*TitleInfo
	*provider.Info
}
