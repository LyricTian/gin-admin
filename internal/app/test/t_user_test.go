package test

import (
	"net/http/httptest"
	"testing"

	"github.com/LyricTian/gin-admin/v6/internal/app/schema"
	"github.com/LyricTian/gin-admin/v6/pkg/unique"
	"github.com/LyricTian/gin-admin/v6/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	const router = apiPrefix + "v1/users"
	var err error

	w := httptest.NewRecorder()

	// post /menus
	addMenuItem := &schema.Menu{
		Name:       unique.MustUUID().String(),
		ShowStatus: 1,
		Status:     1,
	}
	engine.ServeHTTP(w, newPostRequest(apiPrefix+"v1/menus", addMenuItem))
	assert.Equal(t, 200, w.Code)
	var addMenuItemRes ResID
	err = parseReader(w.Body, &addMenuItemRes)
	assert.Nil(t, err)

	// post /roles
	addRoleItem := &schema.Role{
		Name:   unique.MustUUID().String(),
		Status: 1,
		RoleMenus: schema.RoleMenus{
			&schema.RoleMenu{
				MenuID: addMenuItemRes.ID,
			},
		},
	}
	engine.ServeHTTP(w, newPostRequest(apiPrefix+"v1/roles", addRoleItem))
	assert.Equal(t, 200, w.Code)
	var addRoleItemRes ResID
	err = parseReader(w.Body, &addRoleItemRes)
	assert.Nil(t, err)

	// post /users
	addItem := &schema.User{
		UserName: unique.MustUUID().String(),
		RealName: unique.MustUUID().String(),
		Status:   1,
		Password: util.MD5HashString("test"),
		UserRoles: schema.UserRoles{
			&schema.UserRole{
				RoleID: addRoleItemRes.ID,
			},
		},
	}
	engine.ServeHTTP(w, newPostRequest(router, addItem))
	assert.Equal(t, 200, w.Code)
	var addItemRes ResID
	err = parseReader(w.Body, &addItemRes)
	assert.Nil(t, err)

	// get /users/:id
	engine.ServeHTTP(w, newGetRequest("%s/%s", nil, router, addItemRes.ID))
	assert.Equal(t, 200, w.Code)
	var getItem schema.User
	err = parseReader(w.Body, &getItem)
	assert.Nil(t, err)
	assert.Equal(t, addItem.UserName, getItem.UserName)
	assert.Equal(t, addItem.Status, getItem.Status)
	assert.NotEmpty(t, getItem.ID)

	// put /users/:id
	putItem := getItem
	putItem.UserName = unique.MustUUID().String()
	engine.ServeHTTP(w, newPutRequest("%s/%s", putItem, router, getItem.ID))
	assert.Equal(t, 200, w.Code)
	err = parseOK(w.Body)
	assert.Nil(t, err)

	// query /users
	engine.ServeHTTP(w, newGetRequest(router, newPageParam()))
	assert.Equal(t, 200, w.Code)
	var pageItems []*schema.User
	err = parsePageReader(w.Body, &pageItems)
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(pageItems), 1)
	if len(pageItems) > 0 {
		assert.Equal(t, putItem.ID, pageItems[0].ID)
		assert.Equal(t, putItem.UserName, pageItems[0].UserName)
	}

	// delete /users/:id
	engine.ServeHTTP(w, newDeleteRequest("%s/%s", router, addItemRes.ID))
	assert.Equal(t, 200, w.Code)
	err = parseOK(w.Body)
	assert.Nil(t, err)

	// delete /roles/:id
	engine.ServeHTTP(w, newDeleteRequest(apiPrefix+"v1/roles/%s", addRoleItemRes.ID))
	assert.Equal(t, 200, w.Code)
	err = parseOK(w.Body)
	assert.Nil(t, err)

	// delete /menus/:id
	engine.ServeHTTP(w, newDeleteRequest(apiPrefix+"v1/menus/%s", addMenuItemRes.ID))
	assert.Equal(t, 200, w.Code)
	err = parseOK(w.Body)
	assert.Nil(t, err)
}
