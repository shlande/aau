package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/shlande/dmhy-rss/pkg/server/port"
	"net/http"
)

func Start(address string, api port.Api) {
	http.ListenAndServe(address, BuildHandler(api))
}

func BuildHandler(server port.Api) http.Handler {
	router := mux.NewRouter()
	// logger := log.NewEntry("http")

	router.Get("/search").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		keywords := request.URL.Query().Get("keywords")
		if len(keywords) < 3 {
			err := fmt.Errorf("%w %v", ErrBadRequest, "关键词长度必须大于三")
			writeError(writer, err)
			return
		}
		writeSuccess(writer, server.Search(request.Context(), keywords))
	})

	return router
}

func writeError(w http.ResponseWriter, err error) {
	switch errors.Unwrap(err) {
	case ErrBadRequest:
		w.WriteHeader(400)
	default:
		w.WriteHeader(500)
	}
	w.Write([]byte(err.Error()))
}

func writeSuccess(w http.ResponseWriter, data interface{}) {
	bt, _ := json.Marshal(data)
	w.WriteHeader(200)
	w.Write(bt)
}
