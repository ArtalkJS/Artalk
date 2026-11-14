package main

import (
	"encoding/json"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/config/meta"
	"github.com/stretchr/testify/require"
)

func TestGenerateConfigMetadata(t *testing.T) {
	previousLocale, previousFormat, previousOutput := locale, format, outputFile
	t.Cleanup(func() { locale, format, outputFile = previousLocale, previousFormat, previousOutput })
	locale = "en"

	t.Run("Go source", func(t *testing.T) {
		format = "go"
		outputFile = filepath.Join(t.TempDir(), "metadata.go")
		main()
		source, err := os.ReadFile(outputFile)
		require.NoError(t, err)
		parsed, err := parser.ParseFile(token.NewFileSet(), outputFile, source, 0)
		require.NoError(t, err)
		require.Equal(t, "config", parsed.Name.Name)
		require.Contains(t, string(source), `"ATK_SITE_DEFAULT": "site_default"`)
		require.Contains(t, string(source), "var OptionsMetaCache")
	})

	t.Run("JSON for tooling", func(t *testing.T) {
		format = "json"
		outputFile = filepath.Join(t.TempDir(), "metadata.json")
		main()
		data, err := os.ReadFile(outputFile)
		require.NoError(t, err)
		var options []meta.OptionsMeta
		require.NoError(t, json.Unmarshal(data, &options))
		require.Condition(t, func() bool {
			for _, option := range options {
				if option.Env == "ATK_SITE_DEFAULT" && option.Path == "site_default" {
					return option.AllowsSet
				}
			}
			return false
		})
	})

	t.Run("documentation", func(t *testing.T) {
		format = "markdown"
		outputFile = filepath.Join(t.TempDir(), "env.md")
		const header = "# Configuration\n\n"
		const footer = "\n\nSee the deployment guide.\n"
		template := header + start + "\n\n## General options\n\n| Env | Default | Description | Key |\nOld table\n\n" + end + footer
		require.NoError(t, os.WriteFile(outputFile, []byte(template), 0600))
		main()
		data, err := os.ReadFile(outputFile)
		require.NoError(t, err)
		require.Contains(t, string(data), header+start)
		require.Contains(t, string(data), "## General options\n\n| Env | Default | Description | Key |")
		require.Contains(t, string(data), "| **ATK_SITE_DEFAULT** |")
		require.Contains(t, string(data), end+footer)
		require.NotContains(t, string(data), "Old table")
	})
}

func TestInvalidDocumentationMarkersLeaveFileUnchanged(t *testing.T) {
	for _, body := range []string{"No markers", end + "\n" + start, start + "\nIncomplete document"} {
		file := filepath.Join(t.TempDir(), "env.md")
		require.NoError(t, os.WriteFile(file, []byte(body), 0600))
		require.Panics(t, func() { parseMarkdownEnv(file) })
		require.Panics(t, func() { writeMarkdown(file, "New table") })
		contents, err := os.ReadFile(file)
		require.NoError(t, err)
		require.Equal(t, body, string(contents))
	}
	missing := filepath.Join(t.TempDir(), "missing.md")
	require.Panics(t, func() { parseMarkdownEnv(missing) })
	require.Panics(t, func() { writeMarkdown(missing, "New table") })
	require.NoFileExists(t, missing)
}
