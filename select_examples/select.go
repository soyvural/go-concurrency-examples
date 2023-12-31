package selectexamples

import (
	"context"
	"sync"
	"time"
)

type receiver struct {
	messages []string
	mu       sync.RWMutex
	// drainTimeout is only used to drain remaining messages in the channel
	drainTimeout time.Duration
}

func NewReceiver(ctx context.Context, drainTimeout time.Duration) *receiver {
	return &receiver{
		drainTimeout: drainTimeout,
	}
}

func (r *receiver) add(msg string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.messages = append(r.messages, msg)
}

func (r *receiver) getMessages() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.messages
}

func (r *receiver) receive(c chan string, stop, stopped chan struct{}) {
	for {
		select {
		case msg := <-c:
			r.add(msg)
			// do some processing here
			time.Sleep(time.Millisecond)
		case <-stop:
			r.drain(c)

			stopped <- struct{}{}
			return
		}
	}
}

// drain drains the messages in the channel
// always think draining remaning tasks, messsages when to process is about stop in asynchronous programming
func (r *receiver) drain(c chan string) {
	for {
		select {
		case msg := <-c:
			r.add(msg)
		case <-time.After(r.drainTimeout):
			return
		}
	}
}
