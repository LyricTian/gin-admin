package biz

import (
	"context"
	"io/ioutil"
	"os"
	"strings"

	"github.com/LyricTian/gin-admin/v9/internal/module/rbac/dao"
	"github.com/LyricTian/gin-admin/v9/internal/module/rbac/typed"
	"github.com/LyricTian/gin-admin/v9/internal/x/contextx"
	"github.com/LyricTian/gin-admin/v9/internal/x/utilx"
	"github.com/LyricTian/gin-admin/v9/pkg/errors"
	"github.com/LyricTian/gin-admin/v9/pkg/util/json"
	"github.com/LyricTian/gin-admin/v9/pkg/util/xid"
)

type MenuBiz struct {
	TransRepo              utilx.TransRepo
	MenuRepo               dao.MenuRepo
	MenuActionRepo         dao.MenuActionRepo
	MenuActionResourceRepo dao.MenuActionResourceRepo
}

func (a *MenuBiz) InitFromJSON(ctx context.Context, dataFile string) error {
	// If exists menu data in database, skip init menu data.
	menuResult, err := a.MenuRepo.Query(ctx, typed.MenuQueryParam{
		PaginationParam: utilx.PaginationParam{
			OnlyCount: true,
		},
	})
	if err != nil {
		return err
	} else if menuResult.PageResult.Total > 0 {
		return nil
	}

	// Load menu data from json file
	buf, err := ioutil.ReadFile(dataFile)
	if err != nil {
		if err == os.ErrNotExist {
			return nil
		}
		return err
	}

	var menuTrees typed.Menus
	err = json.Unmarshal(buf, &menuTrees)
	if err != nil {
		return err
	}

	return a.createMenus(ctx, "", menuTrees)
}

func (a *MenuBiz) createMenus(ctx context.Context, parentID string, createItems typed.Menus) error {
	return a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		for _, citem := range createItems {
			menu, err := a.Create(ctx, typed.MenuCreate{
				Name:     citem.Name,
				Sequence: citem.Sequence,
				Icon:     citem.GetIcon(),
				Link:     citem.GetLink(),
				ParentID: parentID,
				Remark:   citem.GetRemark(),
				Hide:     citem.Hide,
				Actions:  citem.Actions,
			})
			if err != nil {
				return err
			}

			if citem.Children != nil {
				err := a.createMenus(ctx, menu.ID, *citem.Children)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (a *MenuBiz) Query(ctx context.Context, params typed.MenuQueryParam) (*typed.MenuQueryResult, error) {
	params.PageSize = -1
	result, err := a.MenuRepo.Query(ctx, params, typed.MenuQueryOptions{
		QueryOptions: utilx.QueryOptions{
			OrderFields: []utilx.OrderByParam{
				{Field: "sequence", Direction: utilx.DESC},
				{Field: "id", Direction: utilx.DESC},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	result.Data = result.Data.ToTree()
	return result, nil
}

func (a *MenuBiz) Get(ctx context.Context, id string) (*typed.Menu, error) {
	menu, err := a.MenuRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	} else if menu == nil {
		return nil, errors.NotFound(errors.ErrNotFoundID, "Menu not found")
	}

	resourceResult, err := a.MenuActionResourceRepo.Query(ctx, typed.MenuActionResourceQueryParam{
		MenuID: id,
	})
	if err != nil {
		return nil, err
	}

	actionResult, err := a.MenuActionRepo.Query(ctx, typed.MenuActionQueryParam{
		MenuID: id,
	})
	if err != nil {
		return nil, err
	}

	menu.Actions = actionResult.Data.FillResources(resourceResult.Data.ToActionIDMap())
	return menu, nil
}

func (a *MenuBiz) getParent(ctx context.Context, parentID string) (*typed.Menu, error) {
	parent, err := a.MenuRepo.Get(ctx, parentID, typed.MenuQueryOptions{
		QueryOptions: utilx.QueryOptions{
			SelectFields: []string{"id", "parent_path"},
		},
	})
	if err != nil {
		return nil, err
	} else if parent == nil {
		return nil, errors.NotFound(errors.ErrNotFoundID, "Parent menu not found")
	}
	return parent, nil
}

func (a *MenuBiz) Create(ctx context.Context, createItem typed.MenuCreate) (*typed.Menu, error) {
	menu := &typed.Menu{
		ID:        xid.NewID(),
		Name:      createItem.Name,
		Sequence:  createItem.Sequence,
		Icon:      &createItem.Icon,
		Link:      &createItem.Link,
		Remark:    &createItem.Remark,
		Hide:      createItem.Hide,
		Status:    typed.MenuStatusEnabled,
		CreatedBy: contextx.FromUserID(ctx),
	}

	if createItem.ParentID != "" {
		parent, err := a.getParent(ctx, createItem.ParentID)
		if err != nil {
			return nil, err
		}
		menu.ParentID = &parent.ID
		menu.ParentPath = parent.GenerateParentPath()
	}

	err := a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		if err := a.MenuRepo.Create(ctx, menu); err != nil {
			return err
		}

		for _, action := range createItem.Actions {
			action.ID = xid.NewID()
			action.MenuID = menu.ID
			if err := a.MenuActionRepo.Create(ctx, action); err != nil {
				return err
			}

			for _, resource := range action.Resources {
				resource.ID = xid.NewID()
				resource.MenuID = menu.ID
				resource.ActionID = action.ID
				if err := a.MenuActionResourceRepo.Create(ctx, resource); err != nil {
					return err
				}
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return menu, nil
}

func (a *MenuBiz) Update(ctx context.Context, id string, createItem typed.MenuCreate) error {
	oldMenu, err := a.MenuRepo.Get(ctx, id)
	if err != nil {
		return err
	} else if oldMenu == nil {
		return errors.NotFound(errors.ErrNotFoundID, "Menu not found")
	}
	oldMenu.Name = createItem.Name
	oldMenu.Sequence = createItem.Sequence
	oldMenu.Icon = &createItem.Icon
	oldMenu.Link = &createItem.Link
	oldMenu.Remark = &createItem.Remark
	oldMenu.Hide = createItem.Hide
	oldMenu.UpdatedBy = contextx.FromUserID(ctx)

	oldParentPath := oldMenu.GenerateParentPath()
	if oldMenu.GetParentID() != createItem.ParentID {
		if createItem.ParentID != "" {
			parent, err := a.getParent(ctx, createItem.ParentID)
			if err != nil {
				return err
			}
			oldMenu.ParentID = &parent.ID
			oldMenu.ParentPath = parent.GenerateParentPath()
		} else {
			oldMenu.ParentID = new(string)
			oldMenu.ParentPath = new(string)
		}
	}

	return a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		if err := a.MenuRepo.Update(ctx, oldMenu); err != nil {
			return err
		}

		if oldMenu.GetParentID() != createItem.ParentID {
			err := a.updateChildrenParentPath(ctx, *oldParentPath, *oldMenu.GenerateParentPath())
			if err != nil {
				return err
			}
		}

		return a.updateActions(ctx, oldMenu.ID, createItem.Actions)
	})
}

func (a *MenuBiz) updateChildrenParentPath(ctx context.Context, oldParentPath, newParentPath string) error {
	menuResult, err := a.MenuRepo.Query(ctx, typed.MenuQueryParam{
		ParentPathPrefix: oldParentPath,
	}, typed.MenuQueryOptions{
		QueryOptions: utilx.QueryOptions{
			SelectFields: []string{"id", "parent_path"},
		},
	})
	if err != nil {
		return err
	} else if menuResult.Data.Len() == 0 {
		return nil
	}

	for _, menu := range menuResult.Data {
		parentPath := strings.Replace(*menu.ParentPath, oldParentPath, newParentPath, 1)
		err := a.MenuRepo.UpdateParentPath(ctx, menu.ID, parentPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *MenuBiz) updateActions(ctx context.Context, id string, actions typed.MenuActions) error {
	if len(actions) == 0 {
		if err := a.MenuActionRepo.DeleteByMenuID(ctx, id); err != nil {
			return err
		}
		return a.MenuActionResourceRepo.DeleteByMenuID(ctx, id)
	}

	actionResult, err := a.MenuActionRepo.Query(ctx, typed.MenuActionQueryParam{
		MenuID: id,
	})
	if err != nil {
		return err
	}
	mActionData := actionResult.Data.ToMap()

	resourceResult, err := a.MenuActionResourceRepo.Query(ctx, typed.MenuActionResourceQueryParam{
		MenuID: id,
	})
	if err != nil {
		return err
	}
	mResourceData := resourceResult.Data.ToMap()

	for _, action := range actions {
		if _, ok := mActionData[action.ID]; !ok {
			action.ID = xid.NewID()
		}

		for _, resource := range action.Resources {
			if oldResource, ok := mResourceData[resource.ID]; ok {
				delete(mResourceData, action.ID)

				if oldResource.Method == resource.Method && oldResource.Path == resource.Path {
					continue
				}
				oldResource.Method = resource.Method
				oldResource.Path = resource.Path
				if err := a.MenuActionResourceRepo.Update(ctx, oldResource); err != nil {
					return err
				}
				continue
			}

			resource.ID = xid.NewID()
			resource.ActionID = action.ID
			resource.MenuID = id
			if err := a.MenuActionResourceRepo.Create(ctx, resource); err != nil {
				return err
			}
		}

		if oldAction, ok := mActionData[action.ID]; ok {
			delete(mActionData, action.ID)

			if oldAction.Code == action.Code && oldAction.Name == action.Name {
				continue
			}
			oldAction.Code = action.Code
			oldAction.Name = action.Name
			if err := a.MenuActionRepo.Update(ctx, oldAction); err != nil {
				return err
			}
			continue
		}

		action.MenuID = id
		if err := a.MenuActionRepo.Create(ctx, action); err != nil {
			return err
		}
	}

	for _, action := range mActionData {
		if err := a.MenuActionRepo.Delete(ctx, action.ID); err != nil {
			return err
		}
	}

	for _, resource := range mResourceData {
		if err := a.MenuActionResourceRepo.Delete(ctx, resource.ID); err != nil {
			return err
		}
	}

	return nil
}

func (a *MenuBiz) deleteByID(ctx context.Context, id string) error {
	if err := a.MenuRepo.Delete(ctx, id); err != nil {
		return err
	}

	if err := a.MenuActionRepo.DeleteByMenuID(ctx, id); err != nil {
		return err
	}

	if err := a.MenuActionResourceRepo.DeleteByMenuID(ctx, id); err != nil {
		return err
	}

	return nil
}

func (a *MenuBiz) Delete(ctx context.Context, id string) error {
	oldMenu, err := a.MenuRepo.Get(ctx, id, typed.MenuQueryOptions{
		QueryOptions: utilx.QueryOptions{
			SelectFields: []string{"id", "parent_path"},
		},
	})
	if err != nil {
		return err
	} else if oldMenu == nil {
		return errors.NotFound(errors.ErrNotFoundID, "Menu not found")
	}

	childResult, err := a.MenuRepo.Query(ctx, typed.MenuQueryParam{
		ParentPathPrefix: *oldMenu.GenerateParentPath(),
	}, typed.MenuQueryOptions{
		QueryOptions: utilx.QueryOptions{
			SelectFields: []string{"id"},
		},
	})
	if err != nil {
		return err
	}

	return a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		for _, menu := range childResult.Data {
			if err := a.deleteByID(ctx, menu.ID); err != nil {
				return err
			}
		}
		return a.deleteByID(ctx, id)
	})
}

func (a *MenuBiz) UpdateStatus(ctx context.Context, id string, status typed.MenuStatus) error {
	oldMenu, err := a.MenuRepo.Get(ctx, id, typed.MenuQueryOptions{
		QueryOptions: utilx.QueryOptions{
			SelectFields: []string{"id", "parent_path"},
		},
	})
	if err != nil {
		return err
	} else if oldMenu == nil {
		return errors.NotFound(errors.ErrNotFoundID, "Menu not found")
	}

	childResult, err := a.MenuRepo.Query(ctx, typed.MenuQueryParam{
		ParentPathPrefix: *oldMenu.GenerateParentPath(),
	}, typed.MenuQueryOptions{
		QueryOptions: utilx.QueryOptions{
			SelectFields: []string{"id"},
		},
	})
	if err != nil {
		return err
	}

	return a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		for _, menu := range childResult.Data {
			if err := a.MenuRepo.UpdateStatus(ctx, menu.ID, status); err != nil {
				return err
			}
		}
		return a.MenuRepo.UpdateStatus(ctx, id, status)
	})
}
