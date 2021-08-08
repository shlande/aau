package mission

import (
	"testing"
	"time"
)

func Test_getNextUpdateTime(t *testing.T) {
	tests := []struct {
		name       string
		weekday    time.Weekday
		lastUpdate time.Time
		want       time.Duration
	}{
		{"after", time.Sunday, time.Date(2021, 6, 1, 0, 0, 0, 0, time.Local), time.Hour * 24 * 5},
		{"same", time.Tuesday, time.Date(2021, 6, 1, 0, 0, 0, 0, time.Local), time.Hour * 24 * 7},
		{"day-after", time.Wednesday, time.Date(2021, 6, 1, 0, 0, 0, 0, time.Local), time.Hour * 24},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getNextUpdateTime(tt.weekday, tt.lastUpdate); got != tt.want {
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
