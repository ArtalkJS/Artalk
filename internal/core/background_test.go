package core

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type backgroundCleanupProbe struct{ finished *atomic.Bool }

func (*backgroundCleanupProbe) Init() error { return nil }
func (p *backgroundCleanupProbe) Dispose() error {
	if !p.finished.Load() {
		return fmt.Errorf("service disposed while background work is running")
	}
	return nil
}

func TestResetWaitsForBackgroundWork(t *testing.T) {
	var finished atomic.Bool
	app := &App{service: &map[string]Service{"probe": &backgroundCleanupProbe{finished: &finished}}}
	release := make(chan struct{})
	var once sync.Once
	defer once.Do(func() { close(release) })
	require.True(t, app.RunBackground(func() { <-release; finished.Store(true) }))
	reset := make(chan error, 1)
	go func() { reset <- app.ResetBootstrapState() }()
	require.Eventually(t, func() bool {
		app.backgroundMu.Lock()
		defer app.backgroundMu.Unlock()
		return app.backgroundStopped
	}, time.Second, time.Millisecond)
	require.False(t, app.RunBackground(func() { t.Error("work accepted during cleanup") }))
	select {
	case err := <-reset:
		t.Fatalf("cleanup returned before work finished: %v", err)
	default:
	}
	once.Do(func() { close(release) })
	select {
	case err := <-reset:
		require.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("cleanup did not finish")
	}
	require.False(t, app.RunBackground(func() { t.Error("work accepted after cleanup") }))
}

func TestBackgroundAdmissionDuringReset(t *testing.T) {
	app := &App{}
	var accepted, finished atomic.Int32
	var callers sync.WaitGroup
	start := make(chan struct{})
	for range 64 {
		callers.Go(func() {
			<-start
			if app.RunBackground(func() { finished.Add(1) }) {
				accepted.Add(1)
			}
		})
	}
	close(start)
	require.NoError(t, app.ResetBootstrapState())
	callers.Wait()
	require.Equal(t, accepted.Load(), finished.Load())
}
