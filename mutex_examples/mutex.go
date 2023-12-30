package mutexexamples

import (
	"sync"
)

// counter is a safe counter with mutex protection
type counter struct {
	mu    sync.RWMutex
	value int
}

// Increment increases the counter value safely
func (c *counter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

// Value returns the current value of the counter safely
func (c *counter) val() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.value
}
