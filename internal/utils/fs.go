package utils

import "os"

// EnsureDir ensures that a target directory exists (like `mkdir -p`),
func EnsureDir(dir string) error {
	return os.MkdirAll(dir, 0700)
}

var CheckFileExist = func(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

var CheckDirExist = func(path string) bool {
	return CheckFileExist(path)
}
