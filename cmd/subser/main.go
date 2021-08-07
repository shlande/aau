package main

import (
	server3 "github.com/shlande/dmhy-rss/internal/server"
	"github.com/shlande/dmhy-rss/pkg/conf"
	store2 "github.com/shlande/dmhy-rss/pkg/controller/store"
	bolt2 "github.com/shlande/dmhy-rss/pkg/controller/store/bolt"
	memory2 "github.com/shlande/dmhy-rss/pkg/controller/store/memory"
	subscriber2 "github.com/shlande/dmhy-rss/pkg/controller/subscriber"
	log2 "github.com/shlande/dmhy-rss/pkg/controller/subscriber/log"
	record2 "github.com/shlande/dmhy-rss/pkg/controller/subscriber/record"
	source2 "github.com/shlande/dmhy-rss/pkg/data/source"
	dmhy2 "github.com/shlande/dmhy-rss/pkg/data/source/dmhy"
	"github.com/shlande/dmhy-rss/pkg/parser"
	"github.com/shlande/dmhy-rss/pkg/parser/common"
	"path"
)

func main() {
	server := server3.NewServer(
		buildParser(),
		buildProvider(),
		buildSubs(),
		buildPermStore(),
		memory2.New(),
	)
	go server.StartHttp(conf.Default.Http.Address)
	select {}
}

func buildParser() parser.Parser {
	return common.New()
}

func buildSubs() *subscriber2.Multi {
	return subscriber2.Combine(
		log2.NewLog(),
		record2.NewRecord(record2.NewJsonKVFromFile(path.Join(conf.Default.Data.OutputDir, "record.json"))),
	)
}

func buildProvider() source2.Provider {
	return dmhy2.NewProvider()
}

func buildPermStore() store2.Store {
	db, err := bolt2.New(conf.Default.Data.OutputName())
	if err != nil {
		panic(err)
	}
	return db
}
