package server

import (
	"github.com/shlande/dmhy-rss/pkg/controller/store"
	"github.com/shlande/dmhy-rss/pkg/controller/store/bolt"
)

type StoreConfig struct {
	Path string
}

func BuildStore(config StoreConfig) store.Interface {
	if len(config.Path) == 0 {
		config.Path = "./data.db"
	}
	return bolt.New(config.Path)
}
