package auth

import (
	"strings"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/markbates/goth"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type AuthProviderInfo struct {
	Name  string `json:"name" validate:"required"`
	Label string `json:"label" validate:"required"`
	Path  string `json:"path" validate:"optional"`
	Icon  string `json:"icon" validate:"required"`
}

func GetProviderInfo(conf *config.Config, providers []goth.Provider) []AuthProviderInfo {
	var info []AuthProviderInfo

	// Email
	if conf.Auth.Email.Enabled {
		info = append(info, AuthProviderInfo{
			Name:  "email",
			Label: "Email",
			Icon:  GetProviderIconBase64("email"),
		})
	}

	for _, provider := range providers {
		name := strings.ToLower(provider.Name())
		title := cases.Title(language.Und, cases.NoLower).String(name)
		// Patch for Microsoft
		if title == "Microsoftonline" {
			title = "Microsoft"
		}
		info = append(info, AuthProviderInfo{
			Name:  name,
			Label: title,
			Path:  "/api/v2/auth/" + name,
			Icon:  GetProviderIconBase64(name),
		})
	}

	// Skip
	if conf.Auth.Anonymous {
		info = append(info, AuthProviderInfo{
			Name:  "skip",
			Label: "Skip",
			Icon:  GetProviderIconBase64("skip"),
		})
	}

	return info
}
