package test

import (
	"net/http/httptest"
	"testing"

	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/LyricTian/gin-admin/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestAPIUser(t *testing.T) {
	const router = apiPrefix + "v1/users"
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
	addRoleItem := &schema.Role{
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
	engine.ServeHTTP(w, newPostRequest(apiPrefix+"v1/roles", addRoleItem))
	assert.Equal(t, 200, w.Code)
	var addNewRoleItem schema.Role
	err = parseReader(w.Body, &addNewRoleItem)
	assert.Nil(t, err)

	// post /users
	addItem := &schema.User{
		UserName: util.MustUUID(),
		RealName: util.MustUUID(),
		Password: util.MD5HashString("123"),
		Status:   1,
		Roles: []*schema.UserRole{
			{RoleID: addNewRoleItem.RecordID},
		},
	}
	engine.ServeHTTP(w, newPostRequest(router, addItem))
	assert.Equal(t, 200, w.Code)
	var addNewItem schema.User
	err = parseReader(w.Body, &addNewItem)
	assert.Nil(t, err)
	assert.Equal(t, addItem.UserName, addNewItem.UserName)
	assert.Equal(t, addItem.RealName, addNewItem.RealName)
	assert.Equal(t, addItem.Status, addNewItem.Status)
	assert.Equal(t, len(addItem.Roles), len(addNewItem.Roles))
	assert.Empty(t, addNewItem.Password)
	assert.NotEmpty(t, addNewItem.RecordID)

	// query /users
	engine.ServeHTTP(w, newGetRequest(router, newPageParam()))
	assert.Equal(t, 200, w.Code)
	var pageItems []*schema.User
	err = parsePageReader(w.Body, &pageItems)
	assert.Nil(t, err)
	assert.Equal(t, len(pageItems), 1)
	if len(pageItems) > 0 {
		assert.Equal(t, addNewItem.RecordID, pageItems[0].RecordID)
		assert.Equal(t, addNewItem.UserName, pageItems[0].UserName)
		assert.Equal(t, addNewItem.RealName, pageItems[0].RealName)
		assert.Equal(t, addNewItem.Status, pageItems[0].Status)
	}

	// put /users/:id
	engine.ServeHTTP(w, newGetRequest("%s/%s", nil, router, addNewItem.RecordID))
	assert.Equal(t, 200, w.Code)
	var putItem schema.User
	err = parseReader(w.Body, &putItem)
	assert.Nil(t, err)

	putItem.UserName = util.MustUUID()
	engine.ServeHTTP(w, newPutRequest("%s/%s", putItem, router, addNewItem.RecordID))
	assert.Equal(t, 200, w.Code)
	var putNewItem schema.User
	err = parseReader(w.Body, &putNewItem)
	assert.Nil(t, err)
	assert.Equal(t, putItem.UserName, putNewItem.UserName)
	assert.Equal(t, putItem.RealName, putNewItem.RealName)
	assert.Equal(t, putItem.Status, putNewItem.Status)
	assert.Equal(t, len(putItem.Roles), len(putNewItem.Roles))

	// delete /menus/:id
	engine.ServeHTTP(w, newDeleteRequest("%s/%s", apiPrefix+"v1/menus", addNewMenuItem.RecordID))
	assert.Equal(t, 200, w.Code)
	err = parseOK(w.Body)
	assert.Nil(t, err)

	// delete /users/:id
	engine.ServeHTTP(w, newDeleteRequest("%s/%s", router, addNewItem.RecordID))
	assert.Equal(t, 200, w.Code)
	err = parseOK(w.Body)
	assert.Nil(t, err)

	// delete /roles/:id
	engine.ServeHTTP(w, newDeleteRequest("%s/%s", apiPrefix+"v1/roles", addNewRoleItem.RecordID))
	assert.Equal(t, 200, w.Code)
	err = parseOK(w.Body)
	assert.Nil(t, err)
}

func BenchmarkAPIUserCreateParallel(b *testing.B) {
	const router = apiPrefix + "v1/users"
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

	// post /roles
	addRoleItem := &schema.Role{
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
	engine.ServeHTTP(w, newPostRequest(apiPrefix+"v1/roles", addRoleItem))
	var addNewRoleItem schema.Role
	_ = parseReader(w.Body, &addNewRoleItem)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// post /users
			addItem := &schema.User{
				UserName: util.MustUUID(),
				RealName: util.MustUUID(),
				Password: util.MD5HashString("123"),
				Status:   1,
				Roles: []*schema.UserRole{
					{RoleID: addNewRoleItem.RecordID},
				},
			}

			w := httptest.NewRecorder()
			engine.ServeHTTP(w, newPostRequest(router, addItem))
		}
	})
}
