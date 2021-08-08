package bolt

import (
	"encoding/json"
	"github.com/shlande/dmhy-rss/pkg/controller/mission"
	"github.com/shlande/dmhy-rss/pkg/data"
	"time"
)

// collectionSummary 包含了collection的收录信息
type collectionSummary struct {
	Id        string
	Animation string
	data.Metadata

	// Collection 的信息
	Latest     int
	LastUpdate time.Time
}

func newCollectionSummary(cl *data.Collection) *collectionSummary {
	return &collectionSummary{
		Id:         cl.Id(),
		Animation:  cl.Animation.Name,
		Metadata:   cl.Metadata,
		Latest:     cl.Latest,
		LastUpdate: cl.LastUpdate,
	}
}

type missionSummary struct {
	CollectionId string
	LastUpdate   time.Time
	SkipTime     int
	Status       mission.Status
}

func newMissionSummary(ms *mission.Mission) *missionSummary {
	return &missionSummary{
		CollectionId: ms.Collection.Id(),
		LastUpdate:   ms.LastUpdate,
		SkipTime:     ms.SkipTime,
		Status:       ms.Status,
	}
}

func mustEncode(obj interface{}) []byte {
	data, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return data
}

func decodeAnimation(data []byte, anm *data.Animation) error {
	err := json.Unmarshal(data, anm)
	if err != nil {
		panic(err)
	}
	return nil
}

func decodeMissionSummary(data []byte, ms *missionSummary) error {
	err := json.Unmarshal(data, ms)
	if err != nil {
		panic(err)
	}
	return nil
}

func decodeSummary(data []byte, sum *collectionSummary) error {
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
