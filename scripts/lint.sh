#!/bin/bash
set -e

# Golang
if ! command -v golangci-lint &> /dev/null; then
    echo "Install Golang linter..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
fi

echo "Lint Golang..."
golangci-lint run
