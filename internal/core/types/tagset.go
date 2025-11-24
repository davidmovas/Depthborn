package types

import "sync"

var _ TagSet = (*BaseTagSet)(nil)

type BaseTagSet struct {
	mu   sync.RWMutex
	tags map[string]struct{}
}

func NewTagSet() TagSet {
	return &BaseTagSet{
		tags: make(map[string]struct{}),
	}
}

func (ts *BaseTagSet) Add(tag string) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.tags[tag] = struct{}{}
}

func (ts *BaseTagSet) Remove(tag string) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	delete(ts.tags, tag)
}

func (ts *BaseTagSet) Has(tag string) bool {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	_, exists := ts.tags[tag]
	return exists
}

func (ts *BaseTagSet) Contains(tags ...string) bool {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	for _, tag := range tags {
		if _, exists := ts.tags[tag]; !exists {
			return false
		}
	}
	return true
}

func (ts *BaseTagSet) ContainsAny(tags ...string) bool {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	for _, tag := range tags {
		if _, exists := ts.tags[tag]; exists {
			return true
		}
	}
	return false
}

func (ts *BaseTagSet) All() []string {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	result := make([]string, 0, len(ts.tags))
	for tag := range ts.tags {
		result = append(result, tag)
	}
	return result
}

func (ts *BaseTagSet) Clear() {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.tags = make(map[string]struct{})
}
