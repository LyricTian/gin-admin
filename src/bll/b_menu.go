package bll

import (
	"context"

	"github.com/LyricTian/gin-admin/src/errors"
	"github.com/LyricTian/gin-admin/src/model"
	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/LyricTian/gin-admin/src/util"
)

// Menu 菜单管理
type Menu struct {
	MenuModel model.IMenu `inject:"IMenu"`
	CommonBll *Common     `inject:""`
}

// CheckDataInit 检查数据是否初始化
func (a *Menu) CheckDataInit(ctx context.Context) (bool, error) {
	result, err := a.MenuModel.Query(ctx, schema.MenuQueryParam{}, schema.MenuQueryOptions{
		PageParam: &schema.PaginationParam{PageSize: -1},
	})
	if err != nil {
		return false, err
	}
	return result.PageResult.Total > 0, nil
}

// QueryPage 查询分页数据
func (a *Menu) QueryPage(ctx context.Context, params schema.MenuQueryParam, pp *schema.PaginationParam) ([]*schema.Menu, *schema.PaginationResult, error) {
	result, err := a.MenuModel.Query(ctx, params, schema.MenuQueryOptions{
		PageParam: pp,
	})
	if err != nil {
		return nil, nil, err
	}
	return result.Data, result.PageResult, nil
}

// QueryTree 查询菜单树
func (a *Menu) QueryTree(ctx context.Context, includeActions, includeResources bool) ([]*schema.MenuTree, error) {
	result, err := a.MenuModel.Query(ctx, schema.MenuQueryParam{}, schema.MenuQueryOptions{
		IncludeActions:   includeActions,
		IncludeResources: includeResources,
	})
	if err != nil {
		return nil, err
	}
	return result.Data.ToTrees().ToTree(), nil
}

// Get 查询指定数据
func (a *Menu) Get(ctx context.Context, recordID string) (*schema.Menu, error) {
	item, err := a.MenuModel.Get(ctx, recordID,
		schema.MenuQueryOptions{
			IncludeResources: true,
			IncludeActions:   true,
		},
	)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.ErrNotFound
	}
	return item, nil
}

// 获取父级路径
func (a *Menu) getParentPath(ctx context.Context, parentID string) (string, error) {
	if parentID == "" {
		return "", nil
	}

	pitem, err := a.MenuModel.Get(ctx, parentID)
	if err != nil {
		return "", err
	} else if pitem == nil {
		return "", errors.NewBadRequestError("无效的父级节点")
	}

	var parentPath string
	if v := pitem.ParentPath; v != "" {
		parentPath = v + "/"
	}
	parentPath = parentPath + pitem.RecordID
	return parentPath, nil
}

// Create 创建数据
func (a *Menu) Create(ctx context.Context, item schema.Menu) (*schema.Menu, error) {
	parentPath, err := a.getParentPath(ctx, item.ParentID)
	if err != nil {
		return nil, err
	}

	item.ParentPath = parentPath
	item.RecordID = util.MustUUID()
	item.Creator = a.CommonBll.GetUserID(ctx)
	err = a.MenuModel.Create(ctx, item)
	if err != nil {
		return nil, err
	}

	return a.Get(ctx, item.RecordID)
}

// Update 更新数据
func (a *Menu) Update(ctx context.Context, recordID string, item schema.Menu) (*schema.Menu, error) {
	if recordID == item.ParentID {
		return nil, errors.NewBadRequestError("不允许使用节点自身作为父级节点")
	}

	oldItem, err := a.MenuModel.Get(ctx, recordID)
	if err != nil {
		return nil, err
	} else if oldItem == nil {
		return nil, errors.ErrNotFound
	}
	item.ParentPath = oldItem.ParentPath

	err = a.CommonBll.ExecTrans(ctx, func(ctx context.Context) error {
		// 如果父级更新，需要更新当前节点及节点下级的父级路径
		if item.ParentID != oldItem.ParentID {
			parentPath, err := a.getParentPath(ctx, item.ParentID)
			if err != nil {
				return err
			}
			item.ParentPath = parentPath

			opath := oldItem.ParentPath
			if opath != "" {
				opath += "/"
			}
			opath += oldItem.RecordID

			result, err := a.MenuModel.Query(ctx, schema.MenuQueryParam{
				PrefixParentPath: opath,
			})
			if err != nil {
				return err
			}

			for _, menu := range result.Data {
				npath := item.ParentPath + menu.ParentPath[len(opath):]
				err = a.MenuModel.UpdateParentPath(ctx, menu.RecordID, npath)
				if err != nil {
					return err
				}
			}
		}

		return a.MenuModel.Update(ctx, recordID, item)
	})
	if err != nil {
		return nil, err
	}
	return a.Get(ctx, recordID)
}

// Delete 删除数据
func (a *Menu) Delete(ctx context.Context, recordID string) error {
	result, err := a.MenuModel.Query(ctx, schema.MenuQueryParam{
		ParentID: &recordID,
	}, schema.MenuQueryOptions{PageParam: &schema.PaginationParam{PageSize: -1}})
	if err != nil {
		return err
	} else if result.PageResult.Total > 0 {
		return errors.NewBadRequestError("含有子级菜单，不能删除")
	}

	return a.MenuModel.Delete(ctx, recordID)
}
