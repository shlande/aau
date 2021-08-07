package data

import "time"

// Animation 这里介绍了一个番剧的详细信息
type Animation struct {
	Name       string
	Translated string
	Summary    string
	AirDate    time.Time
	AirWeekday time.Weekday
	// 在weekday 的基础上加AirTime的时间
	AirTime time.Duration
	// 总共集数
	TotalEpisodes int
	// 分类
	Category string
}
