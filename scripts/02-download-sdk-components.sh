#!/bin/bash

################################################################################
# Android SDK Components Download Script
# Downloads SDK packages, platforms, build tools, system images, etc.
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
    echo -e "${YELLOW}Please run 01-download-android-studio.sh first${NC}"
    exit 1
fi

source "$CONFIG_FILE"

DOWNLOAD_DIR="$PROJECT_ROOT/$DOWNLOAD_DIR"
SDK_DIR="$DOWNLOAD_DIR/sdk"
ANDROID_SDK_ROOT="$SDK_DIR/android-sdk"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Android Studio Offline Installer${NC}"
echo -e "${BLUE}Step 2: Download SDK Components${NC}"
echo -e "${BLUE}========================================${NC}"
echo

################################################################################
# Setup SDK Command Line Tools
################################################################################

setup_sdk_tools() {
    echo -e "${GREEN}Setting up SDK command line tools...${NC}"

    mkdir -p "$ANDROID_SDK_ROOT/cmdline-tools"

    # Find and extract cmdline tools
    local cmdline_zip=$(find "$SDK_DIR" -name "commandlinetools-*.zip" | head -1)

    if [ -z "$cmdline_zip" ]; then
        echo -e "${RED}Error: Command line tools not found${NC}"
        echo -e "${YELLOW}Please run 01-download-android-studio.sh first${NC}"
        exit 1
    fi

    echo "Extracting: $cmdline_zip"
    unzip -q -o "$cmdline_zip" -d "$ANDROID_SDK_ROOT/cmdline-tools"

    # Move to proper structure (cmdline-tools/latest/)
    if [ -d "$ANDROID_SDK_ROOT/cmdline-tools/cmdline-tools" ]; then
        mv "$ANDROID_SDK_ROOT/cmdline-tools/cmdline-tools" "$ANDROID_SDK_ROOT/cmdline-tools/latest"
    fi

    echo -e "${GREEN}✓ SDK tools set up${NC}"
    echo
}

################################################################################
# Download SDK Packages
################################################################################

download_sdk_packages() {
    echo -e "${GREEN}Downloading SDK packages...${NC}"
    echo "This may take a while depending on your selections..."
    echo

    export ANDROID_SDK_ROOT="$ANDROID_SDK_ROOT"
    local sdkmanager="$ANDROID_SDK_ROOT/cmdline-tools/latest/bin/sdkmanager"

    if [ ! -f "$sdkmanager" ]; then
        echo -e "${RED}Error: sdkmanager not found at $sdkmanager${NC}"
        exit 1
    fi

    # Accept licenses automatically
    yes | "$sdkmanager" --licenses || true

    # Read packages from config file
    local packages_file="$PROJECT_ROOT/config/sdk-packages.txt"

    if [ -f "$packages_file" ]; then
        echo -e "${YELLOW}Using packages from: $packages_file${NC}"

        # Read packages, skip comments and empty lines
        local packages=$(grep -v '^#' "$packages_file" | grep -v '^$' | tr '\n' ' ')

        if [ -n "$packages" ]; then
            echo "Installing packages: $packages"
            "$sdkmanager" $packages --sdk_root="$ANDROID_SDK_ROOT"
        fi
    fi

    # Download specific API levels from config
    if [ -n "$API_LEVELS" ]; then
        echo -e "${YELLOW}Downloading API levels: $API_LEVELS${NC}"
        for api in $API_LEVELS; do
            "$sdkmanager" "platforms;android-${api}" --sdk_root="$ANDROID_SDK_ROOT"
        done
    fi

    # Download build tools
    if [ -n "$BUILD_TOOLS_VERSIONS" ]; then
        echo -e "${YELLOW}Downloading Build Tools: $BUILD_TOOLS_VERSIONS${NC}"
        for version in $BUILD_TOOLS_VERSIONS; do
            "$sdkmanager" "build-tools;${version}" --sdk_root="$ANDROID_SDK_ROOT"
        done
    fi

    # Download system images
    if [ -n "$SYSTEM_IMAGES" ]; then
        echo -e "${YELLOW}Downloading System Images...${NC}"
        for image in $SYSTEM_IMAGES; do
            echo "Downloading: $image"
            "$sdkmanager" "$image" --sdk_root="$ANDROID_SDK_ROOT"
        done
    fi

    # Download NDK if requested
    if [ "$DOWNLOAD_NDK" = "true" ] && [ -n "$NDK_VERSION" ]; then
        echo -e "${YELLOW}Downloading NDK ${NDK_VERSION}...${NC}"
        "$sdkmanager" "ndk;${NDK_VERSION}" --sdk_root="$ANDROID_SDK_ROOT"
    fi

    echo -e "${GREEN}✓ SDK packages downloaded${NC}"
    echo
}

################################################################################
# Download Platform Tools
################################################################################

download_platform_tools() {
    if [ "$DOWNLOAD_PLATFORM_TOOLS" != "true" ]; then
        return
    fi

    echo -e "${GREEN}Downloading standalone Platform Tools...${NC}"

    local platform_arch=""
    case "$PLATFORM" in
        linux)
            platform_arch="linux"
            ;;
        mac|mac_arm)
            platform_arch="darwin"
            ;;
        windows)
            platform_arch="windows"
            ;;
    esac

    local platform_tools_url="https://dl.google.com/android/repository/platform-tools-latest-${platform_arch}.zip"
    local output_file="$SDK_DIR/platform-tools-${platform_arch}.zip"

    wget -O "$output_file" "$platform_tools_url" || curl -L -o "$output_file" "$platform_tools_url"

    echo -e "${GREEN}✓ Platform tools downloaded${NC}"
    echo
}

################################################################################
# Download Google Maven Repository
################################################################################

download_google_maven() {
    if [ "$DOWNLOAD_GOOGLE_REPOS" != "true" ]; then
        return
    fi

    echo -e "${GREEN}Downloading Google Maven Repository...${NC}"

    local maven_url="https://dl.google.com/android/repository/google_m2repository_r58.zip"
    local output_file="$SDK_DIR/google_m2repository.zip"

    wget -O "$output_file" "$maven_url" || curl -L -o "$output_file" "$maven_url"

    # Extract to SDK
    mkdir -p "$ANDROID_SDK_ROOT/extras/google"
    unzip -q -o "$output_file" -d "$ANDROID_SDK_ROOT/extras/google/"

    echo -e "${GREEN}✓ Google Maven repository downloaded${NC}"
    echo
}

################################################################################
# Create SDK info file
################################################################################

create_sdk_info() {
    echo -e "${GREEN}Creating SDK info file...${NC}"

    local info_file="$SDK_DIR/sdk-info.txt"

    cat > "$info_file" << EOF
Android SDK Components Download Information
============================================

Download Date: $(date)
SDK Root: $ANDROID_SDK_ROOT

Configuration:
--------------
API Levels: $API_LEVELS
Build Tools: $BUILD_TOOLS_VERSIONS
System Images: $SYSTEM_IMAGES
NDK Version: $NDK_VERSION
Platform: $PLATFORM

Installed Packages:
-------------------
EOF

    # List installed packages
    export ANDROID_SDK_ROOT="$ANDROID_SDK_ROOT"
    "$ANDROID_SDK_ROOT/cmdline-tools/latest/bin/sdkmanager" --list_installed --sdk_root="$ANDROID_SDK_ROOT" >> "$info_file"

    echo -e "${GREEN}✓ SDK info saved to: $info_file${NC}"
    echo
}

################################################################################
# Main execution
################################################################################

main() {
    echo "Configuration:"
    echo "  SDK Root: $ANDROID_SDK_ROOT"
    echo "  API Levels: $API_LEVELS"
    echo "  Build Tools: $BUILD_TOOLS_VERSIONS"
    echo

    setup_sdk_tools
    download_sdk_packages
    download_platform_tools
    download_google_maven
    create_sdk_info

    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}Step 2 Complete!${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo
    echo "SDK Size:"
    du -sh "$ANDROID_SDK_ROOT"
    echo
    echo -e "${YELLOW}Next step: Run ./scripts/03-download-gradle.sh${NC}"
}

main "$@"
