package api

import (
	"github.com/LyricTian/gin-admin/v10/internal/mods/rbac/biz"
	"github.com/LyricTian/gin-admin/v10/internal/mods/rbac/schema"
	"github.com/LyricTian/gin-admin/v10/pkg/util"
	"github.com/gin-gonic/gin"
)

// Logger management
type Logger struct {
	LoggerBIZ *biz.Logger
}

// @Tags LoggerAPI
// @Security ApiKeyAuth
// @Summary Query logger list
// @Param current query int true "pagination index" default(1)
// @Param pageSize query int true "pagination size" default(10)
// @Param level query string false "log level"
// @Param traceID query string false "trace ID"
// @Param userName query string false "user name"
// @Param tag query string false "log tag"
// @Param message query string false "log message"
// @Param startTime query string false "start time"
// @Param endTime query string false "end time"
// @Success 200 {object} util.ResponseResult{data=[]schema.Logger}
// @Failure 401 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/loggers [get]
func (a *Logger) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.LoggerQueryParam
	if err := util.ParseQuery(c, &params); err != nil {
		util.ResError(c, err)
		return
	}

	result, err := a.LoggerBIZ.Query(ctx, params)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResPage(c, result.Data, result.PageResult)
}
