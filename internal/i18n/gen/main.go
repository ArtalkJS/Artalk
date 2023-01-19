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
	"strings"
	"time"
)

var (
	workDir    string
	scanDirs   string
	outputFile string
)

func init() {
	flag.StringVar(&workDir, "w", "", "specify work directory")
	flag.StringVar(&scanDirs, "d", "./internal", "specify which directory to scan")
	flag.StringVar(&outputFile, "o", "./i18n/en.yml", "specify output filename")
}

func main() {
	flag.Parse()

	if workDir != "" {
		if err := os.Chdir(workDir); err != nil {
			panic(err)
		}
	}

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

	// Use the ioutil package to write the keys to the YAML file
	content := ""
	for _, key := range keysSlice {
		content += fmt.Sprintf("%s:\n", key)
	}
	err := os.WriteFile(outputFile, []byte(content), 0644)
	if err != nil {
		panic(err)
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

			extKeys := extractFromCode(codeStr)

			if len(extKeys) > 0 {
				result = append(result, extKeys...)
				duration := time.Since(start)
				fmt.Printf("found %d msgs in `%s` (duration: %v)\n", len(extKeys), path, duration)
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
					i18nKeySet = append(i18nKeySet, lit.Value)
					return true
				}
			}
		}
		return true
	})

	return i18nKeySet
}
