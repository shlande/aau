package store

import (
	"github.com/shlande/dmhy-rss/pkg/classify"
	"github.com/shlande/dmhy-rss/pkg/task"
)

// Store 存放收集到的内容,要提供能通过id快速检索collection的方法
type Store interface {
	Save(collection ...*classify.Collection) (string, error)
	Get(id string) (*classify.Collection, error)
	SaveTask(worker ...task.Worker) error
	GetTask(collectionId string) task.Worker
}
