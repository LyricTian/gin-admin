package test

import (
	"net/http/httptest"
	"testing"

	"github.com/LyricTian/gin-admin/v8/internal/app/schema"
	"github.com/LyricTian/gin-admin/v8/pkg/util/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMenu(t *testing.T) {
	const router = apiPrefix + "v1/menus"
	var err error

	w := httptest.NewRecorder()

	// post /menus
	addItem := &schema.Menu{
		Name:   uuid.MustUUID().String(),
		IsShow: 1,
		Status: 1,
	}
	engine.ServeHTTP(w, newPostRequest(router, addItem))
	assert.Equal(t, 200, w.Code)
	var addItemRes ResID
	err = parseReader(w.Body, &addItemRes)
	assert.Nil(t, err)

	// get /menus/:id
	engine.ServeHTTP(w, newGetRequest("%s/%d", nil, router, addItemRes.ID))
	assert.Equal(t, 200, w.Code)
	var getItem schema.Menu
	err = parseReader(w.Body, &getItem)
	assert.Nil(t, err)
	assert.Equal(t, addItem.Name, getItem.Name)
	assert.Equal(t, addItem.Status, getItem.Status)
	assert.NotEmpty(t, getItem.ID)

	// put /menus/:id
	putItem := getItem
	putItem.Name = uuid.MustUUID().String()
	engine.ServeHTTP(w, newPutRequest("%s/%d", putItem, router, getItem.ID))
	assert.Equal(t, 200, w.Code)
	err = parseOK(w.Body)
	assert.Nil(t, err)

	// query /menus
	engine.ServeHTTP(w, newGetRequest(router, newPageParam()))
	assert.Equal(t, 200, w.Code)
	var pageItems []*schema.Menu
	err = parsePageReader(w.Body, &pageItems)
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(pageItems), 1)
	if len(pageItems) > 0 {
		assert.Equal(t, putItem.ID, pageItems[0].ID)
		assert.Equal(t, putItem.Name, pageItems[0].Name)
	}

	// delete /menus/:id
	engine.ServeHTTP(w, newDeleteRequest("%s/%d", router, addItemRes.ID))
	assert.Equal(t, 200, w.Code)
	err = parseOK(w.Body)
	assert.Nil(t, err)
}
