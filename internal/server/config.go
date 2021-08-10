package server

import (
	"github.com/shlande/dmhy-rss/pkg/controller/manager"
	"github.com/shlande/dmhy-rss/pkg/controller/manual"
	"github.com/shlande/dmhy-rss/pkg/data/parser"
	"github.com/shlande/dmhy-rss/pkg/data/provider/bgm"
	"github.com/shlande/dmhy-rss/pkg/data/source/dmhy"
	"github.com/shlande/dmhy-rss/pkg/data/tools"
)

type Config struct {
	StoreConfig
	SubscribeConfig
}

var DefaultConfig = Config{
	StoreConfig:     StoreConfig{Path: "./data.db"},
	SubscribeConfig: SubscribeConfig{RecordPath: "./resource.json"},
}

func BuildServer(config Config) *Server {
	store := BuildStore(config.StoreConfig)
	sub := BuildSubscriber(config.SubscribeConfig)
	anmProvider := bgm.New(store.Animation())
	clp := tools.NewCollectionProvider(parser.New(), dmhy.NewProvider())
	manager := manager.NewManager(clp, sub, store.Mission(), store.Collection(), store.Log())

	return &Server{
		store:   store,
		manager: manager,
		manual:  manual.New(manager.AddChan(), clp),
		pvd:     anmProvider,
		clpd:    clp,
	}
}
