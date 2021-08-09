package mission

import (
	"testing"
	"time"
)

func Test_getNextUpdateTime(t *testing.T) {
	now := time.Date(2021, 8, 9, 11, 27, 0, 0, time.Local)
	tests := []struct {
		name     string
		airTime  time.Time
		airBreak time.Duration
		want     time.Duration
	}{
		{"after", time.Date(2021, 7, 7, 15, 0, 0, 0, time.Local), time.Hour * 24 * 7, time.Hour*51 + time.Minute*33},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getNextUpdateTime(tt.airTime, tt.airBreak, 0, now); got != tt.want {
				t.Errorf("getNextUpdateTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimer(t *testing.T) {
	timer := time.NewTimer(0)
	time.Sleep(1)
	select {
	case <-timer.C:
	default:
		panic("should ok")
	}

	timer = time.NewTimer(0)
	time.Sleep(1)
	select {
	case <-timer.C:
	}

	timer.Reset(time.Minute)
	select {
	case <-timer.C:
		panic("should block")
	default:
	}
}
