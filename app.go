package iface

import (
	"github.com/shlande/dmhy-rss/core"
)

// Server 对外界暴露对服务
type Server interface {
	// Search 添加一个监控
	Search(keyword string) []*core.Collection
	// Add 添加监听
	Add(hash string) error
	// Delete 删除监控
	Delete(hash string) error
	// List 获取当前所有的正在监听的内容
	List() []*core.Collection
}
