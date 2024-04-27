package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	workDir    string
	scanDirs   string
	outputFile string
	updateDir  string
	originFile string
)

func init() {
	flag.StringVar(&workDir, "w", "", "specify work directory")
	flag.StringVar(&scanDirs, "d", "./internal", "specify which directory to scan")
	flag.StringVar(&outputFile, "o", "", "specify output filename")
	flag.StringVar(&updateDir, "u", "", "specify which directory to update")
	flag.StringVar(&originFile, "origin-file", "i18n/en.yml", "specify origin file such as `i18n/en.yml`")
}

func main() {
	flag.Parse()

	if workDir != "" {
		if err := os.Chdir(workDir); err != nil {
			panic(err)
		}
	}

	originFile, _ = filepath.Abs(originFile)

	// Initialize an empty map to store the keys
	keys := make(map[string]bool)

	for _, dir := range strings.Split(scanDirs, ",") {
		for _, k := range scan(strings.TrimSpace(dir)) {
			// Add the key to the map if it doesn't already exist
			if !keys[k] {
				keys[k] = true
			}
		}
	}

	// Use the sort package to sort the keys alphabetically
	var keysSlice []string
	for key := range keys {
		keysSlice = append(keysSlice, key)
	}
	sort.Strings(keysSlice)

	if outputFile != "" {
		saveToFile(outputFile, keysSlice)
	} else if updateDir != "" {
		updateDirectory(updateDir, keysSlice)
	}
}

func scan(dir string) []string {
	result := []string{}

	// Use the os package to search for Go files recursively
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			start := time.Now()

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			buf, err := io.ReadAll(file)
			if err != nil {
				return err
			}

			codeStr := string(buf)
			if !strings.Contains(codeStr, "i18n.T") {
				return nil
			}

			keys := extractFromCode(codeStr)

			if len(keys) > 0 {
				result = append(result, keys...)
				duration := time.Since(start)
				fmt.Printf("found %d msgs in `%s` (duration: %v)\n", len(keys), path, duration)
			}
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	return result
}

func extractFromCode(src string) []string {
	fileSet := token.NewFileSet()
	f, err := parser.ParseFile(fileSet, "", src, 0)
	if err != nil {
		panic(err)
	}

	var i18nKeySet []string
	ast.Inspect(f, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
			if fun.Sel.Name != "T" {
				return true
			}
			if ident, ok := fun.X.(*ast.Ident); ok {
				if ident.Name != "i18n" {
					return true
				}
				if lit, ok := call.Args[0].(*ast.BasicLit); ok {
					// string without quotes
					i18nKeySet = append(i18nKeySet, lit.Value[1:len(lit.Value)-1])
					return true
				}
			}
		}
		return true
	})

	return i18nKeySet
}

func saveToFile(filename string, keys []string) error {
	// write the keys to the YAML file
	content := ""
	for _, key := range keys {
		content += fmt.Sprintf("%s:\n", key)
	}
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		return err
	}
	return nil
}

func updateDirectory(directory string, keys []string) error {
	// check if the directory exists
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		return fmt.Errorf("directory `%s` does not exist", directory)
	}

	// walk through the directory, and find all yaml files
	return filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || (!strings.HasSuffix(path, ".yml") && !strings.HasSuffix(path, ".yaml")) {
			return err
		}

		defer func() {
			fmt.Printf("updated `%s`\n", path)
		}()

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		buf, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		// parse yaml and update the keys
		yamlMap := make(map[string]string)
		err = yaml.Unmarshal(buf, &yamlMap)
		if err != nil {
			return err
		}

		// write the yaml back to the file
		content := ""
		for _, key := range keys {
			value, ok := yamlMap[key]
			if !ok {
				pathAbs, _ := filepath.Abs(path)
				if pathAbs != originFile {
					value = "__MISSING__"
				}
			}
			content += fmt.Sprintf("%s: %s\n", strconv.Quote(key), strconv.Quote(value))
		}

		err = os.WriteFile(path, []byte(content), 0644)
		if err != nil {
			return err
		}

		return nil
	})
}
