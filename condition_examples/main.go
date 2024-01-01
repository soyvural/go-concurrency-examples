package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

const (
	resourceLimit = 10
	numGoroutines = 5
)

// waitForResource waits until the resourceCounter reaches the resourceLimit
func waitForResource(id int, cond *sync.Cond, resourceCounter *atomic.Int32, wg *sync.WaitGroup) {
	defer wg.Done()
	cond.L.Lock()
	defer cond.L.Unlock()
	if resourceCounter.Load() < resourceLimit {
		fmt.Printf("ID: %d is waiting.\n", id)
		cond.Wait()
	}
	fmt.Printf("ID: %d is running.\n", id)
}

// increaseResource increases the resourceCounter and broadcasts when it reaches the resourceLimit
func increaseResource(cond *sync.Cond, resourceCounter *atomic.Int32) {
	cond.L.Lock()
	defer cond.L.Unlock()
	for i := 0; i < resourceLimit; i++ {
		if resourceCounter.Add(1) == resourceLimit {
			cond.Broadcast()
			fmt.Printf("resourceCounter is %d, broadcast.\n", resourceCounter.Load())
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	var mu sync.Mutex
	cond := sync.NewCond(&mu)
	resourceCounter := &atomic.Int32{}
	wg := sync.WaitGroup{}

	for i := 0; i < numGoroutines; i++ {
		i := i
		wg.Add(1)
		go waitForResource(i, cond, resourceCounter, &wg)
	}
	go increaseResource(cond, resourceCounter)

	wg.Wait()

}
