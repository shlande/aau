package tools

import (
	"context"
	"github.com/shlande/dmhy-rss/pkg/data"
	"github.com/shlande/dmhy-rss/pkg/data/parser"
	"github.com/shlande/dmhy-rss/pkg/data/source"
	"github.com/shlande/dmhy-rss/third_part/workqueue"
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
	infos, err := c.pvd.Keywords(ctx, animation.Translated)
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
	for _, v := range cls {
		c.addToCache(v)
	}
	return cls, nil
}

func (c *CollectionProvider) addToCache(collection *data.Collection) {
	c.cache[collection.Id()] = collection
	// 一小时后gc会清理掉cache
	c.gcq.AddAfter(collection.Id(), time.Hour)
}

func (c *CollectionProvider) Get(collectionId string) *data.Collection {
	return c.cache[collectionId]
}
