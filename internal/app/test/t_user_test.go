package test

import (
	"net/http/httptest"
	"testing"

	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/LyricTian/gin-admin/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	const router = apiPrefix + "v1/users"
	var err error

	w := httptest.NewRecorder()

	// post /menus
	addMenuItem := &schema.Menu{
		Name:       util.MustUUID(),
		ShowStatus: 1,
		Status:     1,
	}
	engine.ServeHTTP(w, newPostRequest(apiPrefix+"v1/menus", addMenuItem))
	assert.Equal(t, 200, w.Code)
	var addMenuItemRes ResRecordID
	err = parseReader(w.Body, &addMenuItemRes)
	assert.Nil(t, err)

	// post /roles
	addRoleItem := &schema.Role{
		Name:   util.MustUUID(),
		Status: 1,
		RoleMenus: schema.RoleMenus{
			&schema.RoleMenu{
				MenuID: addMenuItemRes.RecordID,
			},
		},
	}
	engine.ServeHTTP(w, newPostRequest(apiPrefix+"v1/roles", addRoleItem))
	assert.Equal(t, 200, w.Code)
	var addRoleItemRes ResRecordID
	err = parseReader(w.Body, &addRoleItemRes)
	assert.Nil(t, err)

	// post /users
	addItem := &schema.User{
		UserName: util.MustUUID(),
		RealName: util.MustUUID(),
		Status:   1,
		Password: util.MD5HashString("test"),
		UserRoles: schema.UserRoles{
			&schema.UserRole{
				RoleID: addRoleItemRes.RecordID,
			},
		},
	}
	engine.ServeHTTP(w, newPostRequest(router, addItem))
	assert.Equal(t, 200, w.Code)
	var addItemRes ResRecordID
	err = parseReader(w.Body, &addItemRes)
	assert.Nil(t, err)

	// get /users/:id
	engine.ServeHTTP(w, newGetRequest("%s/%s", nil, router, addItemRes.RecordID))
	assert.Equal(t, 200, w.Code)
	var getItem schema.User
	err = parseReader(w.Body, &getItem)
	assert.Nil(t, err)
	assert.Equal(t, addItem.UserName, getItem.UserName)
	assert.Equal(t, addItem.Status, getItem.Status)
	assert.NotEmpty(t, getItem.RecordID)

	// put /users/:id
	putItem := getItem
	putItem.UserName = util.MustUUID()
	engine.ServeHTTP(w, newPutRequest("%s/%s", putItem, router, getItem.RecordID))
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
		assert.Equal(t, putItem.RecordID, pageItems[0].RecordID)
		assert.Equal(t, putItem.UserName, pageItems[0].UserName)
	}

	// delete /users/:id
	engine.ServeHTTP(w, newDeleteRequest("%s/%s", router, addItemRes.RecordID))
	assert.Equal(t, 200, w.Code)
	err = parseOK(w.Body)
	assert.Nil(t, err)

	// delete /roles/:id
	engine.ServeHTTP(w, newDeleteRequest(apiPrefix+"v1/roles/%s", addRoleItemRes.RecordID))
	assert.Equal(t, 200, w.Code)
	err = parseOK(w.Body)
	assert.Nil(t, err)

	// delete /menus/:id
	engine.ServeHTTP(w, newDeleteRequest(apiPrefix+"v1/menus/%s", addMenuItemRes.RecordID))
	assert.Equal(t, 200, w.Code)
	err = parseOK(w.Body)
	assert.Nil(t, err)
}
