package main

import (
	"flag"
	"log"
	"os"
	"path"

	"github.com/mikepianka/drsz"
)

type Config struct {
	RootDir   string
	CreateCsv bool
	CsvPath   string
}

// cli parses command line arguments into validated config.
func cli() Config {
	cfg := Config{}

	csvArg := flag.String("o", "none", "save output to the provided CSV filepath")
	flag.Parse()
	rootArg := flag.Arg(0)

	// validate root arg
	if rootArg == "" {
		log.Fatal(
			"You did not pass a directory argument.\n" +
				"Example usage: drsz ./some_directory\n" +
				"Example usage (creates output file): drsz -o ./results.csv ./some_directory")
	}

	cfg.RootDir = path.Clean(rootArg)

	info, err := os.Stat(cfg.RootDir)
	if err != nil {
		log.Fatal(err)
	}

	if !info.IsDir() {
		log.Fatal("provided root path is not a directory")
	}

	// validate csv arg
	if *csvArg != "none" {
		// csv arg was provided
		cfg.CsvPath = path.Clean(*csvArg)
		cfg.CreateCsv = true

		if !drsz.IsCsvPath(cfg.CsvPath) {
			log.Fatal("provided output path is not a CSV filename")
		}

		_, err := os.Stat(cfg.CsvPath)
		if err == nil {
			log.Fatalf("output file already exists at %s", cfg.CsvPath)
		}
	}

	return cfg
}

func main() {
	cfg := cli()

	root, err := drsz.NewRootDir(cfg.RootDir)
	if err != nil {
		log.Fatal(err)
	}

	err = root.FindTops()
	if err != nil {
		log.Fatal(err)
	}

	err = root.CalcStats()
	if err != nil {
		log.Fatal(err)
	}

	if !cfg.CreateCsv {
		return
	}

	err = root.ExportCSV(cfg.CsvPath)
	if err != nil {
		log.Fatal(err)
	}
}
