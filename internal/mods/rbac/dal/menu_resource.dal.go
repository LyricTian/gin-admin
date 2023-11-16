package dal

import (
	"context"

	"github.com/LyricTian/gin-admin/v10/internal/mods/rbac/schema"
	"github.com/LyricTian/gin-admin/v10/pkg/errors"
	"github.com/LyricTian/gin-admin/v10/pkg/util"
	"gorm.io/gorm"
)

// Get menu resource storage instance
func GetMenuResourceDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDB(ctx, defDB).Model(new(schema.MenuResource))
}

// Menu resource management for RBAC
type MenuResource struct {
	DB *gorm.DB
}

// Query menu resources from the database based on the provided parameters and options.
func (a *MenuResource) Query(ctx context.Context, params schema.MenuResourceQueryParam, opts ...schema.MenuResourceQueryOptions) (*schema.MenuResourceQueryResult, error) {
	var opt schema.MenuResourceQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	db := GetMenuResourceDB(ctx, a.DB)
	if v := params.MenuID; len(v) > 0 {
		db = db.Where("menu_id = ?", v)
	}
	if v := params.MenuIDs; len(v) > 0 {
		db = db.Where("menu_id IN ?", v)
	}

	var list schema.MenuResources
	pageResult, err := util.WrapPageQuery(ctx, db, params.PaginationParam, opt.QueryOptions, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	queryResult := &schema.MenuResourceQueryResult{
		PageResult: pageResult,
		Data:       list,
	}
	return queryResult, nil
}

// Get the specified menu resource from the database.
func (a *MenuResource) Get(ctx context.Context, id string, opts ...schema.MenuResourceQueryOptions) (*schema.MenuResource, error) {
	var opt schema.MenuResourceQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	item := new(schema.MenuResource)
	ok, err := util.FindOne(ctx, GetMenuResourceDB(ctx, a.DB).Where("id=?", id), opt.QueryOptions, item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}
	return item, nil
}

// Exist checks if the specified menu resource exists in the database.
func (a *MenuResource) Exists(ctx context.Context, id string) (bool, error) {
	ok, err := util.Exists(ctx, GetMenuResourceDB(ctx, a.DB).Where("id=?", id))
	return ok, errors.WithStack(err)
}

// ExistsMethodPathByMenuID checks if the specified menu resource exists in the database.
func (a *MenuResource) ExistsMethodPathByMenuID(ctx context.Context, method, path, menuID string) (bool, error) {
	ok, err := util.Exists(ctx, GetMenuResourceDB(ctx, a.DB).Where("method=? AND path=? AND menu_id=?", method, path, menuID))
	return ok, errors.WithStack(err)
}

// Create a new menu resource.
func (a *MenuResource) Create(ctx context.Context, item *schema.MenuResource) error {
	result := GetMenuResourceDB(ctx, a.DB).Create(item)
	return errors.WithStack(result.Error)
}

// Update the specified menu resource in the database.
func (a *MenuResource) Update(ctx context.Context, item *schema.MenuResource) error {
	result := GetMenuResourceDB(ctx, a.DB).Where("id=?", item.ID).Select("*").Omit("created_at").Updates(item)
	return errors.WithStack(result.Error)
}

// Delete the specified menu resource from the database.
func (a *MenuResource) Delete(ctx context.Context, id string) error {
	result := GetMenuResourceDB(ctx, a.DB).Where("id=?", id).Delete(new(schema.MenuResource))
	return errors.WithStack(result.Error)
}

// Deletes the menu resource by menu id.
func (a *MenuResource) DeleteByMenuID(ctx context.Context, menuID string) error {
	result := GetMenuResourceDB(ctx, a.DB).Where("menu_id=?", menuID).Delete(new(schema.MenuResource))
	return errors.WithStack(result.Error)
}
