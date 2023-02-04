package i18n

import (
	"fmt"

	"github.com/ArtalkJS/Artalk/internal/pkged"
	"github.com/ArtalkJS/Artalk/internal/utils"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

//go:generate go run ./gen -w ../../ -d internal,server,cmd -o i18n/en.yml

var Locales map[string]string

func Init(locale string) {
	if locale == "" { // default lang
		locale = "en"
	}

	yamlStr, err := pkged.FS().ReadFile(fmt.Sprintf("i18n/%s.yml", locale))
	if err != nil {
		logrus.Warn("invalid locale config please check, now it is set to `en`")
		Init("en")
		return
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
