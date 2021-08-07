package tools

import (
	"context"
	"github.com/shlande/dmhy-rss/pkg/data"
	"github.com/shlande/dmhy-rss/pkg/data/parser"
	"github.com/shlande/dmhy-rss/pkg/data/source"
	"log"
)

type CollectionProvider struct {
	parser.Parser
	source.Provider
}

func (c *CollectionProvider) Search(ctx context.Context, animation *data.Animation) ([]*data.Collection, error) {
	infos, err := c.Provider.Keywords(ctx, animation.Translated)
	if err != nil && err != data.ErrEpisodeExist {
		return nil, err
	}
	var sources = make([]*data.Source, 0, len(infos))
	for _, v := range infos {
		result, err := c.Parser.Parse(v.Name)
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
	return Classify(animation, sources), nil
}
