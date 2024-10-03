package config

import (
	"cmp"
	"path/filepath"
	"testing"

	"github.com/adrg/xdg"
	"github.com/artalkjs/artalk/v2/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestI18nPatch(t *testing.T) {
	test := func(input string, expected string) {
		t.Run("locale="+input, func(t *testing.T) {
			config := &Config{Locale: input}
			config.i18nPatch()
			assert.Equal(t, expected, config.Locale)
		})
	}

	test("", "en")
	test("en", "en")
	test("en-US", "en")

	test("zh", "zh-CN")
	test("zh-CN", "zh-CN")

	t.Run("case convert", func(t *testing.T) {
		// @see https://www.techonthenet.com/js/language_tags.php
		test("zh-cn", "zh-CN")
		test("ZH-cn", "zh-CN")
		test("EN", "en")
	})
}

func mockCheckFileExist(mockFiles map[string]bool) func(string) bool {
	return func(filename string) bool {
		return mockFiles[filename]
	}
}

func TestRetrieveConfigFile(t *testing.T) {
	orgCheckFileExist := utils.CheckFileExist
	orgConfDefaultFilenames := CONF_DEFAULT_FILENAMES
	orgXDGConfigHome := xdg.ConfigHome
	defer func() {
		utils.CheckFileExist = orgCheckFileExist
		CONF_DEFAULT_FILENAMES = orgConfDefaultFilenames
		xdg.ConfigHome = orgXDGConfigHome
	}()

	tests := []struct {
		name          string
		alias         []string
		mockFiles     map[string]bool
		xdgConfigHome string
		expectedFile  string
	}{
		{
			name:  "work dir config",
			alias: []string{"artalk.yml"},
			mockFiles: map[string]bool{
				"artalk.yml":                    true,
				"data/artalk.yml":               true,
				"/home/user/.config/artalk.yml": true,
				"/etc/artalk/artalk.yml":        true,
			},
			expectedFile: "./artalk.yml",
		},
		{
			name:  "xdg config home",
			alias: []string{"artalk.yml"},
			mockFiles: map[string]bool{
				"/home/user/.config/artalk/artalk.yml": true,
				"/etc/artalk/artalk.yml":               true,
			},
			xdgConfigHome: "/home/user/.config",
			expectedFile:  "/home/user/.config/artalk/artalk.yml",
		},
		{
			name:  "linux packing version config (/etc/artalk/artalk.yml)",
			alias: []string{"artalk.yml"},
			mockFiles: map[string]bool{
				"/etc/artalk/artalk.yml": true,
			},
			expectedFile: "/etc/artalk/artalk.yml",
		},
		{
			name:  "multiple config files",
			alias: []string{"artalk.yml", "artalk.yaml"},
			mockFiles: map[string]bool{
				"artalk.yml":  false,
				"artalk.yaml": true,
			},
			expectedFile: "./artalk.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CONF_DEFAULT_FILENAMES = tt.alias
			utils.CheckFileExist = mockCheckFileExist(tt.mockFiles)
			xdg.ConfigHome = cmp.Or(tt.xdgConfigHome, orgXDGConfigHome)

			result, _ := filepath.Abs(RetrieveConfigFile())
			expected, _ := filepath.Abs(tt.expectedFile)
			assert.Equal(t, expected, result, "Expected to find the config file in %s", tt.expectedFile)
		})
	}
}

func TestRetrieveDataDir(t *testing.T) {
	orgCheckDirExist := utils.CheckDirExist
	orgXDGDataHome := xdg.DataHome
	defer func() {
		utils.CheckDirExist = orgCheckDirExist
		xdg.DataHome = orgXDGDataHome
	}()

	tests := []struct {
		name        string
		mockDirs    map[string]bool
		xdgDataHome string
		expectedDir string
	}{
		{
			name: "xdg data home",
			mockDirs: map[string]bool{
				"/home/user/.local/share/artalk": true,
			},
			xdgDataHome: "/home/user/.local/share",
			expectedDir: "/home/user/.local/share/artalk",
		},
		{
			name: "linux packing version data (/var/lib/artalk)",
			mockDirs: map[string]bool{
				"/home/user/.local/share/artalk": false,
				"/var/lib/artalk":                true,
			},
			expectedDir: "/var/lib/artalk",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			utils.CheckDirExist = func(dir string) bool {
				return tt.mockDirs[dir]
			}
			xdg.DataHome = cmp.Or(tt.xdgDataHome, orgXDGDataHome)

			result, _ := filepath.Abs(RetrieveDataDir())
			expected, _ := filepath.Abs(tt.expectedDir)
			assert.Equal(t, expected, result, "Expected to find the data dir in %s", tt.expectedDir)
		})
	}
}
