package gormdemo

import (
	"context"
	"fmt"

	gcontext "github.com/LyricTian/gin-admin/src/context"
	"github.com/LyricTian/gin-admin/src/errors"
	"github.com/LyricTian/gin-admin/src/logger"
	"github.com/LyricTian/gin-admin/src/model/gorm/common"
	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/LyricTian/gin-admin/src/service/gormplus"
	"github.com/LyricTian/gin-admin/src/util"
	"github.com/jinzhu/gorm"
)

// NewModel 实例化示例存储
func NewModel(db *gormplus.DB) *Model {
	db.AutoMigrate(new(Demo))
	return &Model{db}
}

// Model 示例程序存储
type Model struct {
	db *gormplus.DB
}

func (a *Model) getFuncName(name string) string {
	return fmt.Sprintf("gorm.demo.%s", name)
}

func (a *Model) getDemoDB(ctx context.Context) *gorm.DB {
	return gormcommon.FromTransDB(ctx, a.db).Model(Demo{})
}

// Query 查询数据
func (a *Model) Query(ctx context.Context, params schema.DemoQueryParam, pp *schema.PaginationParam) ([]*schema.Demo, *schema.PaginationResult, error) {
	span := logger.StartSpan(ctx, "查询数据", a.getFuncName("Query"))
	defer span.Finish()

	db := a.getDemoDB(ctx)
	if v := params.Code; v != "" {
		db = db.Where("code LIKE ?", "%"+v+"%")
	}
	if v := params.Name; v != "" {
		db = db.Where("name LIKE ?", "%"+v+"%")
	}
	if v := params.Status; v > 0 {
		db = db.Where("status=?", v)
	}
	db = db.Order("id DESC")

	var items []*Demo
	pageResult, err := gormcommon.WrapPageQuery(db, pp, &items)
	if err != nil {
		span.Errorf(err.Error())
		return nil, nil, errors.New("查询数据发生错误")
	}

	return Demos(items).ToSchemaDemos(), pageResult, nil
}

// Get 查询指定数据
func (a *Model) Get(ctx context.Context, recordID string) (*schema.Demo, error) {
	span := logger.StartSpan(ctx, "查询指定数据", a.getFuncName("Get"))
	defer span.Finish()

	db := a.getDemoDB(ctx).Where("record_id=?", recordID)
	var item Demo
	ok, err := a.db.FindOne(db, &item)
	if err != nil {
		span.Errorf(err.Error())
		return nil, errors.New("查询指定数据发生错误")
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaDemo(), nil
}

// CheckCode 检查编号是否存在
func (a *Model) CheckCode(ctx context.Context, code string) (bool, error) {
	span := logger.StartSpan(ctx, "检查编号是否存在", a.getFuncName("CheckCode"))
	defer span.Finish()

	db := a.getDemoDB(ctx).Where("code=?", code)
	exists, err := a.db.Check(db)
	if err != nil {
		span.Errorf(err.Error())
		return false, errors.New("检查编号是否存在发生错误")
	}
	return exists, nil
}

// Create 创建数据
func (a *Model) Create(ctx context.Context, item schema.Demo) error {
	span := logger.StartSpan(ctx, "创建数据", a.getFuncName("Create"))
	defer span.Finish()

	demo := new(Demo)
	_ = util.FillStruct(item, demo)
	demo.Creator, _ = gcontext.FromUserID(ctx)
	result := a.getDemoDB(ctx).Create(demo)
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("创建数据发生错误")
	}
	return nil
}

// Update 更新数据
func (a *Model) Update(ctx context.Context, recordID string, item schema.Demo) error {
	span := logger.StartSpan(ctx, "更新数据", a.getFuncName("Update"))
	defer span.Finish()

	demo := new(Demo)
	_ = util.FillStruct(item, demo)
	result := a.getDemoDB(ctx).Where("record_id=?", recordID).Omit("record_id", "creator").Updates(demo)
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("更新数据发生错误")
	}
	return nil
}

// Delete 删除数据
func (a *Model) Delete(ctx context.Context, recordID string) error {
	span := logger.StartSpan(ctx, "删除数据", a.getFuncName("Delete"))
	defer span.Finish()

	result := a.getDemoDB(ctx).Where("record_id=?", recordID).Delete(Demo{})
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("删除数据发生错误")
	}
	return nil
}

// UpdateStatus 更新状态
func (a *Model) UpdateStatus(ctx context.Context, recordID string, status int) error {
	span := logger.StartSpan(ctx, "更新状态", a.getFuncName("UpdateStatus"))
	defer span.Finish()

	result := a.getDemoDB(ctx).Where("record_id=?", recordID).Update("status", status)
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("更新状态发生错误")
	}
	return nil
}
