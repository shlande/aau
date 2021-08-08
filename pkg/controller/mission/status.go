package mission

import (
	"time"
)

type Status uint8

func (s Status) String() string {
	return statusString[s]
}

const (
	// Waiting 已经完成当前阶段的更新，等待下一步更新
	Waiting Status = iota
	// Updating 正在等待获取到最新的更新
	Updating
	// Finish 已经完成更新
	Finish
)

type Action uint8

func (a Action) String() string {
	return actionString[a]
}

const (
	_ Action = iota
	UpdateBegin
	UpdateSuccess
	UpdateFailed
	UpdateSkipped
	Terminated
	Finished
)

var (
	actionString = map[Action]string{
		UpdateBegin:   "开始更新",
		UpdateSuccess: "更新成功",
		UpdateFailed:  "更新失败",
		UpdateSkipped: "跳过更新",
		Terminated:    "任务中止",
		Finished:      "任务完成",
	}

	statusString = map[Status]string{
		// Waiting 已经完成当前阶段的更新，等待下一步更新
		Waiting: "等待中",
		// Updating 正在等待获取到最新的更新
		Updating: "正在更新",
		// Finish 已经完成更新
		Finish: "已完成",
	}
)

type Log struct {
	Action
	EmitTime time.Time
	Message  string
}
