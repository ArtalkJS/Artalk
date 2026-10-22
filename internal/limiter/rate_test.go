package limiter

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRateLimiterWindowAndMemory(t *testing.T) {
	now := time.Unix(100, 0)
	l := NewRateLimiter(2, 3, time.Minute)
	l.now = func() time.Time { return now }
	require.True(t, l.Allow("a"))
	require.True(t, l.Allow("a"))
	require.False(t, l.Allow("a"))
	require.True(t, l.Allow("b")) // A rejected attempt must not consume global quota.
	for i := range 1000 {
		require.False(t, l.Allow(fmt.Sprint(i)))
	}
	require.Len(t, l.counts, 2) // Rejected keys allocate no state.
	now = now.Add(time.Minute)
	require.True(t, l.Allow("c"))
	require.Len(t, l.counts, 1)
	require.Equal(t, 1, l.total)
}

func TestRateLimiterConcurrent(t *testing.T) {
	for _, tc := range []struct {
		name                string
		perKey, total, want int
		unique              bool
	}{
		{"per key", 10, 50, 10, false},
		{"instance", 10, 50, 50, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			l := NewRateLimiter(tc.perKey, tc.total, time.Minute)
			now := time.Unix(100, 0)
			l.now = func() time.Time { return now }
			var passed atomic.Int32
			var wg sync.WaitGroup
			for i := range 200 {
				wg.Go(func() {
					key := "shared"
					if tc.unique {
						key = fmt.Sprint(i)
					}
					if l.Allow(key) {
						passed.Add(1)
					}
				})
			}
			wg.Wait()
			require.EqualValues(t, tc.want, passed.Load())
		})
	}
}
