package port

import (
	"context"
	"github.com/shlande/dmhy-rss/pkg/controller/mission"
	"github.com/shlande/dmhy-rss/pkg/data"
)

type Api interface {
	// SearchAnimation 通过关键词去查找番剧
	SearchAnimation(ctx context.Context, keywords string) ([]*data.Animation, error)
	// GetAnimationBySession 获取季度番剧
	GetAnimationBySession(ctx context.Context, year int, session int) ([]*data.Animation, error)
	// SetAnimationKeywords 设置番剧查找资源时使用的关键词
	SetAnimationKeywords(ctx context.Context, anmId, keywords string) error

	// ListCollectionByAnmId 通过番剧id查找对应的资源
	ListCollectionByAnmId(ctx context.Context, anmUniId string) ([]*data.Collection, error)
	// ListCollectionByKeywords 通过关键词查找collection信息,值得注意的是这些collection的信息可能会有问题。
	ListCollectionByKeywords(ctx context.Context, keywords string) ([]*data.Collection, error)

	// CreateMission 添加collection到监控列表中，进行同步更新
	CreateMission(ctx context.Context, collectionId string) error
	// CancelMission 取消任务
	CancelMission(ctx context.Context, collectionId string) error
	// ListMission 列出任务，active用于筛选活跃任务
	ListMission(ctx context.Context, active bool) ([]*mission.Mission, error)

	// GetLogs 获取日志
	GetLogs(ctx context.Context, missionId string) ([]*mission.Log, error)

	// GetCollection 通过id查找collection
	GetCollection(ctx context.Context, collectionId string) (*data.Collection, error)
}
