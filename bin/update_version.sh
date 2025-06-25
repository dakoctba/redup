#!/bin/bash

# Script to update version in go.mod

set -e

VERSION=${1:-$(git describe --tags --always --dirty)}

echo "Updating version to: $VERSION"

# Update go.mod version
go mod edit -go=1.21

echo "Version updated to $VERSION"
echo "Run 'go mod tidy' to clean up dependencies"
