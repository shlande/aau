package provider

import (
	"context"
	"github.com/shlande/dmhy-rss/pkg/data"
)

type Session uint8

const (
	_              = iota
	Spring Session = iota
	Summer
	Fall
	Winter
)

type Provider interface {
	// Session 查找一个季度的番剧
	Session(ctx context.Context, year int, session Session) ([]*data.Animation, error)
	// Search 模糊查询
	Search(ctx context.Context, keywords string) ([]*data.Animation, error)
	// Get 通过id精准获取
	Get(ctx context.Context, id string) (*data.Animation, error)
}
