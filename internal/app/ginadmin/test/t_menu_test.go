package test

import (
	"net/http/httptest"
	"testing"

	"github.com/LyricTian/gin-admin/internal/app/ginadmin/schema"
	"github.com/LyricTian/gin-admin/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestMenu(t *testing.T) {
	const router = "v1/menus"
	var err error

	w := httptest.NewRecorder()

	// post /menus
	addItem := &schema.Menu{
		Name:     util.MustUUID(),
		Sequence: 9999999,
		Router:   "/system/menu",
		Actions: []*schema.MenuAction{
			{Code: "query", Name: "query"},
		},
		Resources: []*schema.MenuResource{
			{Code: "query", Name: "query", Method: "GET", Path: "/test/v1/menus"},
		},
	}
	engine.ServeHTTP(w, newPostRequest(router, addItem))
	assert.Equal(t, 200, w.Code)

	var addNewItem schema.Menu
	err = parseReader(w.Body, &addNewItem)
	assert.Nil(t, err)
	assert.Equal(t, addItem.Name, addNewItem.Name)
	assert.Equal(t, addItem.Router, addNewItem.Router)
	assert.Equal(t, addItem.Sequence, addNewItem.Sequence)
	assert.Equal(t, len(addItem.Actions), len(addNewItem.Actions))
	assert.Equal(t, len(addItem.Resources), len(addNewItem.Resources))
	assert.NotEmpty(t, addNewItem.RecordID)

	// query /menus?q=page
	engine.ServeHTTP(w, newGetRequest(router,
		newPageParam(map[string]string{"q": "page"})))
	assert.Equal(t, 200, w.Code)
	var pageItems []*schema.Menu
	err = parsePageReader(w.Body, &pageItems)
	assert.Nil(t, err)
	assert.Equal(t, len(pageItems), 1)
	if len(pageItems) > 0 {
		assert.Equal(t, addNewItem.RecordID, pageItems[0].RecordID)
		assert.Equal(t, addNewItem.Name, pageItems[0].Name)
	}

	// put /menus/:id
	engine.ServeHTTP(w, newGetRequest("%s/%s", nil, router, addNewItem.RecordID))
	assert.Equal(t, 200, w.Code)
	var putItem schema.Menu
	err = parseReader(w.Body, &putItem)
	putItem.Name = util.MustUUID()
	engine.ServeHTTP(w, newPutRequest("%s/%s", putItem, router, addNewItem.RecordID))
	assert.Equal(t, 200, w.Code)
	var putNewItem schema.Menu
	err = parseReader(w.Body, &putNewItem)
	assert.Nil(t, err)
	assert.Equal(t, putItem.Name, putNewItem.Name)
	assert.Equal(t, len(putItem.Actions), len(putNewItem.Actions))
	assert.Equal(t, len(putItem.Resources), len(putNewItem.Resources))

	// delete /menus/:id
	engine.ServeHTTP(w, newDeleteRequest("%s/%s", router, addNewItem.RecordID))
	assert.Equal(t, 200, w.Code)
	err = parseOK(w.Body)
	assert.Nil(t, err)
}
