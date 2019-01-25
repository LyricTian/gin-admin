package entity

import (
	"github.com/LyricTian/gin-admin/src/schema"
)

// GetUserTableName 获取用户表名
func GetUserTableName() string {
	return User{}.TableName()
}

// GetUserRoleTableName 获取用户角色关联表名
func GetUserRoleTableName() string {
	return UserRole{}.TableName()
}

// SchemaUser 用户对象
type SchemaUser schema.User

// ToUser 转换为用户实体
func (a SchemaUser) ToUser() *User {
	item := &User{
		RecordID: a.RecordID,
		UserName: a.UserName,
		RealName: a.RealName,
		Password: a.Password,
		Status:   a.Status,
	}
	return item
}

// User 用户实体
type User struct {
	Model
	RecordID string `gorm:"column:record_id;size:36;unique_index;"` // 记录内码
	UserName string `gorm:"column:user_name;size:64;index;"`        // 用户名
	RealName string `gorm:"column:real_name;size:32;index;"`        // 真实姓名
	Password string `gorm:"column:password;size:40;"`               // 密码(sha1(md5(明文))加密)
	Status   int    `gorm:"column:status;index;"`                   // 状态(1:启用 2:停用)
	Creator  string `gorm:"column:creator;size:36;"`                // 创建者
}

// TableName 表名
func (a User) TableName() string {
	return a.Model.TableName("user")
}

// ToSchemaUser 转换为用户对象
func (a User) ToSchemaUser(includePassword bool) *schema.User {
	item := &schema.User{
		RecordID:  a.RecordID,
		UserName:  a.UserName,
		RealName:  a.RealName,
		Status:    a.Status,
		CreatedAt: a.CreatedAt,
	}
	if includePassword {
		item.Password = a.Password
	}
	return item
}

// Users 用户列表
type Users []*User

// ToSchemaUsers 转换为用户对象列表
func (a Users) ToSchemaUsers(includePassword bool) []*schema.User {
	list := make([]*schema.User, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaUser(includePassword)
	}
	return list
}

// UserRole 用户角色关联实体
type UserRole struct {
	Model
	UserID string `gorm:"column:user_id;size:36;index;"` // 用户内码
	RoleID string `gorm:"column:role_id;size:36;index;"` // 角色内码
}

// TableName 表名
func (a UserRole) TableName() string {
	return a.Model.TableName("user_role")
}

// UserRoles 用户角色关联列表
type UserRoles []*UserRole

// ToRoleIDs 转换为角色ID列表
func (a UserRoles) ToRoleIDs() []string {
	roleIDs := make([]string, len(a))
	for i, item := range a {
		roleIDs[i] = item.RoleID
	}
	return roleIDs
}
