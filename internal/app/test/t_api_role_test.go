package test

import (
	"net/http/httptest"
	"testing"

	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/LyricTian/gin-admin/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestAPIRole(t *testing.T) {
	const router = apiPrefix + "v1/roles"
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
	engine.ServeHTTP(w, newPostRequest(apiPrefix+"v1/menus", addMenuItem))
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

	// query /roles
	engine.ServeHTTP(w, newGetRequest(router, newPageParam()))
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
	assert.Nil(t, err)
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
	engine.ServeHTTP(w, newDeleteRequest("%s/%s", apiPrefix+"v1/menus", addNewMenuItem.RecordID))
	assert.Equal(t, 200, w.Code)
	err = parseOK(w.Body)
	assert.Nil(t, err)

	// delete /roles/:id
	engine.ServeHTTP(w, newDeleteRequest("%s/%s", router, addNewItem.RecordID))
	assert.Equal(t, 200, w.Code)
	err = parseOK(w.Body)
	assert.Nil(t, err)
}

func BenchmarkAPIRoleCreateParallel(b *testing.B) {
	const router = apiPrefix + "v1/roles"

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
	engine.ServeHTTP(w, newPostRequest(apiPrefix+"v1/menus", addMenuItem))
	var addNewMenuItem schema.Menu
	_ = parseReader(w.Body, &addNewMenuItem)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
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
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, newPostRequest(router, addItem))
			if w.Code != 200 {
				b.Errorf("Expected value: %d, given value: %d", 200, w.Code)
			}
		}
	})
}
