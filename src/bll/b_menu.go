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
		PageParam: &schema.PaginationParam{PageSize: 1},
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
func (a *Menu) QueryTree(ctx context.Context, includeResource bool) ([]*schema.MenuTree, error) {
	types := []int{1, 2}
	if includeResource {
		types = append(types, 3)
	}
	result, err := a.MenuModel.Query(ctx, schema.MenuQueryParam{
		Types: types,
	})
	if err != nil {
		return nil, err
	}
	return result.Data.ToTrees().ToTree(), nil
}

// Get 查询指定数据
func (a *Menu) Get(ctx context.Context, recordID string) (*schema.Menu, error) {
	item, err := a.MenuModel.Get(ctx, recordID)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.ErrNotFound
	}
	return item, nil
}

func (a *Menu) checkAndGetParent(ctx context.Context, item schema.Menu, oldItem *schema.Menu) (*schema.Menu, error) {
	if oldItem == nil || oldItem.Code != item.Code {
		exists, err := a.MenuModel.CheckCodeWithParentID(ctx, item.Code, item.ParentID)
		if err != nil {
			return nil, err
		} else if exists {
			return nil, errors.NewBadRequestError("同一父级下编号不允许重复")
		}
	}

	var parentItem *schema.Menu
	if item.ParentID != "" {
		pitem, err := a.MenuModel.Get(ctx, item.ParentID)
		if err != nil {
			return nil, err
		} else if pitem == nil {
			return nil, errors.NewBadRequestError("无效的父级节点")
		}
		parentItem = pitem
	}

	switch {
	case item.Type == 3 && (parentItem == nil || !(parentItem.Type == 2 || parentItem.Type == 1)):
		return nil, errors.NewBadRequestError("资源类型只能依赖于模块或功能")
	case item.Type == 2 && (parentItem == nil || parentItem.Type != 1):
		return nil, errors.NewBadRequestError("功能类型只能依赖于模块")
	case item.Type == 1 && parentItem != nil && parentItem.Type != 1:
		return nil, errors.NewBadRequestError("模块类型只能依赖于模块")
	}

	return parentItem, nil
}

// 获取父级路径
func (a *Menu) getParentPath(parentItem *schema.Menu) string {
	var parentPath string
	if parentItem != nil {
		if v := parentItem.ParentPath; v != "" {
			parentPath = v + "/"
		}
		parentPath = parentPath + parentItem.RecordID
	}
	return parentPath
}

// Create 创建数据
func (a *Menu) Create(ctx context.Context, item schema.Menu) (*schema.Menu, error) {
	parentItem, err := a.checkAndGetParent(ctx, item, nil)
	if err != nil {
		return nil, err
	}

	item.ParentPath = a.getParentPath(parentItem)
	item.RecordID = util.MustUUID()
	item.Creator = a.CommonBll.GetUserID(ctx)
	err = a.CommonBll.ExecTrans(ctx, func(ctx context.Context) error {
		return a.MenuModel.Create(ctx, item)
	})
	if err != nil {
		return nil, err
	}

	return &item, nil
}

// Update 更新数据
func (a *Menu) Update(ctx context.Context, recordID string, item schema.Menu) error {
	if recordID == item.ParentID {
		return errors.NewBadRequestError("不允许使用节点自身作为父级节点")
	}

	oldItem, err := a.MenuModel.Get(ctx, recordID)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	parentItem, err := a.checkAndGetParent(ctx, item, oldItem)
	if err != nil {
		return err
	}
	item.ParentPath = oldItem.ParentPath

	return a.CommonBll.ExecTrans(ctx, func(ctx context.Context) error {
		// 如果父级更新，需要更新当前节点及节点下级的父级路径
		if item.ParentID != oldItem.ParentID {
			item.ParentPath = a.getParentPath(parentItem)

			opath := oldItem.ParentPath
			if opath != "" {
				opath += "/"
			}
			opath += oldItem.RecordID

			result, err := a.MenuModel.Query(ctx, schema.MenuQueryParam{
				ParentPath: opath,
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
}

// Delete 删除数据
func (a *Menu) Delete(ctx context.Context, recordIDs ...string) error {
	return a.CommonBll.ExecTrans(ctx, func(ctx context.Context) error {
		for _, recordID := range recordIDs {
			exists, err := a.MenuModel.CheckChild(ctx, recordID)
			if err != nil {
				return err
			} else if exists {
				return errors.NewBadRequestError("含有子级菜单，不能删除")
			}

			err = a.MenuModel.Delete(ctx, recordID)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
