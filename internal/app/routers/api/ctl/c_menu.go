package ctl

import (
	"github.com/LyricTian/gin-admin/internal/app/bll"
	"github.com/LyricTian/gin-admin/internal/app/ginplus"
	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/LyricTian/gin-admin/pkg/util"
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
// @Param name query string false "名称"
// @Param hidden query int false "隐藏菜单(0:不隐藏 1:隐藏)"
// @Param parentID query string false "父级ID"
// @Success 200 {array} schema.Menu "查询结果：{list:列表数据,pagination:{current:页索引,pageSize:页大小,total:总数量}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/menus [get]
func (a *Menu) Query(c *gin.Context) {
	params := schema.MenuQueryParam{
		LikeName: c.Query("name"),
	}

	if v := c.Query("parentID"); v != "" {
		params.ParentID = &v
	}

	if v := c.Query("hidden"); v != "" {
		if hidden := util.S(v).DefaultInt(0); hidden > -1 {
			params.Hidden = &hidden
		}
	}

	result, err := a.MenuBll.Query(ginplus.NewContext(c), params, schema.MenuQueryOptions{
		PageParam: ginplus.GetPaginationParam(c),
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
// @Param includeActions query int false "是否包含动作数据(1是)"
// @Param includeResources query int false "是否包含资源数据(1是)"
// @Success 200 {array} schema.MenuTree "查询结果：{list:列表数据}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/menus.tree [get]
func (a *Menu) QueryTree(c *gin.Context) {
	result, err := a.MenuBll.Query(ginplus.NewContext(c), schema.MenuQueryParam{}, schema.MenuQueryOptions{
		IncludeActions:   c.Query("includeActions") == "1",
		IncludeResources: c.Query("includeResources") == "1",
	})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResList(c, result.Data.ToTrees().ToTree())
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
	item, err := a.MenuBll.Get(ginplus.NewContext(c), c.Param("id"), schema.MenuQueryOptions{
		IncludeActions:   true,
		IncludeResources: true,
	})
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
// @Success 200 {object} schema.Menu
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

	nitem, err := a.MenuBll.Update(ginplus.NewContext(c), c.Param("id"), item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nitem)
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
