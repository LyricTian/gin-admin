package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/LyricTian/gin-admin/v9/internal/config"
	"github.com/LyricTian/gin-admin/v9/internal/inject"
	"github.com/LyricTian/gin-admin/v9/internal/middleware"
	"github.com/gin-gonic/gin"
)

var (
	engine *gin.Engine
)

func init() {
	ctx := context.Background()

	config.MustLoad("config_test.toml")
	config.C.Middleware.Auth.Disable = true
	config.C.Middleware.Casbin.Disable = true

	injector, _, err := inject.BuildInjector(ctx)
	if err != nil {
		panic(err)
	}

	err = injector.RBAC.Init(ctx)
	if err != nil {
		panic(err)
	}

	engine = gin.New()
	g := engine.Group("/api")
	{
		g.Use(middleware.AuthWithConfig(middleware.AuthConfig{
			Disable: config.C.Middleware.Auth.Disable,
			DefaultUserID: func(c *gin.Context) string {
				return config.C.Dictionary.RootUser.ID
			},
		}))

		// Register RBAC APIs
		injector.RBAC.RegisterAPI(ctx, g)
	} // end
}

func newPostRequest(v interface{}, formatRouter string, args ...interface{}) *http.Request {
	buf := new(bytes.Buffer)
	_ = json.NewEncoder(buf).Encode(v)
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf(formatRouter, args...), buf)
	return req
}

func newPutRequest(v interface{}, formatRouter string, args ...interface{}) *http.Request {
	buf := new(bytes.Buffer)
	_ = json.NewEncoder(buf).Encode(v)
	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(formatRouter, args...), buf)
	return req
}

func newDeleteRequest(formatRouter string, args ...interface{}) *http.Request {
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf(formatRouter, args...), nil)
	return req
}

func newGetRequest(params map[string]string, formatRouter string, args ...interface{}) *http.Request {
	values := make(url.Values)
	for k, v := range params {
		values.Set(k, v)
	}

	urlStr := fmt.Sprintf(formatRouter, args...)
	if len(values) > 0 {
		urlStr += "?" + values.Encode()
	}

	req, _ := http.NewRequest(http.MethodGet, urlStr, nil)
	return req
}
