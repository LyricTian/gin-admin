package test

import (
	"net/http/httptest"
	"testing"

	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/LyricTian/gin-admin/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestAPIDemo(t *testing.T) {
	const router = apiPrefix + "v1/demos"
	var err error

	w := httptest.NewRecorder()

	// post /demos
	addItem := &schema.Demo{
		Code:   util.MustUUID(),
		Name:   util.MustUUID(),
		Status: 1,
	}
	engine.ServeHTTP(w, newPostRequest(router, addItem))
	assert.Equal(t, 200, w.Code)

	var addNewItem schema.Demo
	err = parseReader(w.Body, &addNewItem)
	assert.Nil(t, err)
	assert.Equal(t, addItem.Code, addNewItem.Code)
	assert.Equal(t, addItem.Code, addNewItem.Code)
	assert.Equal(t, addItem.Status, addNewItem.Status)
	assert.NotEmpty(t, addNewItem.RecordID)

	// query /demos
	engine.ServeHTTP(w, newGetRequest(router, newPageParam()))
	assert.Equal(t, 200, w.Code)
	var pageItems []*schema.Demo
	err = parsePageReader(w.Body, &pageItems)
	assert.Nil(t, err)
	assert.Equal(t, len(pageItems), 1)
	if len(pageItems) > 0 {
		assert.Equal(t, addNewItem.RecordID, pageItems[0].RecordID)
		assert.Equal(t, addNewItem.Name, pageItems[0].Name)
	}

	// put /demos/:id
	engine.ServeHTTP(w, newGetRequest("%s/%s", nil, router, addNewItem.RecordID))
	assert.Equal(t, 200, w.Code)
	var putItem schema.Demo
	err = parseReader(w.Body, &putItem)
	assert.Nil(t, err)

	putItem.Name = util.MustUUID()
	engine.ServeHTTP(w, newPutRequest("%s/%s", putItem, router, addNewItem.RecordID))
	assert.Equal(t, 200, w.Code)

	var putNewItem schema.Demo
	err = parseReader(w.Body, &putNewItem)
	assert.Nil(t, err)
	assert.Equal(t, putItem.Name, putNewItem.Name)

	// delete /demos/:id
	engine.ServeHTTP(w, newDeleteRequest("%s/%s", router, addNewItem.RecordID))
	assert.Equal(t, 200, w.Code)
	err = parseOK(w.Body)
	assert.Nil(t, err)
}
