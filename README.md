# drsz
A small CLI tool for finding the size of subdirectories.

## About
Given a directory path to begin searching from, **drsz** will calculate the total size in bytes of each top-level subdirectory. Handy if you want to quickly find what directories might be using up alot of storage space on your drive.

## Usage
At the command line run `./drsz ./some_directory`. **drsz** will recursively search through each top-level subdirectory in `./some_directory`, add up the file sizes of the contents, and then print the total size to the console.

```
>>> drsz ./docker-geoserver
Calculating top-level subdirectory sizes in ./docker-geoserver/*
.git = 982 kB
build_data = 181 kB
clustering = 1.5 kB
resources = 670 B
scripts = 24 kB
volume = 123.4 MB
Search complete after finding 6 top-level subdirectories.
```

If you want the output saved to a CSV file you can pass an `-o` argument: `./drsz -o ./output.csv ./some_directory`. Note that you need to pass the `-o` argument *before* the directory path. The CSV file will have two columns with the subdirectory paths and total sizes in bytes.

## Installation
Linux and Windows amd64 binaries have been pre-built and are available in Releases. Simply download one and follow the Usage instructions above.

At the command line you can either run the executable by providing the full filepath to it, or make an alias in your `~/.bashrc` on Linux (or your Windows path) to create a shortened `drsz` command for more convenient access.

You can also build from source using the `build.sh` script that has been preconfigured to create Linux and Windows binaries.
