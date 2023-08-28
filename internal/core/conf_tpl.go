package core

import (
	"strings"

	"github.com/ArtalkJS/Artalk/internal/config"
)

func (app *App) ConfTpl() string {
	if app.Conf() == nil {
		return config.Template("en")
	}
	return config.Template(strings.TrimSpace(app.Conf().Locale))
}
