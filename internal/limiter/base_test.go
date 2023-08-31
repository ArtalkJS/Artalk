package limiter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLimiter(t *testing.T) {
	const myIP = "127.0.0.1"

	assertPassAndLog := func(t *testing.T, limiter *Limiter) {
		if assert.True(t, limiter.IsPass(myIP)) {
			limiter.Log(myIP)
		}
	}
	assertPrevent := func(t *testing.T, limiter *Limiter) {
		for i := 0; i < 100; i++ {
			assert.False(t, limiter.IsPass(myIP))
		}
	}

	t.Run("NeverReset", func(t *testing.T) {
		limiter := NewLimiter(&LimiterConf{
			AlwaysMode:          false,
			MaxActionDuringTime: 2,
			ResetTimeout:        0, // never reset
		})

		// comment add *1
		assertPassAndLog(t, limiter)

		// comment add *2
		assertPassAndLog(t, limiter)

		// comment add *n (n > 2)
		for i := 0; i < 100; i++ {
			assertPrevent(t, limiter)
		}

		time.Sleep(1 * time.Second)

		// comment add after 1s sleep
		assertPrevent(t, limiter)

		t.Run("User Verify Pass", func(t *testing.T) {
			assertPrevent(t, limiter)

			limiter.MarkVerifyPassed(myIP)
			assertPassAndLog(t, limiter)
			assert.False(t, limiter.IsPass(myIP), "need verify immediately after user pass verify and commented")
		})

		t.Run("User Verify Fail", func(t *testing.T) {
			assertPrevent(t, limiter)

			limiter.MarkVerifyFailed(myIP)
			assertPrevent(t, limiter)
		})
	})

	t.Run("ResetAfter1s", func(t *testing.T) {
		limiter := NewLimiter(&LimiterConf{
			AlwaysMode:          false,
			MaxActionDuringTime: 2,
			ResetTimeout:        1, // reset after 1s
		})

		// comment add *1
		assertPassAndLog(t, limiter)

		// comment add *2
		assertPassAndLog(t, limiter)

		// comment add *3 immediately
		assertPrevent(t, limiter)

		// delay 1s
		time.Sleep(1 * time.Second)

		// comment add *4 after delay 1s
		assertPassAndLog(t, limiter)
		assertPassAndLog(t, limiter)

		assertPrevent(t, limiter)
	})

	t.Run("AlwaysMode", func(t *testing.T) {
		limiter := NewLimiter(&LimiterConf{
			AlwaysMode: true,
		})

		// always mode let use verify immediately
		assertPrevent(t, limiter)

		// test mark verify passed
		limiter.MarkVerifyPassed(myIP)

		assertPassAndLog(t, limiter)

		for i := 0; i < 100; i++ {
			// only true once before logged
			assertPrevent(t, limiter)
		}

		// test mark verify field
		for i := 0; i < 100; i++ {
			limiter.MarkVerifyFailed(myIP)
			assertPrevent(t, limiter)
			limiter.Log(myIP)
		}

		// test mark verify passed after fail marked
		limiter.MarkVerifyPassed(myIP)
		assertPassAndLog(t, limiter)

		for i := 0; i < 100; i++ {
			// only true once before logged
			assertPrevent(t, limiter)
		}
	})
}
