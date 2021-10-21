package user

import (
	"context"

	"gorm.io/gorm"

	"github.com/LyricTian/gin-admin/v8/internal/app/dao/util"
	"github.com/LyricTian/gin-admin/v8/internal/app/schema"
	"github.com/LyricTian/gin-admin/v8/pkg/util/structure"
)

func GetUserDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDBWithModel(ctx, defDB, new(User))
}

type SchemaUser schema.User

func (a SchemaUser) ToUser() *User {
	item := new(User)
	structure.Copy(a, item)
	return item
}

type User struct {
	util.Model
	UserName string  `gorm:"size:64;uniqueIndex;default:'';not null;"` // 用户名
	RealName string  `gorm:"size:64;index;default:'';"`                // 真实姓名
	Password string  `gorm:"size:40;default:'';"`                      // 密码
	Email    *string `gorm:"size:255;"`                                // 邮箱
	Phone    *string `gorm:"size:20;"`                                 // 手机号
	Status   int     `gorm:"index;default:0;"`                         // 状态(1:启用 2:停用)
	Creator  uint64  `gorm:""`                                         // 创建者
}

func (a User) ToSchemaUser() *schema.User {
	item := new(schema.User)
	structure.Copy(a, item)
	return item
}

type Users []*User

func (a Users) ToSchemaUsers() []*schema.User {
	list := make([]*schema.User, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaUser()
	}
	return list
}
