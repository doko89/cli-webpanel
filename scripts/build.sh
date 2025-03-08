#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

VERSION="v0.1.0"
BINARY_NAME="webpanel"
PLATFORMS=("linux/amd64" "linux/386" "linux/arm" "linux/arm64")

# Create release directory
rm -rf release
mkdir -p release

# Build for each platform
for platform in "${PLATFORMS[@]}"
do
    # Split platform into OS and architecture
    IFS='/' read -r -a array <<< "$platform"
    OS="${array[0]}"
    ARCH="${array[1]}"
    
    echo -e "${YELLOW}Building for $OS/$ARCH...${NC}"
    
    # Set environment variables for cross-compilation
    export GOOS=$OS
    export GOARCH=$ARCH
    
    # Build binary
    OUTPUT="release/${BINARY_NAME}-${OS}-${ARCH}"
    if [ "$ARCH" = "arm" ]; then
        # Build for different ARM versions (v6, v7)
        for VERSION in 6 7; do
            export GOARM=$VERSION
            go build -o "${OUTPUT}-v${VERSION}" cmd/webpanel/main.go
            echo -e "${GREEN}Built ${OUTPUT}-v${VERSION}${NC}"
        done
    else
        go build -o "$OUTPUT" cmd/webpanel/main.go
        echo -e "${GREEN}Built $OUTPUT${NC}"
    fi
    
    # Create checksums
    if [ -f "$OUTPUT" ]; then
        sha256sum "$OUTPUT" > "${OUTPUT}.sha256"
    fi
done

echo -e "\n${GREEN}Build complete!${NC}"
echo -e "Binaries are available in the release directory"