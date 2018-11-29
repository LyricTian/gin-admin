package test

import (
	"net/http/httptest"
	"testing"

	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	const router = "users"
	var err error

	w := httptest.NewRecorder()

	// post /users
	addItem := &schema.User{
		UserName: "test_user_1",
		RealName: "测试用户1",
		Password: "123456",
		Status:   1,
		RoleIDs:  []string{"test_role_1"},
	}
	engine.ServeHTTP(w, newPostRequest(router, addItem))
	assert.Equal(t, 200, w.Code)

	var addNewItem schema.User
	err = parseReader(w.Body, &addNewItem)
	assert.Nil(t, err)
	assert.Equal(t, addItem.UserName, addNewItem.UserName)
	assert.Equal(t, addItem.RealName, addNewItem.RealName)
	assert.Equal(t, addItem.RoleIDs, addNewItem.RoleIDs)
	assert.Empty(t, addNewItem.Password)
	assert.NotEqual(t, addNewItem.ID, 0)
	assert.NotEmpty(t, addNewItem.RecordID)

	// get /users?type=page
	engine.ServeHTTP(w, newGetRequest(router,
		newPageParam(map[string]string{"type": "page"})))
	assert.Equal(t, 200, w.Code)
	var pageItems []*schema.User
	err = parsePageReader(w.Body, &pageItems)
	assert.Nil(t, err)
	assert.Equal(t, len(pageItems), 1)
	assert.Equal(t, pageItems[0].RecordID, addNewItem.RecordID)

	// put /users/:id
	putItem := *pageItems[0]
	putItem.UserName = "test_user_2"
	putItem.RealName = "测试用户2"
	putItem.RoleIDs = []string{"test_role_2"}
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
	assert.Equal(t, getItem.RoleIDs, putItem.RoleIDs)

	// delete /users/:id
	engine.ServeHTTP(w, newDeleteRequest("%s/%s", router, addNewItem.RecordID))
	assert.Equal(t, 200, w.Code)
	releaseReader(w.Body)
}
