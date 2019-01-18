package bll

import (
	"context"
	"sync"

	"github.com/LyricTian/gin-admin/src/errors"
	"github.com/LyricTian/gin-admin/src/model"
	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/LyricTian/gin-admin/src/util"
)

// Menu 菜单管理
type Menu struct {
	lock      sync.RWMutex
	CommonBll *Common     `inject:""`
	MenuModel model.IMenu `inject:"IMenu"`
}

// QueryPage 查询分页数据
func (a *Menu) QueryPage(ctx context.Context, params schema.MenuPageQueryParam, pageIndex, pageSize uint) (int, []*schema.Menu, error) {
	return a.MenuModel.QueryPage(ctx, params, pageIndex, pageSize)
}

// QueryTree 查询菜单树
func (a *Menu) QueryTree(ctx context.Context) ([]*schema.MenuTreeResult, error) {
	items, err := a.MenuModel.QueryList(ctx, schema.MenuListQueryParam{
		Types: []int{10, 20},
	})
	if err != nil {
		return nil, err
	}

	return schema.MenuList(items).ToTreeResult(), nil
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

func (a *Menu) getLevelCode(ctx context.Context, parentID string) (string, error) {
	menuList, err := a.MenuModel.QueryList(ctx, schema.MenuListQueryParam{
		ParentID: parentID,
	})
	if err != nil {
		return "", err
	}

	levelCodes := schema.MenuList(menuList).ToLevelCodes()
	levelCode := util.GetLevelCode(levelCodes)
	if len(levelCode) == 0 {
		return "", errors.NewInternalServerError("分级码生成失败")
	}
	return levelCode, nil
}

// Create 创建数据
func (a *Menu) Create(ctx context.Context, item schema.Menu) (string, error) {
	exists, err := a.MenuModel.CheckCode(ctx, item.Code, item.ParentID)
	if err != nil {
		return "", err
	} else if exists {
		return "", errors.NewBadRequestError("同一父级下编号不允许重复")
	}

	item.RecordID = util.MustUUID()
	err = a.CommonBll.ExecTrans(ctx, func(ctx context.Context) error {
		a.lock.Lock()
		defer a.lock.Unlock()

		levelCode, err := a.getLevelCode(ctx, item.ParentID)
		if err != nil {
			return err
		}
		item.LevelCode = levelCode
		return a.MenuModel.Create(ctx, item)
	})
	if err != nil {
		return "", err
	}

	return item.RecordID, nil
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
	} else if item.Code != oldItem.Code {
		exists, err := a.MenuModel.CheckCode(ctx, item.Code, item.ParentID)
		if err != nil {
			return err
		} else if exists {
			return errors.NewBadRequestError("同一父级下编号不允许重复")
		}
	}
	item.LevelCode = oldItem.LevelCode

	return a.CommonBll.ExecTrans(ctx, func(ctx context.Context) error {
		a.lock.Lock()
		defer a.lock.Unlock()

		// 如果父级更新，需要更新当前节点及节点下级的分级码
		if item.ParentID != oldItem.ParentID {
			levelCode, err := a.getLevelCode(ctx, item.ParentID)
			if err != nil {
				return err
			}
			item.LevelCode = levelCode

			oldLevelCode := oldItem.LevelCode
			menuList, err := a.MenuModel.QueryList(ctx, schema.MenuListQueryParam{
				LevelCode: oldLevelCode,
			})
			if err != nil {
				return err
			}

			for _, menu := range menuList {
				if menu.RecordID == recordID {
					continue
				}
				newLevelCode := levelCode + menu.LevelCode[len(oldLevelCode):]
				err = a.MenuModel.UpdateLevelCode(ctx, menu.RecordID, newLevelCode)
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

			return a.MenuModel.Delete(ctx, recordID)
		}
		return nil
	})
}
