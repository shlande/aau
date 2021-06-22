package task

import (
	"context"
)

// Worker 负责监控更新collection
// 当一个collection需要更新时发送消息
type Worker interface {
	Id() string
	Run(ctx context.Context)
	// Log 输出日志
	Log() []*Log
	Stop()
	// 结束并删除任务
	Terminate()
}
