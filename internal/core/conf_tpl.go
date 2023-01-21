package core

import (
	"fmt"
	"io"
	"strings"

	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/pkged"
)

// GetConfTpl is get config template file with specified locale
func GetConfTpl(localeParam ...string) string {
	var locale string
	if len(localeParam) > 0 {
		locale = strings.TrimSpace(localeParam[0])
	}
	if locale == "" && config.Instance != nil {
		locale = strings.TrimSpace(config.Instance.Locale)
	}

	const CONF_DEFAULT = "conf/artalk.example.yml"

	var filename string
	if locale == "en" || locale == "" {
		filename = CONF_DEFAULT
	} else {
		filename = fmt.Sprintf("conf/artalk.example.%s.yml", locale)
	}

	f, err := pkged.FS().Open(filename)
	if err != nil {
		// locale not found
		f, err = pkged.FS().Open(CONF_DEFAULT)
		if err != nil {
			panic(err)
		}
	}

	fileStr, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	return string(fileStr)
}
