package mission

import (
	"fmt"
	"github.com/shlande/dmhy-rss/pkg/data"
	"time"
)

func NewMission(animation *data.Animation, metadata data.Metadata) *Mission {
	return &Mission{
		Collection: data.NewCollection(animation, metadata),
		SkipTime:   0,
		Status:     Updating,
	}
}

// Mission 只作为一个数据源
type Mission struct {
	*data.Collection

	LastUpdate time.Time
	SkipTime   int
	Status
}

func (m *Mission) GetExpectedLatest() int {
	if m.Collection.Latest == m.TotalEpisodes {
		return m.TotalEpisodes
	}
	switch m.Status {
	case Updating:
		return m.Latest + 1
	default:
		return m.Latest
	}
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

// TODO：完成skip功能
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
		anm := m.Animation
		return getNextUpdateTime(anm.AirDate, anm.AirBreak, anm.AirBreak*time.Duration(m.SkipTime), m.LastUpdate)
	default:
		panic("unreachable")
	}
}

func (m *Mission) Next(val interface{}) *Log {
	m.LastUpdate = time.Now()
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

func getNextUpdateTime(airDate time.Time, airBreak, skip time.Duration, lastUpdate time.Time) time.Duration {
	// 获取当前是第几周
	// 获取当前时间
	if airBreak == 0 {
		return time.Hour
	}

	airDate.Add(skip)

	for airDate.Before(lastUpdate) {
		airDate = airDate.Add(airBreak)
	}

	return airDate.Sub(lastUpdate)
}
