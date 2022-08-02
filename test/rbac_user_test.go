package test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/LyricTian/gin-admin/v9/internal/module/rbac/typed"
	"github.com/LyricTian/gin-admin/v9/pkg/util/hash"
	"github.com/stretchr/testify/assert"
)

func TestRBACUser(t *testing.T) {
	assert := assert.New(t)
	addUser := typed.UserCreate{
		Username: "test",
		Name:     "TestA",
		Password: hash.MD5String("123456"),
	}

	addRes := httptest.NewRecorder()
	engine.ServeHTTP(addRes, newPostRequest(addUser, "/api/rbac/v1/users"))
	assert.Equal(200, addRes.Code)

	var addedUser typed.User
	err := json.Unmarshal(addRes.Body.Bytes(), &addedUser)
	assert.Nil(err)
	assert.NotEmpty(addedUser.ID)
	assert.Equal(addedUser.Status, typed.UserStatusActivated)

	getRes := httptest.NewRecorder()
	engine.ServeHTTP(getRes, newGetRequest(nil, "/api/rbac/v1/users/%s", addedUser.ID))
	assert.Equal(200, getRes.Code)
	var getUser typed.User
	err = json.Unmarshal(getRes.Body.Bytes(), &getUser)
	assert.Nil(err)
	assert.Equal(addedUser.ID, getUser.ID)
	assert.Equal(addUser.Name, getUser.Name)
	assert.Equal(addUser.Username, getUser.Username)

	addUser.Username = "test2"
	updateRes := httptest.NewRecorder()
	engine.ServeHTTP(updateRes, newPutRequest(addUser, "/api/rbac/v1/users/%s", addedUser.ID))
	assert.Equal(200, updateRes.Code)

	getRes = httptest.NewRecorder()
	engine.ServeHTTP(getRes, newGetRequest(nil, "/api/rbac/v1/users/%s", addedUser.ID))
	assert.Equal(200, getRes.Code)
	var updatedUser typed.User
	err = json.Unmarshal(getRes.Body.Bytes(), &updatedUser)
	assert.Nil(err)
	assert.Equal(addedUser.ID, updatedUser.ID)
	assert.Equal(addUser.Username, updatedUser.Username)
	assert.NotEqual(getUser.Username, updatedUser.Username)

	deleteRes := httptest.NewRecorder()
	engine.ServeHTTP(deleteRes, newDeleteRequest("/api/rbac/v1/users/%s", addedUser.ID))
	assert.Equal(200, deleteRes.Code)
}
