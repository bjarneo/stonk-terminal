#!/bin/bash
env GOOS=linux GOARCH=arm go build -v cli.go
