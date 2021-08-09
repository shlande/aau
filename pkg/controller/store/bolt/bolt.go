package bolt

import (
	"errors"
	"github.com/boltdb/bolt"
	"github.com/shlande/dmhy-rss/pkg/controller/mission"
	"github.com/shlande/dmhy-rss/pkg/controller/store"
	"github.com/shlande/dmhy-rss/pkg/data"
	"github.com/sirupsen/logrus"
	"strconv"
)

var (
	animationKey         = []byte("animation")
	collectionSummaryKey = []byte("collectionSummary")
	resourceKey          = []byte("resource")
	missionSummaryKey    = []byte("mission")
	logKey               = []byte("log")
	pinKey               = []byte("pin")
)

func New(path string) *Bolt {
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1})
	tx, err := db.Begin(true)
	defer func() {
		if err == nil {
			err = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()
	if err != nil {
		logrus.Panicln(err)
		panic(err)
	}

	_, err = tx.CreateBucketIfNotExists(animationKey)
	if err != nil {
		panic(err)
	}
	_, err = tx.CreateBucketIfNotExists(collectionSummaryKey)
	if err != nil {
		panic(err)
	}
	_, err = tx.CreateBucketIfNotExists(missionSummaryKey)
	if err != nil {
		panic(err)
	}
	_, err = tx.CreateBucketIfNotExists(resourceKey)
	if err != nil {
		panic(err)
	}
	_, err = tx.CreateBucketIfNotExists(logKey)
	if err != nil {
		panic(err)
	}
	_, err = tx.CreateBucketIfNotExists(pinKey)
	if err != nil {
		panic(err)
	}

	return &Bolt{DB: db}
}

type Bolt struct {
	*logrus.Entry
	*bolt.DB
}

func (b *Bolt) Collection() store.CollectionInterface {
	return &collection{b}
}

func (b *Bolt) Mission() store.MissionInterface {
	return ms{b}
}

func (b *Bolt) Log() store.LogInterface {
	return log{b}
}

func (b *Bolt) Animation() store.AnimationInterface {
	return anm{b}
}

func (b *Bolt) Pin() store.PinInterface {
	return p{b}
}

func setCollectionSummary(tx *bolt.Tx, summary *collectionSummary) (err error) {
	return tx.Bucket(collectionSummaryKey).Put([]byte(summary.Id), mustEncode(summary))
}

func getCollectionSummary(tx *bolt.Tx, key string, info *collectionSummary) error {
	dt := tx.Bucket(collectionSummaryKey).Get([]byte(key))
	if dt == nil {
		return store.ErrNotFound
	}
	return decodeSummary(dt, info)
}

// Animation 动画简介持久存储
func setAnimation(tx *bolt.Tx, anm *data.Animation) error {
	return tx.Bucket(animationKey).Put([]byte(anm.Id), mustEncode(anm))
}

func getAnimation(tx *bolt.Tx, animationId string, animation *data.Animation) error {
	dt := tx.Bucket(animationKey).Get([]byte(animationId))
	if dt == nil {
		return store.ErrNotFound
	}
	return decodeAnimation(dt, animation)
}

// Resource 资源持久存储
func getResourceBucketKey(collectionKey string) []byte {
	return strconv.AppendQuote(resourceKey, "-"+collectionKey)
}

func getResourceKey(resource *data.Source) []byte {
	return strconv.AppendInt(nil, int64(resource.Episode), 10)
}

func getAllSource(tx *bolt.Tx, collectionKey string) (ss []*data.Source, err error) {
	err = tx.Bucket(getResourceBucketKey(collectionKey)).ForEach(func(k, v []byte) error {
		var source = &data.Source{}
		err := decodeResource(v, source)
		if err != nil {
			return err
		}
		ss = append(ss, source)
		return nil
	})
	return ss, err
}

func setSource(tx *bolt.Tx, collectionKey string, source *data.Source) (err error) {
	bk, err := tx.CreateBucketIfNotExists(getResourceBucketKey(collectionKey))
	if err != nil {
		return err
	}
	return bk.Put(getResourceKey(source), mustEncode(source))
}

// Mission 任务执持久化存储
func getMissionSummary(tx *bolt.Tx, collectionId string, ms *missionSummary) error {
	bt := tx.Bucket(missionSummaryKey).Get([]byte(collectionId))
	if bt == nil {
		return store.ErrNotFound
	}
	return decodeMissionSummary(bt, ms)
}

func setMissionSummary(tx *bolt.Tx, ms *missionSummary) error {
	return tx.Bucket(missionSummaryKey).Put([]byte(ms.CollectionId), mustEncode(ms))
}

// Pin 固定任务持久化存储
func getPin(tx *bolt.Tx, key string) (uint8, error) {
	dt := tx.Bucket(pinKey).Get([]byte(key))
	if dt == nil {
		return 0, store.ErrNotFound
	}
	return dt[0], nil
}

func pin(tx *bolt.Tx, animation *data.Animation) error {
	return tx.Bucket(pinKey).Put([]byte(animation.Id), []byte{0})
}

func finish(tx *bolt.Tx, animation *data.Animation) error {
	return tx.Bucket(pinKey).Put([]byte(animation.Id), []byte{1})
}

func unpin(tx *bolt.Tx, animation *data.Animation) error {
	return tx.Bucket(pinKey).Delete([]byte(animation.Id))
}

// Log 日志持久化存储
func getLogBucketKey(missionKey string) []byte {
	return strconv.AppendQuoteToASCII(logKey, "-"+missionKey)
}

func getLogKey(log *mission.Log) []byte {
	return strconv.AppendInt(nil, log.EmitTime.Unix(), 10)
}

func setLog(tx *bolt.Tx, missionKey string, log *mission.Log) error {
	bk, err := tx.CreateBucketIfNotExists(getLogBucketKey(missionKey))
	if err != nil {
		return err
	}
	return bk.Put(getLogKey(log), mustEncode(log))
}

func getAllLog(tx *bolt.Tx, missionKey string) (ls []*mission.Log, err error) {
	err = tx.Bucket(getLogBucketKey(missionKey)).ForEach(func(k, v []byte) error {
		var log = &mission.Log{}
		err := decodeLog(v, log)
		if err != nil {
			return err
		}
		ls = append(ls, log)
		return nil
	})
	return ls, err
}

// cleanup 清理工作
func (b *Bolt) cleanup(tx *bolt.Tx, err error) error {
	if err == nil {
		if err = tx.Commit(); err != nil {
			if !tx.Writable() && errors.Is(bolt.ErrTxNotWritable, err) {
				return nil
			}
			b.Errorln("commit失败: ", err)
		}
	} else {
		if err := tx.Rollback(); err != nil {
			b.Errorln("rollback失败: ", err)
		}
	}
	return err
}

type collection struct {
	*Bolt
}

func (c *collection) Save(collection *data.Collection) (err error) {
	tx, err := c.Bolt.Begin(true)
	if err != nil {
		return err
	}

	defer func() {
		err = c.cleanup(tx, err)
	}()

	err = getCollectionSummary(tx, collection.Id(), &collectionSummary{})
	if err != nil {
		// 首次插入的话就把source也放入
		if err == store.ErrNotFound {
			for _, v := range collection.Items {
				err := setSource(tx, collection.Id(), v)
				if err != nil {
					return err
				}
			}
			err = nil
		} else {
			return err
		}
	}
	err = setCollectionSummary(tx, newCollectionSummary(collection))
	if err != nil {
		return err
	}

	// 保存animation数据
	return setAnimation(tx, collection.Animation)
}

func (c *collection) Get(id string) (cl *data.Collection, err error) {
	err = c.DB.View(func(tx *bolt.Tx) error {
		var summary = &collectionSummary{}
		err := getCollectionSummary(tx, id, summary)
		if err != nil {
			return err
		}

		var anm = &data.Animation{}
		err = getAnimation(tx, summary.Animation, anm)
		if err != nil {
			return err
		}

		ss, err := getAllSource(tx, id)
		if err != nil {
			return err
		}

		cl = data.NewCollection(anm, summary.Metadata)
		for _, v := range ss {
			err := cl.Add(v)
			// 这里的错误一定不可能发生
			if err != nil {
				panic(err)
			}
		}
		return nil
	})
	return cl, err
}

func (c *collection) GetAll() (cls []*data.Collection, err error) {
	var keys = make([]string, 0, 10)
	err = c.DB.View(func(tx *bolt.Tx) error {
		return tx.Bucket(collectionSummaryKey).ForEach(func(k, v []byte) error {
			keys = append(keys, string(k))
			return nil
		})
	})
	for _, v := range keys {
		cl, err := c.Get(v)
		if err != nil {
			return nil, err
		}
		cls = append(cls, cl)
	}
	return cls, err
}

type log struct {
	*Bolt
}

func (l log) Save(missionId string, log *mission.Log) error {
	tx, err := l.DB.Begin(true)
	if err != nil {
		return err
	}

	defer func() {
		err = l.cleanup(tx, err)
	}()

	return setLog(tx, missionId, log)
}

func (l log) GetAll(missionId string) (ls []*mission.Log, err error) {
	err = l.DB.View(func(tx *bolt.Tx) error {
		ls, err = getAllLog(tx, missionId)
		if err != nil {
			return err
		}
		return nil
	})
	return ls, err
}

type ms struct {
	*Bolt
}

func (m ms) Save(mission *mission.Mission) error {
	tx, err := m.DB.Begin(true)
	if err != nil {
		return err
	}

	defer func() { err = m.cleanup(tx, err) }()

	return setMissionSummary(tx, newMissionSummary(mission))
}

func (m ms) Get(id string) (ms *mission.Mission, err error) {
	var mss *missionSummary
	err = m.DB.View(func(tx *bolt.Tx) error {
		mss = &missionSummary{}
		return getMissionSummary(tx, id, mss)
	})
	cl, err := (&collection{m.Bolt}).Get(id)
	if err != nil {
		return nil, err
	}
	return &mission.Mission{
		Collection: cl,
		LastUpdate: mss.LastUpdate,
		SkipTime:   mss.SkipTime,
		Status:     mss.Status,
	}, nil
}

func (m ms) GetAll(active bool) (mss []*mission.Mission, err error) {
	var keys = make([]string, 0, 10)

	err = m.DB.View(func(tx *bolt.Tx) error {
		return tx.Bucket(missionSummaryKey).ForEach(func(k, v []byte) error {
			keys = append(keys, string(k))
			return nil
		})
	})

	for _, v := range keys {
		ms, err := m.Get(v)
		if err != nil {
			return nil, err
		}
		// 选择出活跃的或者不活跃的任务
		if !(ms.Status == mission.Finish) == active {
			mss = append(mss, ms)
		}
	}
	return mss, err
}

type anm struct {
	*Bolt
}

func (a anm) Save(animation *data.Animation) error {
	tx, err := a.DB.Begin(true)
	if err != nil {
		return err
	}

	defer func() { err = a.cleanup(tx, err) }()

	return setAnimation(tx, animation)
}

func (a anm) Get(id string) (anm *data.Animation, err error) {
	err = a.DB.View(func(tx *bolt.Tx) error {
		anm = &data.Animation{}
		return getAnimation(tx, id, anm)
	})
	return anm, err
}

type p struct {
	*Bolt
}

func (p p) Pin(animation *data.Animation) error {
	tx, err := p.DB.Begin(true)
	if err != nil {
		return err
	}

	defer func() { err = p.cleanup(tx, err) }()

	return pin(tx, animation)
}

func (p p) Unpin(animation *data.Animation) error {
	tx, err := p.DB.Begin(true)
	if err != nil {
		return err
	}

	defer func() { err = p.cleanup(tx, err) }()

	status, err := getPin(tx, animation.Id)
	if err != nil {
		return err
	}
	if status == 1 {
		return errors.New("无法取消已经完成的pin")
	}

	return unpin(tx, animation)
}

func (p p) IsPin(animation *data.Animation) (pinned bool, err error) {
	var val uint8
	err = p.DB.View(func(tx *bolt.Tx) error {
		val, err = getPin(tx, animation.Id)
		return err
	})
	if err == store.ErrNotFound {
		return false, nil
	}
	return val == 0, err
}

func (p p) Finish(animation *data.Animation) error {
	tx, err := p.DB.Begin(true)
	if err != nil {
		return err
	}

	defer func() { err = p.cleanup(tx, err) }()

	status, err := getPin(tx, animation.Id)
	if err != nil {
		return err
	}
	switch status {
	case 0:
		return finish(tx, animation)
	case 1:
		return errors.New("pin已经完成")
	default:
		return errors.New("无效的状态")
	}

}

func (p p) IsFinish(animation *data.Animation) (finish bool, err error) {
	var val uint8
	err = p.DB.View(func(tx *bolt.Tx) error {
		val, err = getPin(tx, animation.Id)
		return err
	})
	if err == store.ErrNotFound {
		return false, nil
	}
	return val == 1, err
}

func (p p) GetPinned(active interface{}) (anms []*data.Animation, err error) {
	err = p.DB.View(func(tx *bolt.Tx) error {
		var keys []string

		err = tx.Bucket(pinKey).ForEach(func(k, v []byte) error {
			if active == nil ||
				v[0] == 0 && active.(bool) ||
				v[0] == 1 && !active.(bool) {
				keys = append(keys, string(k))
			}
			return nil
		})

		for _, v := range keys {
			// 获取anm
			var anm = &data.Animation{}
			err := getAnimation(tx, v, anm)
			if err != nil {
				return err
			}
			anms = append(anms, anm)
		}
		return nil
	})
	return anms, err
}
