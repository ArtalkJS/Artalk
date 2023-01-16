package main

import (
	"embed"

	"github.com/ArtalkJS/ArtalkGo/cmd"
	"github.com/ArtalkJS/ArtalkGo/internal/pkged"
)

//go:embed public/*
//go:embed artalk-go.example.yml
var embedFS embed.FS

func main() {
	pkged.SetFS(embedFS)
	cmd.Execute()
}
