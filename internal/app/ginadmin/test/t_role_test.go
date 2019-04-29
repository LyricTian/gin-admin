package test

import (
	"net/http/httptest"
	"testing"

	"github.com/LyricTian/gin-admin/internal/app/ginadmin/schema"
	"github.com/LyricTian/gin-admin/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestRole(t *testing.T) {
	const router = "v1/roles"
	var err error

	w := httptest.NewRecorder()

	// post /menus
	addMenuItem := &schema.Menu{
		Name:     util.MustUUID(),
		Sequence: 9999999,
		Actions: []*schema.MenuAction{
			{Code: "query", Name: "query"},
		},
		Resources: []*schema.MenuResource{
			{Code: "query", Name: "query", Method: "GET", Path: "/test/v1/menus"},
		},
	}
	engine.ServeHTTP(w, newPostRequest("v1/menus", addMenuItem))
	assert.Equal(t, 200, w.Code)
	var addNewMenuItem schema.Menu
	err = parseReader(w.Body, &addNewMenuItem)
	assert.Nil(t, err)

	// post /roles
	addItem := &schema.Role{
		Name:     util.MustUUID(),
		Sequence: 9999999,
		Menus: []*schema.RoleMenu{
			{
				MenuID:    addNewMenuItem.RecordID,
				Actions:   []string{"query"},
				Resources: []string{"query"},
			},
		},
	}
	engine.ServeHTTP(w, newPostRequest(router, addItem))
	assert.Equal(t, 200, w.Code)
	var addNewItem schema.Role
	err = parseReader(w.Body, &addNewItem)
	assert.Nil(t, err)
	assert.Equal(t, addItem.Name, addNewItem.Name)
	assert.Equal(t, addItem.Sequence, addNewItem.Sequence)
	assert.Equal(t, len(addItem.Menus), len(addNewItem.Menus))
	assert.NotEmpty(t, addNewItem.RecordID)

	// query /roles?q=page
	engine.ServeHTTP(w, newGetRequest(router,
		newPageParam(map[string]string{"q": "page"})))
	assert.Equal(t, 200, w.Code)
	var pageItems []*schema.Role
	err = parsePageReader(w.Body, &pageItems)
	assert.Nil(t, err)
	assert.Equal(t, len(pageItems), 1)
	if len(pageItems) > 0 {
		assert.Equal(t, addNewItem.RecordID, pageItems[0].RecordID)
		assert.Equal(t, addNewItem.Name, pageItems[0].Name)
	}

	// put /roles/:id
	engine.ServeHTTP(w, newGetRequest("%s/%s", nil, router, addNewItem.RecordID))
	assert.Equal(t, 200, w.Code)

	var putItem schema.Role
	err = parseReader(w.Body, &putItem)
	assert.Equal(t, len(putItem.Menus), 1)

	putItem.Name = util.MustUUID()
	putItem.Menus = []*schema.RoleMenu{
		{
			MenuID:    addNewMenuItem.RecordID,
			Actions:   []string{},
			Resources: []string{"query"},
		},
	}

	engine.ServeHTTP(w, newPutRequest("%s/%s", putItem, router, addNewItem.RecordID))
	assert.Equal(t, 200, w.Code)
	var putNewItem schema.Role
	err = parseReader(w.Body, &putNewItem)
	assert.Nil(t, err)
	assert.Equal(t, putItem.Name, putNewItem.Name)
	assert.Equal(t, len(putItem.Menus), len(putNewItem.Menus))
	assert.Equal(t, 0, len(putNewItem.Menus[0].Actions))
	assert.Equal(t, 1, len(putNewItem.Menus[0].Resources))

	// delete /menus/:id
	engine.ServeHTTP(w, newDeleteRequest("%s/%s", "v1/menus", addNewMenuItem.RecordID))
	assert.Equal(t, 200, w.Code)
	err = parseOK(w.Body)
	assert.Nil(t, err)

	// delete /roles/:id
	engine.ServeHTTP(w, newDeleteRequest("%s/%s", router, addNewItem.RecordID))
	assert.Equal(t, 200, w.Code)
	err = parseOK(w.Body)
	assert.Nil(t, err)
}
