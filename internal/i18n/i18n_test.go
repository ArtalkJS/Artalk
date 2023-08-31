package i18n

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
