#!/usr/bin/env bash
set -euo pipefail

echo "==> Running go vet..."
go vet ./...

echo "==> Running tests with race detector..."
go test -v -race ./...

echo "==> Building CLI binary..."
go build -v -o bin/bokuao ./cmd/bokuao

echo "==> All verification checks passed successfully!"
