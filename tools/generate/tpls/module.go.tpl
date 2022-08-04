package {{.ModuleLowerName}}

import (
	"context"

	"{{.PkgName}}/internal/module/{{.ModuleLowerName}}/api"
	"{{.PkgName}}/internal/module/{{.ModuleLowerName}}/biz"
	"{{.PkgName}}/internal/module/{{.ModuleLowerName}}/dao"
	"{{.PkgName}}/internal/module/{{.ModuleLowerName}}/typed"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"gorm.io/gorm"
)

// Collection of {{.ModuleName}} wire providers
var Set = wire.NewSet(
	wire.Struct(new({{.ModuleName}}), "*"),
    wire.Struct(new(dao.{{.Name}}Repo), "*"),
	wire.Struct(new(biz.{{.Name}}Biz), "*"),
	wire.Struct(new(api.{{.Name}}API), "*"),
) // end

// {{.ModuleName}} module is a {{.ModuleName}} service
type {{.ModuleName}} struct {
	DB       *gorm.DB
    {{.Name}}API  *api.{{.Name}}API
} // end

func (a *{{.ModuleName}}) Init(ctx context.Context) error {
	// Auto migrate tables for {{.ModuleName}}
	if err := a.autoMigrate(ctx); err != nil {
		return err
	}

	return nil
}

func (a *{{.ModuleName}}) autoMigrate(ctx context.Context) error {
	return a.DB.AutoMigrate(
		&typed.{{.Name}}{},
	) // end
}

func (a *{{.ModuleName}}) RegisterAPI(ctx context.Context, group *gin.RouterGroup) {
	r := group.Group("{{.ModuleLowerName}}")
	v1 := r.Group("v1")
	{
        g{{.Name}} := v1.Group("{{.LowerPluralName}}")
		{
			g{{.Name}}.GET("", a.{{.Name}}API.Query)
			g{{.Name}}.GET(":id", a.{{.Name}}API.Get)
			g{{.Name}}.POST("", a.{{.Name}}API.Create)
			g{{.Name}}.PUT(":id", a.{{.Name}}API.Update)
			g{{.Name}}.DELETE(":id", a.{{.Name}}API.Delete)
		}
	} // end
}
