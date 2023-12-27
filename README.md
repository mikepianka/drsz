# drsz

A small CLI tool for finding the size of subdirectories.

## About

Given a directory path to begin searching from, **drsz** will calculate the total size in bytes of each top-level subdirectory. Handy if you want to quickly find what directories might be using up alot of storage space on your drive.

## Usage

At the command line run `./drsz ./some_directory`. **drsz** will recursively search through each top-level subdirectory in `./some_directory`, add up the file sizes of the contents, and then print the total size to the console.

For example:

```
>>> drsz ./spruce
Found 3 top level directories in /Users/mike/spruce
Calculating... 100% |████████████████████████████████████████| (3/3)
Name     Size       Last_Modified
.git     62 kB      2023-11-30 22:03:12.990783711 -0500 EST
bin      4.9 MB     2023-12-01 04:13:18.166776944 -0500 EST
docs     47 kB      2023-11-30 21:53:10.93936325 -0500 EST
Completed in 1.88025ms
```

### Usage Options

Default usage which outputs results to the terminal:
`drsz ./some_directory`

Write results to an output file:
`drsz -o ./results.csv ./some_directory`

Speed up the program with 5x concurrent directory searches. Using this option will increase the read load on the drive being searched.
`drsz -c 5 ./some_directory`

## Installation

Mac, Linux, and Windows binaries have been pre-built and are available in Releases. Simply download one and follow the Usage instructions above. Alternatively, clone the repo and run `make` to build from source.

At the command line you can either run the executable by providing the full filepath to it, or make an alias in your path to the binary to create a shortened `drsz` command for more convenient access.
