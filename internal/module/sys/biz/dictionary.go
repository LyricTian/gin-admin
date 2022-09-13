package biz

import (
	"context"
	"strings"

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
	params.Pagination = false
	params.PageSize = -1

	if v := params.ParentID; v != nil {
		if *v == "-1" {
			params.ParentID = nil
		} else if *v == "0" {
			params.ParentID = new(string)
		}
	}

	var parentData typed.Dictionaries
	if params.Key != "" {
		var parentID string
		for _, key := range strings.Split(params.Key, typed.DictionaryPathDelimiter) {
			dictionary, err := a.DictionaryRepo.GetByKeyAndParentID(ctx, key, parentID)
			if err != nil {
				return nil, err
			} else if dictionary == nil {
				return nil, errors.NotFound(errors.ErrNotFoundID, "Dictionary not found")
			}
			parentID = dictionary.ID
			params.LikeLeftParentPath = params.LikeLeftParentPath + parentID + typed.DictionaryPathDelimiter
			parentData = append(parentData, dictionary)
		}
	}

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

	if len(parentData) > 0 {
		result.Data = append(result.Data, parentData...)
	}

	if params.QueryValue != "" && len(result.Data) > 0 {
		parentResult, err := a.DictionaryRepo.Query(ctx, typed.DictionaryQueryParam{
			IDs: result.Data.SplitParentIDs(),
		}, typed.DictionaryQueryOptions{
			QueryOptions: queryOpts,
		})
		if err != nil {
			return nil, err
		} else if len(parentResult.Data) > 0 {
			result.Data = append(result.Data, parentResult.Data...)
		}
	}

	result.Data = result.Data.ToTree()
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
		Key:       createItem.Key,
		Value:     &createItem.Value,
		Remark:    &createItem.Remark,
		CreatedBy: contextx.FromUserID(ctx),
	}

	if createItem.ParentID != "" {
		parent, err := a.DictionaryRepo.Get(ctx, createItem.ParentID)
		if err != nil {
			return nil, err
		} else if parent == nil {
			return nil, errors.NotFound(errors.ErrNotFoundID, "Dictionary parent not found")
		}
		dictionary.ParentID = &createItem.ParentID

		parentPath := *parent.ParentPath + parent.ID + typed.DictionaryPathDelimiter
		dictionary.ParentPath = &parentPath
	}

	exists, err := a.DictionaryRepo.GetByKeyAndParentID(ctx, createItem.Key, createItem.ParentID)
	if err != nil {
		return nil, err
	} else if exists != nil {
		return nil, errors.BadRequest(errors.ErrBadRequestID, "Dictionary key already exists")
	}

	err = a.TransRepo.Exec(ctx, func(ctx context.Context) error {
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

func (a *DictionaryBiz) Update(ctx context.Context, id string, updateItem typed.DictionaryUpdate) error {
	oldDictionary, err := a.DictionaryRepo.Get(ctx, id)
	if err != nil {
		return err
	} else if oldDictionary == nil {
		return errors.NotFound(errors.ErrNotFoundID, "Dictionary not found")
	} else if updateItem.Key != oldDictionary.Key {
		exists, err := a.DictionaryRepo.GetByKeyAndParentID(ctx, updateItem.Key, *oldDictionary.ParentID)
		if err != nil {
			return err
		} else if exists != nil {
			return errors.BadRequest(errors.ErrBadRequestID, "Dictionary key already exists")
		}
	}

	oldDictionary.Key = updateItem.Key
	oldDictionary.Value = &updateItem.Value
	oldDictionary.Remark = &updateItem.Remark
	oldDictionary.UpdatedBy = contextx.FromUserID(ctx)

	return a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		if err := a.DictionaryRepo.Update(ctx, oldDictionary); err != nil {
			return err
		}

		return nil
	})
}

func (a *DictionaryBiz) Delete(ctx context.Context, id string) error {
	dictionary, err := a.DictionaryRepo.Get(ctx, id)
	if err != nil {
		return err
	} else if dictionary == nil {
		return errors.NotFound(errors.ErrNotFoundID, "Dictionary not found")
	}

	return a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		if err := a.DictionaryRepo.Delete(ctx, id); err != nil {
			return err
		}

		return a.DictionaryRepo.DeleteByParentPath(ctx, *dictionary.ParentPath+id+typed.DictionaryPathDelimiter)
	})
}
