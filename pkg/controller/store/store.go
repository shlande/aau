package store

import (
	"github.com/shlande/dmhy-rss/pkg/classify"
	worker2 "github.com/shlande/dmhy-rss/pkg/controller/worker"
)

// Store 存放收集到的内容,要提供能通过id快速检索collection的方法
type Store interface {
	Save(collection ...*classify.Collection) error
	Get(id string) (*classify.Collection, error)
	SaveWorker(worker ...*worker2.Worker) error
	GetWorker(collectionId string) (*worker2.RecoverHelper, error)
	ListWorker() ([]*worker2.RecoverHelper, error)
}
