package tools

import (
	"context"
	"github.com/shlande/dmhy-rss/pkg/controller/store"
	"github.com/shlande/dmhy-rss/pkg/data"
	"github.com/shlande/dmhy-rss/pkg/data/parser"
	"github.com/shlande/dmhy-rss/pkg/data/source"
	"github.com/shlande/dmhy-rss/third_part/workqueue"
	"github.com/sirupsen/logrus"
	"log"
	"time"
)

func NewCollectionProvider(p parser.Parser, provider source.Provider) *CollectionProvider {
	cl := &CollectionProvider{
		gcq:    workqueue.NewDelayingQueue(),
		parser: p,
		pvd:    provider,
		cache:  map[string]*data.Collection{},
	}
	go cl.gc()
	return cl
}

type CollectionProvider struct {
	gcq    workqueue.DelayingInterface
	parser parser.Parser
	pvd    source.Provider
	store  store.CollectionInterface
	cache  map[string]*data.Collection
}

func (c *CollectionProvider) gc() {
	for {
		val, shutdown := c.gcq.Get()
		if shutdown {
			break
		}
		delete(c.cache, val.(string))
	}
}

func (c *CollectionProvider) Search(ctx context.Context, animation *data.Animation) ([]*data.Collection, error) {
	cls, err := c.search(ctx, animation.GetKeywords(), animation)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	for _, v := range cls {
		c.addToCache(v)
	}
	return cls, err
}

func (c *CollectionProvider) search(ctx context.Context, keywords string, animation *data.Animation) ([]*data.Collection, error) {
	infos, err := c.pvd.Keywords(ctx, keywords)
	if err != nil && err != data.ErrEpisodeExist {
		return nil, err
	}
	var sources = make([]*data.Source, 0, len(infos))
	for _, v := range infos {
		result, err := c.parser.Parse(v.Name)
		if err != nil {
			log.Println(err)
			continue
		}
		sources = append(sources, &data.Source{
			Name:     v.Name,
			Info:     v,
			Episode:  result.Episode,
			Metadata: result.Metadata,
		})
	}
	cls := Classify(animation, sources)
	return cls, nil
}

// Keywords 通过keywords查找的时候，collection中的animation会是nil，而且该collection无法用于加载任务
func (c *CollectionProvider) Keywords(ctx context.Context, keywords string) ([]*data.Collection, error) {
	return c.search(ctx, keywords, nil)
}

func (c *CollectionProvider) addToCache(collection *data.Collection) {
	c.cache[collection.Id()] = collection
	// 一小时后gc会清理掉cache
	c.gcq.AddAfter(collection.Id(), time.Hour)
}

func (c *CollectionProvider) Get(collectionId string) *data.Collection {
	return c.cache[collectionId]
}
