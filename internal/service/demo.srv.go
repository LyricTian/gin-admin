package service

import (
	"context"
	"fmt"

	"github.com/google/wire"

	"github.com/LyricTian/gin-admin/v9/internal/dao"
	"github.com/LyricTian/gin-admin/v9/internal/module/consts"
	"github.com/LyricTian/gin-admin/v9/internal/schema"

	"github.com/LyricTian/gin-admin/v9/pkg/errors"
	"github.com/LyricTian/gin-admin/v9/pkg/util/xid"
)

var DemoSet = wire.NewSet(wire.Struct(new(DemoSrv), "*"))

type DemoSrv struct {
	TransRepo *dao.TransRepo
	DemoRepo  *dao.DemoRepo
}

func (a *DemoSrv) Query(ctx context.Context, params schema.DemoQueryParam) (*schema.DemoQueryResult, error) {
	params.Pagination = true
	result, err := a.DemoRepo.Query(ctx, params)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *DemoSrv) Get(ctx context.Context, id string) (*schema.Demo, error) {
	item, err := a.DemoRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.NotFound(consts.ErrNotFoundID, fmt.Sprintf("not found: %v", id))
	}

	return item, nil
}

func (a *DemoSrv) Create(ctx context.Context, item schema.Demo) (*schema.Demo, error) {
	item.ID = xid.NewID()

	err := a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		return a.DemoRepo.Create(ctx, &item)
	})
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (a *DemoSrv) Update(ctx context.Context, id string, item schema.Demo) error {
	oldItem, err := a.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.NotFound(consts.ErrNotFoundID, fmt.Sprintf("not found: %v", id))
	}

	item.ID = id
	item.CreatedAt = oldItem.CreatedAt

	return a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		return a.DemoRepo.Update(ctx, &item)
	})
}

func (a *DemoSrv) Delete(ctx context.Context, id string) error {
	oldItem, err := a.DemoRepo.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.NotFound(consts.ErrNotFoundID, fmt.Sprintf("not found: %v", id))
	}

	return a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		return a.DemoRepo.Delete(ctx, id)
	})
}
