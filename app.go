package main

import (
	"context"
	"time"
)

// Collection 包含了单个种类的所有record
type Collection struct {
	Name       string
	Fansub     []string
	LastUpdate *time.Time
	Videos     []*Record
}

type Downloader interface {
	Download(ctx context.Context, magnet string) error
}

type App interface {
	// Add 添加一个监控
	Add(name string) error
	// Delete 删除监控
	Delete(name string) error
	// List 获取当前所有的正在监听的内容
	List() []string
}

type app struct {
	// 保存所有的内容
	name map[string][]*Record
}

func (a *app) Add(name string) error {
	if _, has := a.name[name]; !has {
		a.name[name] = make([]*Record, 0, 10)
	}
	return nil
}

func (a *app) Delete(name string) error {
	panic("implement me")
}

func (a *app) List() []string {
	panic("implement me")
}

func (a *app) sort([]Record) []*Collection {
	return nil
}
