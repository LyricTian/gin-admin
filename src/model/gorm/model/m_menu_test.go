package model_test

import (
	"context"
	"testing"

	"github.com/LyricTian/gin-admin/src/model/gorm/model"
	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/LyricTian/gin-admin/src/util"
	"github.com/stretchr/testify/assert"
)

func TestMenu(t *testing.T) {
	menu := new(model.Menu).Init(gdb)

	addItem := schema.Menu{
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
	err := menu.Create(context.Background(), addItem)
	assert.Nil(t, err)

	updateItem := addItem
	updateItem.Name = util.MustUUID()
	err = menu.Update(context.Background(), addItem.RecordID, updateItem)
	assert.Nil(t, err)

	getItem, err := menu.Get(context.Background(), addItem.RecordID, schema.MenuQueryOptions{
		IncludeActions:   true,
		IncludeResources: true,
	})
	assert.Nil(t, err)
	assert.NotNil(t, getItem)

	assert.Equal(t, updateItem.Name, getItem.Name)
	assert.Equal(t, len(updateItem.Actions), len(getItem.Actions))
	assert.Equal(t, len(updateItem.Resources), len(getItem.Resources))

	err = menu.Delete(context.Background(), addItem.RecordID)
	assert.Nil(t, err)
}
