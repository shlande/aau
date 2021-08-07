package memory

import (
	"github.com/shlande/dmhy-rss/pkg/classify"
	store2 "github.com/shlande/dmhy-rss/pkg/controller/store"
	worker2 "github.com/shlande/dmhy-rss/pkg/controller/worker"
)

func New() *Store {
	return &Store{cls: make(map[string]*classify.Collection)}
}

type Store struct {
	cls map[string]*classify.Collection
}

func (s *Store) SaveWorker(worker ...*worker2.Worker) error {
	panic("memory不支持保存worker信息")
}

func (s *Store) GetWorker(collectionId string) (*worker2.RecoverHelper, error) {
	panic("memory不支持保存worker信息")
}

func (s *Store) ListWorker() ([]*worker2.RecoverHelper, error) {
	panic("memory不支持保存worker信息")
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
		return nil, store2.ErrNotFound
	}
	return cl, nil
}
