package ctl

import (
	"strings"

	"github.com/LyricTian/gin-admin/internal/app/bll"
	"github.com/LyricTian/gin-admin/internal/app/ginplus"
	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/LyricTian/gin-admin/pkg/util"
	"github.com/gin-gonic/gin"
)

// NewUser 创建用户管理控制器
func NewUser(bUser bll.IUser) *User {
	return &User{
		UserBll: bUser,
	}
}

// User 用户管理
type User struct {
	UserBll bll.IUser
}

// Query 查询数据
// @Tags 用户管理
// @Summary 查询数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param current query int true "分页索引" default(1)
// @Param pageSize query int true "分页大小" default(10)
// @Param userName query string false "用户名(模糊查询)"
// @Param realName query string false "真实姓名(模糊查询)"
// @Param roleIDs query string false "角色ID(多个以英文逗号分隔)"
// @Param status query int false "状态(1:启用 2:停用)"
// @Success 200 {array} schema.UserShow "查询结果：{list:列表数据,pagination:{current:页索引,pageSize:页大小,total:总数量}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/users [get]
func (a *User) Query(c *gin.Context) {
	var params schema.UserQueryParam
	params.LikeUserName = c.Query("userName")
	params.LikeRealName = c.Query("realName")
	if v := util.S(c.Query("status")).DefaultInt(0); v > 0 {
		params.Status = v
	}

	if v := c.Query("roleIDs"); v != "" {
		params.RoleIDs = strings.Split(v, ",")
	}

	result, err := a.UserBll.QueryShow(ginplus.NewContext(c), params, schema.UserQueryOptions{
		IncludeRoles: true,
		PageParam:    ginplus.GetPaginationParam(c),
	})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResPage(c, result.Data, result.PageResult)
}

// Get 查询指定数据
// Get 查询指定数据
// @Tags 用户管理
// @Summary 查询指定数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path string true "记录ID"
// @Success 200 {object} schema.User
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 404 {object} schema.HTTPError "{error:{code:0,message:资源不存在}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/users/{id} [get]
func (a *User) Get(c *gin.Context) {
	item, err := a.UserBll.Get(ginplus.NewContext(c), c.Param("id"), schema.UserQueryOptions{
		IncludeRoles: true,
	})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, item.CleanSecure())
}

// Create 创建数据
// @Tags 用户管理
// @Summary 创建数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param body body schema.User true "创建数据"
// @Success 200 {object} schema.User
// @Failure 400 {object} schema.HTTPError "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/users [post]
func (a *User) Create(c *gin.Context) {
	var item schema.User
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	item.Creator = ginplus.GetUserID(c)
	nitem, err := a.UserBll.Create(ginplus.NewContext(c), item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nitem.CleanSecure())
}

// Update 更新数据
// @Tags 用户管理
// @Summary 更新数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path string true "记录ID"
// @Param body body schema.User true "更新数据"
// @Success 200 {object} schema.User
// @Failure 400 {object} schema.HTTPError "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/users/{id} [put]
func (a *User) Update(c *gin.Context) {
	var item schema.User
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	nitem, err := a.UserBll.Update(ginplus.NewContext(c), c.Param("id"), item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nitem.CleanSecure())
}

// Delete 删除数据
// @Tags 用户管理
// @Summary 删除数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path string true "记录ID"
// @Success 200 {object} schema.HTTPStatus "{status:OK}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/users/{id} [delete]
func (a *User) Delete(c *gin.Context) {
	err := a.UserBll.Delete(ginplus.NewContext(c), c.Param("id"))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

// Enable 启用数据
// @Tags 用户管理
// @Summary 启用数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path string true "记录ID"
// @Success 200 {object} schema.HTTPStatus "{status:OK}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/users/{id}/enable [patch]
func (a *User) Enable(c *gin.Context) {
	err := a.UserBll.UpdateStatus(ginplus.NewContext(c), c.Param("id"), 1)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

// Disable 禁用数据
// @Tags 用户管理
// @Summary 禁用数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path string true "记录ID"
// @Success 200 {object} schema.HTTPStatus "{status:OK}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/users/{id}/disable [patch]
func (a *User) Disable(c *gin.Context) {
	err := a.UserBll.UpdateStatus(ginplus.NewContext(c), c.Param("id"), 2)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}
