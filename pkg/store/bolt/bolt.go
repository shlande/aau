package bolt

import (
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/shlande/dmhy-rss/pkg/classify"
	"github.com/shlande/dmhy-rss/pkg/log"
	"github.com/shlande/dmhy-rss/pkg/parser"
	"github.com/shlande/dmhy-rss/pkg/provider"
	"github.com/shlande/dmhy-rss/pkg/store"
	"github.com/shlande/dmhy-rss/pkg/worker"
	"github.com/sirupsen/logrus"
	"strconv"
)

var (
	infoKey    = []byte("info")
	episodeKey = []byte("episode")
	workerKey  = []byte("worker")
	logKey     = []byte("log")
)

func New(path string) (bt *Bolt, err error) {
	logger := log.NewEntry("bolt")
	defer func() {
		if err != nil {
			logger.Panicln(err)
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
		logger.Panicln(err)
		return nil, err
	}
	_, err = tx.CreateBucketIfNotExists(infoKey)
	if err != nil {
		return nil, err
	}
	_, err = tx.CreateBucketIfNotExists(workerKey)
	if err != nil {
		return nil, err
	}
	_, err = tx.CreateBucketIfNotExists(episodeKey)
	if err != nil {
		return nil, err
	}
	return &Bolt{Entry: logger, DB: db}, nil
}

type Bolt struct {
	*logrus.Entry
	*bolt.DB
}

func (b *Bolt) Save(collection ...*classify.Collection) error {
	tx, err := b.DB.Begin(true)
	defer func() {
		err = b.cleanup(tx, err)
	}()
	if err != nil {
		return err
	}
	for _, v := range collection {
		info := NewInfo(v)
		clId := v.Id()
		err = b.setInfo(tx, clId, info)
		if err != nil {
			return err
		}
		for _, detail := range v.Items {
			key := clId + "." + strconv.Itoa(detail.Episode)
			err = b.setEpisode(tx, key, NewEpisode(detail))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *Bolt) setInfo(tx *bolt.Tx, key string, info *Info) (err error) {
	return tx.Bucket(infoKey).Put([]byte(key), info.Encode())
}

func (b *Bolt) getInfo(tx *bolt.Tx, key string, info *Info) error {
	data := tx.Bucket([]byte("info")).Get([]byte(key))
	if data == nil {
		return store.ErrNotFound
	}
	return decodeInfo(data, info)
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

func (b *Bolt) getEpisode(tx *bolt.Tx, key string, episode *Episode) error {
	data := tx.Bucket(episodeKey).Get([]byte(key))
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
	tx, err := b.DB.Begin(false)
	defer func() {
		err = b.cleanup(tx, err)
	}()
	if err != nil {
		err = fmt.Errorf("%w:%v", ErrOpenTx, err)
		b.Errorln(err)
		return nil, err
	}
	err = b.getInfo(tx, id, info)
	if err != nil && !errors.Is(err, store.ErrNotFound) {
		return nil, err
	}
	var detail *parser.Detail
	var episode *Episode
	for i := 0; i <= info.Latest; i++ {
		if err == nil {
			detail = &parser.Detail{TitleInfo: &parser.TitleInfo{}, Info: &provider.Info{}}
			episode = NewEpisode(detail)
		}
		key := id + "." + strconv.FormatInt(int64(i), 10)
		err := b.getEpisode(tx, key, episode)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				continue
			} else {
				return nil, err
			}
		}
		_ = cl.Add(detail)
	}
	return cl, nil
}

func (b *Bolt) ListWorker() (hlps []*worker.RecoverHelper, err error) {
	tx, err := b.DB.Begin(false)
	defer func() {
		err = b.cleanup(tx, err)
	}()
	if err != nil {
		return nil, err
	}
	var ids []string
	// TODO:优化一下
	tx.Bucket(workerKey).ForEach(func(k, v []byte) error {
		ids = append(ids, string(k))
		return nil
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
func (b *Bolt) SaveWorker(worker ...*worker.Worker) (err error) {
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

func (b *Bolt) setLog(tx *bolt.Tx, key string, logs []*worker.Log) error {
	for i, v := range logs {
		key := key + "." + strconv.Itoa(i)
		err := tx.Bucket(logKey).Put([]byte(key), encodeLog(v))
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Bolt) getLog(tx *bolt.Tx, key string, from, to int) (logs []*worker.Log, err error) {
	if to-from > 0 {
		logs = make([]*worker.Log, 0, to-from)
	}
	for i := from; i < to; i++ {
		key := key + "." + strconv.Itoa(i)
		data := tx.Bucket(logKey).Get([]byte(key))
		if data == nil {
			return nil, store.ErrNotFound
		}
		var log = &worker.Log{}
		err := decodeLog(data, log)
		if err != nil {
			return nil, fmt.Errorf("%w 无法加载数据 key:%v err:%v", store.ErrOperation, key, err.Error())
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func (b *Bolt) getWorkerStatus(tx *bolt.Tx, key string, wk *Worker) error {
	data := tx.Bucket(workerKey).Get([]byte(key))
	if data == nil {
		return store.ErrNotFound
	}
	return decodeWorker(data, wk)
}

func (b *Bolt) setWorkerStatus(tx *bolt.Tx, key string, wk *Worker) error {
	return tx.Bucket(workerKey).Put([]byte(key), encodeWorker(wk))
}

func (b *Bolt) GetWorker(workerId string) (wkr *worker.RecoverHelper, err error) {
	cl, err := b.Get(workerId)
	if err != nil {
		return nil, err
	}
	tx, err := b.Begin(false)
	defer func() {
		err = b.cleanup(tx, err)
	}()
	if err != nil {
		return nil, err
	}
	var wk = &Worker{}
	err = b.getWorkerStatus(tx, workerId, wk)
	if err != nil {
		return nil, err
	}
	logs, err := b.getLog(tx, workerId, 0, wk.Logs)
	if err != nil {
		return nil, err
	}
	return worker.Recover(wk.Status, cl, wk.UpdateTime, logs), nil
}
