#!/usr/bin/env bash

rm -rf ./build/

env GOOS=linux GOARCH=amd64 go build -o ./build/ .

env GOOS=windows GOARCH=amd64 go build -o ./build/ .
