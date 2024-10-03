package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnsureDir(t *testing.T) {
	// Test creating a directory
	dir := "test_dir"
	err := EnsureDir(dir)
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	// Test creating an existing directory
	err = EnsureDir(dir)
	assert.NoError(t, err)

	// Test creating a nested directory
	nestedDir := "nested/test_dir"
	err = EnsureDir(nestedDir)
	assert.NoError(t, err)
	defer os.RemoveAll("nested")
}

func TestCheckFileExist(t *testing.T) {
	// Test existing file
	existingFile := "test_file_exist.txt"
	_, createErr := os.Create(existingFile)
	assert.NoError(t, createErr, "Error creating test file")
	defer os.Remove(existingFile)

	exists := CheckFileExist(existingFile)
	assert.True(t, exists, "Expected file to exist, but it doesn't")

	// Test non-existing file
	nonExistingFile := "non_existing.txt"
	exists = CheckFileExist(nonExistingFile)
	assert.False(t, exists, "Expected file to not exist, but it does")
}

func TestCheckDirExist(t *testing.T) {
	// Test existing directory
	existingDir := "test_dir_exist"
	err := os.Mkdir(existingDir, 0700)
	assert.NoError(t, err, "Error creating test directory")
	defer os.Remove(existingDir)

	exists := CheckDirExist(existingDir)
	assert.True(t, exists, "Expected directory to exist, but it doesn't")

	// Test non-existing directory
	nonExistingDir := "non_existing"
	exists = CheckDirExist(nonExistingDir)
	assert.False(t, exists, "Expected directory to not exist, but it does")
}
