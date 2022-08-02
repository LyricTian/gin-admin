package test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/LyricTian/gin-admin/v9/internal/module/rbac/typed"
	"github.com/stretchr/testify/assert"
)

func TestRBACMenu(t *testing.T) {
	assert := assert.New(t)
	addMenu := typed.MenuCreate{
		Name:     "test",
		Sequence: 1,
		Icon:     "icon",
		Link:     "/test",
		Remark:   "test",
		Hide:     true,
		Actions: typed.MenuActions{
			&typed.MenuAction{
				Code: "ro",
				Name: "Readonly",
				Resources: typed.MenuActionResources{
					&typed.MenuActionResource{
						Path:   "/api/v1/test",
						Method: "GET",
					},
				},
			},
		},
	}

	addRes := httptest.NewRecorder()
	engine.ServeHTTP(addRes, newPostRequest(addMenu, "/api/rbac/v1/menus"))
	assert.Equal(200, addRes.Code)

	var addedMenu typed.Menu
	err := json.Unmarshal(addRes.Body.Bytes(), &addedMenu)
	assert.Nil(err)
	assert.NotEmpty(addedMenu.ID)

	getRes := httptest.NewRecorder()
	engine.ServeHTTP(getRes, newGetRequest(nil, "/api/rbac/v1/menus/%s", addedMenu.ID))
	assert.Equal(200, getRes.Code)
	var getMenu typed.Menu
	err = json.Unmarshal(getRes.Body.Bytes(), &getMenu)
	assert.Nil(err)
	assert.Equal(addedMenu.ID, getMenu.ID)
	assert.Equal(addMenu.Name, getMenu.Name)
	assert.Equal(addMenu.Sequence, getMenu.Sequence)
	assert.Equal(len(addMenu.Actions), len(getMenu.Actions))

	addMenu.Name = "test2"
	updateRes := httptest.NewRecorder()
	engine.ServeHTTP(updateRes, newPutRequest(addMenu, "/api/rbac/v1/menus/%s", addedMenu.ID))
	assert.Equal(200, updateRes.Code)

	getRes = httptest.NewRecorder()
	engine.ServeHTTP(getRes, newGetRequest(nil, "/api/rbac/v1/menus/%s", addedMenu.ID))
	assert.Equal(200, getRes.Code)
	var updatedMenu typed.Menu
	err = json.Unmarshal(getRes.Body.Bytes(), &updatedMenu)
	assert.Nil(err)
	assert.Equal(addedMenu.ID, updatedMenu.ID)
	assert.Equal(addMenu.Name, updatedMenu.Name)
	assert.NotEqual(getMenu.Name, updatedMenu.Name)

	deleteRes := httptest.NewRecorder()
	engine.ServeHTTP(deleteRes, newDeleteRequest("/api/rbac/v1/menus/%s", addedMenu.ID))
	assert.Equal(200, deleteRes.Code)
}
