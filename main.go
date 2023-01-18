package main

import (
	"embed"

	"github.com/ArtalkJS/Artalk/cmd"
	"github.com/ArtalkJS/Artalk/internal/pkged"
)

//go:embed public/*
//go:embed i18n/*
//go:embed artalk.example.yml
var embedFS embed.FS

func main() {
	pkged.SetFS(embedFS)
	cmd.Execute()
}
