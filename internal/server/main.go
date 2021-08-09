package main

import (
	"context"
	"github.com/shlande/dmhy-rss/internal/server/port/http"
	"github.com/shlande/dmhy-rss/pkg/controller/manager"
	"github.com/shlande/dmhy-rss/pkg/controller/manual"
	"github.com/shlande/dmhy-rss/pkg/data/parser"
	"github.com/shlande/dmhy-rss/pkg/data/provider/bgm"
	"github.com/shlande/dmhy-rss/pkg/data/source/dmhy"
	"github.com/shlande/dmhy-rss/pkg/data/tools"
)

func main() {
	store := buildStore(defaultConfig.StoreConfig)
	sub := buildSubscribe(defaultConfig.SubscribeConfig)
	anmProvider := bgm.New(store.Animation())
	clp := tools.NewCollectionProvider(parser.New(), dmhy.NewProvider())
	manager := manager.NewManager(clp, sub, store.Mission(), store.Collection(), store.Log())
	go manager.Run(context.TODO())
	server := &Server{
		store:   store,
		manager: manager,
		manual:  manual.New(manager.AddChan(), clp),
		pvd:     anmProvider,
	}
	http.Start(":9090", server)
}

type Config struct {
	StoreConfig
	SubscribeConfig
}

var defaultConfig = Config{
	StoreConfig:     StoreConfig{Path: "./data.db"},
	SubscribeConfig: SubscribeConfig{RecordPath: "./resource.json"},
}
