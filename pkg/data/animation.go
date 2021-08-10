package data

import (
	"errors"
	"time"
)

// Animation 这里介绍了一个番剧的详细信息
type Animation struct {
	Id         string
	Name       string
	Translated string
	// Keywords 指定查找时使用的关键词
	// 默认名称有些时候会无法查找到资源，此时需要用户手动更改关键词
	Keywords string
	Summary  string
	AirDate  time.Time
	// 更新等待周期，可能是month，day，once
	AirBreak time.Duration
	// 总共集数
	TotalEpisodes int
	// 分类
	Category string
}

func (a *Animation) GetKeywords() string {
	if len(a.Keywords) == 0 {
		return a.Translated
	}
	return a.Keywords
}

func (a *Animation) SetKeywords(keywords string) error {
	if len(keywords) < 3 {
		return errors.New("关键词长度不能小于3")
	}
	a.Keywords = keywords
	return nil
}
