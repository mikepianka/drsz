package drsz

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/dustin/go-humanize"
)

// Dir holds information about a directory.
type Dir struct {
	// AbsPath is the absolute path to the directory.
	AbsPath string
	// SizeBytes is the total size of the directory's contents in bytes.
	SizeBytes int64
	// LastModified is the most recent modification time of any file in the directory's contents.
	LastModified time.Time
}

// Name returns the name of the directory.
func (d Dir) Name() string {
	return path.Base(d.AbsPath)
}

// SizeString returns the size of the directory as a human readable string.
func (d Dir) SizeString() string {
	return humanize.Bytes(uint64(d.SizeBytes))
}

// WalkCalc recursively walks through the directory, calculating its total size and the most recent file modification time.
func (d *Dir) WalkCalc() error {
	var size int64
	var lastMod time.Time

	err := filepath.Walk(d.AbsPath, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			// found a file
			// add to total size
			size += info.Size()
			// set last modified time if it's more recent
			mod := info.ModTime()
			if mod.After(lastMod) {
				lastMod = mod
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("error while searching in %s: %v", d.AbsPath, err)
	}

	d.SizeBytes = size
	d.LastModified = lastMod
	return nil
}

// NewDir returns a pointer to a new Dir initialized with dirPath.
func NewDir(dirPath string) (*Dir, error) {
	absPath, err := resolveDirPath(dirPath)
	if err != nil {
		return nil, err
	}

	return &Dir{AbsPath: absPath}, nil
}

// resolveDirPath resolves the absolute path of the provided dirPath and confirms it is an accessible directory.
func resolveDirPath(dirPath string) (string, error) {
	abs, err := filepath.Abs(dirPath)
	if err != nil {
		return "", err
	}

	info, err := os.Stat(abs)
	if err != nil {
		return "", err
	}

	// no issues reading path, make sure it's a dir
	if !info.IsDir() {
		return "", fmt.Errorf("provided path is not a directory")
	}

	// path exists
	return abs, nil
}
