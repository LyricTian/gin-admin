package ctl

import (
	"strings"

	"github.com/LyricTian/gin-admin/src/bll"
	"github.com/LyricTian/gin-admin/src/errors"
	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/LyricTian/gin-admin/src/util"
	"github.com/LyricTian/gin-admin/src/web/context"
)

// Demo demo
// @Name Demo
// @Description demo接口
type Demo struct {
	DemoBll *bll.Demo `inject:""`
}

// Query 查询数据
func (a *Demo) Query(ctx *context.Context) {
	switch ctx.Query("q") {
	case "page":
		a.QueryPage(ctx)
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
// @Param status query int false "状态(1:启用 2:停用)"
// @Success 200 []schema.Demo "查询结果：{list:列表数据,pagination:{current:页索引,pageSize:页大小,total:总数量}}"
// @Failure 400 option.Interface "{error:{code:0,message:未知的查询类型}}"
// @Failure 401 option.Interface "{error:{code:0,message:未授权}}"
// @Failure 500 option.Interface "{error:{code:0,message:服务器错误}}"
// @Router GET /api/v1/demos?q=page
func (a *Demo) QueryPage(ctx *context.Context) {
	var params schema.DemoQueryParam
	params.Code = ctx.Query("code")
	params.Name = ctx.Query("name")
	params.Status = util.S(ctx.Query("status")).Int()

	items, pr, err := a.DemoBll.Query(ctx.CContext(), params, ctx.GetPaginationParam())
	if err != nil {
		ctx.ResError(err)
		return
	}

	ctx.ResPage(items, pr)
}

// Get 查询指定数据
// @Summary 查询指定数据
// @Param Access-Token header string false "访问令牌"
// @Param id path string true "记录ID"
// @Success 200 schema.Demo
// @Failure 401 option.Interface "{error:{code:0,message:未授权}}"
// @Failure 404 string "{error:{code:0,message:资源不存在}}"
// @Failure 500 option.Interface "{error:{code:0,message:服务器错误}}"
// @Router GET /api/v1/demos/{id}
func (a *Demo) Get(ctx *context.Context) {
	item, err := a.DemoBll.Get(ctx.CContext(), ctx.Param("id"))
	if err != nil {
		ctx.ResError(err)
		return
	}
	ctx.ResSuccess(item)
}

// Create 创建数据
// @Summary 创建数据
// @Param Access-Token header string false "访问令牌"
// @Param body body schema.Demo true
// @Success 200 option.Interface "{record_id:记录ID}"
// @Failure 400 option.Interface "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 option.Interface "{error:{code:0,message:未授权}}"
// @Failure 500 option.Interface "{error:{code:0,message:服务器错误}}"
// @Router POST /api/v1/demos
func (a *Demo) Create(ctx *context.Context) {
	var item schema.Demo
	if err := ctx.ParseJSON(&item); err != nil {
		ctx.ResError(err)
		return
	}

	newItem, err := a.DemoBll.Create(ctx.CContext(), item)
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
// @Param body body schema.Demo true
// @Success 200 option.Interface "{status:OK}"
// @Failure 400 option.Interface "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 option.Interface "{error:{code:0,message:未授权}}"
// @Failure 500 option.Interface "{error:{code:0,message:服务器错误}}"
// @Router PUT /api/v1/demos/{id}
func (a *Demo) Update(ctx *context.Context) {
	var item schema.Demo
	if err := ctx.ParseJSON(&item); err != nil {
		ctx.ResError(err)
		return
	}

	err := a.DemoBll.Update(ctx.CContext(), ctx.Param("id"), item)
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
// @Router DELETE /api/v1/demos/{id}
func (a *Demo) Delete(ctx *context.Context) {
	err := a.DemoBll.Delete(ctx.CContext(), ctx.Param("id"))
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
// @Router DELETE /api/v1/demos
func (a *Demo) DeleteMany(ctx *context.Context) {
	ids := strings.Split(ctx.Query("batch"), ",")
	if len(ids) == 0 {
		ctx.ResError(errors.NewBadRequestError("无效的请求参数"))
		return
	}

	for _, id := range ids {
		err := a.DemoBll.Delete(ctx.CContext(), id)
		if err != nil {
			ctx.ResError(err)
			return
		}
	}

	ctx.ResOK()
}

// Enable 启用数据
// @Summary 启用数据
// @Param Access-Token header string false "访问令牌"
// @Param id path string true "记录ID"
// @Success 200 option.Interface "{status:OK}"
// @Failure 401 option.Interface "{error:{code:0,message:未授权}}"
// @Failure 500 option.Interface "{error:{code:0,message:服务器错误}}"
// @Router PATCH /api/v1/demos/{id}/enable
func (a *Demo) Enable(ctx *context.Context) {
	err := a.DemoBll.UpdateStatus(ctx.CContext(), ctx.Param("id"), 1)
	if err != nil {
		ctx.ResError(err)
		return
	}
	ctx.ResOK()
}

// Disable 禁用数据
// @Summary 禁用数据
// @Param Access-Token header string false "访问令牌"
// @Param id path string true "记录ID"
// @Success 200 option.Interface "{status:OK}"
// @Failure 401 option.Interface "{error:{code:0,message:未授权}}"
// @Failure 500 option.Interface "{error:{code:0,message:服务器错误}}"
// @Router PATCH /api/v1/demos/{id}/disable
func (a *Demo) Disable(ctx *context.Context) {
	err := a.DemoBll.UpdateStatus(ctx.CContext(), ctx.Param("id"), 2)
	if err != nil {
		ctx.ResError(err)
		return
	}
	ctx.ResOK()
}
