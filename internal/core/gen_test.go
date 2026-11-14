package core

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/pkged"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestGenerateTemplatePreservesExistingFiles(t *testing.T) {
	root, err := filepath.Abs("../..")
	require.NoError(t, err)
	previousFS := pkged.FS()
	t.Cleanup(func() { pkged.SetFS(previousFS) })
	pkged.SetFS(os.DirFS(root))
	// Catch fatal generation errors without terminating the test process.
	restore := zap.ReplaceGlobals(zap.NewNop().WithOptions(zap.WithFatalHook(zapcore.WriteThenPanic)))
	t.Cleanup(restore)
	dir := t.TempDir()
	resource := "internal/template/email_tpl/default.html"
	want, err := os.ReadFile(filepath.Join(root, resource))
	require.NoError(t, err)
	Gen(resource, dir, false)
	output := filepath.Join(dir, "default.html")
	got, err := os.ReadFile(output)
	require.NoError(t, err)
	require.Equal(t, want, got)
	require.NoError(t, os.WriteFile(output, []byte("custom template"), 0600))
	require.Panics(t, func() { Gen(resource, dir, false) })
	got, err = os.ReadFile(output)
	require.NoError(t, err)
	require.Equal(t, "custom template", string(got))
	Gen(resource, dir, true)
	got, err = os.ReadFile(output)
	require.NoError(t, err)
	require.Equal(t, want, got)
	require.Panics(t, func() { Gen("missing-resource", dir, false) })
	require.Panics(t, func() { Gen(resource, filepath.Join(output, "child"), false) })
}
