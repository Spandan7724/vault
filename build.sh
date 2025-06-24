#!/bin/bash

# Cross-platform build script for vault password manager
# Creates binaries for major OS/architecture combinations

set -e

VERSION=${1:-"dev"}
BUILD_DIR="dist"

echo "Building vault v${VERSION} for multiple platforms..."

# Clean previous builds
rm -rf $BUILD_DIR
mkdir -p $BUILD_DIR

# Build configurations: OS/ARCH
declare -a platforms=(
    "linux/amd64"
    "linux/arm64" 
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
    "windows/arm64"
)

# Build for each platform
for platform in "${platforms[@]}"; do
    IFS='/' read -r -a array <<< "$platform"
    GOOS="${array[0]}"
    GOARCH="${array[1]}"
    
    # Binary name (what goes inside the archive)
    if [ "$GOOS" = "windows" ]; then
        binary_name="vault.exe"
    else
        binary_name="vault"
    fi
    
    # Archive name (the downloadable file)
    archive_name="vault-${VERSION}-${GOOS}-${GOARCH}"
    
    echo "Building for $GOOS/$GOARCH..."
    
    env GOOS="$GOOS" GOARCH="$GOARCH" go build \
        -ldflags="-s -w -X main.version=${VERSION}" \
        -o "$BUILD_DIR/$binary_name" \
        .
    
    # Create compressed archives with correct binary name inside
    if [ "$GOOS" = "windows" ]; then
        # Create zip for Windows
        (cd $BUILD_DIR && zip -q "${archive_name}.zip" "$binary_name")
    else
        # Create tar.gz for Unix-like systems
        (cd $BUILD_DIR && tar -czf "${archive_name}.tar.gz" "$binary_name")
    fi
    
    # Clean up the binary (keep only the archive)
    rm "$BUILD_DIR/$binary_name"
    
    echo "✓ Built ${archive_name}"
done

echo ""
echo "Build complete! Artifacts in $BUILD_DIR:"
ls -la $BUILD_DIR/

echo ""
echo "To create a GitHub release:"
echo "1. Tag your commit: git tag v${VERSION}"
echo "2. Push the tag: git push origin v${VERSION}"
echo "3. Upload the files in $BUILD_DIR/ to the GitHub release"