package main

import (
	"github.com/shlande/dmhy-rss/clawer"
	"github.com/shlande/dmhy-rss/core"
)

// Server 对外界暴露对服务
type Server interface {
	// Search 添加一个监控
	Search(keyword string) []*clawer.Collection
	// Delete 删除监控
	Delete(id string) error
	// List 获取当前所有的正在监听的内容
	List() []*core.Collection
}
