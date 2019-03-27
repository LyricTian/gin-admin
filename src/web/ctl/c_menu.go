package ctl

import (
	"github.com/LyricTian/gin-admin/src/bll"
	"github.com/LyricTian/gin-admin/src/errors"
	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/LyricTian/gin-admin/src/util"
	"github.com/LyricTian/gin-admin/src/web/context"
)

// Menu 菜单管理
// @Name Menu
// @Description 菜单管理接口
type Menu struct {
	MenuBll *bll.Menu `inject:""`
}

// Query 查询数据
func (a *Menu) Query(ctx *context.Context) {
	switch ctx.Query("q") {
	case "page":
		a.QueryPage(ctx)
	case "tree":
		a.QueryTree(ctx)
	default:
		ctx.ResError(errors.NewBadRequestError("未知的查询类型"))
	}
}

// QueryPage 查询分页数据
// @Summary 查询分页数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param current query int true "分页索引" 1
// @Param pageSize query int true "分页大小" 10
// @Param name query string false "名称"
// @Param hidden query int false "隐藏菜单(0:不隐藏 1:隐藏)"
// @Param parent_id query string false "父级ID"
// @Success 200 []schema.Menu "分页查询结果：{list:列表数据,pagination:{current:页索引,pageSize:页大小,total:总数量}}"
// @Failure 400 schema.HTTPError "{error:{code:0,message:未知的查询类型}}"
// @Failure 401 schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router GET /api/v1/menus?q=page
func (a *Menu) QueryPage(ctx *context.Context) {
	params := schema.MenuQueryParam{
		LikeName: ctx.Query("name"),
	}

	if v := ctx.Query("parent_id"); v != "" {
		params.ParentID = &v
	}

	if v := ctx.Query("hidden"); v != "" {
		if hidden := util.S(v).Int(); hidden > -1 {
			params.Hidden = &hidden
		}
	}

	items, pr, err := a.MenuBll.QueryPage(ctx.GetContext(), params, ctx.GetPaginationParam())
	if err != nil {
		ctx.ResError(err)
		return
	}
	ctx.ResPage(items, pr)
}

// QueryTree 查询菜单树
// @Summary 查询菜单树
// @Param Authorization header string false "Bearer 用户令牌"
// @Param include_actions query int false "是否包含动作数据(1是)"
// @Param include_resources query int false "是否包含资源数据(1是)"
// @Success 200 option.Interface "查询结果：{list:菜单树}"
// @Failure 400 schema.HTTPError "{error:{code:0,message:未知的查询类型}}"
// @Failure 401 schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router GET /api/v1/menus?q=tree
func (a *Menu) QueryTree(ctx *context.Context) {
	treeData, err := a.MenuBll.QueryTree(ctx.GetContext(), ctx.Query("include_actions") == "1", ctx.Query("include_resources") == "1")
	if err != nil {
		ctx.ResError(err)
		return
	}

	ctx.ResList(treeData)
}

// Get 查询指定数据
// @Summary 查询指定数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path string true "记录ID"
// @Success 200 schema.Menu
// @Failure 401 schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 404 schema.HTTPError "{error:{code:0,message:资源不存在}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router GET /api/v1/menus/{id}
func (a *Menu) Get(ctx *context.Context) {
	item, err := a.MenuBll.Get(ctx.GetContext(), ctx.Param("id"))
	if err != nil {
		ctx.ResError(err)
		return
	}
	ctx.ResSuccess(item)
}

// Create 创建数据
// @Summary 创建数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param body body schema.Menu true
// @Success 200 schema.Menu
// @Failure 400 schema.HTTPError "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router POST /api/v1/menus
func (a *Menu) Create(ctx *context.Context) {
	var item schema.Menu
	if err := ctx.ParseJSON(&item); err != nil {
		ctx.ResError(err)
		return
	}

	nitem, err := a.MenuBll.Create(ctx.GetContext(), item)
	if err != nil {
		ctx.ResError(err)
		return
	}
	ctx.ResSuccess(nitem)
}

// Update 更新数据
// @Summary 更新数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path string true "记录ID"
// @Param body body schema.Menu true
// @Success 200 schema.Menu
// @Failure 400 schema.HTTPError "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router PUT /api/v1/menus/{id}
func (a *Menu) Update(ctx *context.Context) {
	var item schema.Menu
	if err := ctx.ParseJSON(&item); err != nil {
		ctx.ResError(err)
		return
	}

	nitem, err := a.MenuBll.Update(ctx.GetContext(), ctx.Param("id"), item)
	if err != nil {
		ctx.ResError(err)
		return
	}
	ctx.ResSuccess(nitem)
}

// Delete 删除数据
// @Summary 删除数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path string true "记录ID"
// @Success 200 schema.HTTPStatus "{status:OK}"
// @Failure 401 schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router DELETE /api/v1/menus/{id}
func (a *Menu) Delete(ctx *context.Context) {
	err := a.MenuBll.Delete(ctx.GetContext(), ctx.Param("id"))
	if err != nil {
		ctx.ResError(err)
		return
	}
	ctx.ResOK()
}
