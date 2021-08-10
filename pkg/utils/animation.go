package utils

import "time"

func GetSessionTime(year int, session int) (start time.Time, end time.Time) {
	// 就一般理性而言，番剧都会在某一个月集中开始发布，因此在计算开始的时候，把范围往前面挪动会更加合适
	if session == 0 {
		year = year - 1
	}
	return time.Date(year, time.Month((11+session*3)%12+1), 0, 0, 0, 0, 0, time.Local),
		time.Date(year+session/3, time.Month((11+(session+1)*3)%12+1), 0, 0, 0, 0, 0, time.Local)
}

func GetSession(airTime time.Time) (year, session int) {
	year = airTime.Year()
	for i := 0; i < 3; i++ {
		b, e := GetSessionTime(year, i)
		if airTime.After(b) && airTime.Before(e) {
			return year, i
		}
	}
	panic("unreachable")
}
