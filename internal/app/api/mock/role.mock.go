package mock

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var RoleSet = wire.NewSet(wire.Struct(new(RoleMock), "*"))

// RoleMock 角色管理
type RoleMock struct {
}

// Query 查询数据
// @Tags 角色管理
// @Summary 查询数据
// @Security ApiKeyAuth
// @Param current query int true "分页索引" default(1)
// @Param pageSize query int true "分页大小" default(10)
// @Param queryValue query string false "查询值"
// @Param status query int false "状态(1:启用 2:禁用)"
// @Success 200 {object} schema.ListResult{list=[]schema.Role} "查询结果"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/roles [get]
func (a *RoleMock) Query(c *gin.Context) {
}

// QuerySelect 查询选择数据
// @Tags 角色管理
// @Summary 查询选择数据
// @Security ApiKeyAuth
// @Param queryValue query string false "查询值"
// @Param status query int false "状态(1:启用 2:禁用)"
// @Success 200 {object} schema.ListResult{list=[]schema.Role} "查询结果"
// @Failure 400 {object} schema.ErrorResult "{error:{code:0,message:未知的查询类型}}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/roles.select [get]
func (a *RoleMock) QuerySelect(c *gin.Context) {
}

// Get 查询指定数据
// @Tags 角色管理
// @Summary 查询指定数据
// @Security ApiKeyAuth
// @Param id path int true "唯一标识"
// @Success 200 {object} schema.Role
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 404 {object} schema.ErrorResult "{error:{code:0,message:资源不存在}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/roles/{id} [get]
func (a *RoleMock) Get(c *gin.Context) {
}

// Create 创建数据
// @Tags 角色管理
// @Summary 创建数据
// @Security ApiKeyAuth
// @Param body body schema.Role true "创建数据"
// @Success 200 {object} schema.IDResult
// @Failure 400 {object} schema.ErrorResult "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/roles [post]
func (a *RoleMock) Create(c *gin.Context) {
}

// Update 更新数据
// @Tags 角色管理
// @Summary 更新数据
// @Security ApiKeyAuth
// @Param id path int true "唯一标识"
// @Param body body schema.Role true "更新数据"
// @Success 200 {object} schema.Role
// @Failure 400 {object} schema.ErrorResult "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/roles/{id} [put]
func (a *RoleMock) Update(c *gin.Context) {
}

// Delete 删除数据
// @Tags 角色管理
// @Summary 删除数据
// @Security ApiKeyAuth
// @Param id path int true "唯一标识"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/roles/{id} [delete]
func (a *RoleMock) Delete(c *gin.Context) {
}

// Enable 启用数据
// @Tags 角色管理
// @Summary 启用数据
// @Security ApiKeyAuth
// @Param id path int true "唯一标识"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/roles/{id}/enable [patch]
func (a *RoleMock) Enable(c *gin.Context) {
}

// Disable 禁用数据
// @Tags 角色管理
// @Summary 禁用数据
// @Security ApiKeyAuth
// @Param id path int true "唯一标识"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/roles/{id}/disable [patch]
func (a *RoleMock) Disable(c *gin.Context) {
}
