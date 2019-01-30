package model

import (
	"context"
	"fmt"

	"github.com/LyricTian/gin-admin/src/errors"
	"github.com/LyricTian/gin-admin/src/logger"
	"github.com/LyricTian/gin-admin/src/model/gorm/entity"
	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/LyricTian/gin-admin/src/service/gormplus"
)

// InitResource 初始化资源存储
func InitResource(db *gormplus.DB) *Resource {
	db.AutoMigrate(new(entity.Resource))
	return NewResource(db)
}

// NewResource 实例化资源存储
func NewResource(db *gormplus.DB) *Resource {
	return &Resource{db: db}
}

// Resource 资源存储
type Resource struct {
	db *gormplus.DB
}

func (a *Resource) getFuncName(name string) string {
	return fmt.Sprintf("gorm.model.Resource.%s", name)
}

func (a *Resource) getResourceDB(ctx context.Context) *gormplus.DB {
	return FromDBWithModel(ctx, a.db, entity.Resource{})
}

func (a *Resource) getQueryOption(opts ...schema.ResourceQueryOptions) schema.ResourceQueryOptions {
	if len(opts) > 0 {
		return opts[0]
	}
	return schema.ResourceQueryOptions{}
}

// Query 查询数据
func (a *Resource) Query(ctx context.Context, params schema.ResourceQueryParam, opts ...schema.ResourceQueryOptions) (schema.ResourceQueryResult, error) {
	span := logger.StartSpan(ctx, "查询数据", a.getFuncName("Query"))
	defer span.Finish()

	db := a.getResourceDB(ctx).DB
	if v := params.Name; v != "" {
		db = db.Where("name LIKE ?", "%"+v+"%")
	}
	if v := params.Path; v != "" {
		db = db.Where("path LIKE ?", v+"%")
	}
	db = db.Order("id DESC")

	var qr schema.ResourceQueryResult
	opt := a.getQueryOption(opts...)
	var items entity.Resources
	pr, err := WrapPageQuery(db, opt.PageParam, &items)
	if err != nil {
		span.Errorf(err.Error())
		return qr, errors.New("查询数据发生错误")
	}
	qr.PageResult = pr
	qr.Data = items.ToSchemaResources()

	return qr, nil
}

// Get 查询指定数据
func (a *Resource) Get(ctx context.Context, recordID string) (*schema.Resource, error) {
	span := logger.StartSpan(ctx, "查询指定数据", a.getFuncName("Get"))
	defer span.Finish()

	db := a.getResourceDB(ctx).Where("record_id=?", recordID)
	var item entity.Resource
	ok, err := a.db.FindOne(db, &item)
	if err != nil {
		span.Errorf(err.Error())
		return nil, errors.New("查询指定数据发生错误")
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaResource(), nil
}

// CheckPathAndMethod 检查访问路径和请求方法是否存在
func (a *Resource) CheckPathAndMethod(ctx context.Context, path, method string) (bool, error) {
	span := logger.StartSpan(ctx, "检查访问路径和请求方法是否存在", a.getFuncName("CheckPathAndMethod"))
	defer span.Finish()

	db := a.getResourceDB(ctx).Where("path=? AND method=?", path, method)
	exists, err := a.db.Check(db)
	if err != nil {
		span.Errorf(err.Error())
		return false, errors.New("检查访问路径和请求方法是否存在发生错误")
	}
	return exists, nil
}

// Create 创建数据
func (a *Resource) Create(ctx context.Context, item schema.Resource) error {
	span := logger.StartSpan(ctx, "创建数据", a.getFuncName("Create"))
	defer span.Finish()

	eitem := entity.SchemaResource(item).ToResource()
	eitem.Creator = FromUserID(ctx)
	result := a.getResourceDB(ctx).Create(eitem)
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("创建数据发生错误")
	}
	return nil
}

// Update 更新数据
func (a *Resource) Update(ctx context.Context, recordID string, item schema.Resource) error {
	span := logger.StartSpan(ctx, "更新数据", a.getFuncName("Update"))
	defer span.Finish()

	eitem := entity.SchemaResource(item).ToResource()
	result := a.getResourceDB(ctx).Where("record_id=?", recordID).Omit("record_id", "creator").Updates(eitem)
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("更新数据发生错误")
	}
	return nil
}

// Delete 删除数据
func (a *Resource) Delete(ctx context.Context, recordID string) error {
	span := logger.StartSpan(ctx, "删除数据", a.getFuncName("Delete"))
	defer span.Finish()

	result := a.getResourceDB(ctx).Where("record_id=?", recordID).Delete(entity.Resource{})
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("删除数据发生错误")
	}
	return nil
}
