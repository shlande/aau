package main

import "context"

// Store 用来保存爬到的动漫
type Store interface {
	// Store 保存数据
	Store(ctx context.Context, record *Record) error
	// FindByName 通过名称搜索
	FindByName(ctx context.Context, name string) ([]*Record, error)
}
