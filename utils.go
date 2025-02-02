package drsz

import (
	"path"
	"strings"
)

// IsCsvPath checks that the provided filepath is to a CSV.
func IsCsvPath(filepath string) bool {
	clean := path.Clean(filepath)
	ext := path.Ext(clean)
	return strings.ToLower(ext) == ".csv"
}
