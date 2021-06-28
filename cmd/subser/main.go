package main

import (
	"github.com/shlande/dmhy-rss/pkg/conf"
	"github.com/shlande/dmhy-rss/pkg/parser"
	"github.com/shlande/dmhy-rss/pkg/parser/common"
	"github.com/shlande/dmhy-rss/pkg/provider"
	"github.com/shlande/dmhy-rss/pkg/provider/dmhy"
	server2 "github.com/shlande/dmhy-rss/pkg/server"
	"github.com/shlande/dmhy-rss/pkg/store"
	"github.com/shlande/dmhy-rss/pkg/store/bolt"
	"github.com/shlande/dmhy-rss/pkg/store/memory"
	"github.com/shlande/dmhy-rss/pkg/subscriber"
	"github.com/shlande/dmhy-rss/pkg/subscriber/log"
	"github.com/shlande/dmhy-rss/pkg/subscriber/record"
	"path"
)

func main() {
	server := server2.NewServer(
		buildParser(),
		buildProvider(),
		buildSubs(),
		buildPermStore(),
		memory.New(),
	)
	go server.StartHttp(conf.Default.Http.Address)
	select {}
}

func buildParser() parser.Parser {
	return common.New()
}

func buildSubs() *subscriber.Multi {
	return subscriber.Combine(
		log.NewLog(),
		record.NewRecord(record.NewJsonKVFromFile(path.Join(conf.Default.Data.OutputDir, "record.json"))),
	)
}

func buildProvider() provider.Provider {
	return dmhy.NewProvider()
}

func buildPermStore() store.Store {
	db, err := bolt.New(conf.Default.Data.OutputName())
	if err != nil {
		panic(err)
	}
	return db
}
