package hook_test

import (
	"errors"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/hook"
	"github.com/stretchr/testify/require"
)

func TestTriggerRunsHandlersInOrderAndAppliesRemoval(t *testing.T) {
	var h hook.Hook[string]
	var calls []string
	regular := h.Add(func(value string) error {
		calls = append(calls, "regular:"+value)
		return nil
	})
	h.PreAdd(func(value string) error {
		calls = append(calls, "first:"+value)
		h.Remove(regular)
		return nil
	})

	require.NoError(t, h.Trigger("event", func(value string) error {
		calls = append(calls, "one-off:"+value)
		return nil
	}))
	require.Equal(t, []string{"first:event", "regular:event", "one-off:event"}, calls)
	calls = nil
	require.NoError(t, h.Trigger("next"))
	require.Equal(t, []string{"first:next"}, calls)
	h.RemoveAll()
	require.NoError(t, h.Trigger("cleared"))
	require.Equal(t, []string{"first:next"}, calls)
}

func TestTriggerStopsAfterHandlerError(t *testing.T) {
	failure := errors.New("handler failed")
	for _, tc := range []struct {
		name, result string
		err          error
	}{
		{"intentional stop", "", hook.ErrStopPropagation},
		{"handler failure", "handler failed", failure},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var h hook.Hook[struct{}]
			h.Add(func(struct{}) error { return tc.err })
			h.Add(func(struct{}) error {
				t.Fatal("handler after the stop should not run")
				return nil
			})
			err := h.Trigger(struct{}{})
			if tc.result == "" {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tc.err)
			}
		})
	}
}
