package ctl

import (
	"github.com/LyricTian/gin-admin/internal/app/bll"
	"github.com/LyricTian/gin-admin/internal/app/ginplus"
	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/gin-gonic/gin"
)

// NewRole 创建角色管理控制器
func NewRole(bRole bll.IRole) *Role {
	return &Role{
		RoleBll: bRole,
	}
}

// Role 角色管理
type Role struct {
	RoleBll bll.IRole
}

// Query 查询数据
// @Tags 角色管理
// @Summary 查询数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param current query int true "分页索引" default(1)
// @Param pageSize query int true "分页大小" default(10)
// @Param name query string false "角色名称(模糊查询)"
// @Success 200 {array} schema.Role "查询结果：{list:列表数据,pagination:{current:页索引,pageSize:页大小,total:总数量}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/roles [get]
func (a *Role) Query(c *gin.Context) {
	var params schema.RoleQueryParam
	params.LikeName = c.Query("name")

	result, err := a.RoleBll.Query(ginplus.NewContext(c), params, schema.RoleQueryOptions{
		PageParam: ginplus.GetPaginationParam(c),
	})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResPage(c, result.Data, result.PageResult)
}

// QuerySelect 查询选择数据
// @Tags 角色管理
// @Summary 查询选择数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Success 200 {array} schema.Role "查询结果：{list:角色列表}"
// @Failure 400 {object} schema.HTTPError "{error:{code:0,message:未知的查询类型}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/roles.select [get]
func (a *Role) QuerySelect(c *gin.Context) {
	result, err := a.RoleBll.Query(ginplus.NewContext(c), schema.RoleQueryParam{})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResList(c, result.Data)
}

// Get 查询指定数据
// @Tags 角色管理
// @Summary 查询指定数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path string true "记录ID"
// @Success 200 {object} schema.Role
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 404 {object} schema.HTTPError "{error:{code:0,message:资源不存在}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/roles/{id} [get]
func (a *Role) Get(c *gin.Context) {
	item, err := a.RoleBll.Get(ginplus.NewContext(c), c.Param("id"), schema.RoleQueryOptions{
		IncludeMenus: true,
	})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, item)
}

// Create 创建数据
// @Tags 角色管理
// @Summary 创建数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param body body schema.Role true "创建数据"
// @Success 200 {object} schema.Role
// @Failure 400 {object} schema.HTTPError "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/roles [post]
func (a *Role) Create(c *gin.Context) {
	var item schema.Role
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	item.Creator = ginplus.GetUserID(c)
	nitem, err := a.RoleBll.Create(ginplus.NewContext(c), item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, nitem)
}

// Update 更新数据
// @Tags 角色管理
// @Summary 更新数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path string true "记录ID"
// @Param body body schema.Role true "更新数据"
// @Success 200 {object} schema.Role
// @Failure 400 {object} schema.HTTPError "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/roles/{id} [put]
func (a *Role) Update(c *gin.Context) {
	var item schema.Role
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	nitem, err := a.RoleBll.Update(ginplus.NewContext(c), c.Param("id"), item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nitem)
}

// Delete 删除数据
// @Tags 角色管理
// @Summary 删除数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path string true "记录ID"
// @Success 200 {object} schema.HTTPStatus "{status:OK}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/roles/{id} [delete]
func (a *Role) Delete(c *gin.Context) {
	err := a.RoleBll.Delete(ginplus.NewContext(c), c.Param("id"))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}
