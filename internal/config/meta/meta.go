package meta

import (
	"cmp"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/artalkjs/artalk/v2/internal/i18n"
	"github.com/goccy/go-yaml"
	"github.com/iancoleman/strcase"
	"github.com/knadh/koanf/maps"
)

const EnvPrefix = "ATK_"

type OptionsMeta struct {
	Title      string   `json:"title" yaml:"title"`
	Desc       string   `json:"desc" yaml:"desc"`
	Type       string   `json:"type" yaml:"type"`
	Options    []string `json:"options" yaml:"options"`
	Env        string   `json:"env" yaml:"env"`
	Path       string   `json:"path" yaml:"path"`
	PathText   string   `json:"path_title" yaml:"path_title"`
	Default    any      `json:"default" yaml:"default"`
	IsRoot     bool     `json:"is_root" yaml:"is_root"`
	HasChild   bool     `json:"has_child" yaml:"has_child"`
	AllowsSet  bool     `json:"allows_edit" yaml:"allows_edit"`
	CommentRaw string   `json:"comment_raw" yaml:"comment_raw"`
}

// GetOptionsMetaData returns the metadata of the config options from the YAML config template.
//
// The YAML config template includes comments and default values.
// The format of the comments is as follows:
//
//	# Title (Desc) [Option1, Option2]
//
// The metadata includes the title, description, type, options, environment variable name, path, and default value.
func GetOptionsMetaData(yamlTemplate string) ([]OptionsMeta, error) {
	// parse and get comments in yaml template
	var defaultConf map[string]interface{}
	commentMap := yaml.CommentMap{} // The commentMap is a map, the key has a prefix "$."
	if err := yaml.UnmarshalWithOptions([]byte(yamlTemplate), &defaultConf, yaml.CommentToMap(commentMap)); err != nil {
		return nil, err
	}

	metas := []OptionsMeta{}

	// The defaults is a flattened map of the default config values. (defaults["a.b.c"] = "value")
	// The leafNodePathMap is a flattened map of the paths of the leaf nodes. (leafNodePathMap["a.b.c"] = ["a", "b", "c"])
	defaults, leafNodePathMap := maps.Flatten(defaultConf, nil, ".")

	// Traverse all paths and generate metadata
	for path := range getPathsWithAllRootNodes(leafNodePathMap) {
		c, ok := commentMap["$."+path]
		if !ok || len(c) == 0 {
			c = newCommentNodeByPath(path) // if not found, create a new comment node
		}

		comment := strings.TrimSpace(strings.Join(c[0].Texts, ""))
		title, desc, options := extractComment(comment)

		hasChild := getNodeHasChild(leafNodePathMap, path)
		metas = append(metas, OptionsMeta{
			Title:      title,
			Desc:       desc,
			Type:       fmt.Sprintf("%T", defaults[path]),
			Options:    options,
			Env:        EnvPrefix + strings.ReplaceAll(strings.ToUpper(path), ".", "_"),
			Path:       path,
			PathText:   "",
			Default:    defaults[path],
			IsRoot:     !strings.Contains(path, "."),
			HasChild:   hasChild,
			AllowsSet:  !hasChild,
			CommentRaw: comment,
		})
	}

	metas = generatePathTextForMetas(metas)

	// sort by path alphabetically
	sort.Slice(metas, func(i, j int) bool { return metas[i].Path < metas[j].Path })

	return metas, nil
}

// Create a new comment node by path, using the last part of the path as the title.
func newCommentNodeByPath(path string) []*yaml.Comment {
	arr := strings.Split(path, ".")
	return []*yaml.Comment{{Texts: []string{
		strcase.ToCamel(arr[len(arr)-1]),
	}}}
}

// Check if a node has child nodes
//
// If passing a path, and the path with a dot suffix is found in the pathMap, it has child nodes.
func getNodeHasChild(pathMap map[string][]string, path string) bool {
	for k := range pathMap {
		if strings.HasPrefix(k, path+".") {
			return true
		}
	}
	return false
}

// Get the paths with all root nodes
//
// It will traverse all composite paths which are split by dots.
// For example, if the paths are "a.b.c", it will traverse ["a", "a.b", "a.b.c"].
// The all traversed paths will be append to the return map.
func getPathsWithAllRootNodes(pathMap map[string][]string) map[string][]string {
	temp := map[string][]string{}
	for k, v := range pathMap {
		temp[k] = v
	}
	for k := range temp {
		arr := strings.Split(k, ".")
		if len(arr) == 1 {
			continue
		}
		for i := 1; i < len(arr); i++ {
			pathMap[strings.Join(arr[:i], ".")] = arr[:i]
		}
	}
	return pathMap
}

func extractComment(comment string) (title string, desc string, options []string) {
	title = comment

	// Remove double-bar from title
	barReg := regexp.MustCompile(`(?m)(--.*?--)`)
	title = barReg.ReplaceAllString(title, "")

	// Extract Title and Description from comment
	descRe := regexp.MustCompile(`(?m)(\(.+?\))`)
	title = descRe.ReplaceAllStringFunc(title, func(s string) string {
		match := descRe.FindStringSubmatch(s)[1]
		if match != "" {
			desc = strings.Trim(match, "()")
			return ""
		}
		return s
	})

	// Extract Options from comment
	optionsRe := regexp.MustCompile(`(?m)(\[.+?\])`)
	title = optionsRe.ReplaceAllStringFunc(title, func(s string) string {
		match := optionsRe.FindStringSubmatch(s)[1]
		if match != "" {
			arr := strings.Split(strings.Trim(match, "[]"), ",")
			for _, v := range arr {
				options = append(options, strings.Trim(strings.TrimSpace(v), "\""))
			}
			return ""
		}
		return s
	})

	// Trim spaces
	title = strings.TrimSpace(title)
	desc = strings.TrimSpace(desc)

	// Localize the title
	if title == "Enabled" {
		title = i18n.T("Enabled")
	}

	return title, desc, options
}

func generatePathTextForMetas(metaSlice []OptionsMeta) []OptionsMeta {
	metas := make([]OptionsMeta, len(metaSlice))
	metaTitleMap := map[string]string{}
	for i, m := range metaSlice {
		metas[i] = m
		metaTitleMap[m.Path] = metas[i].Title
	}

	const sep = " > "
	const placeholder = "__"
	for i := range metas {
		split := strings.Split(metas[i].Path, ".")
		pathText := []string{}
		for j := 1; j <= len(split); j++ {
			pathText = append(pathText, cmp.Or(metaTitleMap[strings.Join(split[:j], ".")], placeholder))
		}
		metas[i].PathText = strings.Join(pathText, sep)
	}

	return metas
}
