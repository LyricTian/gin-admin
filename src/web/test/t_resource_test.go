package test

import (
	"net/http/httptest"
	"testing"

	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/LyricTian/gin-admin/src/util"
	"github.com/stretchr/testify/assert"
)

func TestResource(t *testing.T) {
	const router = "v1/resources"
	var err error

	w := httptest.NewRecorder()

	// post /resources
	addItem := &schema.Resource{
		Path:   util.MustUUID(),
		Method: "GET",
		Name:   "测试资源",
		Code:   "test",
	}
	engine.ServeHTTP(w, newPostRequest(router, addItem))
	assert.Equal(t, 200, w.Code)

	var addNewItem schema.Resource
	err = parseReader(w.Body, &addNewItem)
	assert.Nil(t, err)
	assert.NotEmpty(t, addNewItem.RecordID)

	// get /resources/:id
	engine.ServeHTTP(w, newGetRequest("%s/%s", nil, router, addNewItem.RecordID))
	assert.Equal(t, 200, w.Code)

	var addGetItem schema.Resource
	err = parseReader(w.Body, &addGetItem)

	assert.Nil(t, err)
	assert.Equal(t, addGetItem.Code, addItem.Code)
	assert.Equal(t, addGetItem.Name, addItem.Name)
	assert.Equal(t, addGetItem.Path, addItem.Path)
	assert.Equal(t, addGetItem.Method, addItem.Method)

	// get /resources?q=page
	engine.ServeHTTP(w, newGetRequest(router,
		newPageParam(map[string]string{"q": "page"})))
	assert.Equal(t, 200, w.Code)
	var pageItems []*schema.Resource
	err = parsePageReader(w.Body, &pageItems)
	assert.Nil(t, err)
	assert.Equal(t, len(pageItems), 1)
	assert.Equal(t, pageItems[0].RecordID, addNewItem.RecordID)

	// get /resources/:id
	engine.ServeHTTP(w, newGetRequest("%s/%s", nil, router, addNewItem.RecordID))
	assert.Equal(t, 200, w.Code)
	var putItem schema.Resource
	err = parseReader(w.Body, &putItem)
	// put /resources/:id
	putItem.Code = "test2"
	putItem.Name = "测试资源2"
	engine.ServeHTTP(w, newPutRequest("%s/%s", putItem, router, addNewItem.RecordID))
	assert.Equal(t, 200, w.Code)
	releaseReader(w.Body)

	// get /resources/:id
	engine.ServeHTTP(w, newGetRequest("%s/%s", nil, router, addNewItem.RecordID))
	assert.Equal(t, 200, w.Code)

	var getItem schema.Resource
	err = parseReader(w.Body, &getItem)

	assert.Nil(t, err)
	assert.Equal(t, getItem.RecordID, addNewItem.RecordID)
	assert.Equal(t, getItem.Code, putItem.Code)
	assert.Equal(t, getItem.Name, putItem.Name)

	// delete /resources/:id
	engine.ServeHTTP(w, newDeleteRequest("%s/%s", router, addNewItem.RecordID))
	assert.Equal(t, 200, w.Code)
	releaseReader(w.Body)
}
