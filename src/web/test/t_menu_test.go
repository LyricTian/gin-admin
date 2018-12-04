package test

import (
	"net/http/httptest"
	"testing"

	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/stretchr/testify/assert"
)

func TestMenu(t *testing.T) {
	const router = "menus"
	var err error

	w := httptest.NewRecorder()

	// post /menus
	addItem := &schema.Menu{
		Code:     "test_menu_1",
		Name:     "测试菜单",
		Type:     10,
		Sequence: 1,
		Icon:     "test",
		Path:     "/test",
		Status:   1,
		IsHide:   2,
	}
	engine.ServeHTTP(w, newPostRequest(router, addItem))
	assert.Equal(t, 200, w.Code)

	var addNewItem schema.Menu
	err = parseReader(w.Body, &addNewItem)
	assert.Nil(t, err)
	assert.Equal(t, addItem.Code, addNewItem.Code)
	assert.Equal(t, addItem.Name, addNewItem.Name)
	assert.Equal(t, addItem.Type, addNewItem.Type)
	assert.Equal(t, addItem.Sequence, addNewItem.Sequence)
	assert.Equal(t, addItem.Icon, addNewItem.Icon)
	assert.Equal(t, addItem.Path, addNewItem.Path)
	assert.Equal(t, addItem.Status, addNewItem.Status)
	assert.NotEqual(t, addNewItem.ID, 0)
	assert.NotEmpty(t, addNewItem.RecordID)

	// get /menus?type=page
	engine.ServeHTTP(w, newGetRequest(router,
		newPageParam(map[string]string{"type": "page"})))
	assert.Equal(t, 200, w.Code)
	var pageItems []*schema.Menu
	err = parsePageReader(w.Body, &pageItems)
	assert.Nil(t, err)
	assert.Equal(t, len(pageItems), 1)
	assert.Equal(t, pageItems[0].RecordID, addNewItem.RecordID)

	// put /menus/:id
	putItem := *pageItems[0]
	putItem.Code = "test_menu_2"
	putItem.Name = "测试菜单2"
	engine.ServeHTTP(w, newPutRequest("%s/%s", putItem, router, addNewItem.RecordID))
	assert.Equal(t, 200, w.Code)
	releaseReader(w.Body)

	// patch /menus/:id/disable
	engine.ServeHTTP(w, newPatchRequest("%s/%s/disable", router, addNewItem.RecordID))
	assert.Equal(t, 200, w.Code)
	releaseReader(w.Body)

	// get /menus/:id
	engine.ServeHTTP(w, newGetRequest("%s/%s", nil, router, addNewItem.RecordID))
	assert.Equal(t, 200, w.Code)

	var getItem schema.Menu
	err = parseReader(w.Body, &getItem)

	assert.Nil(t, err)
	assert.Equal(t, getItem.RecordID, addNewItem.RecordID)
	assert.Equal(t, getItem.Code, putItem.Code)
	assert.Equal(t, getItem.Name, putItem.Name)
	assert.Equal(t, getItem.Status, 2)

	// delete /menus/:id
	engine.ServeHTTP(w, newDeleteRequest("%s/%s", router, addNewItem.RecordID))
	assert.Equal(t, 200, w.Code)
	releaseReader(w.Body)
}
