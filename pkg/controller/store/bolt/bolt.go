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

func (c collection) Save(collection *data.Collection) error {
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
		}
		return err
	}
	return setCollectionSummary(tx, newCollectionSummary(collection))
}

func (c collection) Get(id string) (cl *data.Collection, err error) {
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

func (c collection) GetAll() ([]*data.Collection, error) {
	c.DB.View(func(tx *bolt.Tx) error {
		tx.Bucket(collectionSummaryKey)
	})
}
