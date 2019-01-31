package auth

import (
	"fmt"
	"time"

	"github.com/LyricTian/gin-admin/src/config"
	"github.com/LyricTian/gin-admin/src/errors"
	"github.com/LyricTian/gin-admin/src/logger"
	"github.com/LyricTian/gin-admin/src/web/context"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// 定义错误
var (
	ErrJWTInvalid = errors.NewBadRequestError("无效的请求")
)

// JwtItem jwt存储数据项
type JwtItem struct {
	UserInfo  string    `json:"user_info,omitempty"` // 用户信息
	ExpiredAt time.Time `json:"expired,omitempty"`   // 过期时间
}

// Valid 参数校验
func (a *JwtItem) Valid() error {
	if a.UserInfo == "" || a.ExpiredAt.IsZero() {
		return ErrJWTInvalid
	}
	return nil
}

// NewJWTAuth 创建jwt授权
func NewJWTAuth() Auther {
	return &jwtAuth{}
}

type jwtAuth struct {
}

func (a *jwtAuth) getFunctionName(name string) string {
	return fmt.Sprintf("auth.jwt.%s", name)
}

func (a *jwtAuth) Entry(skipper SkipperFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func (a *jwtAuth) sign(c *gin.Context, info UserInfo) error {
	cfg := config.GetJWTConfig()

	method := jwt.SigningMethodHS512
	switch cfg.SignMethod {
	case "HS256":
		method = jwt.SigningMethodHS256
	case "HS384":
		method = jwt.SigningMethodHS384
	case "HS512":
		method = jwt.SigningMethodHS512
	}
	token := jwt.NewWithClaims(method, &JwtItem{
		UserInfo:  info.String(),
		ExpiredAt: time.Now().Add(time.Second * time.Duration(cfg.Expired)),
	})
	tokenString, err := token.SignedString([]byte(cfg.Secret))
	if err != nil {
		logger.StartSpan(context.New(c).CContext(), "jwt签名", a.getFunctionName("sign")).Errorf(err.Error())
		return errors.NewInternalServerError("签名发生错误")
	}

	c.Header(cfg.HeaderName, tokenString)
	return nil
}

func (a *jwtAuth) SaveUserInfo(c *gin.Context, info UserInfo) error {
	return a.sign(c, info)
}

func (a *jwtAuth) GetUserInfo(c *gin.Context) (*UserInfo, error) {
	cfg := config.GetJWTConfig()
	tokenString := c.GetHeader(cfg.HeaderName)
	if tokenString == "" {
		return nil, nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &JwtItem{}, func(_ *jwt.Token) (interface{}, error) {
		return []byte(cfg.Secret), nil
	})
	if err != nil {
		if err == ErrJWTInvalid {
			return nil, err
		}
		logger.StartSpan(context.New(c).CContext(), "解析JWT签名", a.getFunctionName("GetUserInfo")).Errorf(err.Error())
		return nil, errors.NewBadRequestError("解析签名发生错误")
	}

	claims := token.Claims.(*JwtItem)
	if claims.ExpiredAt.Before(time.Now()) {
		return nil, nil
	}

	userInfo := parseUserInfo(claims.UserInfo)

	// 校验用户数据后，再把令牌数据重新签名并响应给客户端(为了解决用户持续访问时的令牌过期问题)
	err = a.sign(c, *userInfo)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}

func (a *jwtAuth) Destroy(c *gin.Context) error {
	return nil
}
