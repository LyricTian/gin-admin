package mock

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// DemoSet 注入Demo
var DemoSet = wire.NewSet(wire.Struct(new(Demo), "*"))

// Demo 示例程序
type Demo struct {
}

// Query 查询数据
// @Tags Demo
// @Security ApiKeyAuth
// @Summary 查询数据
// @Param current query int true "分页索引" default(1)
// @Param pageSize query int true "分页大小" default(10)
// @Param queryValue query string false "查询值"
// @Success 200 {array} schema.Demo "查询结果：{list:列表数据,pagination:{current:页索引,pageSize:页大小,total:总数量}}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/demos [get]
func (a *Demo) Query(c *gin.Context) {
}

// Get 查询指定数据
// @Tags Demo
// @Security ApiKeyAuth
// @Summary 查询指定数据
// @Param id path string true "唯一标识"
// @Success 200 {object} schema.Demo
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 404 {object} schema.ErrorResult "{error:{code:0,message:资源不存在}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/demos/{id} [get]
func (a *Demo) Get(c *gin.Context) {
}

// Create 创建数据
// @Tags Demo
// @Security ApiKeyAuth
// @Summary 创建数据
// @Param body body schema.Demo true "创建数据"
// @Success 200 {object} schema.IDResult
// @Failure 400 {object} schema.ErrorResult "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/demos [post]
func (a *Demo) Create(c *gin.Context) {
}

// Update 更新数据
// @Tags Demo
// @Security ApiKeyAuth
// @Summary 更新数据
// @Param id path string true "唯一标识"
// @Param body body schema.Demo true "更新数据"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 400 {object} schema.ErrorResult "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/demos/{id} [put]
func (a *Demo) Update(c *gin.Context) {
}

// Delete 删除数据
// @Tags Demo
// @Security ApiKeyAuth
// @Summary 删除数据
// @Param id path string true "唯一标识"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/demos/{id} [delete]
func (a *Demo) Delete(c *gin.Context) {
}

// Enable 启用数据
// @Tags Demo
// @Security ApiKeyAuth
// @Summary 启用数据
// @Param id path string true "唯一标识"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/demos/{id}/enable [patch]
func (a *Demo) Enable(c *gin.Context) {
}

// Disable 禁用数据
// @Tags Demo
// @Security ApiKeyAuth
// @Summary 禁用数据
// @Param id path string true "唯一标识"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/demos/{id}/disable [patch]
func (a *Demo) Disable(c *gin.Context) {
}
