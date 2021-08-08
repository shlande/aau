package mission

import (
	"fmt"
	"github.com/shlande/dmhy-rss/pkg/data"
	"time"
)

// Mission 只作为一个数据源
type Mission struct {
	*data.Collection

	SkipTime int
	Status
}

func (m *Mission) GetNextExpectedResource() int {
	if m.Collection.Latest == m.TotalEpisodes {
		return m.TotalEpisodes
	}
	return m.Latest + 1
}

func (m *Mission) next() {
	if m.Collection.IsFull() {
		m.Status = Finish
	}
	switch m.Status {
	case Updating:
		m.Status = Waiting
	case Waiting:
		m.Status = Updating
	}
}

func (m *Mission) skip() {
	if m.Status == Updating {
		m.next()
	}
	m.SkipTime++
}

func (m *Mission) GetNextUpdateDelay() time.Duration {
	switch m.Status {
	case Updating:
		should := m.LastUpdate.Add(time.Hour)
		// 如果应该更新的时间比现在早，那么延迟为0
		if should.Before(time.Now()) {
			return 0
		}
		return should.Sub(time.Now())
	case Waiting:
		return getNextUpdateTime(m.Animation.AirWeekday, m.LastUpdate)
	default:
		panic("unreachable")
	}
}

func (m *Mission) Next(val interface{}) *Log {
	// val为nil只可能是从waiting转updating
	if val == nil {
		switch m.Status {
		case Waiting:
			m.next()
			return &Log{Action: UpdateBegin, EmitTime: time.Now()}
		default:
			panic("unreachable")
		}
	}
	// 如果是错误，只能是更新错误
	if err, ok := val.(error); ok {
		switch m.Status {
		case Updating:
			return &Log{Action: UpdateFailed, EmitTime: time.Now(), Message: err.Error()}
		default:
			panic("unreachable")
		}

	}
	// 此时val只能为source数组类型
	ss := val.([]*data.Source)
	// 构造message
	var message string
	for _, v := range ss {
		message += "|" + v.Name
	}
	log := &Log{EmitTime: time.Now(), Message: fmt.Sprintf("总共更新了 %v 条数据:%v", len(ss), message)}
	// 执行下一步，判断是否完成
	m.next()
	switch m.Status {
	case Waiting:
		log.Action = UpdateSuccess
	case Finish:
		log.Action = Finished
	}
	return log
}

func getNextUpdateTime(weekday time.Weekday, lastUpdate time.Time) time.Duration {
	// 获取当前是第几周
	var day int
	if lastUpdate.Weekday() >= weekday {
		day = int(weekday + 7 - lastUpdate.Weekday())
	} else {
		day = int(weekday - lastUpdate.Weekday())
	}
	return time.Date(lastUpdate.Year(), lastUpdate.Month(), lastUpdate.Day()+day, 0, 0, 0, 0, time.Local).Sub(lastUpdate)
}
