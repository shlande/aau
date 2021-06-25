package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/shlande/dmhy-rss/pkg/log"
	"github.com/shlande/dmhy-rss/pkg/server/port"
	"net/http"
	"time"
)

func Start(address string, api port.Api) {
	http.ListenAndServe(address, BuildHandler(api))
}

func BuildHandler(server port.Api) http.Handler {
	router := mux.NewRouter()
	logger := log.NewEntry("http")
	router.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			codeInterceptor := &interceptStatusCode{ResponseWriter: writer}
			handler.ServeHTTP(codeInterceptor, request)
			switch codeInterceptor.code {
			case 200:
				fallthrough
			case 400:
				logger.Infoln(codeInterceptor.code, request.URL.Path, request.RemoteAddr)
			default:
				logger.Errorln(codeInterceptor.code, request.URL.Path, request.RemoteAddr)
			}
		})
	})

	router.Path("/search").Methods(http.MethodGet).HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		keywords := request.URL.Query().Get("keywords")
		if len(keywords) < 3 {
			err := fmt.Errorf("%w %v", ErrBadRequest, "关键词长度必须大于三")
			write(writer, err)
			return
		}
		write(writer, server.Search(request.Context(), keywords))
	})

	router.Path("/watch").Methods(http.MethodPost).HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()
		week, err := parseWeekday(query.Get("time"))
		if err != nil {
			write(w, err)
			return
		}
		write(w, server.Watch(query.Get("id"), week))
	})

	router.Path("/watch/{id}").Methods(http.MethodGet).HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		args := mux.Vars(request)
		write(writer, server.WatchStatus(args["id"]))
		return
	})

	router.Path("/watch/{id}").Methods(http.MethodDelete).HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		args := mux.Vars(request)
		write(writer, server.UnWatch(args["id"]))
		return
	})

	return router
}

func write(w http.ResponseWriter, data interface{}) {
	if data == nil {
		w.WriteHeader(200)
		return
	}
	switch dt := data.(type) {
	case error:
		switch errors.Unwrap(dt) {
		case ErrBadRequest:
			w.WriteHeader(400)
		default:
			w.WriteHeader(500)
		}
		w.Write([]byte(dt.Error()))
	default:
		bt, _ := json.Marshal(data)
		w.WriteHeader(200)
		w.Write(bt)
	}
}

func parseWeekday(weekday string) (week time.Weekday, err error) {
	switch weekday {
	case "1":
		week = time.Monday
	case "2":
		week = time.Tuesday
	case "3":
		week = time.Wednesday
	case "4":
		week = time.Thursday
	case "5":
		week = time.Friday
	case "6":
		week = time.Saturday
	case "7":
		week = time.Sunday
	default:
		err = fmt.Errorf("%w: %v", ErrBadRequest, "无效的日期")
	}
	return
}

type interceptStatusCode struct {
	code int
	http.ResponseWriter
}

func (i *interceptStatusCode) WriteHeader(code int) {
	i.code = code
	i.ResponseWriter.WriteHeader(code)
}
