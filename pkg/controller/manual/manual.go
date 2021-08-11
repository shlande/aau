package manual

import (
	"context"
	"errors"
	"github.com/shlande/dmhy-rss/pkg/controller/mission"
	"github.com/shlande/dmhy-rss/pkg/controller/store"
	"github.com/shlande/dmhy-rss/pkg/data"
	"github.com/shlande/dmhy-rss/pkg/data/tools"
)

func New(msc chan<- *mission.Mission, clp *tools.CollectionProvider, collectionInterface store.CollectionInterface) *Manual {
	return &Manual{
		msc:                msc,
		CollectionProvider: clp,
		store:              collectionInterface,
	}
}

type Manual struct {
	msc chan<- *mission.Mission
	*tools.CollectionProvider
	store store.CollectionInterface
}

func (m *Manual) CreateMission(collectionId string) error {
	collection := m.CollectionProvider.Get(collectionId)
	if collection == nil {
		return errors.New("无效的collectionId，可能已经过期了🤔")
	}
	ms := mission.NewMission(collection.Animation, collection.Metadata)
	err := ms.Valid()
	if err == nil {
		m.msc <- ms
	}
	return err
}

func (m *Manual) Get(ctx context.Context, collectionId string) (cl *data.Collection, err error) {
	cl = m.CollectionProvider.Get(collectionId)
	if cl == nil {
		cl, err = m.store.Get(collectionId)
	}
	return
}
