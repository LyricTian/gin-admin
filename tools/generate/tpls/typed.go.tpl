package typed

import (
	"time"

	"{{.PkgName}}/internal/x/utilx"
)

{{if .Description}}// {{.Description}}{{end}}
type {{.Name}} struct {
	ID        string     `gorm:"size:20;primarykey;" json:"id"`
	{{range .Fields}}{{.Name}} {{if .Optional}}*{{end}}{{.Type}} `{{if .GormTag}}gorm:"{{.GormTag}}" {{end}}json:"{{.LowerUnderlineName}}"`{{if .Comments}} // {{.Comments}}{{end}}{{end}}
	CreatedAt time.Time  `gorm:"index;" json:"created_at"`
	CreatedBy string     `gorm:"size:20;" json:"created_by"`
	UpdatedAt time.Time  `gorm:"index;" json:"updated_at"`
	UpdatedBy string     `gorm:"size:20;" json:"updated_by"`
}

type {{.Name}}QueryParam struct {
	utilx.PaginationParam
	{{range .Fields}}{{if .QueryIfExpression}}{{.Name}} {{.Type}} `form:"{{if .InQuery}}{{.FirstLowerName}}{{else}}-{{end}}"`{{end}}{{end}}
}

type {{.Name}}QueryOptions struct {
	utilx.QueryOptions
}

type {{.Name}}QueryResult struct {
	Data       {{.PluralName}}
	PageResult *utilx.PaginationResult
}

type {{.PluralName}} []*{{.Name}}

type {{.Name}}Create struct {
    {{range .Fields}}{{if .InCreate}}{{.Name}} {{.Type}} `json:"{{.LowerUnderlineName}}"{{if .BindingTag}} binding:"{{.BindingTag}}"{{end}}`{{end}}{{end}}
}
