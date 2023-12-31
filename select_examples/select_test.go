package selectexamples

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestReceive(t *testing.T) {
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
			r := &receiver{
				drainTimeout: 100 * time.Millisecond,
			}
			c := make(chan string, tc.n)
			stop := make(chan struct{})
			done := make(chan struct{})

			go r.receive(c, stop, done)

			var messages []string
			for i := 0; i < tc.n; i++ {
				msg := fmt.Sprintf("msg-%d", i)
				c <- msg
				messages = append(messages, msg)
			}

			stop <- struct{}{}
			<-done

			gotMessages := r.getMessages()

			// short-cut if messages sizes mismatch so no need to make a detail comaprison
			if len(messages) != len(gotMessages) {
				t.Errorf("messages size mismatch: want %d, got %d", len(messages), len(gotMessages))
				return
			}

			if diff := cmp.Diff(messages, r.messages); diff != "" {
				t.Errorf("mismatch messages(-want +got):\n%s", diff)
			}
		})
	}
}
