package test

import (
	"net/http/httptest"
	"testing"

	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/stretchr/testify/assert"
)

func TestRole(t *testing.T) {
	const router = "roles"
	var err error

	w := httptest.NewRecorder()

	// post /menus
	addMenuItem := &schema.Menu{
		Code:     "test_role_menu_1",
		Name:     "测试角色菜单",
		Type:     10,
		Sequence: 1,
		Icon:     "test",
		Path:     "/test",
		Status:   1,
		IsHide:   2,
	}
	engine.ServeHTTP(w, newPostRequest("menus", addMenuItem))
	assert.Equal(t, 200, w.Code)

	var addMenuNewItem schema.Menu
	err = parseReader(w.Body, &addMenuNewItem)

	// post /roles
	addItem := &schema.Role{
		Name:    "测试角色",
		Memo:    "角色备注",
		Status:  1,
		MenuIDs: []string{addMenuNewItem.RecordID},
	}
	engine.ServeHTTP(w, newPostRequest(router, addItem))
	assert.Equal(t, 200, w.Code)

	var addNewItem schema.Role
	err = parseReader(w.Body, &addNewItem)
	assert.Nil(t, err)
	assert.Equal(t, addItem.Name, addNewItem.Name)
	assert.Equal(t, addItem.Memo, addNewItem.Memo)
	assert.Equal(t, addItem.Status, addNewItem.Status)
	assert.Equal(t, addItem.MenuIDs, addNewItem.MenuIDs)
	assert.NotEqual(t, addNewItem.ID, 0)
	assert.NotEmpty(t, addNewItem.RecordID)

	// get /roles?type=page
	engine.ServeHTTP(w, newGetRequest(router,
		newPageParam(map[string]string{"type": "page"})))
	assert.Equal(t, 200, w.Code)
	var pageItems []*schema.Role
	err = parsePageReader(w.Body, &pageItems)
	assert.Nil(t, err)
	assert.Equal(t, len(pageItems), 1)
	assert.Equal(t, pageItems[0].RecordID, addNewItem.RecordID)

	// put /roles/:id
	putItem := *pageItems[0]
	putItem.Name = "测试角色2"
	putItem.Memo = "角色备注2"
	putItem.MenuIDs = []string{addMenuNewItem.RecordID}
	engine.ServeHTTP(w, newPutRequest("%s/%s", putItem, router, addNewItem.RecordID))
	assert.Equal(t, 200, w.Code)
	releaseReader(w.Body)

	// patch /roles/:id/disable
	engine.ServeHTTP(w, newPatchRequest("%s/%s/disable", router, addNewItem.RecordID))
	assert.Equal(t, 200, w.Code)
	releaseReader(w.Body)

	// get /roles/:id
	engine.ServeHTTP(w, newGetRequest("%s/%s", nil, router, addNewItem.RecordID))
	assert.Equal(t, 200, w.Code)

	var getItem schema.Role
	err = parseReader(w.Body, &getItem)

	assert.Nil(t, err)
	assert.Equal(t, getItem.RecordID, addNewItem.RecordID)
	assert.Equal(t, getItem.Name, putItem.Name)
	assert.Equal(t, getItem.MenuIDs, putItem.MenuIDs)
	assert.Equal(t, getItem.Status, 2)

	// delete /roles/:id
	engine.ServeHTTP(w, newDeleteRequest("%s/%s", router, addNewItem.RecordID))
	assert.Equal(t, 200, w.Code)
	releaseReader(w.Body)

	// delete /menus/:id
	engine.ServeHTTP(w, newDeleteRequest("%s/%s", "menus", addMenuNewItem.RecordID))
	assert.Equal(t, 200, w.Code)
	releaseReader(w.Body)
}
