package renderer

import (
	"sync"
	"time"
)

// Debouncer batches rapid calls into single delayed call
// Useful for preventing excessive re-renders
type Debouncer struct {
	delay    time.Duration
	timer    *time.Timer
	callback func()
	mu       sync.Mutex
}

// NewDebouncer creates a new debouncer
// delay: minimum time between callback invocations
// callback: function to call after delay
func NewDebouncer(delay time.Duration, callback func()) *Debouncer {
	return &Debouncer{
		delay:    delay,
		callback: callback,
	}
}

// Call triggers the debounced callback
// If called multiple times rapidly, only the last call executes after delay
func (d *Debouncer) Call() {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Stop existing timer if any
	if d.timer != nil {
		d.timer.Stop()
	}

	// Start new timer
	d.timer = time.AfterFunc(d.delay, func() {
		d.callback()
	})
}

// CallImmediate calls callback immediately and resets debounce timer
func (d *Debouncer) CallImmediate() {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Stop existing timer
	if d.timer != nil {
		d.timer.Stop()
		d.timer = nil
	}

	// Call immediately
	d.callback()
}

// Stop cancels pending debounced call
func (d *Debouncer) Stop() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
		d.timer = nil
	}
}
