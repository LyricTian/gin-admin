package biz

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/LyricTian/gin-admin/v10/internal/mods/rbac/dal"
	"github.com/LyricTian/gin-admin/v10/internal/mods/rbac/schema"
	"github.com/LyricTian/gin-admin/v10/pkg/errors"
	"github.com/LyricTian/gin-admin/v10/pkg/util"
)

// Menu management for RBAC
type Menu struct {
	Trans           *util.Trans
	MenuDAL         *dal.Menu
	MenuResourceDAL *dal.MenuResource
	RoleMenuDAL     *dal.RoleMenu
}

// Query menus from the data access object based on the provided parameters and options.
func (a *Menu) Query(ctx context.Context, params schema.MenuQueryParam) (*schema.MenuQueryResult, error) {
	params.Pagination = false

	result, err := a.MenuDAL.Query(ctx, params, schema.MenuQueryOptions{
		QueryOptions: util.QueryOptions{
			OrderFields: schema.MenusOrderParams,
		},
	})
	if err != nil {
		return nil, err
	}

	if params.LikeName != "" {
		result.Data, err = a.appendChildren(ctx, result.Data)
		if err != nil {
			return nil, err
		}
	}

	result.Data = result.Data.ToTree()
	return result, nil
}

func (a *Menu) appendChildren(ctx context.Context, data schema.Menus) (schema.Menus, error) {
	if len(data) == 0 {
		return data, nil
	}

	existsInData := func(id string) bool {
		for _, item := range data {
			if item.ID == id {
				return true
			}
		}
		return false
	}

	for _, item := range data {
		childResult, err := a.MenuDAL.Query(ctx, schema.MenuQueryParam{
			ParentPathPrefix: item.ParentPath + item.ID + util.TreePathDelimiter,
		})
		if err != nil {
			return nil, err
		}
		for _, child := range childResult.Data {
			if existsInData(child.ID) {
				continue
			}
			data = append(data, child)
		}
	}

	if parentIDs := data.SplitParentIDs(); len(parentIDs) > 0 {
		parentResult, err := a.MenuDAL.Query(ctx, schema.MenuQueryParam{
			InIDs: parentIDs,
		})
		if err != nil {
			return nil, err
		}
		for _, p := range parentResult.Data {
			if existsInData(p.ID) {
				continue
			}
			data = append(data, p)
		}
	}
	sort.Sort(data)

	return data, nil
}

// Get the specified menu from the data access object.
func (a *Menu) Get(ctx context.Context, id string) (*schema.Menu, error) {
	menu, err := a.MenuDAL.Get(ctx, id)
	if err != nil {
		return nil, err
	} else if menu == nil {
		return nil, errors.NotFound("", "Menu not found")
	}

	menuResResult, err := a.MenuResourceDAL.Query(ctx, schema.MenuResourceQueryParam{
		MenuID: menu.ID,
	})
	if err != nil {
		return nil, err
	}
	menu.Resources = menuResResult.Data

	return menu, nil
}

// Create a new menu in the data access object.
func (a *Menu) Create(ctx context.Context, formItem *schema.MenuForm) (*schema.Menu, error) {
	menu := &schema.Menu{
		ID:        util.NewXID(),
		CreatedAt: time.Now(),
	}

	if parentID := formItem.ParentID; parentID != "" {
		parent, err := a.MenuDAL.Get(ctx, parentID)
		if err != nil {
			return nil, err
		} else if parent == nil {
			return nil, errors.NotFound("", "Parent not found")
		}
		menu.ParentPath = parent.ParentPath + parent.ID + util.TreePathDelimiter
	}

	if exists, err := a.MenuDAL.ExistsCodeByParentID(ctx, menu.Code, menu.ParentID); err != nil {
		return nil, err
	} else if exists {
		return nil, errors.BadRequest("", "Menu code already exists at the same level")
	}

	if err := formItem.FillTo(menu); err != nil {
		return nil, err
	}

	err := a.Trans.Exec(ctx, func(ctx context.Context) error {
		if err := a.MenuDAL.Create(ctx, menu); err != nil {
			return err
		}

		for _, res := range formItem.Resources {
			res.ID = util.NewXID()
			res.MenuID = menu.ID
			res.CreatedAt = time.Now()
			if err := a.MenuResourceDAL.Create(ctx, res); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return menu, nil
}

// Update the specified menu in the data access object.
func (a *Menu) Update(ctx context.Context, id string, formItem *schema.MenuForm) error {
	menu, err := a.MenuDAL.Get(ctx, id)
	if err != nil {
		return err
	} else if menu == nil {
		return errors.NotFound("", "Menu not found")
	}

	oldParentPath := menu.ParentPath
	oldStatus := menu.Status
	var childData schema.Menus
	if menu.ParentID != formItem.ParentID {
		if parentID := formItem.ParentID; parentID != "" {
			parent, err := a.MenuDAL.Get(ctx, parentID)
			if err != nil {
				return err
			} else if parent == nil {
				return errors.NotFound("", "Parent not found")
			}
			menu.ParentPath = parent.ParentPath + parent.ID + util.TreePathDelimiter
		} else {
			menu.ParentPath = ""
		}

		childResult, err := a.MenuDAL.Query(ctx, schema.MenuQueryParam{
			ParentPathPrefix: oldParentPath + menu.ID + util.TreePathDelimiter,
		}, schema.MenuQueryOptions{
			QueryOptions: util.QueryOptions{
				SelectFields: []string{"id", "parent_path"},
			},
		})
		if err != nil {
			return err
		}
		childData = childResult.Data
	}

	if menu.Code != formItem.Code {
		if exists, err := a.MenuDAL.ExistsCodeByParentID(ctx, formItem.Code, formItem.ParentID); err != nil {
			return err
		} else if exists {
			return errors.BadRequest("", "Menu code already exists at the same level")
		}
	}

	if err := formItem.FillTo(menu); err != nil {
		return err
	}

	return a.Trans.Exec(ctx, func(ctx context.Context) error {
		if oldStatus != formItem.Status {
			opath := oldParentPath + menu.ID + util.TreePathDelimiter
			if err := a.MenuDAL.UpdateStatusByParentPath(ctx, opath, formItem.Status); err != nil {
				return err
			}
		}

		for _, child := range childData {
			opath := oldParentPath + menu.ID + util.TreePathDelimiter
			npath := menu.ParentPath + menu.ID + util.TreePathDelimiter
			err := a.MenuDAL.UpdateParentPath(ctx, child.ID, strings.Replace(child.ParentPath, opath, npath, 1))
			if err != nil {
				return err
			}
		}

		if err := a.MenuDAL.Update(ctx, menu); err != nil {
			return err
		}

		if err := a.MenuResourceDAL.DeleteByMenuID(ctx, id); err != nil {
			return err
		}
		for _, res := range formItem.Resources {
			if res.ID == "" {
				res.ID = util.NewXID()
			}
			res.MenuID = id
			if res.CreatedAt.IsZero() {
				res.CreatedAt = time.Now()
			}
			res.UpdatedAt = time.Now()
			if err := a.MenuResourceDAL.Create(ctx, res); err != nil {
				return err
			}
		}

		return nil
	})
}

// Delete the specified menu from the data access object.
func (a *Menu) Delete(ctx context.Context, id string) error {
	menu, err := a.MenuDAL.Get(ctx, id)
	if err != nil {
		return err
	} else if menu == nil {
		return errors.NotFound("", "Menu not found")
	}

	childResult, err := a.MenuDAL.Query(ctx, schema.MenuQueryParam{
		ParentPathPrefix: menu.ParentPath + menu.ID + util.TreePathDelimiter,
	}, schema.MenuQueryOptions{
		QueryOptions: util.QueryOptions{
			SelectFields: []string{"id"},
		},
	})
	if err != nil {
		return err
	}

	return a.Trans.Exec(ctx, func(ctx context.Context) error {
		if err := a.delete(ctx, id); err != nil {
			return err
		}

		for _, child := range childResult.Data {
			if err := a.delete(ctx, child.ID); err != nil {
				return err
			}
		}
		return nil
	})
}

func (a *Menu) delete(ctx context.Context, id string) error {
	if err := a.MenuDAL.Delete(ctx, id); err != nil {
		return err
	}
	if err := a.MenuResourceDAL.DeleteByMenuID(ctx, id); err != nil {
		return err
	}
	if err := a.RoleMenuDAL.DeleteByMenuID(ctx, id); err != nil {
		return err
	}
	return nil
}
