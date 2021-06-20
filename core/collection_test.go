package core

import (
	"encoding/json"
	"github.com/shlande/dmhy-rss/pkg/parser"
	"os"
	"reflect"
	"testing"
	"time"
)

func init() {
	f, err := os.Open("../test.json")
	if err != nil {
		panic(err)
	}
	err = json.NewDecoder(f).Decode(cls)
	if err != nil {
		panic(err)
	}
}

var cls []*parser.Collection

func TestNewCollection(t *testing.T) {
	type args struct {
		collection *parser.Collection
		updateTime time.Weekday
		episodes   int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "1", args: args{collection: cls[0], updateTime: time.Monday, episodes: 12}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCollection(tt.args.collection, tt.args.updateTime, tt.args.episodes); !reflect.DeepEqual(got.Hash_, tt.want) {
				t.Errorf("NewCollection() = %v, want %v", got.Hash_, tt.want)
			}
		})
	}
}
