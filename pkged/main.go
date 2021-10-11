package pkged

import (
	"github.com/markbates/pkger"
	"github.com/markbates/pkger/pkging"
)

func Open(p string) (pkging.File, error) {
	return pkger.Open(p)
}
