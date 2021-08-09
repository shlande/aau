package data

import "time"

// Animation 这里介绍了一个番剧的详细信息
type Animation struct {
	Id         string
	Name       string
	Translated string
	Summary    string
	AirDate    time.Time
	// 更新等待周期，可能是month，day，once
	AirBreak time.Duration
	// 总共集数
	TotalEpisodes int
	// 分类
	Category string
}
