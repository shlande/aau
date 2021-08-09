package bgm

import (
	"context"
	"github.com/shlande/dmhy-rss/pkg/controller/store"
	"github.com/shlande/dmhy-rss/pkg/controller/store/bolt"
	"testing"
	"time"
)

func getLocalTestStore() store.AnimationInterface {
	return bolt.New("./test.db").Animation()
}

func TestParseBroadcast(t *testing.T) {
	tests := []struct {
		name         string
		bs           string
		wantAirTime  time.Time
		wantAirBreak time.Duration
	}{
		{
			"R/2021-07-22T15:55:00.000Z/P7D",
			"R/2021-07-22T15:55:00.000Z/P7D",
			time.Date(2021, 7, 22, 23, 55, 0, 0, time.Local),
			time.Hour * 24 * 7,
		},
		{
			"R/2020-10-02T00:00:01.000Z/P0D",
			"R/2020-10-02T00:00:01.000Z/P0D",
			time.Date(2020, 10, 2, 8, 0, 1, 0, time.Local),
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAirTime, gotAirBreak, err := ParseBroadcast(tt.bs)
			if err != nil {
				t.Error(err)
			}
			if !tt.wantAirTime.Equal(gotAirTime) {
				t.Errorf("ParseBroadcast() gotAirTime = %v, want %v", gotAirTime, tt.wantAirTime)
			}
			if gotAirBreak != tt.wantAirBreak {
				t.Errorf("ParseBroadcast() gotAirBreak = %v, want %v", gotAirBreak, tt.wantAirBreak)
			}
		})
	}
}

func TestSearch(t *testing.T) {
	bgm := &Provider{
		lastUpdate: time.Time{},
		data:       nil,
	}
	res, err := bgm.search(context.Background(), "无职转生")
	if err != nil {
		panic(err)
	}
	res = res
}

func TestBgmData(t *testing.T) {
	resp, err := getBgmData()
	if err != nil {
		panic(err)
	}
	resp = resp
}

func TestProvider_Search(t *testing.T) {
	an, err := New(getLocalTestStore()).Search(context.Background(), "无职转生")
	if err != nil {
		panic(err)
	}
	an = an
}

func TestProvider_Session(t *testing.T) {
	anm, err := New(getLocalTestStore()).Session(context.Background(), 2021, 2)
	if err != nil {
		panic(err)
	}
	anm = anm
}
