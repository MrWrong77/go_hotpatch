#!/bin/bash
rm -rf Plugins/**/*.so
rm -rf main

# go run tool/main.go
go generate ./...
./main