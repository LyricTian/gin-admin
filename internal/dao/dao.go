package dao

import (
	"strings"

	"github.com/google/wire"
	"gorm.io/gorm"

	"github.com/LyricTian/gin-admin/v9/internal/config"
	"github.com/LyricTian/gin-admin/v9/internal/dao/repo"
	"github.com/LyricTian/gin-admin/v9/internal/dao/util"
	"github.com/LyricTian/gin-admin/v9/internal/schema"
) // end

// RepoSet repo injection
var RepoSet = wire.NewSet(
	util.TransSet,
	repo.DemoSet,
) // end

// Define repo type alias
type (
	TransRepo = util.Trans
	DemoRepo  = repo.DemoRepo
) // end

// Auto migration for given models
func AutoMigrate(db *gorm.DB) error {
	if dbType := config.C.Gorm.DBType; strings.ToLower(dbType) == "mysql" {
		db = db.Set("gorm:table_options", "ENGINE=InnoDB")
	}

	err := db.AutoMigrate(
		new(schema.Demo),
	) // end
	if err != nil {
		return err
	}

	return nil
}
