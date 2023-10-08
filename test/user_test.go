package test

import (
	"net/http"
	"testing"

	"github.com/LyricTian/gin-admin/v10/internal/mods/rbac/schema"
	"github.com/LyricTian/gin-admin/v10/pkg/crypto/hash"
	"github.com/LyricTian/gin-admin/v10/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	e := tester(t)

	menuFormItem := schema.MenuForm{
		Code:        "user",
		Name:        "User management",
		Description: "User management",
		Sequence:    7,
		Type:        "page",
		Path:        "/system/user",
		Properties:  `{"icon":"user"}`,
		Status:      schema.MenuStatusEnabled,
	}

	var menu schema.Menu
	e.POST(baseAPI + "/menus").WithJSON(menuFormItem).
		Expect().Status(http.StatusOK).JSON().Decode(&util.ResponseResult{Data: &menu})

	assert := assert.New(t)
	assert.NotEmpty(menu.ID)
	assert.Equal(menuFormItem.Code, menu.Code)
	assert.Equal(menuFormItem.Name, menu.Name)
	assert.Equal(menuFormItem.Description, menu.Description)
	assert.Equal(menuFormItem.Sequence, menu.Sequence)
	assert.Equal(menuFormItem.Type, menu.Type)
	assert.Equal(menuFormItem.Path, menu.Path)
	assert.Equal(menuFormItem.Properties, menu.Properties)
	assert.Equal(menuFormItem.Status, menu.Status)

	roleFormItem := schema.RoleForm{
		Code: "user",
		Name: "Normal",
		Menus: schema.RoleMenus{
			{MenuID: menu.ID},
		},
		Description: "Normal",
		Sequence:    8,
		Status:      schema.RoleStatusEnabled,
	}

	var role schema.Role
	e.POST(baseAPI + "/roles").WithJSON(roleFormItem).Expect().Status(http.StatusOK).JSON().Decode(&util.ResponseResult{Data: &role})
	assert.NotEmpty(role.ID)
	assert.Equal(roleFormItem.Code, role.Code)
	assert.Equal(roleFormItem.Name, role.Name)
	assert.Equal(roleFormItem.Description, role.Description)
	assert.Equal(roleFormItem.Sequence, role.Sequence)
	assert.Equal(roleFormItem.Status, role.Status)
	assert.Equal(len(roleFormItem.Menus), len(role.Menus))

	userFormItem := schema.UserForm{
		Username: "test",
		Name:     "Test",
		Password: hash.MD5String("test"),
		Phone:    "0720",
		Email:    "test@gmail.com",
		Remark:   "test user",
		Status:   schema.UserStatusActivated,
		Roles:    schema.UserRoles{{RoleID: role.ID}},
	}

	var user schema.User
	e.POST(baseAPI + "/users").WithJSON(userFormItem).Expect().Status(http.StatusOK).JSON().Decode(&util.ResponseResult{Data: &user})
	assert.NotEmpty(user.ID)
	assert.Equal(userFormItem.Username, user.Username)
	assert.Equal(userFormItem.Name, user.Name)
	assert.Equal(userFormItem.Phone, user.Phone)
	assert.Equal(userFormItem.Email, user.Email)
	assert.Equal(userFormItem.Remark, user.Remark)
	assert.Equal(userFormItem.Status, user.Status)
	assert.Equal(len(userFormItem.Roles), len(user.Roles))

	var users schema.Users
	e.GET(baseAPI+"/users").WithQuery("username", userFormItem.Username).Expect().Status(http.StatusOK).JSON().Decode(&util.ResponseResult{Data: &users})
	assert.GreaterOrEqual(len(users), 1)

	newName := "Test 1"
	newStatus := schema.UserStatusFreezed
	user.Name = newName
	user.Status = newStatus
	e.PUT(baseAPI + "/users/" + user.ID).WithJSON(user).Expect().Status(http.StatusOK)

	var getUser schema.User
	e.GET(baseAPI + "/users/" + user.ID).Expect().Status(http.StatusOK).JSON().Decode(&util.ResponseResult{Data: &getUser})
	assert.Equal(newName, getUser.Name)
	assert.Equal(newStatus, getUser.Status)

	e.DELETE(baseAPI + "/users/" + user.ID).Expect().Status(http.StatusOK)
	e.GET(baseAPI + "/users/" + user.ID).Expect().Status(http.StatusNotFound)

	e.DELETE(baseAPI + "/roles/" + role.ID).Expect().Status(http.StatusOK)
	e.GET(baseAPI + "/roles/" + role.ID).Expect().Status(http.StatusNotFound)

	e.DELETE(baseAPI + "/menus/" + menu.ID).Expect().Status(http.StatusOK)
	e.GET(baseAPI + "/menus/" + menu.ID).Expect().Status(http.StatusNotFound)
}
