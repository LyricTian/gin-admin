package test

import (
	"net/http/httptest"
	"testing"

	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/LyricTian/gin-admin/src/util"
	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	const router = "v1/users"
	var err error

	w := httptest.NewRecorder()

	// post /menus
	addMenuItem := &schema.Menu{
		Code:     util.MustUUID(),
		Name:     "测试角色菜单",
		Type:     1,
		Sequence: 9999,
		Icon:     "test",
		Path:     "/test",
		IsHide:   2,
	}
	engine.ServeHTTP(w, newPostRequest("v1/menus", addMenuItem))
	assert.Equal(t, 200, w.Code)

	var addMenuNewItem schema.Menu
	err = parseReader(w.Body, &addMenuNewItem)

	// post /roles
	addRoleItem := &schema.Role{
		Name:    util.MustUUID(),
		Memo:    "角色备注",
		Status:  1,
		MenuIDs: []string{addMenuNewItem.RecordID},
	}
	engine.ServeHTTP(w, newPostRequest("v1/roles", addRoleItem))
	assert.Equal(t, 200, w.Code)

	var addNewRoleItem schema.Role
	err = parseReader(w.Body, &addNewRoleItem)
	assert.Nil(t, err)

	// post /users
	addItem := &schema.User{
		UserName: util.MustUUID(),
		RealName: "测试用户1",
		Password: util.MD5HashString("123456"),
		Status:   1,
		RoleIDs:  []string{addNewRoleItem.RecordID},
	}
	engine.ServeHTTP(w, newPostRequest(router, addItem))
	assert.Equal(t, 200, w.Code)

	var addNewItem schema.User
	err = parseReader(w.Body, &addNewItem)
	assert.Nil(t, err)

	// get /users/:id
	engine.ServeHTTP(w, newGetRequest("%s/%s", nil, router, addNewItem.RecordID))
	assert.Equal(t, 200, w.Code)

	var addGetItem schema.User
	err = parseReader(w.Body, &addGetItem)
	assert.Equal(t, addItem.UserName, addGetItem.UserName)
	assert.Equal(t, addItem.RealName, addGetItem.RealName)
	assert.Equal(t, addItem.Status, addGetItem.Status)
	assert.Equal(t, addItem.RoleIDs, addGetItem.RoleIDs)
	assert.NotEmpty(t, addGetItem.RecordID)
	assert.Empty(t, addGetItem.Password)

	// get /users?q=page
	engine.ServeHTTP(w, newGetRequest(router,
		newPageParam(map[string]string{"q": "page"})))
	assert.Equal(t, 200, w.Code)
	var pageItems []*schema.User
	err = parsePageReader(w.Body, &pageItems)
	assert.Nil(t, err)
	assert.Equal(t, len(pageItems), 1)
	assert.Equal(t, pageItems[0].RecordID, addNewItem.RecordID)

	// get /users/:id
	engine.ServeHTTP(w, newGetRequest("%s/%s", nil, router, addNewItem.RecordID))
	assert.Equal(t, 200, w.Code)
	var putItem schema.User
	err = parseReader(w.Body, &putItem)

	// put /users/:id
	putItem.UserName = util.MustUUID()
	putItem.RealName = "测试用户2"
	engine.ServeHTTP(w, newPutRequest("%s/%s", putItem, router, addNewItem.RecordID))
	assert.Equal(t, 200, w.Code)
	releaseReader(w.Body)

	// get /users/:id
	engine.ServeHTTP(w, newGetRequest("%s/%s", nil, router, addNewItem.RecordID))
	assert.Equal(t, 200, w.Code)
	var getItem schema.User
	err = parseReader(w.Body, &getItem)

	assert.Nil(t, err)
	assert.Equal(t, getItem.RecordID, addNewItem.RecordID)
	assert.Equal(t, getItem.UserName, putItem.UserName)
	assert.Equal(t, getItem.RealName, putItem.RealName)

	// delete /users/:id
	engine.ServeHTTP(w, newDeleteRequest("%s/%s", router, addNewItem.RecordID))
	assert.Equal(t, 200, w.Code)
	releaseReader(w.Body)

	// delete /roles/:id
	engine.ServeHTTP(w, newDeleteRequest("%s/%s", "v1/roles", addNewRoleItem.RecordID))
	assert.Equal(t, 200, w.Code)
	releaseReader(w.Body)

	// delete /menus/:id
	engine.ServeHTTP(w, newDeleteRequest("%s/%s", "v1/menus", addMenuNewItem.RecordID))
	assert.Equal(t, 200, w.Code)
	releaseReader(w.Body)
}
