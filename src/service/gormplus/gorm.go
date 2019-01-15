package gormplus

import (
	"time"

	"github.com/jinzhu/gorm"

	// gorm存储注入
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Config 配置参数
type Config struct {
	Debug        bool
	DBType       string
	DSN          string
	MaxLifetime  int
	MaxOpenConns int
	MaxIdleConns int
}

// New 创建DB实例
func New(c Config) (*DB, error) {
	db, err := gorm.Open(c.DBType, c.DSN)
	if err != nil {
		return nil, err
	}

	if c.Debug {
		db = db.Debug()
	}

	err = db.DB().Ping()
	if err != nil {
		return nil, err
	}

	db.DB().SetMaxIdleConns(c.MaxIdleConns)
	db.DB().SetMaxOpenConns(c.MaxOpenConns)
	db.DB().SetConnMaxLifetime(time.Duration(c.MaxLifetime) * time.Second)
	return &DB{db}, nil
}

// Wrap 包装gorm
func Wrap(db *gorm.DB) *DB {
	return &DB{db}
}

// DB gorm扩展DB
type DB struct {
	*gorm.DB
}

// FindPage 查询分页数据
func (d *DB) FindPage(db *gorm.DB, pageIndex, pageSize uint, out interface{}) (int, error) {
	var count int
	result := db.Count(&count)
	if err := result.Error; err != nil {
		return 0, err
	} else if count == 0 {
		return 0, nil
	}

	result = db.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(out)
	if err := result.Error; err != nil {
		return 0, err
	}

	return count, nil
}

// FindOne 查询单条数据
func (d *DB) FindOne(db *gorm.DB, out interface{}) error {
	result := db.First(out)
	if err := result.Error; err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	return nil
}
