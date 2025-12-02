package character

import "sync"

var _ FlagSet = (*BaseFlagSet)(nil)

// BaseFlagSet implements FlagSet interface
type BaseFlagSet struct {
	mu    sync.RWMutex
	flags map[string]bool
}

// NewFlagSet creates a new flag set
func NewFlagSet() *BaseFlagSet {
	return &BaseFlagSet{
		flags: make(map[string]bool),
	}
}

func (f *BaseFlagSet) Set(flag string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.flags[flag] = true
}

func (f *BaseFlagSet) Unset(flag string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.flags, flag)
}

func (f *BaseFlagSet) Has(flag string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.flags[flag]
}

func (f *BaseFlagSet) Toggle(flag string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.flags[flag] {
		delete(f.flags, flag)
	} else {
		f.flags[flag] = true
	}
}

func (f *BaseFlagSet) GetAll() []string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	result := make([]string, 0, len(f.flags))
	for flag := range f.flags {
		result = append(result, flag)
	}
	return result
}

func (f *BaseFlagSet) Clear() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.flags = make(map[string]bool)
}

// Clone creates a copy of the flag set
func (f *BaseFlagSet) Clone() *BaseFlagSet {
	f.mu.RLock()
	defer f.mu.RUnlock()

	clone := NewFlagSet()
	for flag := range f.flags {
		clone.flags[flag] = true
	}
	return clone
}

// Count returns number of set flags
func (f *BaseFlagSet) Count() int {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return len(f.flags)
}
