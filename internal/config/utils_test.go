package config

import (
	"fmt"
	"strings"
	"testing"
)

func TestGetConfigEnvNameMapping(t *testing.T) {
	m := GetConfigEnvNameMapping()

	// print the mapping with pretty format
	for k, v := range m {
		fmt.Printf("%s: %s\n", k, v)
	}

	containsArrayFiledInEnv := false
	containsArrayFiledInPath := false
	supportOldVersion := false // contains double underscore in env name

	for env, path := range m {
		if strings.HasSuffix(path, ".$$") && !strings.HasSuffix(env, "_$$") {
			t.Errorf("if path suffix is `.$$`, the env should also have `_$$` suffix, but got %s", env)
		}
		if !strings.HasPrefix(env, "ATK_") {
			t.Errorf("env should start with ATK_, but got %s", env)
		}
		if strings.Contains(env, "_$$") {
			containsArrayFiledInEnv = true
		}
		if strings.Contains(path, ".$$") {
			containsArrayFiledInPath = true
		}
		if strings.Contains(env, "__") { // double underscore
			supportOldVersion = true

			// check if path same
			path2 := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(
				strings.ReplaceAll(strings.TrimPrefix(env, "ATK_"),
					"__", "#"), "_", "."), "#", "_"))
			if path2 != path {
				t.Errorf("old version env name: %s, should map to path: %s, but got %s", env, path, path2)
			}
		}
	}

	if !supportOldVersion {
		t.Error("should support old version env name (eg. ATK_IP__REGION_DB__PATH), but not found")
	}
	if !containsArrayFiledInEnv {
		t.Error("should contains array field env name (eg. ATK_TRUSTED_DOMAINS_$$), but not found")
	}
	if !containsArrayFiledInPath {
		t.Error("should contains array field path (eg. trusted_domains.$$), but not found")
	}
}
