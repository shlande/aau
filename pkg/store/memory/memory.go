package memory

import (
	"github.com/shlande/dmhy-rss/pkg/classify"
	"github.com/shlande/dmhy-rss/pkg/store"
	"github.com/shlande/dmhy-rss/pkg/worker"
)

func New() *Store {
	return &Store{cls: make(map[string]*classify.Collection)}
}

type Store struct {
	cls map[string]*classify.Collection
}

func (s *Store) Save(collection ...*classify.Collection) error {
	for _, v := range collection {
		s.cls[v.Id()] = v
	}
	return nil
}

func (s *Store) Get(id string) (*classify.Collection, error) {
	cl, has := s.cls[id]
	if !has {
		return nil, store.ErrNotFound
	}
	return cl, nil
}

func (s *Store) SaveWorker(worker ...worker.Worker) error {
	panic("memory不支持保存worker信息")
}

func (s *Store) GetWorker(collectionId string) (worker.Worker, error) {
	panic("memory不支持保存worker信息")
}
