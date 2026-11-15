package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/adrg/xdg"
	"github.com/artalkjs/artalk/v2/internal/pkged"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveConfigFileBeforeDataDir(t *testing.T) {
	originalDir, err := os.Getwd()
	require.NoError(t, err)
	originalDataHome := xdg.DataHome
	defer func() {
		require.NoError(t, os.Chdir(originalDir))
		xdg.DataHome = originalDataHome
	}()

	invocationDir := t.TempDir()
	dataHome := t.TempDir()
	invocationDir, err = filepath.EvalSymlinks(invocationDir)
	require.NoError(t, err)
	dataHome, err = filepath.EvalSymlinks(dataHome)
	require.NoError(t, err)
	dataDir := filepath.Join(dataHome, "artalk")
	require.NoError(t, os.Mkdir(dataDir, 0o755))
	require.NoError(t, os.WriteFile(
		filepath.Join(invocationDir, "artalk.yml"),
		[]byte("timezone: UTC\n"),
		0o600,
	))
	require.NoError(t, os.Chdir(invocationDir))
	xdg.DataHome = dataHome

	cfgFile, err := resolveConfigFileBeforeDataDir("", "")
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(invocationDir, "artalk.yml"), cfgFile)

	workDir, err := initDataDir("")
	require.NoError(t, err)
	assert.Equal(t, dataDir, workDir)

	currentDir, err := os.Getwd()
	require.NoError(t, err)
	assert.Equal(t, dataDir, currentDir)

	conf, err := getConfig(cfgFile)
	require.NoError(t, err)
	assert.Equal(t, "UTC", conf.TimeZone)
}

func TestResolveConfigFileWithExplicitPaths(t *testing.T) {
	originalDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { require.NoError(t, os.Chdir(originalDir)) }()

	invocationDir := t.TempDir()
	invocationDir, err = filepath.EvalSymlinks(invocationDir)
	require.NoError(t, err)
	require.NoError(t, os.Chdir(invocationDir))

	cfgFile, err := resolveConfigFileBeforeDataDir("custom.yml", "")
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(invocationDir, "custom.yml"), cfgFile)

	cfgFile, err = resolveConfigFileBeforeDataDir("custom.yml", "./runtime")
	require.NoError(t, err)
	assert.Equal(t, "custom.yml", cfgFile)
}

func TestLaunchReadOnlyCommands(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	previousArgs := os.Args
	t.Cleanup(func() { os.Args = previousArgs })
	cfg := filepath.Join(dir, "artalk.yml")
	require.NoError(t, os.WriteFile(cfg, []byte("site_default: CLI fixture\ntimezone: Local\n"), 0600))
	for _, args := range [][]string{nil, {"version"}, {"--version"}, {"config"}, {"unknown-command"}} {
		os.Args = append([]string{"artalk", "--workdir", dir, "--config", cfg}, args...)
		app := New()
		app.RootCmd.SetArgs(os.Args[1:])
		require.NoError(t, app.Launch())
		require.Nil(t, app.App, "informational commands should not start a database")
	}
}

func TestDefaultDataDirectoryAndConfigGeneration(t *testing.T) {
	root, err := filepath.Abs("..")
	require.NoError(t, err)
	oldFS := pkged.FS()
	t.Cleanup(func() { pkged.SetFS(oldFS) })
	pkged.SetFS(os.DirFS(root))
	originalDataHome, originalConfigHome := xdg.DataHome, xdg.ConfigHome
	t.Cleanup(func() { xdg.DataHome = originalDataHome; xdg.ConfigHome = originalConfigHome })
	xdg.DataHome = t.TempDir()
	xdg.ConfigHome = t.TempDir()
	t.Chdir(t.TempDir())
	workdir, err := initDataDir("")
	require.NoError(t, err)
	require.Equal(t, filepath.Join(xdg.DataHome, "artalk"), workdir)
	conf, err := getConfig("")
	require.NoError(t, err)
	require.NotEmpty(t, conf.AppKey)
	again, err := getConfig("")
	require.NoError(t, err)
	require.Equal(t, conf.AppKey, again.AppKey, "do not regenerate an existing signing key")
	_, err = initDataDir(filepath.Join(t.TempDir(), "missing"))
	require.ErrorContains(t, err, "working directory change error")
}
