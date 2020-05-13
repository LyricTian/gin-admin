package mock

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// UserSet 注入User
var UserSet = wire.NewSet(wire.Struct(new(User), "*"))

// User 用户管理
type User struct {
}

// Query 查询数据
// @Tags 用户管理
// @Summary 查询数据
// @Security ApiKeyAuth
// @Param current query int true "分页索引" default(1)
// @Param pageSize query int true "分页大小" default(10)
// @Param queryValue query string false "查询值"
// @Param roleIDs query string false "角色ID(多个以英文逗号分隔)"
// @Param status query int false "状态(1:启用 2:停用)"
// @Success 200 {array} schema.UserShow "查询结果：{list:列表数据,pagination:{current:页索引,pageSize:页大小,total:总数量}}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/users [get]
func (a *User) Query(c *gin.Context) {
}

// Get 查询指定数据
// Get 查询指定数据
// @Tags 用户管理
// @Summary 查询指定数据
// @Security ApiKeyAuth
// @Param id path string true "唯一标识"
// @Success 200 {object} schema.User
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 404 {object} schema.ErrorResult "{error:{code:0,message:资源不存在}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/users/{id} [get]
func (a *User) Get(c *gin.Context) {
}

// Create 创建数据
// @Tags 用户管理
// @Summary 创建数据
// @Security ApiKeyAuth
// @Param body body schema.User true "创建数据"
// @Success 200 {object} schema.IDResult
// @Failure 400 {object} schema.ErrorResult "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/users [post]
func (a *User) Create(c *gin.Context) {
}

// Update 更新数据
// @Tags 用户管理
// @Summary 更新数据
// @Security ApiKeyAuth
// @Param id path string true "唯一标识"
// @Param body body schema.User true "更新数据"
// @Success 200 {object} schema.User
// @Failure 400 {object} schema.ErrorResult "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/users/{id} [put]
func (a *User) Update(c *gin.Context) {
}

// Delete 删除数据
// @Tags 用户管理
// @Summary 删除数据
// @Security ApiKeyAuth
// @Param id path string true "唯一标识"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/users/{id} [delete]
func (a *User) Delete(c *gin.Context) {
}

// Enable 启用数据
// @Tags 用户管理
// @Summary 启用数据
// @Security ApiKeyAuth
// @Param id path string true "唯一标识"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/users/{id}/enable [patch]
func (a *User) Enable(c *gin.Context) {
}

// Disable 禁用数据
// @Tags 用户管理
// @Summary 禁用数据
// @Security ApiKeyAuth
// @Param id path string true "唯一标识"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/users/{id}/disable [patch]
func (a *User) Disable(c *gin.Context) {
}
