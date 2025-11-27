package utils

import (
	"context"
	"sync"
	"time"
)

// AnimationController manages a single animation ticker
type AnimationController struct {
	frame    int
	mu       sync.RWMutex
	cancel   context.CancelFunc
	ctx      context.Context
	onChange func() // Callback для обновления UI
}

// NewAnimationController creates a new animation controller
func NewAnimationController(onChange func()) *AnimationController {
	ctx, cancel := context.WithCancel(context.Background())
	ac := &AnimationController{
		frame:    0,
		cancel:   cancel,
		ctx:      ctx,
		onChange: onChange,
	}

	// Start ticker in goroutine
	go ac.run()

	return ac
}

// run starts the animation ticker
func (ac *AnimationController) run() {
	ticker := time.NewTicker(16 * time.Millisecond) // ~60 FPS
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ac.mu.Lock()
			ac.frame++
			ac.mu.Unlock()

			// Trigger UI update
			if ac.onChange != nil {
				ac.onChange()
			}

		case <-ac.ctx.Done():
			return
		}
	}
}

// GetFrame returns current frame number (thread-safe)
func (ac *AnimationController) GetFrame() int {
	ac.mu.RLock()
	defer ac.mu.RUnlock()
	return ac.frame
}

// Stop stops the animation
func (ac *AnimationController) Stop() {
	ac.cancel()
}
