package parser

import (
	"github.com/shlande/dmhy-rss/pkg/data"
)

// Parser is a tools used to parse useful info from common
type Parser interface {
	Parse(name string) (*Result, error)
}

func ParseCategory(category string) data.Type {
	if category == "動畫" {
		return data.Animate
	}
	if category == "季度全集" {
		return data.FullSession
	}
	return data.UnknownCategory
}

type Result struct {
	Name string
	// 分类
	data.Metadata
	// 集数
	Episode int
}
