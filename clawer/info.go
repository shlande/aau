package clawer

import "time"

// Detail 一些详细的信息，需要从标题中提取出来的信息
type Detail struct {
	// 番剧名称
	Name string
	// 简体还是繁体
	Language
	// 画质
	Quality
	// 集数
	Episode int
	// 字幕类型
	SubType
}

// Item 爬到的单个记录
type Item struct {
	// 字幕组
	Fansub []string
	// 原始标题
	Title string
	// 类型
	Category
	// 详细信息
	*Detail
	// 磁力链接地址
	MagnetUrl string
	// 原文链接
	Link string
	// 发布时间
	PubDate *time.Time
}

// Collection 包含了单个种类的所有record
type Collection struct {
	// 重复保存了的内容
	Name   string
	Fansub []string
	Quality
	Category
	SubType
	Language
	// Collection 的信息
	Latest     int
	LastUpdate *time.Time

	Items []*Item
}

// TODO:判断是否有重复
func (c *Collection) AddItem(item *Item) {
	c.Items = append(c.Items, item)
}
