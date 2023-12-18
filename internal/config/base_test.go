package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestI18nPatch(t *testing.T) {
	test := func(input string, expected string) {
		t.Run("locale="+input, func(t *testing.T) {
			config := &Config{Locale: input}
			config.i18nPatch()
			assert.Equal(t, expected, config.Locale)
		})
	}

	test("", "en")
	test("en", "en")
	test("en-US", "en")

	test("zh", "zh-CN")
	test("zh-CN", "zh-CN")

	t.Run("case convert", func(t *testing.T) {
		// @see https://www.techonthenet.com/js/language_tags.php
		test("zh-cn", "zh-CN")
		test("ZH-cn", "zh-CN")
		test("EN", "en")
	})
}
