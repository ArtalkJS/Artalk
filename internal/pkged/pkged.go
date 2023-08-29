package pkged

import "io/fs"

var f fs.FS

func SetFS(embedFs fs.FS) {
	f = embedFs
}

func FS() fs.FS {
	return f
}
