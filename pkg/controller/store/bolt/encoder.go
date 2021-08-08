package bolt

import (
	"encoding/json"
	"github.com/shlande/dmhy-rss/pkg/controller/mission"
	"github.com/shlande/dmhy-rss/pkg/data"
	"time"
)

type summary struct {
	Id        string
	Animation string
	data.Metadata

	// Collection 的信息
	Latest     int
	LastUpdate time.Time
}

func truncateCollection(cl *data.Collection) *summary {
	return &summary{
		Id:         cl.Id(),
		Animation:  cl.Animation.Name,
		Metadata:   cl.Metadata,
		Latest:     cl.Latest,
		LastUpdate: cl.LastUpdate,
	}
}

func mustEncode(obj interface{}) []byte {
	data, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return data
}

func decodeSummary(data []byte, sum *summary) error {
	err := json.Unmarshal(data, sum)
	if err != nil {
		panic(err)
	}
	return nil
}

func decodeResource(data []byte, ep *data.Source) error {
	err := json.Unmarshal(data, ep)
	if err != nil {
		panic(err)
	}
	return nil
}

func decodeMission(data []byte, mission *mission.Mission) error {
	err := json.Unmarshal(data, mission)
	if err != nil {
		panic(err)
	}
	return nil
}

func decodeLog(data []byte, log *mission.Log) error {
	err := json.Unmarshal(data, log)
	if err != nil {
		panic(err)
	}
	return nil
}
