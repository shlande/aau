package store

import (
	"github.com/shlande/dmhy-rss/pkg/classify"
	"github.com/shlande/dmhy-rss/pkg/worker"
)

// Store 存放收集到的内容,要提供能通过id快速检索collection的方法
type Store interface {
	Save(collection ...*classify.Collection) error
	Get(id string) (*classify.Collection, error)
	SaveWorker(worker ...*worker.Worker) error
	GetWorker(collectionId string) (*worker.RecoverHelper, error)
	ListWorker() ([]*worker.RecoverHelper, error)
}
