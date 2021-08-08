package store

import (
	"github.com/shlande/dmhy-rss/pkg/controller/mission"
	"github.com/shlande/dmhy-rss/pkg/data"
)

// Interface 门面模式
type Interface interface {
	Collection() CollectionInterface
	Mission() MissionInterface
	Log() LogInterface
	Animation() AnimationInterface
	Pin() PinInterface
}

// CollectionInterface 存放收集到的内容,要提供能通过id快速检索collection的方法
type CollectionInterface interface {
	Save(collection *data.Collection) error
	Get(id string) (*data.Collection, error)
	GetAll() ([]*data.Collection, error)
}

type MissionInterface interface {
	Save(mission *mission.Mission) error
	Get(id string) (*mission.Mission, error)
	GetAll(active bool) ([]*mission.Mission, error)
}

type LogInterface interface {
	Save(missionId string, log *mission.Log) error
	GetAll(missionId string) ([]*mission.Log, error)
}

type AnimationInterface interface {
	Save(animation *data.Animation) error
	Get(id string) (*data.Animation, error)
}

type PinInterface interface {
	AnimationInterface

	Pin(animation *data.Animation) error
	Unpin(animation *data.Animation) error
	IsPin(animation *data.Animation) (bool, error)

	Finish(animation *data.Animation) error
	IsFinish(animation *data.Animation) (bool, error)

	// GetPinned 如果是nil，那么就返回全部，如果是true，则返回活跃的
	GetPinned(active interface{}) ([]*data.Animation, error)
}
