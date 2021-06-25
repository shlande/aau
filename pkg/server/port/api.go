package port

import (
	"context"
	"github.com/shlande/dmhy-rss/pkg/classify"
	"github.com/shlande/dmhy-rss/pkg/worker"
	"time"
)

type Api interface {
	// Search 通过关键词去查找信息
	Search(ctx context.Context, words string) []*CollectionSummary
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

func NewCollectionSummary(collection *classify.Collection) *CollectionSummary {
	if collection == nil {
		return nil
	}
	var episode = make([]*EpisodeSummary, 0, len(collection.Items))
	for _, v := range collection.Items {
		episode = append(episode, &EpisodeSummary{
			Title:   v.Title,
			Episode: v.Episode,
		})
	}

	return &CollectionSummary{
		Id:         collection.Id(),
		Name:       collection.Name,
		Fansub:     collection.Fansub,
		Quality:    collection.Quality.String(),
		Category:   collection.Category.String(),
		SubType:    collection.SubType.String(),
		Language:   collection.Language.String(),
		Latest:     collection.Latest,
		LastUpdate: collection.LastUpdate,
		Episodes:   episode,
	}
}

type CollectionSummary struct {
	Id       string
	Name     string
	Fansub   []string
	Quality  string
	Category string
	SubType  string
	Language string
	// Collection 的信息
	Latest     int
	LastUpdate time.Time
	Episodes   []*EpisodeSummary
}

type EpisodeSummary struct {
	Title   string
	Episode int
}
