package download

import (
	"github.com/shlande/dmhy-rss/pkg/data"
	"testing"
)

// TODO: 添加测试代码
func Test_path_Path(t *testing.T) {
	tests := []struct {
		name       string
		base       string
		collection *data.Collection
		want       string
	}{
		{name: "test", base: "/downloads"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &path{
				base: tt.base,
			}
			if got := p.Path(tt.collection); got != tt.want {
				t.Errorf("Path() = %v, want %v", got, tt.want)
			}
		})
	}
}
