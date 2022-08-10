package api

import (
	"github.com/LyricTian/gin-admin/v9/internal/module/sys/biz"
	"github.com/LyricTian/gin-admin/v9/internal/module/sys/typed"
	"github.com/LyricTian/gin-admin/v9/internal/x/utilx"
	"github.com/gin-gonic/gin"
)

type DictionaryAPI struct {
	DictionaryBiz *biz.DictionaryBiz
}

// @Tags DictionaryAPI
// @Security ApiKeyAuth
// @Summary Query dictionary list
// @Param current query int true "pagination index" default(1)
// @Param pageSize query int true "pagination size" default(10)
// @Param namespace query string false "query namespace"
// @Param key query string false "query key"
// @Success 200 {object} utilx.ResponseResult{data=[]typed.Dictionary} "query result"
// @Failure 401 {object} utilx.ResponseResult
// @Failure 500 {object} utilx.ResponseResult
// @Router /api/sys/v1/dictionaries [get]
func (a *DictionaryAPI) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params typed.DictionaryQueryParam
	if err := utilx.ParseQuery(c, &params); err != nil {
		utilx.ResError(c, err)
		return
	}

	result, err := a.DictionaryBiz.Query(ctx, params)
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResPage(c, result.Data, result.PageResult)
}

// @Tags DictionaryAPI
// @Security ApiKeyAuth
// @Summary Get single dictionary by id
// @Param id path string true "unique id"
// @Success 200 {object} utilx.ResponseResult{data=typed.Dictionary}
// @Failure 401 {object} utilx.ResponseResult
// @Failure 500 {object} utilx.ResponseResult
// @Router /api/sys/v1/dictionaries/{id} [get]
func (a *DictionaryAPI) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.DictionaryBiz.Get(ctx, c.Param("id"))
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResSuccess(c, item)
}

// @Tags DictionaryAPI
// @Security ApiKeyAuth
// @Summary Create dictionary
// @Param body body typed.DictionaryCreate true "request body"
// @Success 200 {object} utilx.ResponseResult{data=typed.Dictionary}
// @Failure 400 {object} utilx.ResponseResult
// @Failure 401 {object} utilx.ResponseResult
// @Failure 500 {object} utilx.ResponseResult
// @Router /api/sys/v1/dictionaries [post]
func (a *DictionaryAPI) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var item typed.DictionaryCreate
	if err := utilx.ParseJSON(c, &item); err != nil {
		utilx.ResError(c, err)
		return
	}

	result, err := a.DictionaryBiz.Create(ctx, item)
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResSuccess(c, result)
}

// @Tags DictionaryAPI
// @Security ApiKeyAuth
// @Summary Update dictionary by id
// @Param id path int true "unique id"
// @Param body body typed.DictionaryCreate true "request body"
// @Success 200 {object} utilx.ResponseResult
// @Failure 400 {object} utilx.ResponseResult
// @Failure 401 {object} utilx.ResponseResult
// @Failure 500 {object} utilx.ResponseResult
// @Router /api/sys/v1/dictionaries/{id} [put]
func (a *DictionaryAPI) Update(c *gin.Context) {
	ctx := c.Request.Context()
	var item typed.DictionaryCreate
	if err := utilx.ParseJSON(c, &item); err != nil {
		utilx.ResError(c, err)
		return
	}

	err := a.DictionaryBiz.Update(ctx, c.Param("id"), item)
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResOK(c)
}

// @Tags DictionaryAPI
// @Security ApiKeyAuth
// @Summary Delete single dictionary by id
// @Param id path string true "unique id"
// @Success 200 {object} utilx.ResponseResult
// @Failure 401 {object} utilx.ResponseResult
// @Failure 500 {object} utilx.ResponseResult
// @Router /api/sys/v1/dictionaries/{id} [delete]
func (a *DictionaryAPI) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.DictionaryBiz.Delete(ctx, c.Param("id"))
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResOK(c)
}
