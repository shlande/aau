package pinnner

import (
	"github.com/shlande/dmhy-rss/pkg/data"
	"github.com/shlande/dmhy-rss/pkg/data/tools"
	"reflect"
	"testing"
)

func Test_pinner_tryFindBest(t *testing.T) {
	type fields struct {
		CollectionProvider tools.CollectionProvider
		finished           map[string]struct{}
		pined              map[string]*data.Animation
		stg                map[string]Strategy
	}
	type args struct {
		animation *data.Animation
		strategy  Strategy
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *data.Collection
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &pinner{
				CollectionProvider: tt.fields.CollectionProvider,
				finished:           tt.fields.finished,
				pined:              tt.fields.pined,
				stg:                tt.fields.stg,
			}
			if got := p.tryFindBest(tt.args.animation, tt.args.strategy); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tryFindBest() = %v, want %v", got, tt.want)
			}
		})
	}
}
