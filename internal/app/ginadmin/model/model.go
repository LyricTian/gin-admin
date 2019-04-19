package model

// Common 提供统一的存储接口
type Common struct {
	Trans ITrans
	Demo  IDemo
	Menu  IMenu
	Role  IRole
	User  IUser
}
