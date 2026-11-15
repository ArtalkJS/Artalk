package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestGenerateTranslations(t *testing.T) {
	oldScan, oldOutput, oldUpdate, oldOrigin := scanDirs, outputFile, updateDir, originFile
	t.Cleanup(func() {
		scanDirs, outputFile, updateDir, originFile = oldScan, oldOutput, oldUpdate, oldOrigin
	})

	scanDirs = t.TempDir()
	source := `package example
func messages(key string) {
	i18n.T("World")
	i18n.T("Hello")
	i18n.T("Hello")
	i18n.T(key)
	other.T("Ignored")
	i18n.Other("Ignored too")
}
`
	require.NoError(t, os.WriteFile(filepath.Join(scanDirs, "messages.go"), []byte(source), 0600))
	require.NoError(t, os.WriteFile(filepath.Join(scanDirs, "unrelated.go"), []byte("package example\n"), 0600))

	outputFile = filepath.Join(t.TempDir(), "keys.yml")
	main()
	keys, err := os.ReadFile(outputFile)
	require.NoError(t, err)
	require.Equal(t, "Hello:\nWorld:\n", string(keys))

	outputFile, updateDir = "", t.TempDir()
	originFile = filepath.Join(updateDir, "en.yml")
	localeFile := filepath.Join(updateDir, "fr.yaml")
	require.NoError(t, os.WriteFile(originFile, []byte("Hello: Hello\n"), 0600))
	require.NoError(t, os.WriteFile(localeFile, []byte("Hello: Bonjour\nOld: Obsolète\n"), 0600))
	main()

	for path, expected := range map[string]map[string]string{
		originFile: {"Hello": "Hello", "World": ""},
		localeFile: {"Hello": "Bonjour", "World": "__MISSING__"},
	} {
		data, err := os.ReadFile(path)
		require.NoError(t, err)
		var messages map[string]string
		require.NoError(t, yaml.Unmarshal(data, &messages))
		require.Equal(t, expected, messages)
	}
}
