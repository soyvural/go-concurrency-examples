package mutexexamples

import (
	"sync"
	"testing"
)

func TestCounterIncrement(t *testing.T) {
	var wg sync.WaitGroup
	expected := 1000
	counter := &counter{}

	for i := 0; i < expected; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}
	wg.Wait()

	if counter.val() != expected {
		t.Errorf("Expected counter value %d, got %d", expected, counter.val())
	}
}
