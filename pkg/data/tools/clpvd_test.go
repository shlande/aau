package tools

import (
	"context"
	"encoding/json"
	"github.com/shlande/dmhy-rss/pkg/data"
	"github.com/shlande/dmhy-rss/pkg/data/parser"
	"github.com/shlande/dmhy-rss/pkg/data/source/dmhy"
	"os"
	"testing"
)

func TestCollectionProvider_Search(t *testing.T) {
	ctx := context.Background()
	c := &CollectionProvider{
		parser: parser.New(),
		pvd:    dmhy.NewProvider(),
	}
	tests := []struct {
		name      string
		animation *data.Animation
		wantErr   bool
	}{
		{
			name: "test",
			animation: &data.Animation{
				Name:       "無職転生～異世界行ったら本気だす～",
				Translated: "无职转生～到了异世界就拿出真本事～",
				Keywords:   "无职转生",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.Search(ctx, tt.animation)
			if (err != nil) != tt.wantErr {
				t.Errorf("Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			file, _ := os.Create("无职转生.json")
			defer file.Close()
			data, err := json.Marshal(got)
			if err != nil {
				panic(err)
			}
			file.Write(data)
		})
	}
}

func TestCollectionProvider_Keywords(t *testing.T) {
	ctx := context.Background()
	c := &CollectionProvider{
		parser: parser.New(),
		pvd:    dmhy.NewProvider(),
	}
	tests := []struct {
		name     string
		keywords string
		wantErr  bool
	}{
		{
			name: "test", keywords: "无职转生",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.Keywords(ctx, tt.keywords)
			if (err != nil) != tt.wantErr {
				t.Errorf("Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			file, _ := os.Create("无职转生-keywords.json")
			defer file.Close()
			data, err := json.Marshal(got)
			if err != nil {
				panic(err)
			}
			file.Write(data)
		})
	}
}
