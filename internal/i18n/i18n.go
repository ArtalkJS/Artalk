package i18n

import (
	"github.com/artalkjs/artalk/v2/internal/log"
	"github.com/artalkjs/artalk/v2/internal/utils"
	"gopkg.in/yaml.v3"
)

//go:generate go run ./gen -w ../../ -d internal,server,cmd -u i18n/

var Locales map[string]string

func Load(locale string, localeProvider func(locale string) ([]byte, error)) {
	if locale == "" { // default lang
		locale = "en"
	}

	yamlStr, err := localeProvider(locale)
	if err != nil {
		log.Warn("invalid locale config please check, now it is set to `en`")
		yamlStr, _ = localeProvider("en")
	}

	clear(Locales) // clear map before unmarshal

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
