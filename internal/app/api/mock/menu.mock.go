package mock

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var MenuSet = wire.NewSet(wire.Struct(new(MenuMock), "*"))

type MenuMock struct{}

// @Tags MenuAPI
// @Summary 查询数据
// @Security ApiKeyAuth
// @Param current query int true "分页索引" default(1)
// @Param pageSize query int true "分页大小" default(10)
// @Param queryValue query string false "查询值"
// @Param status query int false "状态(1:启用 2:禁用)"
// @Param isShow query int false "是否显示(1:显示 2:隐藏)"
// @Param parentID query int false "父级ID"
// @Success 200 {object} schema.ListResult{list=[]schema.Menu} "查询结果"
// @Failure 400 {object} schema.ErrorResult "{error:{code:0,message:bad request}}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:9999,message:invalid signature}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:internal server error}}"
// @Router /api/v1/menus [get]
func (a *MenuMock) Query(c *gin.Context) {
}

// @Tags MenuAPI
// @Summary 查询菜单树
// @Security ApiKeyAuth
// @Param status query int false "状态(1:启用 2:禁用)"
// @Param parentID query int false "父级ID"
// @Success 200 {object} schema.ListResult{list=[]schema.MenuTree} "查询结果"
// @Failure 401 {object} schema.ErrorResult "{error:{code:9999,message:invalid signature}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:internal server error}}"
// @Router /api/v1/menus.tree [get]
func (a *MenuMock) QueryTree(c *gin.Context) {
}

// @Tags MenuAPI
// @Summary 查询指定数据
// @Security ApiKeyAuth
// @Param id path int true "唯一标识"
// @Success 200 {object} schema.Menu
// @Failure 401 {object} schema.ErrorResult "{error:{code:9999,message:invalid signature}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:internal server error}}"
// @Router /api/v1/menus/{id} [get]
func (a *MenuMock) Get(c *gin.Context) {
}

// @Tags MenuAPI
// @Summary 创建数据
// @Security ApiKeyAuth
// @Param body body schema.Menu true "创建数据"
// @Success 200 {object} schema.IDResult
// @Failure 400 {object} schema.ErrorResult "{error:{code:0,message:bad request}}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:9999,message:invalid signature}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:internal server error}}"
// @Router /api/v1/menus [post]
func (a *MenuMock) Create(c *gin.Context) {
}

// @Tags MenuAPI
// @Summary 更新数据
// @Security ApiKeyAuth
// @Param id path int true "唯一标识"
// @Param body body schema.Menu true "更新数据"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 400 {object} schema.ErrorResult "{error:{code:0,message:bad request}}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:9999,message:invalid signature}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:internal server error}}"
// @Router /api/v1/menus/{id} [put]
func (a *MenuMock) Update(c *gin.Context) {
}

// @Tags MenuAPI
// @Summary 删除数据
// @Security ApiKeyAuth
// @Param id path int true "唯一标识"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:9999,message:invalid signature}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:internal server error}}"
// @Router /api/v1/menus/{id} [delete]
func (a *MenuMock) Delete(c *gin.Context) {
}

// @Tags MenuAPI
// @Summary 启用数据
// @Security ApiKeyAuth
// @Param id path int true "唯一标识"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:9999,message:invalid signature}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:internal server error}}"
// @Router /api/v1/menus/{id}/enable [patch]
func (a *MenuMock) Enable(c *gin.Context) {
}

// @Tags MenuAPI
// @Summary 禁用数据
// @Security ApiKeyAuth
// @Param id path int true "唯一标识"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:9999,message:invalid signature}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:internal server error}}"
// @Router /api/v1/menus/{id}/disable [patch]
func (a *MenuMock) Disable(c *gin.Context) {
}
