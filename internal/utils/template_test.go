package utils

import (
	"testing"
)

func TestRenderMustaches(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		data := "Hello, {{ name }}! Your age is {{ age }}."
		dict := map[string]interface{}{
			"name": "Alice",
			"age":  30,
		}

		expected := "Hello, Alice! Your age is 30."

		result := RenderMustaches(data, dict)
		if result != expected {
			t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
		}
	})

	t.Run("WithCustomValueGetter", func(t *testing.T) {
		data := "Hello, {{ name }}!"
		dict := map[string]interface{}{
			"name": "Bob",
		}

		valueGetter := func(key string, val interface{}) string {
			if key == "name" {
				return "Mr. " + val.(string)
			}
			return ""
		}

		expected := "Hello, Mr. Bob!"

		result := RenderMustaches(data, dict, valueGetter)
		if result != expected {
			t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
		}
	})

	t.Run("WithMissingKeys", func(t *testing.T) {
		data := "Hello, {{ name }}!"
		dict := map[string]interface{}{
			"age": 25,
		}

		expected := "Hello, {{ name }}!"

		result := RenderMustaches(data, dict)
		if result != expected {
			t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
		}
	})
}
