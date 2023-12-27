package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/mikepianka/drsz"
)

type Config struct {
	RootDir   string
	CreateCsv bool
	CsvPath   string
	ConcLimit uint8
}

// cli parses command line arguments into validated config.
func cli() Config {
	cfg := Config{}

	csvArg := flag.String("o", "none", "save output to the provided CSV filepath")
	concArg := flag.Int("c", 0, "number of concurrent directory size searches")
	flag.Parse()
	rootArg := flag.Arg(0)

	// validate root arg
	if rootArg == "" {
		log.Fatal(
			"You did not pass a directory argument.\n" +
				"Example usage: drsz ./some_directory\n" +
				"Example usage (creates output file): drsz -o ./results.csv ./some_directory\n" +
				"Example usage (5x concurrent searches): drsz -c 5 ./some_directory")
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

	if *concArg < 0 {
		// can't have negative conc limit, set to 0
		fmt.Println("WARNING: Negative concurrency limit was provided; setting to 0.")
		cfg.ConcLimit = 0
	} else if *concArg > 255 {
		// clamp to 0-255
		fmt.Println("WARNING: Large concurrency limit was provided; clamping to 255.")
		cfg.ConcLimit = 255
	} else {
		cfg.ConcLimit = uint8(*concArg)
	}

	return cfg
}

func main() {
	// collect input via CLI
	cfg := cli()

	// start timer
	now := time.Now()

	// run search
	err := drsz.Run(cfg.RootDir, cfg.ConcLimit, cfg.CreateCsv, cfg.CsvPath)
	if err != nil {
		log.Fatal(err)
	}

	// print elapsed time
	elapsed := time.Since(now)
	fmt.Printf("Completed in %s\n", elapsed)
}
