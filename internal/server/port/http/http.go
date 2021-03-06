package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/shlande/dmhy-rss/internal/server/port"
	"github.com/shlande/dmhy-rss/pkg/controller/mission"
	"github.com/shlande/dmhy-rss/pkg/data"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

func Start(address string, api port.Api) {
	err := http.ListenAndServe(address, BuildHandler(api))
	logrus.Error("http服务异常退出", err)
}

func BuildHandler(server port.Api) http.Handler {
	router := mux.NewRouter()
	router.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			codeInterceptor := &interceptStatusCode{ResponseWriter: writer}
			writer.Header().Set("content-type", "application/json")
			handler.ServeHTTP(codeInterceptor, request)
			switch codeInterceptor.code {
			case 200:
				fallthrough
			case 400:
				logrus.Infoln(codeInterceptor.code, request.Method, request.RequestURI, request.RemoteAddr)
			default:
				logrus.Errorln(codeInterceptor.code, request.Method, request.RequestURI, request.RemoteAddr)
			}
		})
	})

	router.Path("/search/animation/{keywords}").Methods(http.MethodGet).HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		keywords := mux.Vars(request)["keywords"]
		if len(keywords) < 3 {
			err := fmt.Errorf("%w %v", ErrBadRequest, "关键词长度必须大于三")
			write(writer, err)
			return
		}
		res, err := server.SearchAnimation(request.Context(), keywords)
		if err != nil {
			write(writer, err)
			return
		}
		write(writer, res)
	})

	router.Path("/search/collection/{keywords}").Methods(http.MethodGet).HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		keywords := mux.Vars(request)["keywords"]
		if len(keywords) < 3 {
			err := fmt.Errorf("%w %v", ErrBadRequest, "关键词长度必须大于三")
			write(writer, err)
			return
		}
		res, err := server.ListCollectionByKeywords(request.Context(), keywords)
		if err != nil {
			write(writer, err)
			return
		}
		var resp = make([]*collectionSummary, 0, len(res))
		for _, v := range res {
			resp = append(resp, newCollectionSummary(v))
		}
		write(writer, resp)
	})

	router.Path("/search/session/{year}/{session}").Methods(http.MethodGet).HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		y, s := mux.Vars(request)["year"], mux.Vars(request)["session"]
		year, err1 := strconv.ParseInt(y, 10, 64)
		session, err2 := strconv.ParseInt(s, 10, 64)
		if err1 != nil || err2 != nil {
			write(writer, "无效的参数")
			return
		}
		anm, err := server.GetAnimationBySession(request.Context(), int(year), int(session))
		if err != nil {
			write(writer, err)
			return
		}
		write(writer, anm)
	})

	router.Path("/animation/{animationId}/keywords/{keywords}").Methods(http.MethodGet).HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		keywords, anmId := mux.Vars(request)["keywords"], mux.Vars(request)["animationId"]
		write(writer, server.SetAnimationKeywords(request.Context(), anmId, keywords))
	})

	router.Path("/collection/list/{animationId}").Methods(http.MethodGet).HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		anmId := mux.Vars(request)["animationId"]
		cls, err := server.ListCollectionByAnmId(request.Context(), anmId)
		if err != nil {
			write(writer, err)
			return
		}
		var clss = make([]*collectionSummary, 0, len(cls))
		for _, v := range cls {
			clss = append(clss, newCollectionSummary(v))
		}
		write(writer, clss)
	})

	router.Path("/collection/{id}").Methods(http.MethodGet).HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		clId := mux.Vars(request)["id"]
		cl, err := server.GetCollection(request.Context(), clId)
		if err != nil {
			write(writer, err)
			return
		}
		write(writer, newCollectionSummary(cl))
	})

	router.Path("/mission/create/{collectionId}").Methods(http.MethodGet).HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		id := mux.Vars(request)["collectionId"]
		write(w, server.CreateMission(request.Context(), id))
	})

	router.Path("/mission/list/{status}").Methods(http.MethodGet).HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var err error
		status := mux.Vars(request)["status"]
		var active, all bool
		all = all
		switch status {
		case "active":
			active = true
		case "inactive":
			active = false
		case "":
			all = true
		default:
			write(writer, errors.New("位置的请求类型"))
			return
		}
		ms, err := server.ListMission(request.Context(), active)
		if err != nil {
			write(writer, err)
			return
		}
		var mss = make([]*missionSummary, 0, len(ms))
		for _, v := range ms {
			mss = append(mss, newMissionSummary(v))
		}
		write(writer, mss)
		return
	})

	router.Path("/mission/{id}/log").Methods(http.MethodGet).HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		msId := mux.Vars(request)["id"]
		logs, err := server.GetLogs(request.Context(), msId)
		if err != nil {
			write(writer, err)
			return
		}
		write(writer, logs)
	})

	return router
}

type collectionSummary struct {
	Id          string
	Name        string
	Translated  string
	AnimationId string
	Metadata    *metadata

	// Collection 的信息
	Latest     int
	LastUpdate time.Time

	Items []*data.Source
}

func newCollectionSummary(collection *data.Collection) *collectionSummary {
	return &collectionSummary{
		Id:          collection.Id(),
		AnimationId: collection.Animation.Id,
		Name:        collection.Name,
		Translated:  collection.Translated,
		Metadata:    newMetadata(&collection.Metadata),
		Latest:      collection.Latest,
		LastUpdate:  collection.LastUpdate,
		Items:       collection.Items,
	}
}

type missionSummary struct {
	*collectionSummary

	LastUpdate time.Time
	SkipTime   int
	Status     string
}

func newMissionSummary(mission *mission.Mission) *missionSummary {
	return &missionSummary{
		collectionSummary: newCollectionSummary(mission.Collection),
		LastUpdate:        mission.LastUpdate,
		SkipTime:          mission.SkipTime,
		Status:            mission.Status.String(),
	}
}

type metadata struct {
	Fansub   []string
	Quality  string
	Type     string
	Language string
	SubType  string
}

func newMetadata(md *data.Metadata) *metadata {
	return &metadata{
		Fansub:   md.Fansub,
		Quality:  md.Quality.String(),
		Type:     md.Type.String(),
		Language: md.Language.String(),
		SubType:  md.SubType.String(),
	}
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
		bt, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		w.WriteHeader(200)
		w.Write(bt)
	}
}

type interceptStatusCode struct {
	code int
	http.ResponseWriter
}

func (i *interceptStatusCode) WriteHeader(code int) {
	i.code = code
	i.ResponseWriter.WriteHeader(code)
}
