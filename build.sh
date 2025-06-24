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
    
    output_name="vault-${VERSION}-${GOOS}-${GOARCH}"
    
    # Add .exe extension for Windows
    if [ "$GOOS" = "windows" ]; then
        output_name="${output_name}.exe"
    fi
    
    echo "Building for $GOOS/$GOARCH..."
    
    env GOOS="$GOOS" GOARCH="$GOARCH" go build \
        -ldflags="-s -w -X main.version=${VERSION}" \
        -o "$BUILD_DIR/$output_name" \
        .
    
    # Create compressed archives
    if [ "$GOOS" = "windows" ]; then
        # Create zip for Windows
        (cd $BUILD_DIR && zip -q "${output_name%.exe}.zip" "$output_name")
        rm "$BUILD_DIR/$output_name"
    else
        # Create tar.gz for Unix-like systems
        (cd $BUILD_DIR && tar -czf "${output_name}.tar.gz" "$output_name")
        rm "$BUILD_DIR/$output_name"
    fi
    
    echo "✓ Built $output_name"
done

echo ""
echo "Build complete! Artifacts in $BUILD_DIR:"
ls -la $BUILD_DIR/

echo ""
echo "To create a GitHub release:"
echo "1. Tag your commit: git tag v${VERSION}"
echo "2. Push the tag: git push origin v${VERSION}"
echo "3. Upload the files in $BUILD_DIR/ to the GitHub release"