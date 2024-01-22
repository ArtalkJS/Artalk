package template

import (
	"bytes"
	"embed"
	"fmt"
	"os"
)

// -------------------------------------------------------------------
//  Template Loader
// -------------------------------------------------------------------

type TemplateLoader interface {
	Load(tplType RenderType) string
}

// Creates a new default template loader
//
// tplName parameter is the template name,
// if tplName is empty, it will load the default template,
// or you can specify a custom template name.
// it will load the template from internal or external, first look up external, then internal.
func NewFileLoader(tplName string) TemplateLoader {
	return &DefaultTemplateLoader{tplName: tplName}
}

// The default template loader
var _ TemplateLoader = (*DefaultTemplateLoader)(nil)

type DefaultTemplateLoader struct {
	tplName string
}

func (l *DefaultTemplateLoader) Load(tplType RenderType) string {
	// retrieve external template
	if tpl, err := getExternalTpl(l.tplName); err == nil {
		return tpl
	}

	// retrieve internal template
	return getInternalTpl(tplType, l.tplName)
}

// -------------------------------------------------------------------
//  Internal Template Loader
// -------------------------------------------------------------------

//go:embed email_tpl/*
//go:embed notify_tpl/*
var internalTpl embed.FS

// Template base path for different render type
var tplType2BasePath = map[RenderType]string{
	TYPE_EMAIL:  "email_tpl",
	TYPE_NOTIFY: "notify_tpl",
}

// Get internal template
func getInternalTpl(tplType RenderType, tplName string) string {
	if tplName == "" {
		tplName = "default"
	}

	var basePath string
	if path, ok := tplType2BasePath[tplType]; ok {
		basePath = path
	} else {
		basePath = tplType2BasePath[TYPE_EMAIL] // default
	}

	filename := fmt.Sprintf("%s/%s.html", basePath, tplName)
	f, err := internalTpl.Open(filename)
	if err != nil {
		return ""
	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(f); err != nil {
		return ""
	}
	contents := buf.String()

	return contents
}

// -------------------------------------------------------------------
//  External Template Loader
// -------------------------------------------------------------------

// Get external template
func getExternalTpl(filename string) (string, error) {
	// TODO 反复文件 IO 操作会导致性能下降，之后优化可以改成程序启动时加载模板文件到内存中
	// TODO 安全问题：注意这里允许离开工作目录读取文件，
	//      由于 filename 是动态可配置项，若管理员凭证泄露，可能会造成文件泄露
	//		（需注意检测 filename 的合法性，建议程序在低权限账户下运行，或丢到 docker 运行）
	buf, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}
