package utils

import "os"

// EnsureDir ensures that a target directory exists (like `mkdir -p`),
func EnsureDir(dir string) error {
	return os.MkdirAll(dir, 0700)
}

func CheckFileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
