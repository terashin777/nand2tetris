package utils

import (
	"path/filepath"
	"strings"
)

type fp struct{}

var Filepath = fp{}

func (p fp) FileNameWithoutExt(fn string) string {
	return strings.TrimSuffix(fn, filepath.Ext(fn))
}
