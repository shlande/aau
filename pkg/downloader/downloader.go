package downloader

import "context"

// Downloader 外部下载接口
type Downloader interface {
	// Add 添加下载记录
	Add(ctx context.Context, magnet string) error
	// Check 查找状态
	Check(ctx context.Context, magnet string) (float64, error)
	// Delete 删除
	Delete(ctx context.Context, magnet string) error
}
