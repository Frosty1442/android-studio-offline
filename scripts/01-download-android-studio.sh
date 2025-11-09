#!/bin/bash

################################################################################
# Android Studio Download Script
# Downloads Android Studio IDE and JDK for offline installation
################################################################################

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Load configuration
CONFIG_FILE="$PROJECT_ROOT/config/download.config"
if [ ! -f "$CONFIG_FILE" ]; then
    echo -e "${YELLOW}Warning: Config file not found. Creating from template...${NC}"
    cp "$PROJECT_ROOT/config/download.config.template" "$CONFIG_FILE"
    echo -e "${YELLOW}Please edit $CONFIG_FILE and run this script again.${NC}"
    exit 1
fi

source "$CONFIG_FILE"

# Create download directories
DOWNLOAD_DIR="$PROJECT_ROOT/$DOWNLOAD_DIR"
mkdir -p "$DOWNLOAD_DIR/android-studio"
mkdir -p "$DOWNLOAD_DIR/jdk"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Android Studio Offline Installer${NC}"
echo -e "${BLUE}Step 1: Download Android Studio & JDK${NC}"
echo -e "${BLUE}========================================${NC}"
echo

################################################################################
# Download Android Studio
################################################################################

download_android_studio() {
    echo -e "${GREEN}Downloading Android Studio...${NC}"

    local base_url="https://redirector.gvt1.com/edgedl/android/studio/ide-zips"
    local filename=""
    local url=""

    case "$PLATFORM" in
        linux)
            filename="android-studio-${ANDROID_STUDIO_VERSION}-linux.tar.gz"
            ;;
        mac)
            filename="android-studio-${ANDROID_STUDIO_VERSION}-mac.dmg"
            ;;
        mac_arm)
            filename="android-studio-${ANDROID_STUDIO_VERSION}-mac_arm.dmg"
            ;;
        windows)
            filename="android-studio-${ANDROID_STUDIO_VERSION}-windows.exe"
            ;;
        *)
            echo -e "${RED}Error: Unknown platform: $PLATFORM${NC}"
            exit 1
            ;;
    esac

    url="${base_url}/${ANDROID_STUDIO_VERSION}/${filename}"
    local output_file="$DOWNLOAD_DIR/android-studio/$filename"

    if [ -f "$output_file" ] && [ "$RESUME_DOWNLOADS" = "true" ]; then
        echo -e "${YELLOW}File exists. Resuming download...${NC}"
        wget -c -O "$output_file" "$url" || curl -C - -L -o "$output_file" "$url"
    else
        wget -O "$output_file" "$url" || curl -L -o "$output_file" "$url"
    fi

    echo -e "${GREEN}✓ Android Studio downloaded: $filename${NC}"
    echo
}

################################################################################
# Download JDK
################################################################################

download_jdk() {
    if [ "$DOWNLOAD_JDK" != "true" ]; then
        echo -e "${YELLOW}Skipping JDK download (disabled in config)${NC}"
        return
    fi

    echo -e "${GREEN}Downloading JDK ${JDK_VERSION}...${NC}"

    # Using Eclipse Temurin (previously AdoptOpenJDK)
    local jdk_arch=""
    case "$PLATFORM" in
        linux)
            jdk_arch="linux-x64"
            jdk_ext="tar.gz"
            ;;
        mac)
            jdk_arch="mac-x64"
            jdk_ext="tar.gz"
            ;;
        mac_arm)
            jdk_arch="mac-aarch64"
            jdk_ext="tar.gz"
            ;;
        windows)
            jdk_arch="windows-x64"
            jdk_ext="zip"
            ;;
    esac

    # Eclipse Temurin JDK download
    local jdk_url="https://api.adoptium.net/v3/binary/latest/${JDK_VERSION}/ga/${jdk_arch}/jdk/hotspot/normal/eclipse"
    local jdk_filename="jdk-${JDK_VERSION}-${jdk_arch}.${jdk_ext}"
    local output_file="$DOWNLOAD_DIR/jdk/$jdk_filename"

    echo "Downloading from: $jdk_url"
    echo "Saving to: $output_file"

    if [ -f "$output_file" ] && [ "$RESUME_DOWNLOADS" = "true" ]; then
        echo -e "${YELLOW}File exists. Resuming download...${NC}"
        wget -c -O "$output_file" "$jdk_url" || curl -C - -L -o "$output_file" "$jdk_url"
    else
        wget -O "$output_file" "$jdk_url" || curl -L -o "$output_file" "$jdk_url"
    fi

    echo -e "${GREEN}✓ JDK downloaded: $jdk_filename${NC}"
    echo
}

################################################################################
# Download Android SDK Command Line Tools
################################################################################

download_sdk_cmdline_tools() {
    echo -e "${GREEN}Downloading Android SDK Command Line Tools...${NC}"

    local sdk_arch=""
    case "$PLATFORM" in
        linux)
            sdk_arch="linux"
            ;;
        mac|mac_arm)
            sdk_arch="mac"
            ;;
        windows)
            sdk_arch="win"
            ;;
    esac

    local sdk_tools_url="https://dl.google.com/android/repository/commandlinetools-${sdk_arch}-${SDK_TOOLS_VERSION}_latest.zip"
    local sdk_tools_filename="commandlinetools-${sdk_arch}-${SDK_TOOLS_VERSION}.zip"
    local output_file="$DOWNLOAD_DIR/sdk/$sdk_tools_filename"

    mkdir -p "$DOWNLOAD_DIR/sdk"

    if [ -f "$output_file" ] && [ "$RESUME_DOWNLOADS" = "true" ]; then
        echo -e "${YELLOW}File exists. Resuming download...${NC}"
        wget -c -O "$output_file" "$sdk_tools_url" || curl -C - -L -o "$output_file" "$sdk_tools_url"
    else
        wget -O "$output_file" "$sdk_tools_url" || curl -L -o "$output_file" "$sdk_tools_url"
    fi

    echo -e "${GREEN}✓ SDK Command Line Tools downloaded${NC}"
    echo
}

################################################################################
# Main execution
################################################################################

main() {
    echo "Configuration:"
    echo "  Platform: $PLATFORM"
    echo "  Android Studio Version: $ANDROID_STUDIO_VERSION"
    echo "  SDK Tools Version: $SDK_TOOLS_VERSION"
    echo "  Download JDK: $DOWNLOAD_JDK"
    echo "  Download Directory: $DOWNLOAD_DIR"
    echo

    download_android_studio
    download_jdk
    download_sdk_cmdline_tools

    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}Step 1 Complete!${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo
    echo "Downloaded files:"
    ls -lh "$DOWNLOAD_DIR/android-studio/"
    [ "$DOWNLOAD_JDK" = "true" ] && ls -lh "$DOWNLOAD_DIR/jdk/" || true
    ls -lh "$DOWNLOAD_DIR/sdk/" 2>/dev/null || true
    echo
    echo -e "${YELLOW}Next step: Run ./scripts/02-download-sdk-components.sh${NC}"
}

main "$@"
