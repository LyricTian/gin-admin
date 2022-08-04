package dao

import (
	"context"

	"{{.PkgName}}/internal/module/{{.ModuleLowerName}}/typed"
	"{{.PkgName}}/internal/x/utilx"
	"{{.PkgName}}/pkg/errors"
	"gorm.io/gorm"
)

func Get{{.Name}}DB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return utilx.GetDB(ctx, defDB).Model(new(typed.{{.Name}}))
}

type {{.Name}}Repo struct {
	DB *gorm.DB
}

func (a *{{.Name}}Repo) Query(ctx context.Context, params typed.{{.Name}}QueryParam, opts ...typed.{{.Name}}QueryOptions) (*typed.{{.Name}}QueryResult, error) {
	var opt typed.{{.Name}}QueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	db := Get{{.Name}}DB(ctx, a.DB)
    {{range .Fields}}
        {{if .QueryIfExpression}}
        if {{.QueryIfExpression}} {
            db = db.Where({{.QueryWhereCondition}})
        }
        {{end}}
	{{end}}

	var list typed.{{.PluralName}}
	pr, err := utilx.WrapPageQuery(ctx, db, params.PaginationParam, opt.QueryOptions, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	qr := &typed.{{.Name}}QueryResult{
		PageResult: pr,
		Data:       list,
	}

	return qr, nil
}

func (a *{{.Name}}Repo) Get(ctx context.Context, id string, opts ...typed.{{.Name}}QueryOptions) (*typed.{{.Name}}, error) {
	var opt typed.{{.Name}}QueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	item := new(typed.{{.Name}})
	ok, err := utilx.FindOne(ctx, Get{{.Name}}DB(ctx, a.DB).Where("id=?", id), opt.QueryOptions, item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item, nil
}

func (a *{{.Name}}Repo) Exists(ctx context.Context, id string) (bool, error) {
	exists, err := utilx.Exists(ctx, Get{{.Name}}DB(ctx, a.DB).Where("id=?", id))
	return exists, errors.WithStack(err)
}

func (a *{{.Name}}Repo) Create(ctx context.Context, item *typed.{{.Name}}) error {
	result := Get{{.Name}}DB(ctx, a.DB).Create(item)
	return errors.WithStack(result.Error)
}

func (a *{{.Name}}Repo) Update(ctx context.Context, item *typed.{{.Name}}) error {
	result := Get{{.Name}}DB(ctx, a.DB).Where("id=?", item.ID).Omit("created_at", "created_by").Updates(item)
	return errors.WithStack(result.Error)
}

func (a *{{.Name}}Repo) Delete(ctx context.Context, id string) error {
	result := Get{{.Name}}DB(ctx, a.DB).Where("id=?", id).Delete(new(typed.{{.Name}}))
	return errors.WithStack(result.Error)
}
