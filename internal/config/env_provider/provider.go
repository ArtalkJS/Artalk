// Package env implements a koanf.Provider that reads environment
// variables as conf maps.
package env_provider

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"
)

// Env implements an environment variables provider.
type Env struct {
	// Prefix is used to filter environment variables from the `os.Environ()` list.
	prefix string

	// The map of environment variable names to their corresponding config paths.
	//
	// The symbol ".$$" is used to represent numbers in the path (array index, simple type elem).
	//
	// For example, the config path `moderator.keywords.files.0` will be converted to
	// `ATK_MODERATOR_KEYWORDS_FILES_$$`.
	// It in map form is `moderator.keywords.files.$$`: `ATK_MODERATOR_KEYWORDS_FILES_$$`
	//
	// May appears ".$$." to represent nested struct elem array, eg. `admin_users.$$.badge_name`
	// It in map form is `admin_users.$$.badge_name`: `ATK_ADMIN_USERS_$$_BADGE_NAME`
	envPathMap map[string]string

	osEnvironFunc func() []string
}

// Provider returns an environment variables provider
// that reads environment variables with the given prefix.
func Provider(prefix string, envPathMap map[string]string) *Env {
	e := &Env{
		prefix:        prefix,
		envPathMap:    handleEnvPathMap(envPathMap),
		osEnvironFunc: func() []string { return os.Environ() },
	}

	return e
}

// ReadBytes is not supported by the env provider.
func (e *Env) ReadBytes() ([]byte, error) {
	return nil, errors.New("env provider does not support this method")
}

// Read reads all available environment variables into a key:value map
// and returns it.
func (e *Env) Read() (map[string]interface{}, error) {
	// Collect the environment variable keys.
	var keys []string

	for _, k := range e.osEnvironFunc() {
		if strings.HasPrefix(k, e.prefix) {
			keys = append(keys, k)
		}
	}

	simpleElemArrayPaths := getSimpleElemArrayPaths(e.envPathMap)

	mp := make(map[string]interface{})
	for _, k := range keys {
		parts := strings.SplitN(k, "=", 2)
		key, value := parts[0], parts[1]
		if key == "" {
			continue
		}

		// Retrieve the path from the envPathMap.
		keyForPath, numbersInPath := getKeyInEnvPathMap(key)
		hasNumbersInPath := len(numbersInPath) > 0

		path, ok := e.envPathMap[keyForPath]
		if !ok {
			continue
		}

		// Handle get the final path and value
		var (
			finalPath  string = path
			finalValue any    = value
		)

		shouldConvertArray := slices.Contains(simpleElemArrayPaths, path)
		if shouldConvertArray {
			// Only convert simple elem array, eg. "trusted_domains.$$"
			// Convert "ATK_TRUSTED_DOMAINS=1 2 3" to "trusted_domains": ["1", "2", "3"]

			valueSlice := strings.Split(value, " ")

			for i, v := range valueSlice {
				// keep the original numbers
				// for example, "admin_users.$$.xx_items.$$" to "admin_users.1.xx_items.2"
				numbers := append([]string{}, numbersInPath...)
				numbers = append(numbers, fmt.Sprint(i))
				mp[recoverNumbersInPath(path, numbers)] = v
			}

			continue
		}

		if hasNumbersInPath {
			// Recover the numbers in path, eg. "admin_users.$$.badge_name" to "admin_users.1.badge_name"
			finalPath = recoverNumbersInPath(path, numbersInPath)
		}

		mp[finalPath] = finalValue
	}

	return Unflatten(mp), nil
}

func handleEnvPathMap(envPathMap map[string]string) map[string]string {
	// in case for convenience, to support both "ATK_TRUSTED_DOMAINS_0=1" and "ATK_TRUSTED_DOMAINS=1 2 3"
	// append "ATK_TRUSTED_DOMAINS_$$" and "ATK_TRUSTED_DOMAINS" to envPathMap if not exists if path "trusted_domains.$$" exists
	m := make(map[string]string)
	for env, path := range envPathMap {
		if strings.HasSuffix(path, ".$$") {
			envWithoutSuffix := strings.TrimSuffix(env, "_$$")
			m[envWithoutSuffix] = path
			m[envWithoutSuffix+"_$$"] = path
		}
	}
	for k, v := range m {
		envPathMap[k] = v
	}
	return envPathMap
}

// Get the key in envPathMap by envName and the numbers in envName
// It will replace the number in envName with "$$" and return the numbers
func getKeyInEnvPathMap(envName string) (keyInMap string, numbers []string) {
	numbers = []string{}
	reg := regexp.MustCompile(`_(\d+)`)
	envName = reg.ReplaceAllStringFunc(envName, func(s string) string {
		matchNum := reg.FindStringSubmatch(s)[1]
		if matchNum != "" {
			numbers = append(numbers, matchNum)
			return "_$$"
		}
		return s
	})
	return envName, numbers
}

// Recover the numbers in path, eg. "admin_users.$$.badge_name" to "admin_users.1.badge_name"
//
// The input pathWithPlaceholder is the path with placeholder, eg. "admin_users.$$.badge_name"
// The input numbersInPath is the numbers in path, eg. ["1"]
// The return value is the path with numbers, eg. "admin_users.1.badge_name"
func recoverNumbersInPath(pathWithPlaceholder string, numbersInPath []string) string {
	i := 0
	reg := regexp.MustCompile(`.(\$\$)`)
	path := reg.ReplaceAllStringFunc(pathWithPlaceholder, func(s string) string {
		if reg.FindStringSubmatch(s)[1] != "" {
			if i >= len(numbersInPath) {
				return s
			}
			r := "." + numbersInPath[i]
			i++
			return r
		}
		return s
	})
	return path
}

// Get the simple elem array paths from the envPathMap
//
// Simple elem array is: the elem of array with simple type such as
// string, int, bool, etc. eg. "trusted_domains.$$"
//
// The un-simple elem array is the elem array with nested struct,
// eg. "admin_users.$$.badge_name"
func getSimpleElemArrayPaths(m map[string]string) []string {
	paths := make(map[string]bool) // set to avoid duplicate
	for _, path := range m {
		// if has ".$$" suffix is not a struct elem array (simple type elem), eg. "trusted_domains.$$", "admin_users.$$.emails.$$"
		// if only ".$$." with double dot around dollar sign in path (not suffix) is a nested struct elem array, eg. "admin_users.$$.badge_name"
		if strings.HasSuffix(path, ".$$") {
			paths[path] = true
		}
	}
	pathSlice := make([]string, 0, len(paths))
	for k := range paths {
		pathSlice = append(pathSlice, k)
	}
	return pathSlice
}
