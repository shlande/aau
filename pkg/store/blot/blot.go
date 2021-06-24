package blot

import (
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/shlande/dmhy-rss/pkg/classify"
	"github.com/shlande/dmhy-rss/pkg/log"
	"github.com/shlande/dmhy-rss/pkg/parser"
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

func New(path string) (*Blot, error) {
	logger := log.NewEntry("bolt")
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1})
	if err != nil {
		logger.Panicln(err)
		return nil, err
	}
	return &Blot{Entry: logger, DB: db}, nil
}

type Blot struct {
	*logrus.Entry
	*bolt.DB
}

func (b Blot) Save(collection ...*classify.Collection) error {
	tx, err := b.DB.Begin(false)
	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				b.Errorln("rollback失败: ", err)
			}
		}
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
			err = b.setEpisode(tx, clId+"."+strconv.FormatInt(int64(detail.Episode), 10), NewEpisode(detail))
			if err != nil {
				return err
			}
		}
	}
	if err = tx.Commit(); err != nil {
		b.Errorln("commit失败: ", err)
	}
	return nil
}

func (b *Blot) setInfo(tx *bolt.Tx, key string, info *Info) (err error) {
	return tx.Bucket(infoKey).Put([]byte(key), info.Encode())
}

func (b *Blot) getInfo(tx *bolt.Tx, key string, info *Info) error {
	data := tx.Bucket([]byte("info")).Get([]byte(key))
	if data == nil {
		return store.ErrNotFound
	}
	return decodeInfo(data, info)
}

func (b *Blot) getEpisode(tx *bolt.Tx, key string, episode *Episode) error {
	data := tx.Bucket(episodeKey).Get([]byte(key))
	if data == nil {
		return store.ErrNotFound
	}
	return decodeEpisode(data, episode)
}

func (b *Blot) setEpisode(tx *bolt.Tx, key string, episode *Episode) (err error) {
	return tx.Bucket([]byte("episode")).Put([]byte(key), episode.Encode())
}

func (b Blot) Get(id string) (*classify.Collection, error) {
	var cl = &classify.Collection{}
	info := NewInfo(cl)
	tx, err := b.DB.Begin(true)
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
			detail = &parser.Detail{}
			episode = NewEpisode(detail)
		}
		err := b.getEpisode(tx, id+"."+strconv.FormatInt(int64(i), 10), episode)
		if err != nil && !errors.Is(err, store.ErrNotFound) {
			return nil, err
		}
		_ = cl.Add(detail)
	}
	return cl, nil
}

func (b Blot) SaveWorker(worker ...*worker.Worker) error {
	panic("implement me")
}

func (b Blot) GetWorker(workerId string) *worker.Worker {
	panic("implement me")
}
