package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path"
	"regexp"
	"runtime"
	"strings"

	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/config/meta"
	"github.com/ArtalkJS/Artalk/internal/pkged"
	"github.com/samber/lo"
)

var (
	locale string
	format string
	file   string
)

func init() {
	flag.StringVar(&locale, "l", "zh-cn", "specify locale")
	flag.StringVar(&format, "format", "markdown", "specify output format")
	flag.StringVar(&file, "f", "", "specify save file")
}

func main() {
	flag.Parse()

	// change WorkDir run on project root
	{
		_, filename, _, _ := runtime.Caller(0)
		rootDir := path.Join(path.Dir(filename), "../../../..")
		if err := os.Chdir(rootDir); err != nil {
			panic(err)
		}
	}

	// load assets fs
	dirFS := os.DirFS("./")
	pkged.SetFS(dirFS)

	confTemplate := config.Template(locale)
	metas, err := meta.GetOptionsMetaData(confTemplate)
	if err != nil {
		panic(err)
	}

	// output metadata
	result := ""
	if format == "json" {
		result = getJSON(metas)
	} else if format == "markdown" {
		result = getMarkdown(metas)
	} else {
		panic("invalid output format")
	}

	if file == "" {
		fmt.Println(result)
	} else {
		// read file content
		content, err := os.ReadFile(file)
		if err != nil {
			panic(err)
		}

		// find `<!-- env-variables --><!-- /env-variables -->`
		start := "<!-- env-variables -->"
		end := "<!-- /env-variables -->"
		re := regexp.MustCompile(start + `(?s:(.*?))` + end)
		matches := re.FindAllStringSubmatch(string(content), -1)
		if len(matches) == 0 {
			panic("<!-- env-variables --> not found")
		}

		// replace
		newContent := strings.ReplaceAll(string(content), matches[0][0], start+"\n\n"+result+"\n\n"+end)

		// write file
		if err := os.WriteFile(file, []byte(newContent), 0644); err != nil {
			panic(err)
		}

		fmt.Println("File updated successfully")
	}
}

func getJSON(metas []meta.OptionsMeta) string {
	str, err := json.MarshalIndent(metas, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(str)
}

func getMarkdown(metas []meta.OptionsMeta) string {
	md := ""

	// Get root nodes
	rootNodes := []meta.OptionsMeta{
		{
			Title:  "通用配置",
			Desc:   "通用配置项",
			Path:   "general",
			IsRoot: true,
		},
	}
	rootNodes = append(rootNodes, lo.Filter(metas, func(m meta.OptionsMeta, _ int) bool {
		return m.IsRoot && m.HasChild
	})...)

	printGroup := func(root meta.OptionsMeta, items []meta.OptionsMeta) {
		md += fmt.Sprintf("## %s\n\n", root.Title)
		md += fmt.Sprintln("| 环境变量 | 默认值 | 描述 | 路径 |")
		md += fmt.Sprintln("| --- | --- | --- | --- |")
		for _, m := range items {
			if m.Title == "Enabled" {
				m.Title = "启用"
			}
			options := lo.If(len(m.Options) > 0, "[可选：`"+strings.Join(m.Options, ", ")+"`] ").Else("")
			desc := lo.If(m.Desc != "", "("+m.Desc+") ").Else("")
			md += fmt.Sprintf("| **%s** | `%v` | %s %s%s| %s (%s) |\n", m.Env, m.Default, m.Title, desc, options, m.Path, m.PathText)
		}
		md += "\n\n"
	}

	for _, root := range rootNodes {
		printGroup(root, lo.Filter(metas, func(m meta.OptionsMeta, _ int) bool {
			if root.Path == "general" {
				return m.AllowsSet && m.IsRoot
			} else {
				return m.AllowsSet && !m.IsRoot && strings.HasPrefix(m.Path, root.Path)
			}
		}))
	}

	return strings.TrimSpace(md)
}
