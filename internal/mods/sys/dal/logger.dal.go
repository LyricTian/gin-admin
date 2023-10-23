package dal

import (
	"context"
	"fmt"

	rbacSchema "github.com/LyricTian/gin-admin/v10/internal/mods/rbac/schema"
	"github.com/LyricTian/gin-admin/v10/internal/mods/sys/schema"
	"github.com/LyricTian/gin-admin/v10/pkg/errors"
	"github.com/LyricTian/gin-admin/v10/pkg/util"
	"gorm.io/gorm"
)

// Get logger storage instance
func GetLoggerDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDB(ctx, defDB).Model(new(schema.Logger))
}

// Logger management
type Logger struct {
	DB *gorm.DB
}

// Query loggers from the database based on the provided parameters and options.
func (a *Logger) Query(ctx context.Context, params schema.LoggerQueryParam, opts ...schema.LoggerQueryOptions) (*schema.LoggerQueryResult, error) {
	var opt schema.LoggerQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	db := a.DB.Table(fmt.Sprintf("%s AS a", new(schema.Logger).TableName()))
	db = db.Joins(fmt.Sprintf("left join %s b on a.user_id=b.id", new(rbacSchema.User).TableName()))
	db = db.Select("a.*,b.name as user_name,b.username as login_name")

	if v := params.Level; v != "" {
		db = db.Where("a.level = ?", v)
	}
	if v := params.LikeMessage; len(v) > 0 {
		db = db.Where("a.message LIKE ?", "%"+v+"%")
	}
	if v := params.TraceID; v != "" {
		db = db.Where("a.trace_id = ?", v)
	}
	if v := params.LikeUserName; v != "" {
		db = db.Where("b.username LIKE ?", "%"+v+"%")
	}
	if v := params.Tag; v != "" {
		db = db.Where("a.tag = ?", v)
	}
	if start, end := params.StartTime, params.EndTime; start != "" && end != "" {
		db = db.Where("a.created_at BETWEEN ? AND ?", start, end)
	}

	var list schema.Loggers
	pageResult, err := util.WrapPageQuery(ctx, db, params.PaginationParam, opt.QueryOptions, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	queryResult := &schema.LoggerQueryResult{
		PageResult: pageResult,
		Data:       list,
	}
	return queryResult, nil
}
