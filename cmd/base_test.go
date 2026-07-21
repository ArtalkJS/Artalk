package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/adrg/xdg"
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
