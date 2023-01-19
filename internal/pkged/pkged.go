package pkged

import "embed"

var fs embed.FS

func SetFS(embedFs embed.FS) {
	fs = embedFs
}

func FS() embed.FS {
	return fs
}
