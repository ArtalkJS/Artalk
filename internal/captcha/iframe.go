package captcha

import (
	"bytes"
	"embed"
	"html/template"
	"io"
)

//go:embed pages/*
var pages embed.FS

func RenderIFrame(htmlFilename string, params Map) ([]byte, error) {
	f, err := pages.Open("pages/" + htmlFilename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fb, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var tplBuf bytes.Buffer

	tpl := template.New("")
	tpl.Parse(string(fb))
	tpl.Execute(&tplBuf, params)

	return tplBuf.Bytes(), nil
}
