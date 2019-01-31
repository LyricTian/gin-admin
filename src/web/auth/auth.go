package auth

import (
	"sync"

	"github.com/LyricTian/gin-admin/src/util"
	"github.com/gin-gonic/gin"
)

// UserInfo 存储用户信息
type UserInfo struct {
	UserID string `json:"user_id,omitempty"` // 用户ID
}

func (a UserInfo) String() string {
	return util.JSONMarshalToString(a)
}

func parseUserInfo(v string) *UserInfo {
	var info UserInfo
	util.JSONUnmarshal([]byte(v), &info)
	return &info
}

// SkipperFunc 定义中间件跳过函数
type SkipperFunc func(*gin.Context) bool

// Auther 授权接口
type Auther interface {
	// 授权入口中间件
	Entry(skipper SkipperFunc) gin.HandlerFunc

	// 保存用户信息
	SaveUserInfo(c *gin.Context, info UserInfo) error

	// 获取用户信息
	GetUserInfo(c *gin.Context) (*UserInfo, error)

	// 销毁授权信息
	Destroy(c *gin.Context) error
}

var (
	globalAuther Auther
	once         sync.Once
)

// SetGlobalAuther 设定全局的授权模式
func SetGlobalAuther(a Auther) {
	internalAuther(a)
}

func internalAuther(a ...Auther) Auther {
	once.Do(func() {
		if len(a) > 0 && a[0] != nil {
			globalAuther = a[0]
			return
		}
		globalAuther = NewJWTAuth()
	})
	return globalAuther
}

// Entry 授权入口中间件
func Entry(skipper SkipperFunc) gin.HandlerFunc {
	return internalAuther().Entry(skipper)
}

// SaveUserInfo 保存用户信息
func SaveUserInfo(c *gin.Context, info UserInfo) error {
	return internalAuther().SaveUserInfo(c, info)
}

// GetUserInfo 获取用户信息
func GetUserInfo(c *gin.Context) (*UserInfo, error) {
	return internalAuther().GetUserInfo(c)
}

// Destroy 销毁授权信息
func Destroy(c *gin.Context) error {
	return internalAuther().Destroy(c)
}
