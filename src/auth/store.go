package auth

import (
	"os"
	"path/filepath"
	"time"

	"github.com/tidwall/buntdb"
)

// BlackStorer 黑名单存储接口
type BlackStorer interface {
	// 放入黑名单，指定到期时间
	Set(tokenString string, expiration time.Duration) error
	// 检查令牌是否存在
	Check(tokenString string) (bool, error)
	// 关闭存储
	Close() error
}

// NewFileBlackStore 创建基于文件存储的黑名单存储实例
func NewFileBlackStore(path string) BlackStorer {
	os.MkdirAll(filepath.Dir(path), 0777)
	db, err := buntdb.Open(path)
	if err != nil {
		panic(err)
	}
	return &blackStore{
		db: db,
	}
}

type blackStore struct {
	db *buntdb.DB
}

func (a *blackStore) Set(tokenString string, expiration time.Duration) error {
	return a.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(tokenString, "1", &buntdb.SetOptions{Expires: true, TTL: expiration})
		return err
	})
}

func (a *blackStore) Check(tokenString string) (bool, error) {
	var exists bool
	err := a.db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(tokenString)
		if err != nil && err != buntdb.ErrNotFound {
			return err
		}
		exists = val == "1"
		return nil
	})
	return exists, err
}

func (a *blackStore) Close() error {
	return a.db.Close()
}
