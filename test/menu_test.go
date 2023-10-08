package test

import (
	"net/http"
	"testing"

	"github.com/LyricTian/gin-admin/v10/internal/mods/rbac/schema"
	"github.com/LyricTian/gin-admin/v10/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestMenu(t *testing.T) {
	e := tester(t)

	menuFormItem := schema.MenuForm{
		Code:        "menu",
		Name:        "Menu management",
		Description: "Menu management",
		Sequence:    9,
		Type:        "page",
		Path:        "/system/menu",
		Properties:  `{"icon":"menu"}`,
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

	var menus schema.Menus
	e.GET(baseAPI + "/menus").Expect().Status(http.StatusOK).JSON().Decode(&util.ResponseResult{Data: &menus})
	assert.GreaterOrEqual(len(menus), 1)

	newName := "Menu management 1"
	newStatus := schema.MenuStatusDisabled
	menu.Name = newName
	menu.Status = newStatus
	e.PUT(baseAPI + "/menus/" + menu.ID).WithJSON(menu).Expect().Status(http.StatusOK)

	var getMenu schema.Menu
	e.GET(baseAPI + "/menus/" + menu.ID).Expect().Status(http.StatusOK).JSON().Decode(&util.ResponseResult{Data: &getMenu})
	assert.Equal(newName, getMenu.Name)
	assert.Equal(newStatus, getMenu.Status)

	e.DELETE(baseAPI + "/menus/" + menu.ID).Expect().Status(http.StatusOK)
	e.GET(baseAPI + "/menus/" + menu.ID).Expect().Status(http.StatusNotFound)
}
