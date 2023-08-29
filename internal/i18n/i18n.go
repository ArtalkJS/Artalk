package i18n

import (
	"fmt"
	"io"

	"github.com/ArtalkJS/Artalk/internal/log"
	"github.com/ArtalkJS/Artalk/internal/pkged"
	"github.com/ArtalkJS/Artalk/internal/utils"
	"gopkg.in/yaml.v3"
)

//go:generate go run ./gen -w ../../ -d internal,server,cmd -o i18n/en.yml

var Locales map[string]string

func Init(locale string) {
	if locale == "" { // default lang
		locale = "en"
	}

	read := func(l string) ([]byte, error) {
		file, err := pkged.FS().Open(fmt.Sprintf("i18n/%s.yml", l))
		if err != nil {
			return nil, err
		}
		return io.ReadAll(file)
	}

	yamlStr, err := read(locale)
	if err != nil {
		log.Warn("invalid locale config please check, now it is set to `en`")
		yamlStr, _ = read("en")
	}

	yaml.Unmarshal(yamlStr, &Locales)
}

func T(msg string, params ...map[string]interface{}) string {
	v, ok := Locales[msg]
	if !ok || v == "" {
		v = msg
	}
	return msgParams(v, params...)
}

func msgParams(msg string, params ...map[string]interface{}) string {
	if len(params) > 0 {
		return utils.RenderMustaches(msg, params[0])
	}

	return msg
}
