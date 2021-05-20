package gormx

import (
	"gorm.io/gorm/schema"
	"strings"
	"time"

	"github.com/LyricTian/gin-admin/v7/internal/app/config"
	"github.com/LyricTian/gin-admin/v7/internal/app/model/gormx/entity"
	"github.com/LyricTian/gin-admin/v7/pkg/logger"
	"gorm.io/gorm"
)

// Config 配置参数
type Config struct {
	Debug        bool
	Dialector    *gorm.Dialector
	MaxLifetime  int
	MaxOpenConns int
	MaxIdleConns int
	TablePrefix  string
}

// NewDB 创建DB实例
func NewDB(c *Config) (*gorm.DB, func(), error) {
	db, err := gorm.Open(*c.Dialector, &gorm.Config{ // https://gorm.io/zh_CN/docs/gorm_config.html
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   c.TablePrefix, // 全局表前缀
			SingularTable: true,          // 禁用表名复数
		},
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用自动创建外键约束
		SkipDefaultTransaction:                   true, // 跳过默认事务
	})
	if err != nil {
		return nil, nil, err
	}

	if c.Debug {
		db = db.Debug()
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}
	cleanFunc := func() {
		err := sqlDB.Close()
		if err != nil {
			logger.Errorf("Gorm db close error: %s", err.Error())
		}
	}

	//在完成初始化后，GORMV2 会自动ping数据库以检查数据库的可用性，以下代码注释
	//err = sqlDB.Ping()
	//if err != nil {
	//	return nil, cleanFunc, err
	//}

	sqlDB.SetMaxIdleConns(c.MaxIdleConns)
	sqlDB.SetMaxOpenConns(c.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(c.MaxLifetime) * time.Second)
	return db, cleanFunc, nil
}

// AutoMigrate 自动映射数据表
func AutoMigrate(db *gorm.DB) error {
	if dbType := config.C.Gorm.DBType; strings.ToLower(dbType) == "mysql" {
		db = db.Set("gorm:table_options", "ENGINE=InnoDB")
	}

	return db.AutoMigrate(
		new(entity.Demo),
		new(entity.MenuAction),
		new(entity.MenuActionResource),
		new(entity.Menu),
		new(entity.RoleMenu),
		new(entity.Role),
		new(entity.UserRole),
		new(entity.User),
	)
}
