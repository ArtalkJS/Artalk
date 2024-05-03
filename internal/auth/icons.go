package auth

import (
	"embed"
	"encoding/base64"
)

//go:embed icons/*
var iconsFS embed.FS

func GetProviderIconBase64(provider string) string {
	buf, err := iconsFS.ReadFile("icons/" + provider + ".svg")
	if err != nil {
		return ""
	}
	return "data:image/svg+xml;base64," + base64.StdEncoding.EncodeToString(buf)
}
