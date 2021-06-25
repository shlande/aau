package port

import (
	"context"
	"github.com/shlande/dmhy-rss/pkg/classify"
	"github.com/shlande/dmhy-rss/pkg/worker"
	"time"
)

type Api interface {
	// Search 通过关键词去查找信息
	Search(ctx context.Context, words string) []*classify.Collection
	// Watch 添加collection到监控列表中，进行同步更新
	Watch(collectionId string, updateTime time.Weekday) error
	// GetCollection 通过id查找collection
	GetCollection(collectinoId string) *classify.Collection
	// GetWorker 获取worker信息
	GetWorker(workerId string) *WorkerInfo
	// UnWatch 取消更新
	UnWatch(collectionId string) error
	// WatchList 列出所有监控的collection
	// TODO: impl
	WatchList() []*WorkerInfo
}

type WorkerInfo struct {
	Id string
	worker.Status
	UpdateTime time.Weekday
	Logs       []*worker.Log
}

func NewWorkerInfo(worker *worker.Worker) *WorkerInfo {
	if worker == nil {
		return nil
	}
	return &WorkerInfo{
		Id:         worker.Id,
		Status:     worker.Status,
		UpdateTime: worker.UpdateTime,
		Logs:       worker.Log(),
	}
}
