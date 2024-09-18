package config

import (
	"fmt"
	"io"

	"github.com/artalkjs/artalk/v2/internal/pkged"
)

const CONF_TPL_DEFAULT_PATH = "conf/artalk.example.yml"

// Template is get config template file with specified locale
func Template(locale string) string {
	var filename string
	if locale == "en" || locale == "" {
		filename = CONF_TPL_DEFAULT_PATH
	} else {
		filename = fmt.Sprintf("conf/artalk.example.%s.yml", locale)
	}

	f, err := pkged.FS().Open(filename)
	if err != nil {
		// locale not found
		f, err = pkged.FS().Open(CONF_TPL_DEFAULT_PATH)
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
