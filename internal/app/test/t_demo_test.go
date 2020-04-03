package test

import (
	"net/http/httptest"
	"testing"

	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/LyricTian/gin-admin/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestDemo(t *testing.T) {
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
	var addItemRes ResRecordID
	err = parseReader(w.Body, &addItemRes)
	assert.Nil(t, err)

	// get /demos/:id
	engine.ServeHTTP(w, newGetRequest("%s/%s", nil, router, addItemRes.RecordID))
	assert.Equal(t, 200, w.Code)
	var getItem schema.Demo
	err = parseReader(w.Body, &getItem)
	assert.Nil(t, err)
	assert.Equal(t, addItem.Code, getItem.Code)
	assert.Equal(t, addItem.Name, getItem.Name)
	assert.Equal(t, addItem.Status, getItem.Status)
	assert.NotEmpty(t, getItem.RecordID)

	// put /demos/:id
	putItem := getItem
	putItem.Name = util.MustUUID()
	engine.ServeHTTP(w, newPutRequest("%s/%s", putItem, router, getItem.RecordID))
	assert.Equal(t, 200, w.Code)
	err = parseOK(w.Body)
	assert.Nil(t, err)

	// query /demos
	engine.ServeHTTP(w, newGetRequest(router, newPageParam()))
	assert.Equal(t, 200, w.Code)
	var pageItems []*schema.Demo
	err = parsePageReader(w.Body, &pageItems)
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(pageItems), 1)
	if len(pageItems) > 0 {
		assert.Equal(t, putItem.RecordID, pageItems[0].RecordID)
		assert.Equal(t, putItem.Name, pageItems[0].Name)
	}

	// delete /demos/:id
	engine.ServeHTTP(w, newDeleteRequest("%s/%s", router, addItemRes.RecordID))
	assert.Equal(t, 200, w.Code)
	err = parseOK(w.Body)
	assert.Nil(t, err)
}
