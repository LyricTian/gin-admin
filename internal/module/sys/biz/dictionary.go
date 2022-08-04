package biz

import (
	"context"

	"github.com/LyricTian/gin-admin/v9/internal/module/sys/dao"
	"github.com/LyricTian/gin-admin/v9/internal/module/sys/typed"
	"github.com/LyricTian/gin-admin/v9/internal/x/contextx"
	"github.com/LyricTian/gin-admin/v9/internal/x/utilx"
	"github.com/LyricTian/gin-admin/v9/pkg/errors"
	"github.com/LyricTian/gin-admin/v9/pkg/util/xid"
)

type DictionaryBiz struct {
	TransRepo      utilx.TransRepo
	DictionaryRepo dao.DictionaryRepo
}

func (a *DictionaryBiz) Query(ctx context.Context, params typed.DictionaryQueryParam) (*typed.DictionaryQueryResult, error) {
	params.Pagination = true
	queryOpts := utilx.QueryOptions{
		OrderFields: []utilx.OrderByParam{
			{Field: "created_at", Direction: utilx.DESC},
		},
	}

	result, err := a.DictionaryRepo.Query(ctx, params, typed.DictionaryQueryOptions{
		QueryOptions: queryOpts,
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *DictionaryBiz) Get(ctx context.Context, id string) (*typed.Dictionary, error) {
	dictionary, err := a.DictionaryRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	} else if dictionary == nil {
		return nil, errors.NotFound(errors.ErrNotFoundID, "Dictionary not found")
	}

	return dictionary, nil
}

func (a *DictionaryBiz) Create(ctx context.Context, createItem typed.DictionaryCreate) (*typed.Dictionary, error) {
	dictionary := &typed.Dictionary{
		ID:        xid.NewID(),
		Ns:        createItem.Ns,
		Key:       createItem.Key,
		Value:     &createItem.Value,
		Remark:    &createItem.Remark,
		CreatedBy: contextx.FromUserID(ctx),
	}

	err := a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		if err := a.DictionaryRepo.Create(ctx, dictionary); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return dictionary, nil
}

func (a *DictionaryBiz) Update(ctx context.Context, id string, createItem typed.DictionaryCreate) error {
	oldDictionary, err := a.DictionaryRepo.Get(ctx, id)
	if err != nil {
		return err
	} else if oldDictionary == nil {
		return errors.NotFound(errors.ErrNotFoundID, "Dictionary not found")
	}
	oldDictionary.Ns = createItem.Ns
	oldDictionary.Key = createItem.Key
	oldDictionary.Value = &createItem.Value
	oldDictionary.Remark = &createItem.Remark
	oldDictionary.UpdatedBy = contextx.FromUserID(ctx)

	return a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		if err := a.DictionaryRepo.Update(ctx, oldDictionary); err != nil {
			return err
		}

		return nil
	})
}

func (a *DictionaryBiz) Delete(ctx context.Context, id string) error {
	exists, err := a.DictionaryRepo.Exists(ctx, id)
	if err != nil {
		return err
	} else if !exists {
		return errors.NotFound(errors.ErrNotFoundID, "Dictionary not found")
	}

	return a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		if err := a.DictionaryRepo.Delete(ctx, id); err != nil {
			return err
		}

		return nil
	})
}
