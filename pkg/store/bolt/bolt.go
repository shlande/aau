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
			if tx.Writable() && errors.Is(bolt.ErrTxNotWritable, err) {
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
	if err := b.getInfo(tx, id, info); err != nil {
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

func (b *Bolt) SaveWorker(worker ...*worker.Worker) error {
	panic("implement me")
}

func (b *Bolt) GetWorker(workerId string) *worker.Worker {
	panic("implement me")
}
