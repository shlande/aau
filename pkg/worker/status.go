package worker

import (
	"time"
)

type Status uint8

const (
	// Wait 已经完成当前阶段的更新，等待下一步更新
	Wait Status = iota
	// Update 正在等待获取到最新的更新
	Update
	// Download 正在下载
	Download
	// Finish 已经完成更新
	Finish
)

type Action uint8

const (
	UpdateFinish Action = iota
	UpdateFail
	UpdateSkip
	DownloadFinish
	DownloadCancel
	DownloadTimeout
	Terminate
)

type Log struct {
	Action
	EmitTime time.Time
	Message  string
}
