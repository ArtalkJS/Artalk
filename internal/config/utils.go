package config

import (
	"strings"

	"github.com/artalkjs/artalk/v2/internal/utils"
)

// Get the hash function according to the frontend gravatar.params configuration
func GetHashFuncByFrontendConf(conf *Config) func(string) string {
	if gravatar, ok := conf.Frontend["gravatar"].(map[string]interface{}); ok {
		if params, ok := gravatar["params"].(string); ok {
			if strings.Contains(params, "sha256=1") {
				return utils.GetSha256Hash
			}
		}
	}
	return utils.GetMD5Hash
}
