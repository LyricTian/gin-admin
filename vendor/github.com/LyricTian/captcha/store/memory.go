package store

import (
	"container/list"
	"sync"
	"time"
)

// NewMemoryStore An internal store for captcha ids and their values.
func NewMemoryStore(gcInterval, expiration time.Duration) Store {
	mstore := &memoryStore{
		data:       make(map[string]*dataItem),
		list:       list.New(),
		ticker:     time.NewTicker(gcInterval),
		expiration: expiration,
	}

	go mstore.gc()

	return mstore
}

type dataItem struct {
	id        string
	value     []byte
	expiredAt time.Time
}

// memoryStore memory store
type memoryStore struct {
	sync.RWMutex
	data       map[string]*dataItem
	list       *list.List
	ticker     *time.Ticker
	expiration time.Duration
}

func (s *memoryStore) gc() {
	for range s.ticker.C {
		s.RLock()
		e := s.list.Front()
		s.RUnlock()

		for e != nil {
			item := e.Value.(*dataItem)
			if item.expiredAt.Before(time.Now()) {
				s.Lock()
				s.list.Remove(e)
				delete(s.data, item.id)
				e = e.Next()
				s.Unlock()
			} else {
				break
			}
		}
	}
}

func (s *memoryStore) Set(id string, digits []byte) {
	s.Lock()
	defer s.Unlock()

	expiredAt := time.Now().Add(s.expiration)
	if _, ok := s.data[id]; ok {
		e := s.list.Front()
		for e != nil {
			item := e.Value.(*dataItem)
			if item.id == id {
				item.value = digits
				item.expiredAt = expiredAt
				s.list.MoveToBack(e)
				break
			}
			e = e.Next()
		}
		return
	}

	item := &dataItem{
		id:        id,
		value:     digits,
		expiredAt: expiredAt,
	}
	s.data[id] = item
	s.list.PushBack(item)
}

func (s *memoryStore) Get(id string, clear bool) []byte {
	s.RLock()
	item, ok := s.data[id]
	s.RUnlock()

	if !ok {
		return nil
	}

	if clear {
		s.remove(id)
	}

	return item.value
}

func (s *memoryStore) remove(id string) {
	s.RLock()
	e := s.list.Front()
	for e != nil {
		item := e.Value.(*dataItem)
		if item.id == id {
			break
		}
		e = e.Next()
	}
	s.RUnlock()

	s.Lock()
	s.list.Remove(e)
	delete(s.data, id)
	s.Unlock()
}
