package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/dustin/go-humanize"
	"github.com/schollz/progressbar/v3"
)

type Dir struct {
	Name       string
	FullPath   string
	SizeBytes  int64
	SizeString string
}

func dirSize(path string) int64 {
	var size int64

	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			size += info.Size()
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Error searching in %q: %v\n", path, err)
	}

	return size
}

func createOutputFile(dirs []Dir, filepath string) {
	csvFile, err := os.Create(filepath)

	if err != nil {
		log.Fatalf("Failed to create output file: %q: %v\n", filepath, err)
	}

	csvWriter := csv.NewWriter(csvFile)

	defer func() {
		csvWriter.Flush()
		csvFile.Close()
	}()

	header := []string{"directory", "bytes"}
	err = csvWriter.Write(header)

	if err != nil {
		log.Fatalf("Failed to write header to file: %q: %v\n", filepath, err)
	}

	for _, dir := range dirs {
		row := []string{dir.FullPath, fmt.Sprintf("%d", dir.SizeBytes)}
		err = csvWriter.Write(row)

		if err != nil {
			log.Fatalf("Failed to write row %q to file: %q: %v\n", dir.FullPath, filepath, err)
		}
	}

	fmt.Printf("Created output file: %q\n", filepath)
}

func main() {
	saveResults := flag.String("o", "none", "save output to the provided CSV filepath")

	flag.Parse()
	root := flag.Arg(0)

	if root == "" {
		log.Fatal(
			"You did not pass a directory argument.\n" +
				"Example usage: drsz ./some_directory\n" +
				"Example usage (creates output file): drsz -o ./results.csv ./some_directory")
	}

	fmt.Printf("Calculating top-level subdirectory sizes in %s*\n", root)

	contents, err := os.ReadDir(root)

	if err != nil {
		log.Fatal(err)
	}

	// find the top-level subdirectories
	var dirNames []string

	for _, item := range contents {
		if item.IsDir() {
			dirNames = append(dirNames, item.Name())
		}
	}

	// setup progress bar based on number of subdirectories
	bar := progressbar.Default(int64(len(dirNames)))

	// find the total size of each subdirectory
	var dirs []Dir

	for _, d := range dirNames {
		// filepath.Join ensures we get a valid path whether or not root has a trailing slash
		path := filepath.Join(root, d)
		size := dirSize(path)
		sizeReadable := humanize.Bytes(uint64(size))
		dir := Dir{d, path, size, sizeReadable}
		dirs = append(dirs, dir)
		bar.Add(1)
	}

	fmt.Printf("Search complete after finding %d top-level subdirectories:\n", len(dirNames))

	// print results
	for _, dir := range dirs {
		fmt.Printf("%s = %s\n", dir.Name, dir.SizeString)
	}

	// write optional output file
	if *saveResults != "none" {
		createOutputFile(dirs, *saveResults)
	}
}
