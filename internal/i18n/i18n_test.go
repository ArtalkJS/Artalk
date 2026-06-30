package i18n

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestLocaleCatalogKeys(t *testing.T) {
	localeDir := filepath.Join("..", "..", "i18n")
	expected := readLocaleKeys(t, filepath.Join(localeDir, "en.yml"))
	localeFiles, err := filepath.Glob(filepath.Join(localeDir, "*.yml"))
	if err != nil {
		t.Fatal(err)
	}

	for _, localeFile := range localeFiles {
		t.Run(filepath.Base(localeFile), func(t *testing.T) {
			assert.Equal(t, expected, readLocaleKeys(t, localeFile))
		})
	}
}

func readLocaleKeys(t *testing.T, filename string) []string {
	t.Helper()

	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}

	locale := map[string]string{}
	if err := yaml.Unmarshal(data, &locale); err != nil {
		t.Fatal(err)
	}

	keys := make([]string, 0, len(locale))
	for key := range locale {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func TestT(t *testing.T) {
	t.Run("ValidTranslate", func(t *testing.T) {
		Locales = map[string]string{
			"Hello World, {{name}}": "{{name}}：你好世界",
		}
		str := T("Hello World, {{name}}", map[string]interface{}{
			"name": "Kirito",
		})
		assert.Equal(t, "Kirito：你好世界", str)
	})

	t.Run("InvalidTranslate", func(t *testing.T) {
		Locales = map[string]string{}
		str := T("Hello World, {{name}}", map[string]interface{}{
			"name": "Kirito",
		})
		assert.Equal(t, "Hello World, Kirito", str)
	})

	t.Run("NoRenderParams", func(t *testing.T) {
		Locales = map[string]string{}
		str := T("Hello World")
		assert.Equal(t, "Hello World", str)
	})
}

func TestLoad(t *testing.T) {
	t.Run("DefaultLocale", func(t *testing.T) {
		Load("", func(locale string) ([]byte, error) {
			assert.Equal(t, locale, "en") // default is en locale
			return []byte("hello: Hello"), nil
		})
		assert.Equal(t, map[string]string{"hello": "Hello"}, Locales)
	})

	t.Run("ValidLocale", func(t *testing.T) {
		Load("zh-CN", func(locale string) ([]byte, error) {
			assert.Equal(t, locale, "zh-CN")
			return []byte("world: 世界"), nil
		})
		assert.Equal(t, map[string]string{"world": "世界"}, Locales)
	})

	t.Run("InvalidLocale", func(t *testing.T) {
		Load("xxxx", func(locale string) ([]byte, error) {
			if locale == "en" {
				return []byte("Bonjour: Hello"), nil
			}

			assert.Equal(t, locale, "xxxx")
			return nil, fmt.Errorf("locale file not found")
		})

		assert.Equal(t, map[string]string{"Bonjour": "Hello"}, Locales)
	})
}
