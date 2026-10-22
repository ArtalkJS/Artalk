package limiter

import (
	"sync"
	"time"
)

// RateLimiter limits admissions per key and across the instance in a fixed
// window. Unlike the captcha limiter, checking and consuming quota are atomic.
// At most totalLimit keys are retained; the next window discards all old keys.
// No background goroutine or per-key timers are needed.
type RateLimiter struct {
	mu                      sync.Mutex
	perKeyLimit, totalLimit int
	window                  time.Duration
	expires                 time.Time
	counts                  map[string]int
	total                   int
	now                     func() time.Time
}

func NewRateLimiter(perKeyLimit, totalLimit int, window time.Duration) *RateLimiter {
	if perKeyLimit <= 0 || totalLimit <= 0 || window <= 0 {
		panic("rate limiter requires positive limits and window")
	}
	return &RateLimiter{perKeyLimit: perKeyLimit, totalLimit: totalLimit, window: window, now: time.Now}
}

func (l *RateLimiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()
	if !now.Before(l.expires) {
		l.counts = make(map[string]int)
		l.total = 0
		l.expires = now.Add(l.window)
	}
	if l.total >= l.totalLimit || l.counts[key] >= l.perKeyLimit {
		return false
	}
	l.counts[key]++
	l.total++
	return true
}
