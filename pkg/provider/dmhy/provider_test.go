package dmhy

import (
	"bytes"
	"context"
	"encoding/json"
	pvd "github.com/shlande/dmhy-rss/pkg/provider"
	"io"
	"os"
	"testing"
)

func Test_provider_Keywords(t *testing.T) {
	provider := NewProvider()
	ctx := context.Background()
	tests := []struct {
		name      string
		keywords  string
		wantInfos []*pvd.Info
		wantErr   bool
	}{
		{"物质转身.json", "无职转生", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInfos, err := provider.Keywords(ctx, tt.keywords)
			if (err != nil) != tt.wantErr {
				t.Errorf("Keywords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// 输出到文件中去
			f, err := os.Create(tt.name)
			if err != nil {
				panic(err)
			}
			data, err := json.Marshal(gotInfos)
			if err != nil {
				panic(err)
			}
			_, err = io.Copy(f, bytes.NewReader(data))
			if err != nil {
				panic(err)
			}
		})
	}
}
