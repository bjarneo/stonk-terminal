#!/bin/bash

# Docs
# https://golang.org/doc/install/source#environment

rm -rf build/

# Linux
env GOOS=linux GOARCH=arm go build -o build/cli-linux-arm -v cli.go  
env GOOS=linux GOARCH=386 go build -o build/cli-linux-386 -v cli.go  
env GOOS=linux GOARCH=386 go build -o build/cli-linux-amd64 -v cli.go  

# macos
env GOOS=darwin GOARCH=amd64 go build -o build/cli-darwin -v cli.go 

# windows
env GOOS=windows GOARCH=amd64 go build -o build/cli-windows-amd64.exe -v cli.go 

chmod +x build/