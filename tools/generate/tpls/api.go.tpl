package api

import (
	"{{.PkgName}}/internal/module/{{.ModuleLowerName}}/biz"
	"{{.PkgName}}/internal/module/{{.ModuleLowerName}}/typed"
	"{{.PkgName}}/internal/x/utilx"
	"github.com/gin-gonic/gin"
)

type {{.Name}}API struct {
	{{.Name}}Biz *biz.{{.Name}}Biz
}

// @Tags {{.Name}}API
// @Security ApiKeyAuth
// @Summary Query {{.LowerSpaceName}} list
// @Param current query int true "pagination index" default(1)
// @Param pageSize query int true "pagination size" default(10)
{{range .Fields}}{{if .InQuery}}// @Param {{.FirstLowerName}} query {{.Type}} false "{{.QueryComments}}"{{end}}{{end}}
// @Success 200 {object} utilx.ListResult{list=[]typed.{{.Name}}} "query result"
// @Failure 401 {object} utilx.ErrorResult
// @Failure 500 {object} utilx.ErrorResult
// @Router /api/{{.ModuleLowerName}}/v1/{{.LowerPluralName}} [get]
func (a *{{.Name}}API) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params typed.{{.Name}}QueryParam
	if err := utilx.ParseQuery(c, &params); err != nil {
		utilx.ResError(c, err)
		return
	}

	result, err := a.{{.Name}}Biz.Query(ctx, params)
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResList(c, result.Data)
}

// @Tags {{.Name}}API
// @Security ApiKeyAuth
// @Summary Get single {{.LowerSpaceName}} by id
// @Param id path string true "unique id"
// @Success 200 {object} typed.{{.Name}}
// @Failure 401 {object} utilx.ErrorResult
// @Failure 500 {object} utilx.ErrorResult
// @Router /api/{{.ModuleLowerName}}/v1/{{.LowerPluralName}}/{id} [get]
func (a *{{.Name}}API) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.{{.Name}}Biz.Get(ctx, c.Param("id"))
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResSuccess(c, item)
}

// @Tags {{.Name}}API
// @Security ApiKeyAuth
// @Summary Create {{.LowerSpaceName}}
// @Param body body typed.{{.Name}}Create true "request body"
// @Success 200 {object} typed.{{.Name}}
// @Failure 400 {object} utilx.ErrorResult
// @Failure 401 {object} utilx.ErrorResult
// @Failure 500 {object} utilx.ErrorResult
// @Router /api/{{.ModuleLowerName}}/v1/{{.LowerPluralName}} [post]
func (a *{{.Name}}API) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var item typed.{{.Name}}Create
	if err := utilx.ParseJSON(c, &item); err != nil {
		utilx.ResError(c, err)
		return
	}

	result, err := a.{{.Name}}Biz.Create(ctx, item)
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResSuccess(c, result)
}

// @Tags {{.Name}}API
// @Security ApiKeyAuth
// @Summary Update {{.LowerSpaceName}} by id
// @Param id path int true "unique id"
// @Param body body typed.{{.Name}}Create true "request body"
// @Success 200 {object} utilx.OkResult "ok=true"
// @Failure 400 {object} utilx.ErrorResult
// @Failure 401 {object} utilx.ErrorResult
// @Failure 500 {object} utilx.ErrorResult
// @Router /api/{{.ModuleLowerName}}/v1/{{.LowerPluralName}}/{id} [put]
func (a *{{.Name}}API) Update(c *gin.Context) {
	ctx := c.Request.Context()
	var item typed.{{.Name}}Create
	if err := utilx.ParseJSON(c, &item); err != nil {
		utilx.ResError(c, err)
		return
	}

	err := a.{{.Name}}Biz.Update(ctx, c.Param("id"), item)
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResOK(c)
}

// @Tags {{.Name}}API
// @Security ApiKeyAuth
// @Summary Delete single {{.LowerSpaceName}} by id
// @Param id path string true "unique id"
// @Success 200 {object} utilx.OkResult "ok=true"
// @Failure 401 {object} utilx.ErrorResult
// @Failure 500 {object} utilx.ErrorResult
// @Router /api/{{.ModuleLowerName}}/v1/{{.LowerPluralName}}/{id} [delete]
func (a *{{.Name}}API) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.{{.Name}}Biz.Delete(ctx, c.Param("id"))
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResOK(c)
}
