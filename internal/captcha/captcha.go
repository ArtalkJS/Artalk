package captcha

import (
	"embed"
	"io/fs"
)

//go:embed pages/*
var pages embed.FS

func GetPage(name string) (fs.File, error) {
	return pages.Open("pages/" + name)
}
