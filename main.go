package main

import (
	libEmbed "embed"

	"github.com/ArtalkJS/ArtalkGo/cmd"
	"github.com/ArtalkJS/ArtalkGo/internal/pkged"
)

//go:embed frontend/*
//go:embed artalk-go.example.yml
var embedFS libEmbed.FS

func main() {
	pkged.SetFS(embedFS)
	cmd.Execute()
}
