#!/bin/bash

################################################################################
# Package Creation Script
# Creates a portable package of all downloaded components
################################################################################

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Load configuration
CONFIG_FILE="$PROJECT_ROOT/config/download.config"
if [ ! -f "$CONFIG_FILE" ]; then
    echo -e "${RED}Error: Config file not found at $CONFIG_FILE${NC}"
    exit 1
fi

source "$CONFIG_FILE"

DOWNLOAD_DIR="$PROJECT_ROOT/$DOWNLOAD_DIR"
PACKAGE_NAME="${OUTPUT_PACKAGE:-android-studio-offline-$(date +%Y%m%d).tar.gz}"
PACKAGE_DIR="$PROJECT_ROOT/packages"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Android Studio Offline Installer${NC}"
echo -e "${BLUE}Step 5: Create Package${NC}"
echo -e "${BLUE}========================================${NC}"
echo

mkdir -p "$PACKAGE_DIR"

################################################################################
# Verify downloads
################################################################################

verify_downloads() {
    echo -e "${GREEN}Verifying downloads...${NC}"

    local errors=0

    # Check Android Studio
    if ! ls "$DOWNLOAD_DIR/android-studio/"*.* &>/dev/null; then
        echo -e "${RED}  ✗ Android Studio not found${NC}"
        ((errors++))
    else
        echo -e "${GREEN}  ✓ Android Studio found${NC}"
    fi

    # Check SDK
    if [ ! -d "$DOWNLOAD_DIR/sdk/android-sdk" ]; then
        echo -e "${YELLOW}  ⚠ Android SDK not found (optional)${NC}"
    else
        echo -e "${GREEN}  ✓ Android SDK found${NC}"
    fi

    # Check Gradle
    if ! ls "$DOWNLOAD_DIR/gradle/distributions/"*.zip &>/dev/null; then
        echo -e "${YELLOW}  ⚠ Gradle distributions not found (optional)${NC}"
    else
        echo -e "${GREEN}  ✓ Gradle distributions found${NC}"
    fi

    if [ $errors -gt 0 ]; then
        echo -e "${RED}Verification failed. Please run the download scripts first.${NC}"
        exit 1
    fi

    echo
}

################################################################################
# Create installation scripts for package
################################################################################

create_install_script() {
    echo -e "${GREEN}Creating installation script for package...${NC}"

    # Copy the offline installation script
    cp "$SCRIPT_DIR/install-offline.sh" "$DOWNLOAD_DIR/" 2>/dev/null || \
        echo -e "${YELLOW}install-offline.sh will be created next${NC}"

    # Create README for the package
    cat > "$DOWNLOAD_DIR/README.txt" << 'EOF'
Android Studio Offline Installer Package
=========================================

This package contains everything needed for Android development without
internet access.

Contents:
---------
- Android Studio IDE
- Android SDK
- JDK (Java Development Kit)
- Gradle distributions
- Maven dependencies
- Installation scripts

Quick Start:
-----------
1. Extract this package to a directory with at least 50GB free space
2. Run: ./install-offline.sh
3. Follow the prompts

For detailed instructions, see INSTALL.md

System Requirements:
-------------------
- 50+ GB free disk space
- 8+ GB RAM (16 GB recommended)
- 64-bit operating system
- Linux, macOS, or Windows (with WSL)

EOF

    echo -e "${GREEN}✓ Installation script prepared${NC}"
    echo
}

################################################################################
# Create manifest file
################################################################################

create_manifest() {
    echo -e "${GREEN}Creating package manifest...${NC}"

    local manifest="$DOWNLOAD_DIR/MANIFEST.txt"

    cat > "$manifest" << EOF
Android Studio Offline Installer Package
=========================================

Created: $(date)
Platform: $PLATFORM
Package Version: 1.0

Components:
-----------
EOF

    # List Android Studio files
    echo "" >> "$manifest"
    echo "Android Studio:" >> "$manifest"
    ls -lh "$DOWNLOAD_DIR/android-studio/" >> "$manifest" 2>/dev/null || echo "  Not included" >> "$manifest"

    # List JDK
    echo "" >> "$manifest"
    echo "JDK:" >> "$manifest"
    ls -lh "$DOWNLOAD_DIR/jdk/" >> "$manifest" 2>/dev/null || echo "  Not included" >> "$manifest"

    # List SDK info
    echo "" >> "$manifest"
    echo "Android SDK:" >> "$manifest"
    if [ -f "$DOWNLOAD_DIR/sdk/sdk-info.txt" ]; then
        cat "$DOWNLOAD_DIR/sdk/sdk-info.txt" >> "$manifest"
    else
        echo "  Not included" >> "$manifest"
    fi

    # List Gradle
    echo "" >> "$manifest"
    echo "Gradle Distributions:" >> "$manifest"
    ls -lh "$DOWNLOAD_DIR/gradle/distributions/" >> "$manifest" 2>/dev/null || echo "  Not included" >> "$manifest"

    # Size information
    echo "" >> "$manifest"
    echo "Size Information:" >> "$manifest"
    echo "----------------" >> "$manifest"
    du -sh "$DOWNLOAD_DIR"/* >> "$manifest" 2>/dev/null || true
    echo "" >> "$manifest"
    echo "Total Package Size:" >> "$manifest"
    du -sh "$DOWNLOAD_DIR" >> "$manifest"

    echo -e "${GREEN}✓ Manifest created${NC}"
    echo
}

################################################################################
# Create checksum file
################################################################################

create_checksums() {
    if [ "$VERIFY_CHECKSUMS" != "true" ]; then
        echo -e "${YELLOW}Skipping checksum creation${NC}"
        return
    fi

    echo -e "${GREEN}Creating checksums...${NC}"
    echo "This may take a while for large files..."

    local checksum_file="$DOWNLOAD_DIR/SHA256SUMS.txt"

    cd "$DOWNLOAD_DIR"

    # Create checksums for important files
    find . -type f \( -name "*.zip" -o -name "*.tar.gz" -o -name "*.dmg" -o -name "*.exe" \) \
        -exec sha256sum {} \; > "$checksum_file" 2>/dev/null || \
        find . -type f \( -name "*.zip" -o -name "*.tar.gz" -o -name "*.dmg" -o -name "*.exe" \) \
        -exec shasum -a 256 {} \; > "$checksum_file" 2>/dev/null || \
        echo "Could not create checksums" > "$checksum_file"

    cd "$PROJECT_ROOT"

    echo -e "${GREEN}✓ Checksums created${NC}"
    echo
}

################################################################################
# Create compressed package
################################################################################

create_compressed_package() {
    echo -e "${GREEN}Creating compressed package...${NC}"
    echo "This will take a while depending on the size..."
    echo

    local package_path="$PACKAGE_DIR/$PACKAGE_NAME"

    # Create tar.gz package
    echo "Compressing to: $package_path"

    cd "$(dirname "$DOWNLOAD_DIR")"

    tar -czf "$package_path" \
        --exclude='.git' \
        --exclude='.gradle-temp' \
        --exclude='.temp-project' \
        "$(basename "$DOWNLOAD_DIR")" \
        2>&1 | grep -v "Removing leading" || true

    cd "$PROJECT_ROOT"

    # Create split archives if package is very large (>4GB)
    local size=$(stat -f%z "$package_path" 2>/dev/null || stat -c%s "$package_path" 2>/dev/null || echo 0)

    if [ "$size" -gt 4294967296 ]; then
        echo -e "${YELLOW}Package is larger than 4GB, creating split archives...${NC}"

        split -b 4000m "$package_path" "$package_path.part"
        echo -e "${GREEN}Created split archives: $package_path.part*${NC}"

        # Create join script
        cat > "$PACKAGE_DIR/join-parts.sh" << 'EOF'
#!/bin/bash
# Script to join split package parts

PACKAGE_NAME="$1"

if [ -z "$PACKAGE_NAME" ]; then
    echo "Usage: ./join-parts.sh <package-name>"
    echo "Example: ./join-parts.sh android-studio-offline-20241109.tar.gz"
    exit 1
fi

cat ${PACKAGE_NAME}.part* > ${PACKAGE_NAME}
echo "Joined package: ${PACKAGE_NAME}"
echo "Verifying..."
tar -tzf ${PACKAGE_NAME} > /dev/null && echo "Package verified successfully!"
EOF

        chmod +x "$PACKAGE_DIR/join-parts.sh"
    fi

    echo -e "${GREEN}✓ Package created: $package_path${NC}"
    echo
}

################################################################################
# Create package info
################################################################################

create_package_info() {
    echo -e "${GREEN}Creating package information...${NC}"

    local info_file="$PACKAGE_DIR/PACKAGE-INFO.txt"

    cat > "$info_file" << EOF
Android Studio Offline Installer Package
=========================================

Package: $PACKAGE_NAME
Created: $(date)
Platform: $PLATFORM

Package Location:
----------------
$PACKAGE_DIR/$PACKAGE_NAME

Package Size:
------------
EOF

    du -sh "$PACKAGE_DIR/$PACKAGE_NAME" >> "$info_file" 2>/dev/null || echo "N/A" >> "$info_file"

    cat >> "$info_file" << 'EOF'

Installation Instructions:
-------------------------
1. Transfer this package to your offline machine
2. Extract: tar -xzf PACKAGE_NAME
3. Run: cd downloads && ./install-offline.sh
4. Follow the on-screen instructions

For split packages:
1. Transfer all .partXX files to the same directory
2. Run: ./join-parts.sh PACKAGE_NAME
3. Continue with normal installation

Support:
--------
For issues, see docs/04-TROUBLESHOOTING.md in the repository

EOF

    echo -e "${GREEN}✓ Package info created${NC}"
    echo
}

################################################################################
# Main execution
################################################################################

main() {
    verify_downloads
    create_install_script
    create_manifest
    create_checksums
    create_compressed_package
    create_package_info

    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}Step 5 Complete!${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo
    echo "Package Details:"
    echo "---------------"
    cat "$PACKAGE_DIR/PACKAGE-INFO.txt"
    echo
    echo -e "${GREEN}Package ready for transfer to offline machine!${NC}"
    echo
    echo -e "${YELLOW}Next: Transfer the package and run install-offline.sh${NC}"
}

main "$@"
