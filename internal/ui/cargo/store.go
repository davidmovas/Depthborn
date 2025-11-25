package cargo

import (
	"fmt"
	"reflect"
	"sync"
)

var cargo = &cargoStore{
	data: make(map[string]any),
	mu:   sync.RWMutex{},
}

type cargoStore struct {
	data map[string]any
	mu   sync.RWMutex
}

// get retrieves value by key
func (s *cargoStore) get(key string) (any, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	return val, ok
}

// set stores value by key
func (s *cargoStore) set(key string, value any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

// delete removes value by key
func (s *cargoStore) delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}

// has checks if key exists
func (s *cargoStore) has(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.data[key]
	return ok
}

// clear removes all values
func (s *cargoStore) clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = make(map[string]any)
}

// typeKey generates key from type T for global instances
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
