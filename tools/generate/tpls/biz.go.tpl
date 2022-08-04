package biz

import (
	"context"

	"{{.PkgName}}/internal/module/{{.ModuleLowerName}}/dao"
	"{{.PkgName}}/internal/module/{{.ModuleLowerName}}/typed"
	"{{.PkgName}}/internal/x/contextx"
	"{{.PkgName}}/internal/x/utilx"
	"{{.PkgName}}/pkg/errors"
	"{{.PkgName}}/pkg/logger"
	"{{.PkgName}}/pkg/util/xid"
	"{{.PkgName}}/pkg/x/cachex"
	"go.uber.org/zap"
)

type {{.Name}}Biz struct {
	TransRepo    utilx.TransRepo
	{{.Name}}Repo     dao.{{.Name}}Repo
}

func (a *{{.Name}}Biz) Query(ctx context.Context, params typed.{{.Name}}QueryParam) (*typed.{{.Name}}QueryResult, error) {
	params.Pagination = true
	queryOpts := utilx.QueryOptions{
		OrderFields: []utilx.OrderByParam{
			{Field: "created_at", Direction: utilx.DESC},
		},
	}

	result, err := a.{{.Name}}Repo.Query(ctx, params, typed.{{.Name}}QueryOptions{
		QueryOptions: queryOpts,
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *{{.Name}}Biz) Get(ctx context.Context, id string) (*typed.{{.Name}}, error) {
	{{.LowerName}}, err := a.{{.Name}}Repo.Get(ctx, id)
	if err != nil {
		return nil, err
	} else if {{.LowerName}} == nil {
		return nil, errors.NotFound(errors.ErrNotFoundID, "{{.Name}} not found")
	}

	return {{.LowerName}}, nil
}

func (a *{{.Name}}Biz) Create(ctx context.Context, createItem typed.{{.Name}}Create) (*typed.{{.Name}}, error) {
	{{.LowerName}} := &typed.{{.Name}}{
		ID:        xid.NewID(),
        {{range .Fields}}{{if .InCreate}}{{.Name}}: {{if .Optional}}&{{end}}createItem.{{.Name}},{{end}}
		{{end}}
		CreatedBy: contextx.FromUserID(ctx),
	}

	err := a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		if err := a.{{.Name}}Repo.Create(ctx, {{.LowerName}}); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return {{.LowerName}}, nil
}

func (a *{{.Name}}Biz) Update(ctx context.Context, id string, createItem typed.{{.Name}}Create) error {
	old{{.Name}}, err := a.{{.Name}}Repo.Get(ctx, id)
	if err != nil {
		return err
	} else if old{{.Name}} == nil {
		return errors.NotFound(errors.ErrNotFoundID, "{{.Name}} not found")
	}
    {{range .Fields}}{{if .InCreate}}old{{.Name}} = {{if .Optional}}&{{end}}createItem.{{.Name}}{{end}}
	{{end}}
	old{{.Name}}.UpdatedBy = contextx.FromUserID(ctx)

	return a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		if err := a.{{.Name}}Repo.Update(ctx, old{{.Name}}); err != nil {
			return err
		}

		return nil
	})
}

func (a *{{.Name}}Biz) Delete(ctx context.Context, id string) error {
	exists, err := a.{{.Name}}Repo.Exists(ctx, id)
	if err != nil {
		return err
	} else if !exists {
		return errors.NotFound(errors.ErrNotFoundID, "{{.Name}} not found")
	}

	err = a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		if err := a.{{.Name}}Repo.Delete(ctx, id); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
