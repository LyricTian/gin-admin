package data

import (
	"context"
	"os"

	"github.com/LyricTian/gin-admin/internal/app/bll"
	"github.com/LyricTian/gin-admin/internal/app/config"
	"github.com/LyricTian/gin-admin/internal/app/model"
	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/LyricTian/gin-admin/pkg/util"
	"github.com/google/wire"
)

// MenuSet 注入Menu
var MenuSet = wire.NewSet(wire.Struct(new(Menu), "*"))

// Menu 菜单数据
type Menu struct {
	TransModel model.ITrans
	MenuBll    bll.IMenu
}

// Load 加载菜单数据
func (a *Menu) Load() error {
	c := config.C.Menu
	if !c.Enable || c.Data == "" {
		return nil
	}

	ctx := context.Background()
	result, err := a.MenuBll.Query(ctx, schema.MenuQueryParam{
		PaginationParam: schema.PaginationParam{OnlyCount: true},
	})
	if err != nil {
		return err
	} else if result.PageResult.Total > 0 {
		return nil
	}

	data, err := a.readData(c.Data)
	if err != nil {
		return err
	}

	return a.createMenus(ctx, "", data)
}

func (a *Menu) readData(name string) (schema.MenuTrees, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data schema.MenuTrees
	d := util.YAMLNewDecoder(file)
	d.SetStrict(true)
	err = d.Decode(&data)
	return data, err
}

func (a *Menu) createMenus(ctx context.Context, parentID string, list schema.MenuTrees) error {
	return a.TransModel.Exec(ctx, func(ctx context.Context) error {
		for _, item := range list {
			sitem := schema.Menu{
				Name:       item.Name,
				Sequence:   item.Sequence,
				Icon:       item.Icon,
				Router:     item.Router,
				ParentID:   parentID,
				Status:     1,
				ShowStatus: 1,
				Actions:    item.Actions,
			}
			if v := item.ShowStatus; v > 0 {
				sitem.ShowStatus = v
			}
			nsitem, err := a.MenuBll.Create(ctx, sitem)
			if err != nil {
				return err
			}

			if item.Children != nil && len(*item.Children) > 0 {
				err := a.createMenus(ctx, nsitem.RecordID, *item.Children)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}
