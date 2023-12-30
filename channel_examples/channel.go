package channelexamples

import "fmt"

// fibonacci produces the sequence of Fibonacci numbers, f(0=0, f(1)=1, f(n)=f(n-1)+f(n-2).
func fibonacci(n int, c chan int) {
	x, y := 0, 1
	for i := 0; i < n; i++ {
		x, y = y, x+y
	}
	c <- x
}

func fibonacciWithBufferedChannel(n int, c chan int) {
	x, y := 0, 1
	for i := 0; i < n; i++ {
		c <- x
		x, y = y, x+y
	}
	// close channel to indicate that no more values will be sent
	close(c)
}

func fibonacciWithQuitChan(res, quit chan int) {
	x, y := 0, 1
	for {
		select {
		case res <- x:
			x, y = y, x+y
		case <-quit:
			return
		}
	}
}

// worker is an example for how to execute Goroutines in ordered way.
// The workers execute in a chain of Goroutines.
func worker(id int, signal, nextSignal chan struct{}, processFn func(int)) {
	fmt.Printf("worker %d started.\n", id)

	// wait for signal
	<-signal

	// process data
	processFn(id)

	// signal next worker
	nextSignal <- struct{}{}
}
