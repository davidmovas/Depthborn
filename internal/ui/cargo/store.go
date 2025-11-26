package cargo

import (
	"fmt"
	"reflect"
	"sync"
)

var cargo = &cargoStore{
	data:      make(map[string]any),
	observers: make([]ObserverFunc, 0),
	mu:        sync.RWMutex{},
}

// ObserverFunc is called when cargo state changes
// key: the key that was modified
// value: the new value
type ObserverFunc func(key string, value any)

// AddObserver registers global observer for cargo state changes
// Observer is called whenever any key is modified via Set()
func AddObserver(obs ObserverFunc) {
	cargo.addObserver(obs)
}

type cargoStore struct {
	data      map[string]any
	observers []ObserverFunc
	mu        sync.RWMutex
}

func (s *cargoStore) get(key string) (any, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	return val, ok
}

func (s *cargoStore) set(key string, value any) {
	s.mu.Lock()
	s.data[key] = value
	observers := make([]ObserverFunc, len(s.observers))
	copy(observers, s.observers)
	s.mu.Unlock()

	for _, obs := range observers {
		obs(key, value)
	}
}

func (s *cargoStore) delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}

func (s *cargoStore) has(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.data[key]
	return ok
}

func (s *cargoStore) clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = make(map[string]any)
}

func (s *cargoStore) addObserver(obs ObserverFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.observers = append(s.observers, obs)
}

func typeKey[T any]() string {
	var zero T
	t := reflect.TypeOf(zero)
	if t == nil {
		return "nil"
	}
	// Handle pointers
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return fmt.Sprintf("__global_%s_%s", t.PkgPath(), t.Name())
}
