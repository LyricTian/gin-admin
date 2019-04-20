package jwtauth

import (
	"time"
)

// Storer 黑名单存储接口
type Storer interface {
	// 放入令牌，指定到期时间
	Set(tokenString string, expiration time.Duration) error
	// 检查令牌是否存在
	Check(tokenString string) (bool, error)
	// 关闭存储
	Close() error
}
