package limiter

import "path"

func isProtectPath(pathTarget string, protectPaths []string) bool {
	for _, p := range protectPaths {
		if path.Clean(pathTarget) == path.Clean(p) {
			return true
		}
	}
	return false
}
