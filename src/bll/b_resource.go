package bll

import (
	"context"

	"github.com/LyricTian/gin-admin/src/errors"
	"github.com/LyricTian/gin-admin/src/model"
	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/LyricTian/gin-admin/src/util"
)

// 定义错误
var (
	ErrResourcePathAndMethodExists = errors.NewBadRequestError("访问路径和请求方式已经存在")
)

// Resource 资源管理
type Resource struct {
	ResourceModel model.IResource `inject:"IResource"`
	CommonBll     *Common         `inject:""`
}

// Query 查询数据
func (a *Resource) Query(ctx context.Context, params schema.ResourceQueryParam, pp *schema.PaginationParam) ([]*schema.Resource, *schema.PaginationResult, error) {
	result, err := a.ResourceModel.Query(ctx, params, schema.ResourceQueryOptions{PageParam: pp})
	if err != nil {
		return nil, nil, err
	}
	return result.Data, result.PageResult, nil
}

// Get 查询指定数据
func (a *Resource) Get(ctx context.Context, recordID string) (*schema.Resource, error) {
	item, err := a.ResourceModel.Get(ctx, recordID)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.ErrNotFound
	}

	return item, nil
}

func (a *Resource) check(ctx context.Context, item schema.Resource, oldItem *schema.Resource) error {
	if oldItem == nil || oldItem.Path != item.Path || oldItem.Method != item.Method {
		exists, err := a.ResourceModel.CheckPathAndMethod(ctx, item.Path, item.Method)
		if err != nil {
			return err
		} else if exists {
			return ErrResourcePathAndMethodExists
		}
	}
	return nil
}

// Create 创建数据
func (a *Resource) Create(ctx context.Context, item schema.Resource) (*schema.Resource, error) {
	err := a.check(ctx, item, nil)
	if err != nil {
		return nil, err
	}

	item.RecordID = util.MustUUID()
	err = a.ResourceModel.Create(ctx, item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// Update 更新数据
func (a *Resource) Update(ctx context.Context, recordID string, item schema.Resource) error {
	oldItem, err := a.ResourceModel.Get(ctx, recordID)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	err = a.check(ctx, item, oldItem)
	if err != nil {
		return err
	}

	return a.ResourceModel.Update(ctx, recordID, item)
}

// Delete 删除数据
func (a *Resource) Delete(ctx context.Context, recordIDs ...string) error {
	return a.CommonBll.ExecTrans(ctx, func(ctx context.Context) error {
		for _, recordID := range recordIDs {
			err := a.ResourceModel.Delete(ctx, recordID)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
