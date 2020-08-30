package entity

import (
	"context"
	"time"

	"github.com/LyricTian/gin-admin/v7/internal/app/schema"
	"github.com/LyricTian/gin-admin/v7/pkg/util/structure"
	"github.com/jinzhu/gorm"
)

// GetUserDB 获取用户存储
func GetUserDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return GetDBWithModel(ctx, defDB, new(User))
}

// SchemaUser 用户对象
type SchemaUser schema.User

// ToUser 转换为用户实体
func (a SchemaUser) ToUser() *User {
	item := new(User)
	structure.Copy(a, item)
	return item
}

// User 用户实体
type User struct {
	ID        string     `gorm:"column:id;primary_key;size:36;"`
	UserName  string     `gorm:"column:user_name;size:64;index;default:'';not null;"` // 用户名
	RealName  string     `gorm:"column:real_name;size:64;index;default:'';not null;"` // 真实姓名
	Password  string     `gorm:"column:password;size:40;default:'';not null;"`        // 密码(sha1(md5(明文))加密)
	Email     *string    `gorm:"column:email;size:255;index;"`                        // 邮箱
	Phone     *string    `gorm:"column:phone;size:20;index;"`                         // 手机号
	Status    int        `gorm:"column:status;index;default:0;not null;"`             // 状态(1:启用 2:停用)
	Creator   string     `gorm:"column:creator;size:36;"`                             // 创建者
	CreatedAt time.Time  `gorm:"column:created_at;index;"`
	UpdatedAt time.Time  `gorm:"column:updated_at;index;"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index;"`
}

// ToSchemaUser 转换为用户对象
func (a User) ToSchemaUser() *schema.User {
	item := new(schema.User)
	structure.Copy(a, item)
	return item
}

// Users 用户实体列表
type Users []*User

// ToSchemaUsers 转换为用户对象列表
func (a Users) ToSchemaUsers() []*schema.User {
	list := make([]*schema.User, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaUser()
	}
	return list
}
