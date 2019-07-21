package middleware

import (
	"time"

	"github.com/LyricTian/gin-admin/internal/app/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware 跨域请求中间件
func CORSMiddleware() gin.HandlerFunc {
	cfg := config.GetGlobalConfig().CORS
	return cors.New(cors.Config{
		AllowOrigins:     cfg.AllowOrigins,
		AllowMethods:     cfg.AllowMethods,
		AllowHeaders:     cfg.AllowHeaders,
		AllowCredentials: cfg.AllowCredentials,
		MaxAge:           time.Second * time.Duration(cfg.MaxAge),
	})
}
