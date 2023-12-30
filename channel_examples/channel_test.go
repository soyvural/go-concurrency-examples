package channelexamples

import (
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFibonacci(t *testing.T) {
	tests := []struct {
		desc string
		n    int
		want int
	}{
		{
			desc: "n = 0",
			n:    0,
			want: 0,
		},
		{
			desc: "n = 5",
			n:    5,
			want: 5,
		},
		{
			desc: "n = 10",
			n:    10,
			want: 55,
		},
		{
			desc: "n = 15",
			n:    15,
			want: 610,
		},
		{
			desc: "n = 20",
			n:    20,
			want: 6765,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			c := make(chan int)
			go fibonacci(tc.n, c)

			if got := <-c; tc.want != got {
				t.Errorf("got fibonacci(%d) = %d; want %d", tc.n, got, tc.want)
			}
		})

	}
}

func TestFibonacciWithBufferedChannel(t *testing.T) {
	tests := []struct {
		desc       string
		n          int
		bufferSize int
		want       []int
	}{
		{
			desc:       "up to n = 5 and bufferSize = 3",
			n:          5,
			bufferSize: 3,
			want:       []int{0, 1, 1, 2, 3},
		},
		{
			desc:       "up to n = 5 and bufferSize = 10",
			n:          5,
			bufferSize: 10,
			want:       []int{0, 1, 1, 2, 3},
		},
		{
			desc:       "up to n = 10 and bufferSize = 10",
			n:          10,
			bufferSize: 10,
			want:       []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34},
		},
		{
			desc:       "up to n = 10 and bufferSize = 100",
			n:          10,
			bufferSize: 100,
			want:       []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			c := make(chan int, tc.bufferSize)
			go fibonacciWithBufferedChannel(tc.n, c)

			var result []int
			// range loop over the channel waits until the channel is closed
			for n := range c {
				result = append(result, n)
			}

			if diff := cmp.Diff(tc.want, result); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestFibonacciWithQuitChan(t *testing.T) {
	want := []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55}

	res := make(chan int)
	quit := make(chan int)

	go fibonacciWithQuitChan(res, quit)

	// read results one by one from unbuffered channe
	var result []int
	for i := 0; i < len(want); i++ {
		result = append(result, <-res)
	}
	// send quit signal
	quit <- 0

	if diff := cmp.Diff(want, result); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func TestWorker(t *testing.T) {
	tests := []struct {
		desc      string
		workerIDs []int
	}{
		{
			desc:      "size 2",
			workerIDs: []int{1, 2, 3},
		},
		{
			desc:      "size 5",
			workerIDs: []int{47, 48, 49, 50, 51},
		},
		{
			desc:      "size 10",
			workerIDs: []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			var mu sync.Mutex
			var gotWorkerIDs []int
			size := len(tc.workerIDs)

			signals := make([]chan struct{}, size)

			for i := 0; i < len(signals); i++ {
				signals[i] = make(chan struct{})
			}

			for i := 0; i < size; i++ {
				nextIndex := (i + 1) % size
				go worker(tc.workerIDs[i], signals[i], signals[nextIndex], func(id int) {
					mu.Lock()
					defer mu.Unlock()
					gotWorkerIDs = append(gotWorkerIDs, id)
				})
			}

			// let start the chain of workers
			signals[0] <- struct{}{}

			// wait to completed all workers
			<-signals[0]

			if diff := cmp.Diff(tc.workerIDs, gotWorkerIDs); diff != "" {
				t.Errorf("mismatch in worker IDs (-want +got):\n%s", diff)
			}
		})
	}
}
