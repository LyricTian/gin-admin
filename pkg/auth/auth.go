package auth

import (
	"context"
	"errors"
)

// 定义错误
var (
	ErrInvalidToken = errors.New("invalid token")
)

// TokenInfo 令牌信息
type TokenInfo interface {
	// 获取访问令牌
	GetAccessToken() string
	// 获取令牌类型
	GetTokenType() string
	// 获取令牌到期时间戳
	GetExpiresAt() int64
	// JSON编码
	EncodeToJSON() ([]byte, error)
}

// Auther 认证接口
type Auther interface {
	// 生成令牌
	GenerateToken(ctx context.Context, userID string) (TokenInfo, error)

	// 销毁令牌
	DestroyToken(ctx context.Context, accessToken string) error

	// 解析用户ID
	ParseUserID(ctx context.Context, accessToken string) (string, error)

	// 释放资源
	Release() error
}
