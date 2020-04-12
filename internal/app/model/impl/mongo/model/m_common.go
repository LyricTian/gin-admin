package model

import (
	"context"
	"time"

	"github.com/LyricTian/gin-admin/internal/app/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TransFunc 定义事务执行函数
type TransFunc func(context.Context) error

// ExecTrans 执行事务
func ExecTrans(ctx context.Context, cli *mongo.Client, fn TransFunc) error {
	transModel := &Trans{Client: cli}
	return transModel.Exec(ctx, fn)
}

// WrapPageQuery 包装带有分页的查询
func WrapPageQuery(ctx context.Context, c *mongo.Collection, pp schema.PaginationParam, filter interface{}, out interface{}, opts ...*options.FindOptions) (*schema.PaginationResult, error) {
	if pp.OnlyCount {
		count, err := c.CountDocuments(ctx, filter)
		if err != nil {
			return nil, err
		}
		return &schema.PaginationResult{Total: int(count)}, nil
	} else if !pp.Pagination {
		cursor, err := c.Find(ctx, filter, opts...)
		if err != nil {
			return nil, err
		}
		err = cursor.All(ctx, out)
		return nil, err
	}

	total, err := FindPage(ctx, c, pp, filter, out, opts...)
	if err != nil {
		return nil, err
	}

	return &schema.PaginationResult{
		Total:    total,
		Current:  pp.GetCurrent(),
		PageSize: pp.GetPageSize(),
	}, nil
}

// FindPage 查询分页数据
func FindPage(ctx context.Context, c *mongo.Collection, pp schema.PaginationParam, filter interface{}, out interface{}, opts ...*options.FindOptions) (int, error) {
	count, err := c.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	} else if count == 0 {
		return 0, nil
	}

	current, pageSize := pp.GetCurrent(), pp.GetPageSize()
	opt := new(options.FindOptions)
	if len(opts) > 0 {
		opt = opts[0]
	}
	opt.SetSkip(int64((current - 1) * pageSize))
	opt.SetLimit(int64(pageSize))

	cursor, err := c.Find(ctx, filter, opt)
	if err != nil {
		return 0, err
	}
	err = cursor.All(ctx, out)
	return int(count), err
}

// FindOne 查询单条数据
func FindOne(ctx context.Context, c *mongo.Collection, filter, out interface{}) (bool, error) {
	result := c.FindOne(ctx, filter)
	if err := result.Err(); err != nil {
		if err == mongo.ErrNilDocument {
			return false, nil
		}
		return false, err
	}
	err := result.Decode(out)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Insert 插入数据
func Insert(ctx context.Context, c *mongo.Collection, doc interface{}) error {
	_, err := c.InsertOne(ctx, doc)
	return err
}

// InsertMany 插入多条数据
func InsertMany(ctx context.Context, c *mongo.Collection, docs ...interface{}) error {
	_, err := c.InsertMany(ctx, docs)
	return err
}

// UpdateFields 更新指定字段数据
func UpdateFields(ctx context.Context, c *mongo.Collection, filter, doc interface{}) error {
	return Update(ctx, c, filter, doc)
}

// UpdateManyFields 更新多条指定字段的数据
func UpdateManyFields(ctx context.Context, c *mongo.Collection, filter, doc interface{}) error {
	return UpdateMany(ctx, c, filter, doc)
}

// Update 更新数据
func Update(ctx context.Context, c *mongo.Collection, filter, doc interface{}) error {
	_, err := c.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: doc}})
	return err
}

// UpdateMany 更新多条数据
func UpdateMany(ctx context.Context, c *mongo.Collection, filter, doc interface{}) error {
	_, err := c.UpdateMany(ctx, filter, bson.D{{Key: "$set", Value: doc}})
	return err
}

// Delete 删除数据
func Delete(ctx context.Context, c *mongo.Collection, filter interface{}) error {
	_, err := c.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: bson.M{"deleted_at": time.Now()}}})
	return err
}

// DeleteMany 删除多条数据
func DeleteMany(ctx context.Context, c *mongo.Collection, filter interface{}) error {
	_, err := c.UpdateMany(ctx, filter, bson.D{{Key: "$set", Value: bson.M{"deleted_at": time.Now()}}})
	return err
}

// DefaultFilter 默认的查询参数
func DefaultFilter(ctx context.Context, params ...bson.E) bson.D {
	var d bson.D
	if len(params) > 0 {
		d = append(d, params...)
	}
	d = append(d, Filter("deleted_at", bson.M{"$exists": 0}))
	return d
}

// RegexFilter 正则过滤
func RegexFilter(key, value string) bson.E {
	return bson.E{
		Key: key,
		Value: bson.M{
			"$regex":   value,
			"$options": "i",
		},
	}
}

// OrRegexFilter 正则过滤($or)
func OrRegexFilter(key, value string) bson.M {
	return bson.M{
		key: bson.M{
			"$regex":   value,
			"$options": "i",
		},
	}
}

// Filter 过滤
func Filter(key string, value interface{}) bson.E {
	return bson.E{
		Key:   key,
		Value: value,
	}
}

// OrderFieldFunc 排序字段转换函数
type OrderFieldFunc func(string) string

// ParseOrder 解析排序字段
func ParseOrder(items []*schema.OrderField, handle ...OrderFieldFunc) bson.D {
	d := make(bson.D, 0)
	for _, item := range items {
		key := item.Key
		if len(handle) > 0 {
			key = handle[0](key)
		}

		direction := 1
		if item.Direction == schema.OrderByDESC {
			direction = -1
		}
		d = append(d, bson.E{Key: key, Value: direction})
	}

	return d
}
