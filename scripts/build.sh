#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

VERSION=${GITHUB_REF_NAME:-v0.1.0}
BINARY_NAME="webpanel"

# Create release directory
rm -rf release
mkdir -p release

echo -e "${YELLOW}Building for linux/amd64...${NC}"
GOOS=linux GOARCH=amd64 go build -o "release/${BINARY_NAME}-linux-amd64" cmd/webpanel/main.go
sha256sum "release/${BINARY_NAME}-linux-amd64" > "release/${BINARY_NAME}-linux-amd64.sha256"

echo -e "${YELLOW}Building for linux/386...${NC}"
GOOS=linux GOARCH=386 go build -o "release/${BINARY_NAME}-linux-386" cmd/webpanel/main.go
sha256sum "release/${BINARY_NAME}-linux-386" > "release/${BINARY_NAME}-linux-386.sha256"

echo -e "${YELLOW}Building for linux/arm64...${NC}"
GOOS=linux GOARCH=arm64 go build -o "release/${BINARY_NAME}-linux-arm64" cmd/webpanel/main.go
sha256sum "release/${BINARY_NAME}-linux-arm64" > "release/${BINARY_NAME}-linux-arm64.sha256"

echo -e "${YELLOW}Building for linux/arm v6...${NC}"
GOOS=linux GOARCH=arm GOARM=6 go build -o "release/${BINARY_NAME}-linux-arm-v6" cmd/webpanel/main.go
sha256sum "release/${BINARY_NAME}-linux-arm-v6" > "release/${BINARY_NAME}-linux-arm-v6.sha256"

echo -e "${YELLOW}Building for linux/arm v7...${NC}"
GOOS=linux GOARCH=arm GOARM=7 go build -o "release/${BINARY_NAME}-linux-arm-v7" cmd/webpanel/main.go
sha256sum "release/${BINARY_NAME}-linux-arm-v7" > "release/${BINARY_NAME}-linux-arm-v7.sha256"

# Create release notes
cat > release/release-notes.md <<EOF
## CLI Web Panel ${VERSION}

### Supported Architectures:
- linux/amd64 (64-bit x86)
- linux/386 (32-bit x86)
- linux/arm64 (64-bit ARM)
- linux/arm-v6 (Raspberry Pi 1, Zero)
- linux/arm-v7 (Raspberry Pi 2, 3)

### Installation:
\`\`\`bash
curl -fsSL https://raw.githubusercontent.com/doko89/cli-webpanel/main/scripts/install.sh | sudo bash
\`\`\`

### SHA256 Checksums:
\`\`\`
$(cat release/*.sha256)
\`\`\`
EOF

echo -e "${GREEN}Build complete! Binaries and checksums available in release/ directory${NC}"
ls -l release/