package ctl

import (
	"strings"

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
// @Param Access-Token header string false "访问令牌"
// @Param current query int true "分页索引" 1
// @Param pageSize query int true "分页大小" 10
// @Param code query string false "编号"
// @Param name query string false "名称"
// @Param parent_id query string false "父级ID"
// @Param type query int false "菜单类型(1：模块 2：功能 3：资源)"
// @Success 200 []schema.Menu "分页查询结果：{list:列表数据,pagination:{current:页索引,pageSize:页大小,total:总数量}}"
// @Failure 400 option.Interface "{error:{code:0,message:未知的查询类型}}"
// @Failure 401 option.Interface "{error:{code:0,message:未授权}}"
// @Failure 500 option.Interface "{error:{code:0,message:服务器错误}}"
// @Router GET /api/v1/menus?q=page
func (a *Menu) QueryPage(ctx *context.Context) {
	params := schema.MenuQueryParam{
		Code: ctx.Query("code"),
		Name: ctx.Query("name"),
	}

	if v := ctx.Query("parent_id"); v != "" {
		params.ParentID = &v
	}

	if v := ctx.Query("type"); v != "" {
		params.Types = []int{util.S(v).Int()}
	}

	items, pr, err := a.MenuBll.QueryPage(ctx.CContext(), params, ctx.GetPaginationParam())
	if err != nil {
		ctx.ResError(err)
		return
	}

	ctx.ResPage(items, pr)
}

// QueryTree 查询菜单树
// @Summary 查询菜单树
// @Param Access-Token header string false "访问令牌"
// @Success 200 option.Interface "查询结果：{list:[{record_id:记录ID,name:名称,parent_id:父级ID,children:子级树}]}"
// @Failure 400 option.Interface "{error:{code:0,message:未知的查询类型}}"
// @Failure 401 option.Interface "{error:{code:0,message:未授权}}"
// @Failure 500 option.Interface "{error:{code:0,message:服务器错误}}"
// @Router GET /api/v1/menus?q=tree
func (a *Menu) QueryTree(ctx *context.Context) {
	treeData, err := a.MenuBll.QueryTree(ctx.CContext())
	if err != nil {
		ctx.ResError(err)
		return
	}

	ctx.ResList(treeData)
}

// Get 查询指定数据
// @Summary 查询指定数据
// @Param Access-Token header string false "访问令牌"
// @Param id path string true "记录ID"
// @Success 200 schema.Menu
// @Failure 401 option.Interface "{error:{code:0,message:未授权}}"
// @Failure 404 string "{error:{code:0,message:资源不存在}}"
// @Failure 500 option.Interface "{error:{code:0,message:服务器错误}}"
// @Router GET /api/v1/menus/{id}
func (a *Menu) Get(ctx *context.Context) {
	item, err := a.MenuBll.Get(ctx.CContext(), ctx.Param("id"))
	if err != nil {
		ctx.ResError(err)
		return
	}
	ctx.ResSuccess(item)
}

// Create 创建数据
// @Summary 创建数据
// @Param Access-Token header string false "访问令牌"
// @Param body body schema.Menu true
// @Success 200 option.Interface "{record_id:记录ID}"
// @Failure 400 option.Interface "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 option.Interface "{error:{code:0,message:未授权}}"
// @Failure 500 option.Interface "{error:{code:0,message:服务器错误}}"
// @Router POST /api/v1/menus
func (a *Menu) Create(ctx *context.Context) {
	var item schema.Menu
	if err := ctx.ParseJSON(&item); err != nil {
		ctx.ResError(err)
		return
	}

	newItem, err := a.MenuBll.Create(ctx.CContext(), item)
	if err != nil {
		ctx.ResError(err)
		return
	}

	ctx.ResSuccess(context.HTTPNewItem{RecordID: newItem.RecordID})
}

// Update 更新数据
// @Summary 更新数据
// @Param Access-Token header string false "访问令牌"
// @Param id path string true "记录ID"
// @Param body body schema.Menu true
// @Success 200 option.Interface "{status:OK}"
// @Failure 400 option.Interface "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 option.Interface "{error:{code:0,message:未授权}}"
// @Failure 500 option.Interface "{error:{code:0,message:服务器错误}}"
// @Router PUT /api/v1/menus/{id}
func (a *Menu) Update(ctx *context.Context) {
	var item schema.Menu
	if err := ctx.ParseJSON(&item); err != nil {
		ctx.ResError(err)
		return
	}

	err := a.MenuBll.Update(ctx.CContext(), ctx.Param("id"), item)
	if err != nil {
		ctx.ResError(err)
		return
	}
	ctx.ResOK()
}

// Delete 删除数据
// @Summary 删除数据
// @Param Access-Token header string false "访问令牌"
// @Param id path string true "记录ID"
// @Success 200 option.Interface "{status:OK}"
// @Failure 401 option.Interface "{error:{code:0,message:未授权}}"
// @Failure 500 option.Interface "{error:{code:0,message:服务器错误}}"
// @Router DELETE /api/v1/menus/{id}
func (a *Menu) Delete(ctx *context.Context) {
	err := a.MenuBll.Delete(ctx.CContext(), ctx.Param("id"))
	if err != nil {
		ctx.ResError(err)
		return
	}
	ctx.ResOK()
}

// DeleteMany 删除多条数据
// @Summary 删除多条数据
// @Param Access-Token header string false "访问令牌"
// @Param batch query string true "记录ID（多个以,分隔）"
// @Success 200 option.Interface "{status:OK}"
// @Failure 400 option.Interface "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 option.Interface "{error:{code:0,message:未授权}}"
// @Failure 500 option.Interface "{error:{code:0,message:服务器错误}}"
// @Router DELETE /api/v1/menus
func (a *Menu) DeleteMany(ctx *context.Context) {
	ids := strings.Split(ctx.Query("batch"), ",")
	if len(ids) == 0 {
		ctx.ResError(errors.NewBadRequestError("无效的请求参数"))
		return
	}

	err := a.MenuBll.Delete(ctx.CContext(), ids...)
	if err != nil {
		ctx.ResError(err)
		return
	}

	ctx.ResOK()
}
