package jwtauth

import (
	"time"

	"github.com/LyricTian/gin-admin/pkg/auth"
	jwt "github.com/dgrijalva/jwt-go"
)

const defaultKey = "GINADMIN"

var defaultOptions = options{
	tokenType:     "Bearer",
	expired:       7200,
	signingMethod: jwt.SigningMethodHS512,
	signingKey:    defaultKey,
	keyfunc: func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, auth.ErrInvalidToken
		}
		return []byte(defaultKey), nil
	},
}

type options struct {
	signingMethod jwt.SigningMethod
	signingKey    interface{}
	keyfunc       jwt.Keyfunc
	expired       int
	tokenType     string
}

// Option 定义参数项
type Option func(*options)

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
func New(store Storer, opts ...Option) *JWTAuth {
	o := defaultOptions
	for _, opt := range opts {
		opt(&o)
	}

	return &JWTAuth{
		opts: &o,
	}
}

// JWTAuth jwt认证
type JWTAuth struct {
	opts  *options
	store Storer
}

// GenerateToken 生成令牌
func (a *JWTAuth) GenerateToken(userID string) (auth.TokenInfo, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(a.opts.expired) * time.Second).Unix()

	token := jwt.NewWithClaims(a.opts.signingMethod, &jwt.StandardClaims{
		IssuedAt:  now.Unix(),
		ExpiresAt: expiresAt,
		NotBefore: now.Unix(),
		Subject:   userID,
	})

	tokenString, err := token.SignedString(a.opts.signingKey)
	if err != nil {
		return nil, err
	}

	tokenInfo := &tokenInfo{
		ExpiresAt:   expiresAt,
		TokenType:   a.opts.tokenType,
		AccessToken: tokenString,
	}
	return tokenInfo, nil
}

// 解析令牌
func (a *JWTAuth) parseToken(tokenString string) (*jwt.StandardClaims, error) {
	token, _ := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, a.opts.keyfunc)
	if !token.Valid {
		return nil, auth.ErrInvalidToken
	}

	return token.Claims.(*jwt.StandardClaims), nil
}

func (a *JWTAuth) callStore(fn func(Storer) error) error {
	if store := a.store; store != nil {
		return fn(store)
	}
	return nil
}

// DestroyToken 销毁令牌
func (a *JWTAuth) DestroyToken(tokenString string) error {
	claims, err := a.parseToken(tokenString)
	if err != nil {
		return err
	}

	// 如果设定了存储，则将未过期的令牌放入
	return a.callStore(func(store Storer) error {
		expired := time.Unix(claims.ExpiresAt, 0).Sub(time.Now())
		return store.Set(tokenString, expired)
	})
}

// ParseUserID 解析用户ID
func (a *JWTAuth) ParseUserID(tokenString string) (string, error) {
	claims, err := a.parseToken(tokenString)
	if err != nil {
		return "", err
	}

	err = a.callStore(func(store Storer) error {
		exists, err := store.Check(tokenString)
		if err != nil {
			return err
		} else if exists {
			return auth.ErrInvalidToken
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	return claims.Subject, nil
}

// Release 释放资源
func (a *JWTAuth) Release() error {
	return a.callStore(func(store Storer) error {
		return store.Close()
	})
}
