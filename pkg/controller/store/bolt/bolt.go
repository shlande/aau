package bolt

import (
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	worker2 "github.com/shlande/dmhy-rss/pkg/controller/manager"
	"github.com/shlande/dmhy-rss/pkg/controller/mission"
	"github.com/shlande/dmhy-rss/pkg/controller/store"
	"github.com/shlande/dmhy-rss/pkg/data"
	source2 "github.com/shlande/dmhy-rss/pkg/data/source"
	"github.com/sirupsen/logrus"
	"strconv"
)

var (
	animationKey = []byte("animation")
	resourceKey  = []byte("episode")
	missionKey   = []byte("manager")
	logKey       = []byte("log")
)

func New(path string) (bt *Bolt, err error) {
	defer func() {
		if err != nil {
			logrus.Panicln(err)
		}
	}()
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
		return nil, err
	}
	_, err = tx.CreateBucketIfNotExists(animationKey)
	if err != nil {
		return nil, err
	}
	_, err = tx.CreateBucketIfNotExists(missionKey)
	if err != nil {
		return nil, err
	}
	_, err = tx.CreateBucketIfNotExists(resourceKey)
	if err != nil {
		return nil, err
	}
	return &Bolt{DB: db}, nil
}

type Bolt struct {
	*logrus.Entry
	*bolt.DB
}

func (b *Bolt) setSummary(tx *bolt.Tx, summary *summary) (err error) {
	return tx.Bucket(animationKey).Put([]byte(summary.Id), mustEncode(summary))
}

func (b *Bolt) getSummary(tx *bolt.Tx, key string, info *summary) error {
	dt := tx.Bucket([]byte("info")).Get([]byte(key))
	if dt == nil {
		return store.ErrNotFound
	}
	return decodeSummary(dt, info)
}

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

func (b *collection) Save(collection *data.Collection) error {
	tx, err := b.DB.Begin(true)
	defer func() {
		err = b.cleanup(tx, err)
	}()
	if err != nil {
		return err
	}
	clId := collection.Id()
	err = b.set(tx, clId, info)
	if err != nil {
		return err
	}
	for _, detail := range collection.Items {
		key := clId + "." + strconv.Itoa(detail.Episode)
		err = b.setEpisode(tx, key, NewEpisode(detail))
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Bolt) getEpisode(tx *bolt.Tx, key string, episode *Episode) error {
	data := tx.Bucket(resourceKey).Get([]byte(key))
	if data == nil {
		return store.ErrNotFound
	}
	return decodeEpisode(data, episode)
}

func (b *Bolt) setEpisode(tx *bolt.Tx, key string, episode *Episode) (err error) {
	return tx.Bucket([]byte("episode")).Put([]byte(key), episode.Encode())
}

func (b *Bolt) Get(id string) (*classify.Collection, error) {
	var cl = &classify.Collection{}
	info := NewInfo(cl)
	var detail *parser.Detail
	var episode *Episode
	err := b.DB.View(func(tx *bolt.Tx) (err error) {
		for i := 0; i <= info.Latest; i++ {
			if err == nil {
				detail = &parser.Detail{TitleInfo: &parser.TitleInfo{}, Info: &source2.Info{}}
				episode = NewEpisode(detail)
			}
			key := id + "." + strconv.FormatInt(int64(i), 10)
			err = b.getEpisode(tx, key, episode)
			if err != nil {
				if errors.Is(err, store.ErrNotFound) {
					continue
				} else {
					return err
				}
			}
			_ = cl.Add(detail)
		}
		return b.getInfo(tx, id, info)
	})
	if err != nil && !errors.Is(err, store.ErrNotFound) {
		return nil, err
	}
	return cl, nil
}

func (b *Bolt) ListWorker() (hlps []*worker2.RecoverHelper, err error) {
	var ids []string
	err = b.DB.View(func(tx *bolt.Tx) error {
		// TODO:优化一下
		return tx.Bucket(missionKey).ForEach(func(k, v []byte) error {
			ids = append(ids, string(k))
			return nil
		})
	})
	for _, v := range ids {
		wk, err := b.GetWorker(v)
		if err != nil {
			return nil, err
		}
		hlps = append(hlps, wk)
	}
	return
}

func (b *Bolt) SaveWorker(worker ...*worker2.Misson) (err error) {
	tx, err := b.Begin(true)
	defer func() {
		err = b.cleanup(tx, err)
	}()
	if err != nil {
		return err
	}
	for _, v := range worker {
		err = b.setWorkerStatus(tx, v.Id, NewWorker(v))
		if err != nil {
			return err
		}
		// FIXME:虽然worker每次更新的时候都同步了信息，但是最好还是检查一下
	}
	return nil
}

func (b *Bolt) setLog(tx *bolt.Tx, key string, logs []*mission.Log) error {
	for i, v := range logs {
		key := key + "." + strconv.Itoa(i)
		err := tx.Bucket(logKey).Put([]byte(key), encodeLog(v))
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Bolt) getLog(tx *bolt.Tx, key string, from, to int) (logs []*mission.Log, err error) {
	if to-from > 0 {
		logs = make([]*mission.Log, 0, to-from)
	}
	for i := from; i < to; i++ {
		key := key + "." + strconv.Itoa(i)
		data := tx.Bucket(logKey).Get([]byte(key))
		if data == nil {
			return nil, store.ErrNotFound
		}
		var log = &mission.Log{}
		err := decodeLog(data, log)
		if err != nil {
			return nil, fmt.Errorf("%w 无法加载数据 key:%v err:%v", store.ErrOperation, key, err.Error())
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func (b *Bolt) getWorkerStatus(tx *bolt.Tx, key string, wk *Worker) error {
	data := tx.Bucket(missionKey).Get([]byte(key))
	if data == nil {
		return store.ErrNotFound
	}
	return decodeWorker(data, wk)
}

func (b *Bolt) setWorkerStatus(tx *bolt.Tx, key string, wk *Worker) error {
	return tx.Bucket(missionKey).Put([]byte(key), encodeWorker(wk))
}

func (b *Bolt) GetWorker(workerId string) (wkr *worker2.RecoverHelper, err error) {
	cl, err := b.Get(workerId)
	if err != nil {
		return nil, err
	}
	var logs []*mission.Log
	var wk = &Worker{}
	err = b.View(func(tx *bolt.Tx) error {
		err = b.getWorkerStatus(tx, workerId, wk)
		if err != nil {
			return err
		}
		logs, err = b.getLog(tx, workerId, 0, wk.Logs)
		return err
	})
	return worker2.Recover(wk.Status, cl, wk.UpdateTime, logs), nil
}
