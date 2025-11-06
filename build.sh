#!/bin/bash
# build.sh - Build itamae with version information

set -e

# Get version from git tag or use "dev"
VERSION=${VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo "dev")}

# Get git commit hash
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Get build date
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Create bin directory if it doesn't exist
mkdir -p bin

# Build with version information
echo "Building itamae ${VERSION} (${COMMIT})..."

go build -ldflags "\
    -X github.com/yjmrobert/itamae/cmd.Version=${VERSION} \
    -X github.com/yjmrobert/itamae/cmd.GitCommit=${COMMIT} \
    -X github.com/yjmrobert/itamae/cmd.BuildDate=${BUILD_DATE}" \
    -o bin/itamae

echo "✅ Build complete: bin/itamae"
echo "   Version: ${VERSION}"
echo "   Commit:  ${COMMIT}"
echo "   Date:    ${BUILD_DATE}"
