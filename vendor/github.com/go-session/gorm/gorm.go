package gorm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/go-session/session"
	"github.com/jinzhu/gorm"
)

var (
	_             session.ManagerStore = &managerStore{}
	_             session.Store        = &store{}
	jsonMarshal                        = json.Marshal
	jsonUnmarshal                      = json.Unmarshal
)

// SessionItem Data items stored in mysql
type SessionItem struct {
	ID        string    `gorm:"column:id;size:255;primary_key;"`
	Value     string    `gorm:"column:value;size:2048;"`
	CreatedAt time.Time `gorm:"column:created_at;"`
	ExpiredAt time.Time `gorm:"column:expired_at;"`
}

// NewDefaultStore Create an instance of a gorm store
func NewDefaultStore(db *gorm.DB) session.ManagerStore {
	return NewStoreWithDB(db, "", 0)
}

// NewStoreWithDB Create an instance of a gorm store,
// tableName Specify the stored table name (default go_session),
// gcInterval Time interval for executing GC (in seconds, default 600)
func NewStoreWithDB(db *gorm.DB, tableName string, gcInterval int) session.ManagerStore {
	store := &managerStore{
		tableName: "session",
		stdout:    os.Stderr,
	}

	if tableName != "" {
		store.tableName = tableName
	}
	store.db = db.Table(store.tableName)

	if !db.HasTable(store.tableName) {
		err := store.db.CreateTable(&SessionItem{}).Error
		if err != nil {
			panic(err)
		}
		store.db.AddIndex("idx_expired_at", "expired_at")
	}

	interval := 600
	if gcInterval > 0 {
		interval = gcInterval
	}
	store.ticker = time.NewTicker(time.Second * time.Duration(interval))

	go store.gc()
	return store
}

type managerStore struct {
	ticker    *time.Ticker
	db        *gorm.DB
	tableName string
	stdout    io.Writer
}

func (s *managerStore) gc() {
	for range s.ticker.C {
		s.clean()
	}
}

func (s *managerStore) clean() {
	db := s.db.Where("expired_at<=?", time.Now())

	var count int
	err := db.Count(&count).Error
	if err != nil || count == 0 {
		if err != nil {
			s.errorf(err.Error())
		}
		return
	}

	err = db.Delete(nil).Error
	if err != nil {
		s.errorf(err.Error())
	}
}

func (s *managerStore) errorf(format string, args ...interface{}) {
	if s.stdout != nil {
		buf := fmt.Sprintf("[GORM-SESSION-ERROR]: "+format, args...)
		s.stdout.Write([]byte(buf))
	}
}

func (s *managerStore) getValue(sid string) (string, error) {
	var item SessionItem
	err := s.db.Where("id=?", sid).First(&item).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil
		}
	} else if item.ExpiredAt.Before(time.Now()) {
		return "", nil
	}
	return item.Value, nil
}

func (s *managerStore) parseValue(value string) (map[string]interface{}, error) {
	var values map[string]interface{}
	if len(value) > 0 {
		err := jsonUnmarshal([]byte(value), &values)
		if err != nil {
			return nil, err
		}
	}

	return values, nil
}

func (s *managerStore) GetExpired(expired int64) time.Time {
	return time.Now().Add(time.Duration(expired) * time.Second)
}

func (s *managerStore) Check(_ context.Context, sid string) (bool, error) {
	var count int
	result := s.db.Where("id=?", sid).Count(&count)
	if err := result.Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *managerStore) Create(ctx context.Context, sid string, expired int64) (session.Store, error) {
	return newStore(ctx, s, sid, expired, nil), nil
}

func (s *managerStore) Update(ctx context.Context, sid string, expired int64) (session.Store, error) {
	value, err := s.getValue(sid)
	if err != nil {
		return nil, err
	} else if value == "" {
		return newStore(ctx, s, sid, expired, nil), nil
	}

	result := s.db.Where("id=?", sid).Update("expired_at", s.GetExpired(expired))
	if err := result.Error; err != nil {
		return nil, err
	}

	values, err := s.parseValue(value)
	if err != nil {
		return nil, err
	}

	return newStore(ctx, s, sid, expired, values), nil
}

func (s *managerStore) Delete(_ context.Context, sid string) error {
	result := s.db.Where("id=?", sid).Delete(nil)
	return result.Error
}

func (s *managerStore) Refresh(ctx context.Context, oldsid, sid string, expired int64) (session.Store, error) {
	value, err := s.getValue(oldsid)
	if err != nil {
		return nil, err
	} else if value == "" {
		return newStore(ctx, s, sid, expired, nil), nil
	}

	item := &SessionItem{
		ID:        sid,
		Value:     value,
		CreatedAt: time.Now(),
		ExpiredAt: s.GetExpired(expired),
	}
	result := s.db.Create(item)
	if err := result.Error; err != nil {
		return nil, err
	}

	err = s.Delete(nil, oldsid)
	if err != nil {
		return nil, err
	}

	values, err := s.parseValue(value)
	if err != nil {
		return nil, err
	}

	return newStore(ctx, s, sid, expired, values), nil
}

func (s *managerStore) Close() error {
	s.ticker.Stop()
	return nil
}

func newStore(ctx context.Context, s *managerStore, sid string, expired int64, values map[string]interface{}) *store {
	if values == nil {
		values = make(map[string]interface{})
	}

	return &store{
		ctx:     ctx,
		mstore:  s,
		sid:     sid,
		expired: expired,
		values:  values,
	}
}

type store struct {
	sync.RWMutex
	ctx     context.Context
	mstore  *managerStore
	sid     string
	expired int64
	values  map[string]interface{}
}

func (s *store) Context() context.Context {
	return s.ctx
}

func (s *store) SessionID() string {
	return s.sid
}

func (s *store) Set(key string, value interface{}) {
	s.Lock()
	s.values[key] = value
	s.Unlock()
}

func (s *store) Get(key string) (interface{}, bool) {
	s.RLock()
	val, ok := s.values[key]
	s.RUnlock()
	return val, ok
}

func (s *store) Delete(key string) interface{} {
	s.RLock()
	v, ok := s.values[key]
	s.RUnlock()
	if ok {
		s.Lock()
		delete(s.values, key)
		s.Unlock()
	}
	return v
}

func (s *store) Flush() error {
	s.Lock()
	s.values = make(map[string]interface{})
	s.Unlock()
	return s.Save()
}

func (s *store) Save() error {
	var value string

	s.RLock()
	if len(s.values) > 0 {
		buf, err := jsonMarshal(s.values)
		if err != nil {
			s.RUnlock()
			return err
		}
		value = string(buf)
	}
	s.RUnlock()

	exists, err := s.mstore.Check(nil, s.sid)
	if err != nil {
		return err
	} else if !exists {
		item := &SessionItem{
			ID:        s.sid,
			Value:     value,
			CreatedAt: time.Now(),
			ExpiredAt: s.mstore.GetExpired(s.expired),
		}
		result := s.mstore.db.Create(item)
		if err := result.Error; err != nil {
			return err
		}
	}

	result := s.mstore.db.Where("id=?", s.sid).Updates(map[string]interface{}{
		"value":      value,
		"expired_at": s.mstore.GetExpired(s.expired),
	})
	return result.Error
}
