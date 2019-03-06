package ctl

import (
	"strings"

	"github.com/LyricTian/gin-admin/src/bll"
	"github.com/LyricTian/gin-admin/src/errors"
	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/LyricTian/gin-admin/src/web/context"
)

// Resource 资源管理
// @Name Resource
// @Description 资源管理接口
type Resource struct {
	ResourceBll *bll.Resource `inject:""`
}

// Query 查询数据
func (a *Resource) Query(ctx *context.Context) {
	switch ctx.Query("q") {
	case "page":
		a.QueryPage(ctx)
	default:
		ctx.ResError(errors.NewBadRequestError("未知的查询类型"))
	}
}

// QueryPage 查询分页数据
// @Summary 查询分页数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param current query int true "分页索引" 1
// @Param pageSize query int true "分页大小" 10
// @Param name query string false "资源名称(模糊查询)"
// @Param path query string false "访问路径(前缀匹配)"
// @Success 200 []schema.Resource "查询结果：{list:列表数据,pagination:{current:页索引,pageSize:页大小,total:总数量}}"
// @Failure 400 option.Interface "{error:{code:0,message:未知的查询类型}}"
// @Failure 401 option.Interface "{error:{code:0,message:未授权}}"
// @Failure 500 option.Interface "{error:{code:0,message:服务器错误}}"
// @Router GET /api/v1/resources?q=page
func (a *Resource) QueryPage(ctx *context.Context) {
	var params schema.ResourceQueryParam
	params.Name = ctx.Query("name")
	params.Path = ctx.Query("path")

	items, pr, err := a.ResourceBll.Query(ctx.GetContext(), params, ctx.GetPaginationParam())
	if err != nil {
		ctx.ResError(err)
		return
	}

	ctx.ResPage(items, pr)
}

// Get 查询指定数据
// @Summary 查询指定数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path string true "记录ID"
// @Success 200 schema.Resource
// @Failure 401 option.Interface "{error:{code:0,message:未授权}}"
// @Failure 404 string "{error:{code:0,message:资源不存在}}"
// @Failure 500 option.Interface "{error:{code:0,message:服务器错误}}"
// @Router GET /api/v1/resources/{id}
func (a *Resource) Get(ctx *context.Context) {
	item, err := a.ResourceBll.Get(ctx.GetContext(), ctx.Param("id"))
	if err != nil {
		ctx.ResError(err)
		return
	}
	ctx.ResSuccess(item)
}

// Create 创建数据
// @Summary 创建数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param body body schema.Resource true
// @Success 200 option.Interface "{record_id:记录ID}"
// @Failure 400 option.Interface "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 option.Interface "{error:{code:0,message:未授权}}"
// @Failure 500 option.Interface "{error:{code:0,message:服务器错误}}"
// @Router POST /api/v1/resources
func (a *Resource) Create(ctx *context.Context) {
	var item schema.Resource
	if err := ctx.ParseJSON(&item); err != nil {
		ctx.ResError(err)
		return
	}

	newItem, err := a.ResourceBll.Create(ctx.GetContext(), item)
	if err != nil {
		ctx.ResError(err)
		return
	}

	ctx.ResSuccess(schema.HTTPNewItem{RecordID: newItem.RecordID})
}

// Update 更新数据
// @Summary 更新数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path string true "记录ID"
// @Param body body schema.Resource true
// @Success 200 option.Interface "{status:OK}"
// @Failure 400 option.Interface "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 option.Interface "{error:{code:0,message:未授权}}"
// @Failure 500 option.Interface "{error:{code:0,message:服务器错误}}"
// @Router PUT /api/v1/resources/{id}
func (a *Resource) Update(ctx *context.Context) {
	var item schema.Resource
	if err := ctx.ParseJSON(&item); err != nil {
		ctx.ResError(err)
		return
	}

	err := a.ResourceBll.Update(ctx.GetContext(), ctx.Param("id"), item)
	if err != nil {
		ctx.ResError(err)
		return
	}
	ctx.ResOK()
}

// Delete 删除数据
// @Summary 删除数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path string true "记录ID"
// @Success 200 option.Interface "{status:OK}"
// @Failure 401 option.Interface "{error:{code:0,message:未授权}}"
// @Failure 500 option.Interface "{error:{code:0,message:服务器错误}}"
// @Router DELETE /api/v1/resources/{id}
func (a *Resource) Delete(ctx *context.Context) {
	err := a.ResourceBll.Delete(ctx.GetContext(), ctx.Param("id"))
	if err != nil {
		ctx.ResError(err)
		return
	}
	ctx.ResOK()
}

// DeleteMany 删除多条数据
// @Summary 删除多条数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param batch query string true "记录ID（多个以,分隔）"
// @Success 200 option.Interface "{status:OK}"
// @Failure 400 option.Interface "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 option.Interface "{error:{code:0,message:未授权}}"
// @Failure 500 option.Interface "{error:{code:0,message:服务器错误}}"
// @Router DELETE /api/v1/resources
func (a *Resource) DeleteMany(ctx *context.Context) {
	ids := strings.Split(ctx.Query("batch"), ",")
	if len(ids) == 0 {
		ctx.ResError(errors.NewBadRequestError("无效的请求参数"))
		return
	}

	for _, id := range ids {
		err := a.ResourceBll.Delete(ctx.GetContext(), id)
		if err != nil {
			ctx.ResError(err)
			return
		}
	}

	ctx.ResOK()
}
