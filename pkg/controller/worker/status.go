package worker

import (
	"time"
)

type Status uint8

func (s Status) String() string {
	return statusString[s]
}

//func (s Status) MarshalJSON() ([]byte, error) {
//	return []byte(s.String()),nil
//}

const (
	// Wait 已经完成当前阶段的更新，等待下一步更新
	Wait Status = iota
	// Update 正在等待获取到最新的更新
	Update
	// Finish 已经完成更新
	Finish
)

type Action uint8

func (a Action) String() string {
	return actionString[a]
}

//func (a Action) MarshalJSON() ([]byte, error) {
//	return []byte(a.String()),nil
//}

const (
	UpdateFinish Action = iota
	UpdateFail
	UpdateSkip
	Terminate
)

var (
	actionString = map[Action]string{
		UpdateFinish: "更新完成",
		UpdateFail:   "更新失败",
		UpdateSkip:   "跳过更新",
		Terminate:    "用户删除",
	}

	statusString = map[Status]string{
		// Wait 已经完成当前阶段的更新，等待下一步更新
		Wait: "等待中",
		// Update 正在等待获取到最新的更新
		Update: "正在更新",
		// Finish 已经完成更新
		Finish: "已完成",
	}
)

func newLog(action Action, message string) *Log {
	return &Log{
		Action:   action,
		EmitTime: time.Now(),
		Message:  message,
	}
}

type Log struct {
	Action
	EmitTime time.Time
	Message  string
}
