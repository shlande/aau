package manual

import (
	"errors"
	"github.com/shlande/dmhy-rss/pkg/controller/mission"
	"github.com/shlande/dmhy-rss/pkg/data/tools"
)

func New(msc chan<- *mission.Mission, clp *tools.CollectionProvider) *Manual {
	return &Manual{
		msc:                msc,
		CollectionProvider: clp,
	}
}

type Manual struct {
	msc chan<- *mission.Mission
	*tools.CollectionProvider
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
