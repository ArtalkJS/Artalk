package config

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestLocalizedConfigTemplatesContainDefaultKeys(t *testing.T) {
	confDir := filepath.Join("..", "..", "conf")
	defaultKeys := readTemplateKeys(t, filepath.Join(confDir, "artalk.example.yml"))
	localizedTemplates, err := filepath.Glob(filepath.Join(confDir, "artalk.example.*.yml"))
	if err != nil {
		t.Fatal(err)
	}

	for _, template := range localizedTemplates {
		t.Run(filepath.Base(template), func(t *testing.T) {
			localizedKeys := readTemplateKeys(t, template)
			missing := []string{}
			for key := range defaultKeys {
				if _, ok := localizedKeys[key]; !ok {
					missing = append(missing, key)
				}
			}
			sort.Strings(missing)
			if len(missing) != 0 {
				t.Errorf("localized config template is missing keys: %v", missing)
			}
		})
	}
}

func readTemplateKeys(t *testing.T, filename string) map[string]struct{} {
	t.Helper()

	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}

	var template map[string]any
	if err := yaml.Unmarshal(data, &template); err != nil {
		t.Fatal(err)
	}

	keys := map[string]struct{}{}
	collectTemplateKeys(template, "", keys)
	return keys
}

func collectTemplateKeys(value map[string]any, prefix string, keys map[string]struct{}) {
	for key, child := range value {
		path := key
		if prefix != "" {
			path = prefix + "." + key
		}

		if nested, ok := child.(map[string]any); ok {
			collectTemplateKeys(nested, path, keys)
			continue
		}
		keys[path] = struct{}{}
	}
}
