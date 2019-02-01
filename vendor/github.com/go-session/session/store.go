package session

import (
	"container/list"
	"context"
	"sync"
	"time"
)

var (
	_   ManagerStore = &memoryStore{}
	_   Store        = &store{}
	now              = time.Now
)

// ManagerStore Management of session storage, including creation, update, and delete operations
type ManagerStore interface {
	// Check the session store exists
	Check(ctx context.Context, sid string) (bool, error)
	// Create a session store and specify the expiration time (in seconds)
	Create(ctx context.Context, sid string, expired int64) (Store, error)
	// Update a session store and specify the expiration time (in seconds)
	Update(ctx context.Context, sid string, expired int64) (Store, error)
	// Delete a session store
	Delete(ctx context.Context, sid string) error
	// Use sid to replace old sid and return session store
	Refresh(ctx context.Context, oldsid, sid string, expired int64) (Store, error)
	// Close storage, release resources
	Close() error
}

// Store A session id storage operation
type Store interface {
	// Get a session storage context
	Context() context.Context
	// Get the current session id
	SessionID() string
	// Set session value, call save function to take effect
	Set(key string, value interface{})
	// Get session value
	Get(key string) (interface{}, bool)
	// Delete session value, call save function to take effect
	Delete(key string) interface{}
	// Save session data
	Save() error
	// Clear all session data
	Flush() error
}

// NewMemoryStore create an instance of a memory store
func NewMemoryStore() ManagerStore {
	mstore := &memoryStore{
		data:   make(map[string]*dataItem),
		list:   list.New(),
		ticker: time.NewTicker(time.Second),
	}

	go mstore.gc()

	return mstore
}

type dataItem struct {
	sid       string
	expiredAt time.Time
	values    map[string]interface{}
}

func newDataItem(sid string, values map[string]interface{}, expired int64) *dataItem {
	return &dataItem{
		sid:       sid,
		expiredAt: now().Add(time.Duration(expired) * time.Second),
		values:    values,
	}
}

type memoryStore struct {
	sync.RWMutex
	data   map[string]*dataItem
	list   *list.List
	ticker *time.Ticker
}

func (s *memoryStore) gc() {
	for range s.ticker.C {
		s.RLock()
		e := s.list.Front()
		s.RUnlock()

		for e != nil {
			item := e.Value.(*dataItem)
			if item.expiredAt.Before(now()) {
				s.Lock()
				s.list.Remove(e)
				delete(s.data, item.sid)
				e = e.Next()
				s.Unlock()
			} else {
				break
			}
		}
	}
}

func (s *memoryStore) save(sid string, values map[string]interface{}, expired int64) {
	s.Lock()
	defer s.Unlock()

	if item, ok := s.data[sid]; ok {
		item.values = values
		return
	}

	item := newDataItem(sid, values, expired)
	s.data[sid] = item
	s.list.PushBack(item)
}

func (s *memoryStore) Check(_ context.Context, sid string) (bool, error) {
	s.RLock()
	item, ok := s.data[sid]
	s.RUnlock()

	if ok && item.expiredAt.After(now()) {
		return true, nil
	}
	return false, nil
}

func (s *memoryStore) Create(ctx context.Context, sid string, expired int64) (Store, error) {
	return newStore(ctx, s, sid, expired, nil), nil
}

func (s *memoryStore) Update(ctx context.Context, sid string, expired int64) (Store, error) {
	s.Lock()
	defer s.Unlock()

	item, ok := s.data[sid]
	if !ok {
		return newStore(ctx, s, sid, expired, nil), nil
	}

	item.expiredAt = now().Add(time.Duration(expired) * time.Second)
	for e := s.list.Front(); e != nil; e = e.Next() {
		if e.Value.(*dataItem).sid == sid {
			s.list.MoveToBack(e)
			break
		}
	}

	return newStore(ctx, s, sid, expired, item.values), nil
}

func (s *memoryStore) delete(sid string) {
	delete(s.data, sid)

	for e := s.list.Front(); e != nil; e = e.Next() {
		if e.Value.(*dataItem).sid == sid {
			s.list.Remove(e)
			break
		}
	}
}

func (s *memoryStore) Delete(_ context.Context, sid string) error {
	s.Lock()
	defer s.Unlock()

	s.delete(sid)
	return nil
}

func (s *memoryStore) Refresh(ctx context.Context, oldsid, sid string, expired int64) (Store, error) {
	s.Lock()
	defer s.Unlock()

	item, ok := s.data[oldsid]
	if !ok {
		return newStore(ctx, s, sid, expired, nil), nil
	}

	newItem := newDataItem(sid, item.values, expired)
	s.data[sid] = newItem
	s.list.PushBack(newItem)
	s.delete(oldsid)

	return newStore(ctx, s, sid, expired, newItem.values), nil
}

func (s *memoryStore) Close() error {
	s.ticker.Stop()
	return nil
}

func newStore(ctx context.Context, mstore *memoryStore, sid string, expired int64, values map[string]interface{}) *store {
	if values == nil {
		values = make(map[string]interface{})
	}

	return &store{
		mstore:  mstore,
		ctx:     ctx,
		sid:     sid,
		expired: expired,
		values:  values,
	}
}

type store struct {
	sync.RWMutex
	mstore  *memoryStore
	ctx     context.Context
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
	s.RLock()
	values := s.values
	s.RUnlock()

	s.mstore.save(s.sid, values, s.expired)
	return nil
}
