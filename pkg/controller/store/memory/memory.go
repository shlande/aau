package memory

import (
	"github.com/shlande/dmhy-rss/pkg/controller/mission"
	"github.com/shlande/dmhy-rss/pkg/controller/store"
	"github.com/shlande/dmhy-rss/pkg/data"
)

func New() memory {
	return memory{
		cls:  make(map[string]*data.Collection),
		logs: make(map[string][]*mission.Log),
		mss:  make(map[string]*mission.Mission),
	}
}

type memory struct {
	cls  map[string]*data.Collection
	logs map[string][]*mission.Log
	mss  map[string]*mission.Mission
}

func (s memory) Resource() store.ResourceInterface {
	return r(s)
}

func (s memory) Log() store.LogInterface {
	return l(s)
}

func (s memory) Animation() store.AnimationInterface {
	return a(s)
}

func (s memory) Pin() store.PinInterface {
	return p(s)
}

func (s memory) Collection() store.CollectionInterface {
	return collection(s)
}
func (s memory) Mission() store.MissionInterface {
	return m(s)
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

type l memory

func (l l) Save(missionId string, log *mission.Log) error {
	mp := l.logs[missionId]
	mp = append(mp, &*log)
	l.logs[missionId] = mp
	return nil
}

func (l l) GetAll(missionId string) (ms []*mission.Log, err error) {
	mp := l.logs[missionId]
	for _, v := range mp {
		ms = append(ms, v)
	}
	return
}

type m memory

func (m m) Save(mission *mission.Mission) error {
	return nil
}

func (m m) Get(id string) (*mission.Mission, error) {
	return nil, nil
}

func (m m) GetAll(active bool) ([]*mission.Mission, error) {
	return nil, nil
}

type a memory

func (a a) Save(animation *data.Animation) error {
	return nil
}

func (a a) Get(id string) (*data.Animation, error) {
	return nil, store.ErrNotFound
}

type p memory

func (p p) Pin(animation *data.Animation) error {
	return nil
}

func (p p) Unpin(animation *data.Animation) error {
	return nil
}

func (p p) IsPin(animation *data.Animation) (bool, error) {
	return false, nil
}

func (p p) Finish(animation *data.Animation) error {
	return nil
}

func (p p) IsFinish(animation *data.Animation) (bool, error) {
	return false, nil
}

func (p p) GetPinned(active interface{}) ([]*data.Animation, error) {
	return nil, nil
}

type r memory

func (r r) Save(collectionId string, source *data.Source) error {
	return nil
}

func (r r) GetAll(collectionId string) ([]*data.Source, error) {
	return nil, nil
}
