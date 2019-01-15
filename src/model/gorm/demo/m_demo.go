package demo

import (
	"context"

	"github.com/LyricTian/gin-admin/src/logger"
	"github.com/LyricTian/gin-admin/src/model/gorm/common"
	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/LyricTian/gin-admin/src/service/gormplus"
	"github.com/LyricTian/gin-admin/src/util"
	"github.com/pkg/errors"
)

// NewModel 实例化示例存储
func NewModel(db *gormplus.DB) *Model {
	db.AutoMigrate(&Demo{})

	return &Model{db}
}

// Model 示例程序存储
type Model struct {
	db *gormplus.DB
}

// QueryPage 查询分页数据
func (a *Model) QueryPage(ctx context.Context, params schema.DemoQueryParam, pageIndex, pageSize uint) (int, []schema.DemoQueryResult, error) {
	span := logger.StartSpan(ctx, "查询分页数据", "demo.QueryPage")
	defer span.Finish()

	db := a.db.Model(Demo{})
	if v := params.Code; v != "" {
		db = db.Where("code LIKE ?", "%"+v+"%")
	}
	if v := params.Name; v != "" {
		db = db.Where("name LIKE ?", "%"+v+"%")
	}
	if v := params.Status; v > 0 {
		db = db.Where("status=?", v)
	}

	var items []Demo
	count, err := a.db.FindPage(db, pageIndex, pageSize, &items)
	if err != nil {
		span.Errorf(err.Error())
		return 0, nil, errors.New("查询分页数据条数发生错误")
	}

	dataItems := make([]schema.DemoQueryResult, len(items))
	util.FillStructs(items, dataItems)
	return count, dataItems, nil
}

// Get 查询指定数据
func (a *Model) Get(ctx context.Context, recordID string) (*schema.Demo, error) {
	span := logger.StartSpan(ctx, "查询指定数据", "demo.Get")
	defer span.Finish()

	db := a.db.Where("record_id=?", recordID)
	var item Demo
	err := a.db.FindOne(db, &item)
	if err != nil {
		span.Errorf(err.Error())
		return nil, errors.New("查询指定数据发生错误")
	}

	dataItem := new(schema.Demo)
	util.FillStruct(item, dataItem)
	return dataItem, nil
}

// Check 检查数据是否存在
func (a *Model) Check(ctx context.Context, recordID string) (bool, error) {
	span := logger.StartSpan(ctx, "检查数据是否存在", "demo.Check")
	defer span.Finish()

	var count int
	result := a.db.Model(Demo{}).Where("record_id=?", recordID).Count(&count)
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return false, errors.New("检查数据是否存在发生错误")
	}

	return count > 0, nil
}

// CheckCode 检查编号是否存在
func (a *Model) CheckCode(ctx context.Context, code string) (bool, error) {
	span := logger.StartSpan(ctx, "检查编号是否存在", "demo.CheckCode")
	defer span.Finish()

	var count int
	result := a.db.Model(Demo{}).Where("code=?", code).Count(&count)
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return false, errors.New("检查编号是否存在发生错误")
	}

	return count > 0, nil
}

// Create 创建数据
func (a *Model) Create(ctx context.Context, item schema.Demo) error {
	span := logger.StartSpan(ctx, "创建数据", "demo.Create")
	defer span.Finish()

	demo := new(Demo)
	util.FillStruct(item, demo)
	result := a.db.Create(demo)
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("创建数据发生错误")
	}
	return nil
}

// Update 更新数据
func (a *Model) Update(ctx context.Context, recordID string, item schema.Demo) error {
	span := logger.StartSpan(ctx, "更新数据", "demo.Update")
	defer span.Finish()

	demo := new(Demo)
	util.FillStruct(item, demo)
	result := a.db.Model(Demo{}).Where("record_id=?", recordID).Omit("record_id").Updates(demo)
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("更新数据发生错误")
	}
	return nil
}

// Delete 删除数据
func (a *Model) Delete(ctx context.Context, trans interface{}, recordID string) error {
	span := logger.StartSpan(ctx, "删除数据", "demo.Delete")
	defer span.Finish()

	db := a.db
	if trans != nil {
		db = common.GetTrans(trans)
	}
	result := db.Where("record_id=?", recordID).Delete(Demo{})
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("删除数据发生错误")
	}
	return nil
}

// UpdateStatus 更新状态
func (a *Model) UpdateStatus(ctx context.Context, recordID string, status int) error {
	span := logger.StartSpan(ctx, "更新状态", "demo.UpdateStatus")
	defer span.Finish()

	result := a.db.Model(Demo{}).Where("record_id=?", recordID).Update("status", status)
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("更新状态发生错误")
	}
	return nil
}
