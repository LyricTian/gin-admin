package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/go-session/session"
)

var (
	_             session.ManagerStore = &managerStore{}
	_             session.Store        = &store{}
	jsonMarshal                        = json.Marshal
	jsonUnmarshal                      = json.Unmarshal
)

// NewConfig create mysql configuration instance
func NewConfig(dsn string) *Config {
	return &Config{
		DSN:             dsn,
		ConnMaxLifetime: time.Hour * 2,
		MaxOpenConns:    50,
		MaxIdleConns:    25,
	}
}

// Config mysql configuration
type Config struct {
	DSN             string
	ConnMaxLifetime time.Duration
	MaxOpenConns    int
	MaxIdleConns    int
}

// NewStore Create an instance of a mysql store,
// tableName Specify the stored table name (default go_session),
// gcInterval Time interval for executing GC (in seconds, default 600)
func NewStore(config *Config, tableName string, gcInterval int) session.ManagerStore {
	db, err := sql.Open("mysql", config.DSN)
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)

	return NewStoreWithDB(db, tableName, gcInterval)
}

// NewStoreWithDB Create an instance of a mysql store,
// tableName Specify the stored table name (default go_session),
// gcInterval Time interval for executing GC (in seconds, default 600)
func NewStoreWithDB(db *sql.DB, tableName string, gcInterval int) session.ManagerStore {
	store := &managerStore{
		db:        db,
		tableName: "go_session",
		stdout:    os.Stderr,
	}

	if tableName != "" {
		store.tableName = tableName
	}

	interval := 600
	if gcInterval > 0 {
		interval = gcInterval
	}
	store.ticker = time.NewTicker(time.Second * time.Duration(interval))

	store.db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` (`id` VARCHAR(255) NOT NULL PRIMARY KEY, `value` VARCHAR(2048), `expired_at` bigint) engine=InnoDB charset=UTF8;", store.tableName))
	store.db.Exec(fmt.Sprintf("CREATE INDEX `idx_expired_at` ON %s (`expired_at`);", store.tableName))

	go store.gc()
	return store
}

type managerStore struct {
	ticker    *time.Ticker
	db        *sql.DB
	tableName string
	stdout    io.Writer
}

func (s *managerStore) errorf(format string, args ...interface{}) {
	if s.stdout != nil {
		buf := fmt.Sprintf(format, args...)
		s.stdout.Write([]byte(buf))
	}
}

func (s *managerStore) gc() {
	for range s.ticker.C {
		now := time.Now().Unix()

		var count int
		row := s.db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM `%s` WHERE `expired_at`<=?", s.tableName), now)
		err := row.Scan(&count)
		if err != nil {
			s.errorf("[ERROR]:%s", err.Error())
			return
		} else if count > 0 {
			_, err = s.db.Exec(fmt.Sprintf("DELETE FROM `%s` WHERE `expired_at`<=?", s.tableName), now)
			if err != nil {
				s.errorf("[ERROR]:%s", err.Error())
			}
		}
	}
}

func (s *managerStore) getValue(sid string) (string, error) {
	var item SessionItem

	row := s.db.QueryRow(fmt.Sprintf("SELECT `id`,`value`,`expired_at` FROM `%s` WHERE `id`=?", s.tableName), sid)
	err := row.Scan(&item.ID, &item.Value, &item.ExpiredAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", nil
	} else if time.Now().Unix() >= item.ExpiredAt {
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

func (s *managerStore) Check(_ context.Context, sid string) (bool, error) {
	val, err := s.getValue(sid)
	if err != nil {
		return false, err
	}
	return val != "", nil
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

	_, err = s.db.Exec(fmt.Sprintf("UPDATE `%s` SET `expired_at`=? WHERE `id`=?", s.tableName),
		time.Now().Add(time.Duration(expired)*time.Second).Unix(),
		sid)
	if err != nil {
		return nil, err
	}

	values, err := s.parseValue(value)
	if err != nil {
		return nil, err
	}

	return newStore(ctx, s, sid, expired, values), nil
}

func (s *managerStore) Delete(_ context.Context, sid string) error {
	_, err := s.db.Exec(fmt.Sprintf("DELETE FROM `%s` WHERE `id`=?", s.tableName), sid)
	return err
}

func (s *managerStore) Refresh(ctx context.Context, oldsid, sid string, expired int64) (session.Store, error) {
	value, err := s.getValue(oldsid)
	if err != nil {
		return nil, err
	} else if value == "" {
		return newStore(ctx, s, sid, expired, nil), nil
	}

	query := fmt.Sprintf("INSERT INTO `%s` (`id`,`value`,`expired_at`) VALUES (?,?,?);", s.tableName)
	_, err = s.db.Exec(query, sid, value, time.Now().Add(time.Duration(expired)*time.Second).Unix())
	if err != nil {
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
	s.db.Close()
	return nil
}

func newStore(ctx context.Context, s *managerStore, sid string, expired int64, values map[string]interface{}) *store {
	if values == nil {
		values = make(map[string]interface{})
	}

	return &store{
		db:        s.db,
		tableName: s.tableName,
		ctx:       ctx,
		sid:       sid,
		expired:   expired,
		values:    values,
	}
}

type store struct {
	sync.RWMutex
	ctx       context.Context
	db        *sql.DB
	tableName string
	sid       string
	expired   int64
	values    map[string]interface{}
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

	var count int
	row := s.db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE id=?", s.tableName), s.sid)
	err := row.Scan(&count)
	if err != nil {
		return err
	} else if count == 0 {
		query := fmt.Sprintf("INSERT INTO `%s` (`id`,`value`,`expired_at`) VALUES (?,?,?);", s.tableName)
		_, err = s.db.Exec(query, s.sid, value, time.Now().Add(time.Duration(s.expired)*time.Second).Unix())
		return err
	}

	_, err = s.db.Exec(fmt.Sprintf("UPDATE `%s` SET `value`=?,`expired_at`=? WHERE `id`=?", s.tableName),
		value,
		time.Now().Add(time.Duration(s.expired)*time.Second).Unix(),
		s.sid)

	return err
}

// SessionItem Data items stored in mysql
type SessionItem struct {
	ID        string
	Value     string
	ExpiredAt int64
}
