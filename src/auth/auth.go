package auth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// 定义错误
var (
	ErrInvalidToken = errors.New("invalid token")
)

var defaultKey = []byte("GINADMIN")
var defaultOptions = options{
	expired:       7200,
	signingMethod: jwt.SigningMethodHS512,
	signingKey:    defaultKey,
	keyfunc: func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return defaultKey, nil
	},
}

type options struct {
	blackStore    BlackStorer
	signingMethod jwt.SigningMethod
	signingKey    interface{}
	keyfunc       jwt.Keyfunc
	expired       int
}

// Option 定义参数项
type Option func(*options)

// SetBlackStore 设定黑名单存储
func SetBlackStore(store BlackStorer) Option {
	return func(o *options) {
		o.blackStore = store
	}
}

// SetSigningMethod 设定签名方式
func SetSigningMethod(method jwt.SigningMethod) Option {
	return func(o *options) {
		o.signingMethod = method
	}
}

// SetSigningKey 设定签名key
func SetSigningKey(key interface{}) Option {
	return func(o *options) {
		o.signingKey = key
	}
}

// SetKeyfunc 设定验证key的回调函数
func SetKeyfunc(keyFunc jwt.Keyfunc) Option {
	return func(o *options) {
		o.keyfunc = keyFunc
	}
}

// SetExpired 设定令牌过期时长(单位秒，默认7200)
func SetExpired(expired int) Option {
	return func(o *options) {
		o.expired = expired
	}
}

// New 创建认证实例
func New(opts ...Option) *Auth {
	o := defaultOptions
	for _, opt := range opts {
		opt(&o)
	}

	return &Auth{
		opts: &o,
	}
}

// Auth jwt认证
type Auth struct {
	opts *options
}

// GenerateToken 生成令牌
func (a *Auth) GenerateToken(userID string) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(a.opts.signingMethod, &jwt.StandardClaims{
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(time.Duration(a.opts.expired) * time.Second).Unix(),
		NotBefore: now.Unix(),
		Subject:   userID,
	})
	return token.SignedString(a.opts.signingKey)
}

// 解析令牌
func (a *Auth) parseToken(tokenString string) (*jwt.StandardClaims, error) {
	token, _ := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, a.opts.keyfunc)
	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return token.Claims.(*jwt.StandardClaims), nil
}

func (a *Auth) callStore(fn func(BlackStorer) error) error {
	if store := a.opts.blackStore; store != nil {
		return fn(store)
	}
	return nil
}

// DestroyToken 销毁令牌
func (a *Auth) DestroyToken(tokenString string) error {
	claims, err := a.parseToken(tokenString)
	if err != nil {
		return err
	}

	// 如果设定了黑名单存储，则将未过期的令牌放入黑名单
	return a.callStore(func(store BlackStorer) error {
		expired := time.Unix(claims.ExpiresAt, 0).Sub(time.Now())
		return store.Set(tokenString, expired)
	})
}

// ParseUserID 解析用户ID
func (a *Auth) ParseUserID(tokenString string) (string, error) {
	claims, err := a.parseToken(tokenString)
	if err != nil {
		return "", err
	}

	err = a.callStore(func(store BlackStorer) error {
		exists, err := store.Check(tokenString)
		if err != nil {
			return err
		} else if exists {
			return ErrInvalidToken
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	return claims.Subject, nil
}

// Release 释放资源
func (a *Auth) Release() error {
	return a.callStore(func(store BlackStorer) error {
		return store.Close()
	})
}
