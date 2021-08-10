package main

import (
	"context"
	"github.com/shlande/dmhy-rss/internal/server"
	"github.com/shlande/dmhy-rss/internal/server/port/http"
)

func main() {
	server := server.BuildServer(server.DefaultConfig)
	go server.Run(context.TODO())
	http.Start(":9090", server)
}
