package state

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

func New(store2 store2.Store, wks func(id string) (*worker2.Misson, error)) *Updater {
	return &Updater{
		Entry: log.NewEntry("updater"),
		Store: store2,
		wks:   wks,
	}
}

// Updater 负责监控worker的变化，同步更新内容
type Updater struct {
	*logrus.Entry
	store2.Store
	wks func(id string) (*worker2.Misson, error)
}

func (u *Updater) Created(ctx context.Context, collection *classify.Collection) {
	u.save(collection.Id())
}

func (u *Updater) save(id string) {
	w, err := u.wks(id)
	if err != nil {
		time.Sleep(time.Second)
		w, err = u.wks(id)
		if err != nil {
			u.Errorln("无法找到worker id:", id)
			return
		}
	}
	err = u.SaveWorker(w)
	if err != nil {
		u.Errorln("保存worker失败:", err)
	}
	u.Infoln("保存worker信息 id:", w.Id)
}

func (u *Updater) Added(ctx context.Context, detail *parser.Detail) {
	u.save(classify.NewCollection(detail).Id())
}
