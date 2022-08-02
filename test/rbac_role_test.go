package test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/LyricTian/gin-admin/v9/internal/module/rbac/typed"
	"github.com/stretchr/testify/assert"
)

func TestRBACRole(t *testing.T) {
	assert := assert.New(t)
	addRole := typed.RoleCreate{
		Name:     "test",
		Sequence: 1,
	}

	addRes := httptest.NewRecorder()
	engine.ServeHTTP(addRes, newPostRequest(addRole, "/api/rbac/v1/roles"))
	assert.Equal(200, addRes.Code)

	var addedRole typed.Role
	err := json.Unmarshal(addRes.Body.Bytes(), &addedRole)
	assert.Nil(err)
	assert.NotEmpty(addedRole.ID)

	getRes := httptest.NewRecorder()
	engine.ServeHTTP(getRes, newGetRequest(nil, "/api/rbac/v1/roles/%s", addedRole.ID))
	assert.Equal(200, getRes.Code)
	var getRole typed.Role
	err = json.Unmarshal(getRes.Body.Bytes(), &getRole)
	assert.Nil(err)
	assert.Equal(addedRole.ID, getRole.ID)
	assert.Equal(addRole.Name, getRole.Name)
	assert.Equal(addRole.Sequence, getRole.Sequence)

	addRole.Name = "test2"
	updateRes := httptest.NewRecorder()
	engine.ServeHTTP(updateRes, newPutRequest(addRole, "/api/rbac/v1/roles/%s", addedRole.ID))
	assert.Equal(200, updateRes.Code)

	getRes = httptest.NewRecorder()
	engine.ServeHTTP(getRes, newGetRequest(nil, "/api/rbac/v1/roles/%s", addedRole.ID))
	assert.Equal(200, getRes.Code)
	var updatedRole typed.Role
	err = json.Unmarshal(getRes.Body.Bytes(), &updatedRole)
	assert.Nil(err)
	assert.Equal(addedRole.ID, updatedRole.ID)
	assert.Equal(addRole.Name, updatedRole.Name)
	assert.NotEqual(getRole.Name, updatedRole.Name)

	deleteRes := httptest.NewRecorder()
	engine.ServeHTTP(deleteRes, newDeleteRequest("/api/rbac/v1/roles/%s", addedRole.ID))
	assert.Equal(200, deleteRes.Code)
}
