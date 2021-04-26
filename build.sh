#!/bin/bash

# Docs
# https://golang.org/doc/install/source#environment

# Linux
env GOOS=linux GOARCH=arm go build -o build/cli-linux -v cli.go  

# macos
env GOOS=darwin GOARCH=amd64 go build -o build/cli-darwin -v cli.go 