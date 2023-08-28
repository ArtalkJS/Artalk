package main

import (
	"embed"
	"log"

	"github.com/ArtalkJS/Artalk/cmd"
	"github.com/ArtalkJS/Artalk/internal/pkged"
)

//go:embed public/*
//go:embed i18n/*
//go:embed conf/*
var embedFS embed.FS

func main() {
	pkged.SetFS(embedFS)

	app := cmd.New()

	if err := app.Launch(); err != nil {
		log.Fatal(err)
	}
}
