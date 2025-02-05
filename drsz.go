package drsz

import (
	"encoding/csv"
	"fmt"
	"os"
	"path"
	"sync"
	"text/tabwriter"

	"github.com/schollz/progressbar/v3"
)

// RootDir holds information about the top level directories it contains.
type RootDir struct {
	Dir
	TopDirs []*Dir
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
		row := csvRow(dir)
		err = csvWriter.Write(row)

		if err != nil {
			return fmt.Errorf("failed to write row %q to file: %s: %v", dir.AbsPath, csvPath, err)
		}
	}

	fmt.Printf("Exported CSV file %s\n", csvPath)
	return nil
}

// csvRow returns a row of directory info for a CSV file.
func csvRow(dir *Dir) []string {
	size := fmt.Sprintf("%d", dir.SizeBytes)
	timestamp := dir.LastModified.Local().String()
	return []string{dir.AbsPath, size, timestamp}
}

// FindTops finds the top level directories within the provided root dir.
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

// CalcStats calculates the top level directory stats for the provided root dir by recursively walking through.
func (r *RootDir) CalcStats(concLimit uint8) error {
	bar := progressbar.NewOptions64(int64(len(r.TopDirs)), progressbar.OptionSetDescription("Calculating..."), progressbar.OptionSetPredictTime(true), progressbar.OptionShowCount()) // setup progress bar based on number of dirs
	var wg sync.WaitGroup                                                                                                                                                             // setup wait group for tracking dir calc worker progress
	var mu sync.Mutex                                                                                                                                                                 // setup mutex to protect errors slice
	var errors []error                                                                                                                                                                // slice to hold any errors encountered

	if concLimit == 0 {
		concLimit = 1 // if concLimit is zero, only run goroutines one at a time
	}

	// Implement semaphore to limit concurrency
	sem := make(chan struct{}, concLimit) // concLimit is the max number of concurrent goroutines

	for _, d := range r.TopDirs {
		wg.Add(1) // increment wait group
		go func(d *Dir) {
			defer wg.Done()   // decrement wait group once work complete
			sem <- struct{}{} // acquire a concurrency token when performing intensive i/o
			err := d.WalkCalc()
			if err != nil {
				mu.Lock()
				errors = append(errors, err) // collect error
				mu.Unlock()
			}
			fmt.Printf("Processed %s\n", d.AbsPath)
			<-sem      // release token
			bar.Add(1) // increment progress bar
		}(d)
	}

	wg.Wait() // wait for goroutines to finish

	// print results using tabwriter
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', 0)
	// add blank row
	fmt.Fprintln(tw, "")
	// add header
	fmt.Fprintf(tw, "Name\tSize\tLast_Modified\n")
	// add info
	for _, d := range r.TopDirs {
		fmt.Fprintf(tw, "%s\t%s\t%s\n", d.Name(), d.SizeString(), d.LastModified.Local().String())
	}

	// print errors and their associated directories
	if len(errors) != 0 {
		fmt.Fprintln(tw, "")
		fmt.Fprintln(tw, "***WARN*** Processing failed on the following directories ***WARN***")
		for i, err := range errors {
			fmt.Fprintf(tw, "%d: %v\n", i+1, err)
		}
	}

	tw.Flush()

	return nil
}

// NewRootDir returns a pointer to a new RootDir initialized with dirPath.
func NewRootDir(dirPath string) (*RootDir, error) {
	absPath, err := resolveDirPath(dirPath)
	if err != nil {
		return nil, err
	}
	r := &RootDir{Dir: Dir{AbsPath: absPath}}
	return r, nil
}

// Run will execute a drsz search of the provided root dir; optionally creating an output file with the results.
func Run(rootDir string, concLimit uint8, createFile bool, outputFile string) error {
	// initialize root directory
	root, err := NewRootDir(rootDir)
	if err != nil {
		return err
	}

	// find the top-level dirs within root dir
	err = root.FindTops()
	if err != nil {
		return err
	}

	// calculate stats for each top-level dir
	err = root.CalcStats(concLimit)
	if err != nil {
		return err
	}

	// export CSV if requested
	if createFile {
		err = root.ExportCSV(outputFile)
		if err != nil {
			return err
		}
	}

	return nil
}
