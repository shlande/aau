package memory

import (
	"github.com/shlande/dmhy-rss/pkg/controller/store"
	"github.com/shlande/dmhy-rss/pkg/data"
)

func New() memory {
	return memory{cls: make(map[string]*data.Collection)}
}

type memory struct {
	cls map[string]*data.Collection
}

func (s memory) Collection() store.CollectionInterface {
	return collection(s)
}
func (s memory) Mission() store.MissionInterface {
	panic("implement me")
}

type collection memory

func (c collection) Save(collection *data.Collection) error {
	c.cls[collection.Id()] = collection
	return nil
}

func (c collection) Get(id string) (*data.Collection, error) {
	cl, has := c.cls[id]
	if !has {
		return nil, store.ErrNotFound
	}
	return cl, nil
}

func (c collection) GetAll() ([]*data.Collection, error) {
	var cls = make([]*data.Collection, 0, len(c.cls))
	for _, cl := range c.cls {
		cls = append(cls, &*cl)
	}
	return cls, nil
}
