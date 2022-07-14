#!/bin/bash
set -e

# Golang
if ! command -v golangci-lint &> /dev/null; then
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
fi

golangci-lint run
