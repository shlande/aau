package api

import (
	"github.com/shlande/dmhy-rss/pkg/classify"
	"github.com/shlande/dmhy-rss/pkg/task"
	"time"
)

type Api interface {
	// Keywords 通过关键词去查找信息
	Keywords(words string) []*classify.Collection

	// Watch 添加collection到监控列表中，进行同步更新
	Watch(collectionId string, updateTime *time.Time) error
	// UnWatch 取消更新
	UnWatch(collectionId string)
	// WatchList 列出所有监控的collection
	// TODO: impl
	WatchList()

	// Log 获取更新信息
	Log(collectionId string) []*task.Log
}
