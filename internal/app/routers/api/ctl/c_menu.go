package ctl

import (
	"github.com/LyricTian/gin-admin/internal/app/bll"
	"github.com/LyricTian/gin-admin/internal/app/ginplus"
	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/gin-gonic/gin"
)

// NewMenu 创建菜单管理控制器
func NewMenu(bMenu bll.IMenu) *Menu {
	return &Menu{
		MenuBll: bMenu,
	}
}

// Menu 菜单管理
type Menu struct {
	MenuBll bll.IMenu
}

// Query 查询数据
// @Tags 菜单管理
// @Summary 查询数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param current query int true "分页索引" default(1)
// @Param pageSize query int true "分页大小" default(10)
// @Param likeName query string false "名称(模糊查询)"
// @Param status query int false "状态(1:正常 2:隐藏)"
// @Param parentID query string false "父级ID"
// @Success 200 {array} schema.Menu "查询结果：{list:列表数据,pagination:{current:页索引,pageSize:页大小,total:总数量}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/menus [get]
func (a *Menu) Query(c *gin.Context) {
	var params schema.MenuQueryParam
	if err := ginplus.ParseQuery(c, &params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	result, err := a.MenuBll.Query(ginplus.NewContext(c), params, schema.MenuQueryOptions{
		PageParam:   ginplus.GetPaginationParam(c),
		OrderFields: schema.NewOrderFields(map[string]schema.OrderDirection{"sequence": schema.OrderByDESC}),
	})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResPage(c, result.Data, result.PageResult)
}

// QueryTree 查询菜单树
// @Tags 菜单管理
// @Summary 查询菜单树
// @Param Authorization header string false "Bearer 用户令牌"
// @Param likeName query string false "名称(模糊查询)"
// @Param status query int false "状态(1:正常 2:隐藏)"
// @Param parentID query string false "父级ID"
// @Success 200 {array} schema.MenuTree "查询结果：{list:列表数据}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/menus.tree [get]
func (a *Menu) QueryTree(c *gin.Context) {
	var params schema.MenuQueryParam
	if err := ginplus.ParseQuery(c, &params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	result, err := a.MenuBll.Query(ginplus.NewContext(c), params, schema.MenuQueryOptions{
		OrderFields: schema.NewOrderFields(map[string]schema.OrderDirection{"sequence": schema.OrderByDESC}),
	})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResList(c, result.Data.ToTree())
}

// Get 查询指定数据
// @Tags 菜单管理
// @Summary 查询指定数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path string true "记录ID"
// @Success 200 {object} schema.Menu
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 404 {object} schema.HTTPError "{error:{code:0,message:资源不存在}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/menus/{id} [get]
func (a *Menu) Get(c *gin.Context) {
	item, err := a.MenuBll.Get(ginplus.NewContext(c), c.Param("id"))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, item)
}

// Create 创建数据
// @Tags 菜单管理
// @Summary 创建数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param body body schema.Menu true "创建数据"
// @Success 200 {object} schema.Menu
// @Failure 400 {object} schema.HTTPError "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/menus [post]
func (a *Menu) Create(c *gin.Context) {
	var item schema.Menu
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	item.Creator = ginplus.GetUserID(c)
	nitem, err := a.MenuBll.Create(ginplus.NewContext(c), item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, nitem)
}

// Update 更新数据
// @Tags 菜单管理
// @Summary 更新数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path string true "记录ID"
// @Param body body schema.Menu true "更新数据"
// @Success 200 {object} schema.HTTPStatus "{status:OK}"
// @Failure 400 {object} schema.HTTPError "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/menus/{id} [put]
func (a *Menu) Update(c *gin.Context) {
	var item schema.Menu
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	err := a.MenuBll.Update(ginplus.NewContext(c), c.Param("id"), item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResOK(c)
}

// Delete 删除数据
// @Tags 菜单管理
// @Summary 删除数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path string true "记录ID"
// @Success 200 {object} schema.HTTPStatus "{status:OK}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/menus/{id} [delete]
func (a *Menu) Delete(c *gin.Context) {
	err := a.MenuBll.Delete(ginplus.NewContext(c), c.Param("id"))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}
