package source

import "context"

// Provider 提供根据关键词提供信息的功能
type Provider interface {
	Keywords(ctx context.Context, keywords string) ([]*Info, error)
}
