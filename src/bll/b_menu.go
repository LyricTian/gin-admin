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
func (a *Menu) QueryPage(ctx context.Context, params schema.MenuQueryParam, pp *schema.PaginationParam) ([]*schema.Menu, *schema.PaginationResult, error) {
	return a.MenuModel.Query(ctx, params, pp)
}

// QueryTree 查询菜单树
func (a *Menu) QueryTree(ctx context.Context) ([]*schema.MenuTreeResult, error) {
	items, _, err := a.MenuModel.Query(ctx, schema.MenuQueryParam{
		Types: []int{1, 2},
	}, nil)
	if err != nil {
		return nil, err
	}

	return schema.Menus(items).ToTreeResult(), nil
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

func (a *Menu) getLevelCode(ctx context.Context, parentItem *schema.Menu) (string, error) {
	params := schema.MenuQueryParam{}
	if parentItem == nil {
		var value string
		params.ParentID = &value
	} else {
		params.LevelCode = parentItem.LevelCode
	}
	menuList, _, err := a.MenuModel.Query(ctx, params, nil)
	if err != nil {
		return "", err
	}

	levelCodes := schema.Menus(menuList).ToLevelCodes()
	levelCode := util.GetLevelCode(levelCodes)
	if len(levelCode) == 0 {
		return "", errors.NewInternalServerError("分级码生成失败")
	}
	return levelCode, nil
}

func (a *Menu) checkAndGetParent(ctx context.Context, item schema.Menu, oldItem *schema.Menu) (*schema.Menu, error) {
	if oldItem == nil || oldItem.Code != item.Code {
		exists, err := a.MenuModel.CheckCode(ctx, item.Code, item.ParentID)
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
	case item.Type == 3 && (parentItem == nil || parentItem.Type != 2):
		return nil, errors.NewBadRequestError("资源类型只能依赖于功能")
	case item.Type == 2 && (parentItem == nil || parentItem.Type != 1):
		return nil, errors.NewBadRequestError("功能类型只能依赖于模块")
	case item.Type == 1 && parentItem != nil && parentItem.Type != 1:
		return nil, errors.NewBadRequestError("模块类型只能依赖于模块")
	}

	return parentItem, nil
}

// Create 创建数据
func (a *Menu) Create(ctx context.Context, item schema.Menu) (*schema.Menu, error) {
	parentItem, err := a.checkAndGetParent(ctx, item, nil)
	if err != nil {
		return nil, err
	}

	a.lock.Lock()
	defer a.lock.Unlock()

	item.RecordID = util.MustUUID()
	err = a.CommonBll.ExecTrans(ctx, func(ctx context.Context) error {
		levelCode, err := a.getLevelCode(ctx, parentItem)
		if err != nil {
			return err
		}
		item.LevelCode = levelCode
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
	item.LevelCode = oldItem.LevelCode

	a.lock.Lock()
	defer a.lock.Unlock()

	return a.CommonBll.ExecTrans(ctx, func(ctx context.Context) error {
		// 如果父级更新，需要更新当前节点及节点下级的分级码
		if item.ParentID != oldItem.ParentID {
			levelCode, err := a.getLevelCode(ctx, parentItem)
			if err != nil {
				return err
			}
			item.LevelCode = levelCode
			oldLevelCode := oldItem.LevelCode
			menuList, _, err := a.MenuModel.Query(ctx, schema.MenuQueryParam{
				LevelCode: oldLevelCode,
			}, nil)
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

			err = a.MenuModel.Delete(ctx, recordID)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
