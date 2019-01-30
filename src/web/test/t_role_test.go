package test

import (
	"net/http/httptest"
	"testing"

	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/LyricTian/gin-admin/src/util"
	"github.com/stretchr/testify/assert"
)

func TestRole(t *testing.T) {
	const router = "v1/roles"
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
	addItem := &schema.Role{
		Name:     util.MustUUID(),
		Sequence: 9999,
		Memo:     "角色备注",
		MenuIDs:  []string{addMenuNewItem.RecordID},
	}
	engine.ServeHTTP(w, newPostRequest(router, addItem))
	assert.Equal(t, 200, w.Code)

	var addNewItem schema.Role
	err = parseReader(w.Body, &addNewItem)
	assert.Nil(t, err)

	// get /roles/:id
	engine.ServeHTTP(w, newGetRequest("%s/%s", nil, router, addNewItem.RecordID))
	assert.Equal(t, 200, w.Code)

	var addGetItem schema.Role
	err = parseReader(w.Body, &addGetItem)
	assert.Equal(t, addItem.Name, addGetItem.Name)
	assert.Equal(t, addItem.Memo, addGetItem.Memo)
	assert.Equal(t, addItem.MenuIDs, addGetItem.MenuIDs)
	assert.NotEmpty(t, addGetItem.RecordID)

	// get /roles?q=page
	engine.ServeHTTP(w, newGetRequest(router,
		newPageParam(map[string]string{"q": "page"})))
	assert.Equal(t, 200, w.Code)
	var pageItems []*schema.Role
	err = parsePageReader(w.Body, &pageItems)
	assert.Nil(t, err)
	assert.Equal(t, len(pageItems), 1)
	assert.Equal(t, pageItems[0].RecordID, addNewItem.RecordID)

	// get /roles/:id
	engine.ServeHTTP(w, newGetRequest("%s/%s", nil, router, addNewItem.RecordID))
	assert.Equal(t, 200, w.Code)
	var putItem schema.Role
	err = parseReader(w.Body, &putItem)
	// put /roles/:id
	putItem.Name = util.MustUUID()
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

	// delete /roles/:id
	engine.ServeHTTP(w, newDeleteRequest("%s/%s", router, addNewItem.RecordID))
	assert.Equal(t, 200, w.Code)
	releaseReader(w.Body)

	// delete /menus/:id
	engine.ServeHTTP(w, newDeleteRequest("%s/%s", "v1/menus", addMenuNewItem.RecordID))
	assert.Equal(t, 200, w.Code)
	releaseReader(w.Body)
}
