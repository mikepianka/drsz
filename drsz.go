package drsz

import (
	"encoding/csv"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/schollz/progressbar/v3"
)

// NewDir returns a pointer to a new Dir initialized with dirPath.
func NewDir(dirPath string) (*Dir, error) {
	d := &Dir{}
	err := d.SetPath(dirPath)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// Dir holds information about a directory.
type Dir struct {
	AbsPath      string
	SizeBytes    int64
	LastModified time.Time
}

// SizeString returns the size of the directory as a human readable string.
func (d Dir) SizeString() string {
	return humanize.Bytes(uint64(d.SizeBytes))
}

// Name returns the name of the directory.
func (d Dir) Name() string {
	return path.Base(d.AbsPath)
}

// SetPath resolves an absolute path, confirms it is an accessible directory, and sets it in the struct.
func (d *Dir) SetPath(dirPath string) error {
	abs, err := filepath.Abs(dirPath)
	if err != nil {
		return err
	}

	info, err := os.Stat(abs)
	if err != nil {
		return err
	}

	// no issues reading path, make sure it's a dir
	if !info.IsDir() {
		return fmt.Errorf("provided path is not a directory")
	}

	// path exists
	d.AbsPath = abs
	return nil
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

// ExportCSV creates an output CSV file containing directory information at the provided path.
func (r RootDir) ExportCSV(csvPath string) error {
	if !IsCsvPath(csvPath) {
		return fmt.Errorf("provided filepath is not a CSV file")
	}

	csvFile, err := os.Create(csvPath)

	if err != nil {
		return fmt.Errorf("failed to create output file: %s: %v", csvPath, err)
	}

	csvWriter := csv.NewWriter(csvFile)

	defer func() {
		// write buff to file and close it before completion
		csvWriter.Flush()
		csvFile.Close()
	}()

	header := []string{"directory", "bytes", "lastModified"}
	err = csvWriter.Write(header)

	if err != nil {
		return fmt.Errorf("failed to write header to file: %s: %v", csvPath, err)
	}

	for _, dir := range r.TopDirs {
		row := []string{dir.AbsPath, fmt.Sprintf("%d", dir.SizeBytes), dir.LastModified.Local().String()}
		err = csvWriter.Write(row)

		if err != nil {
			return fmt.Errorf("failed to write row %q to file: %s: %v", dir.AbsPath, csvPath, err)
		}
	}

	fmt.Printf("Exported CSV file %s\n", csvPath)
	return nil
}

type RootDir struct {
	Dir
	TopDirs []*Dir
}

func (r *RootDir) FindTops() error {
	contents, err := os.ReadDir(r.AbsPath)
	if err != nil {
		return err
	}

	var topDirs []*Dir
	for _, item := range contents {
		if item.IsDir() {
			dirPath := path.Join(r.AbsPath, item.Name())
			d, err := NewDir(dirPath)
			if err != nil {
				return err
			}
			topDirs = append(topDirs, d)
		}
	}

	r.TopDirs = topDirs
	fmt.Printf("Found %d top level directories in %s\n", len(r.TopDirs), r.AbsPath)

	return nil
}

func (r *RootDir) CalcStats() error {
	bar := progressbar.Default(int64(len(r.TopDirs))) // setup progress bar based on number of dirs
	var wg sync.WaitGroup                             // setup wait group for tracking dir calc worker progress
	var mu sync.Mutex                                 // setup mutex to protect errors slice
	var errors []error                                // slice to hold any errors encountered
	MAX_CONCURR := 4

	// Implement semaphore to limit concurrency
	sem := make(chan struct{}, MAX_CONCURR) // MAX_CONCURR is the max number of concurrent goroutines

	for i, d := range r.TopDirs {
		wg.Add(1) // increment wait group
		go func(d *Dir, i int) {
			defer wg.Done()             // decrement wait group once work complete
			sem <- struct{}{}           // acquire a concurrency token when performing intensive i/o
			time.Sleep(5 * time.Second) // add a synthetic wait to simulate work
			_ = d.WalkCalc()
			<-sem      // release token
			if i > 2 { // switch with err != nil when done testing
				mu.Lock()
				errors = append(errors, fmt.Errorf("uh oh %d", i)) // collect error
				mu.Unlock()
			}
			bar.Add(1) // increment progress bar
		}(d, i)
	}

	wg.Wait() // wait for goroutines to finish

	if len(errors) != 0 {
		// errors encountered, just return first one for simplicity for now
		return fmt.Errorf("encountered %d errors, the first being: %v", len(errors), errors[0])
	}

	// print results
	for _, d := range r.TopDirs {
		fmt.Printf("%s = %s\n", d.Name(), d.SizeString())
	}

	return nil
}

// IsCsvPath checks that the provided filepath is to a CSV.
func IsCsvPath(filepath string) bool {
	clean := path.Clean(filepath)
	ext := strings.ToLower(path.Ext(clean))

	if ext != ".csv" {
		return false
	}

	return true
}

// NewRootDir returns a pointer to a new RootDir initialized with dirPath.
func NewRootDir(dirPath string) (*RootDir, error) {
	r := &RootDir{}
	err := r.SetPath(dirPath)
	if err != nil {
		return nil, err
	}
	return r, nil
}
