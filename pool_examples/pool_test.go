package poolexamples

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestWorkerPoolExecute(t *testing.T) {
	tests := []struct {
		desc string
		n    int
	}{
		{
			desc: "n = 1",
			n:    1,
		},
		{
			desc: "n = 10",
			n:    10,
		},
		{
			desc: "n = 100",
			n:    100,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			pool := newPool(3)
			var callCount atomic.Int32
			for i := 0; i < tc.n; i++ {
				pool.execute(func() {
					callCount.Add(1)
					// simulate some work
					time.Sleep(10 * time.Millisecond)
				})
			}

			pool.wait()

			if callCount.Load() != int32(tc.n) {
				t.Errorf("Expected call count to be %d, got %d", tc.n, callCount.Load())
				return
			}

			// Assert that all tasks have been executed
			if pool.activeCount() != 0 {
				t.Errorf("Expected active workers count to be 0, got %d", pool.activeCount())
			}
		})
	}
}
