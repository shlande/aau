package source

import "time"

type Info struct {
	Name string
	// 发布的时间，可选
	CreateTime time.Time
	// 下载相关的信息，两者中必须有一个
	MagnetUrl  string
	TorrentUrl string
}
