package utils

import "os"

// EnsureDir ensures that a target directory exists (like `mkdir -p`),
func EnsureDir(dir string) error {
	return os.MkdirAll(dir, 0700)
}
