package biz

import (
	"context"
	"time"

	"github.com/LyricTian/gin-admin/v10/internal/library/utils"
	"github.com/LyricTian/gin-admin/v10/internal/mods/rbac/dal"
	"github.com/LyricTian/gin-admin/v10/internal/mods/rbac/schema"
	"github.com/LyricTian/gin-admin/v10/pkg/errors"
	"github.com/LyricTian/gin-admin/v10/pkg/util/xid"
)

type Resource struct {
	Trans       *utils.Trans
	ResourceDAL *dal.Resource
}

func (a *Resource) Query(ctx context.Context, params schema.ResourceQueryParam) (*schema.ResourceQueryResult, error) {
	params.Pagination = true

	result, err := a.ResourceDAL.Query(ctx, params, schema.ResourceQueryOptions{
		QueryOptions: utils.QueryOptions{
			OrderFields: []utils.OrderByParam{
				{Field: "created_at", Direction: utils.DESC},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (a *Resource) Get(ctx context.Context, id string) (*schema.Resource, error) {
	resource, err := a.ResourceDAL.Get(ctx, id)
	if err != nil {
		return nil, err
	} else if resource == nil {
		return nil, errors.NotFound("", "Resource not found")
	}

	return resource, nil
}

func (a *Resource) Create(ctx context.Context, citem schema.ResourceSave) (*schema.Resource, error) {
	if exists, err := a.ResourceDAL.ExistsCode(ctx, citem.Code); err != nil {
		return nil, err
	} else if exists {
		return nil, errors.BadRequest("", "Resource code already exists")
	}

	resource := &schema.Resource{
		ID:          xid.NewID(),
		Code:        citem.Code,
		Object:      citem.Object,
		Action:      citem.Action,
		Description: citem.Description,
		Status:      citem.Status,
		CreatedAt:   time.Now(),
	}

	err := a.Trans.Exec(ctx, func(ctx context.Context) error {
		if err := a.ResourceDAL.Create(ctx, resource); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return resource, nil
}

func (a *Resource) Update(ctx context.Context, id string, uitem schema.ResourceSave) error {
	oldResource, err := a.ResourceDAL.Get(ctx, id)
	if err != nil {
		return err
	} else if oldResource == nil {
		return errors.NotFound("", "Resource not found")
	} else if oldResource.Code != uitem.Code {
		if exists, err := a.ResourceDAL.ExistsCode(ctx, uitem.Code); err != nil {
			return err
		} else if exists {
			return errors.BadRequest("", "Resource code already exists")
		}
	}

	oldResource.Code = uitem.Code
	oldResource.Object = uitem.Object
	oldResource.Action = uitem.Action
	oldResource.Description = uitem.Description
	oldResource.Status = uitem.Status
	return a.Trans.Exec(ctx, func(ctx context.Context) error {
		if err := a.ResourceDAL.Update(ctx, oldResource); err != nil {
			return err
		}
		return nil
	})
}

func (a *Resource) Delete(ctx context.Context, id string) error {
	exists, err := a.ResourceDAL.Exists(ctx, id)
	if err != nil {
		return err
	} else if !exists {
		return errors.NotFound("", "Resource not found")
	}

	return a.Trans.Exec(ctx, func(ctx context.Context) error {
		if err := a.ResourceDAL.Delete(ctx, id); err != nil {
			return err
		}
		return nil
	})
}
