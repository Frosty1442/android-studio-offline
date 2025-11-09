#!/bin/bash

################################################################################
# Gradle Download Script
# Downloads Gradle distributions for offline use
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
GRADLE_DIR="$DOWNLOAD_DIR/gradle"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Android Studio Offline Installer${NC}"
echo -e "${BLUE}Step 3: Download Gradle${NC}"
echo -e "${BLUE}========================================${NC}"
echo

mkdir -p "$GRADLE_DIR/distributions"
mkdir -p "$GRADLE_DIR/wrapper"

################################################################################
# Download Gradle Distributions
################################################################################

download_gradle_versions() {
    echo -e "${GREEN}Downloading Gradle distributions...${NC}"

    local gradle_versions_file="$PROJECT_ROOT/config/gradle-versions.txt"
    local versions=""

    # Get versions from config file if it exists
    if [ -f "$gradle_versions_file" ]; then
        echo -e "${YELLOW}Reading versions from: $gradle_versions_file${NC}"
        versions=$(grep -v '^#' "$gradle_versions_file" | grep -v '^$' | tr '\n' ' ')
    fi

    # Also add versions from config
    if [ -n "$GRADLE_VERSIONS" ]; then
        versions="$versions $GRADLE_VERSIONS"
    fi

    # Remove duplicates
    versions=$(echo "$versions" | tr ' ' '\n' | sort -u | tr '\n' ' ')

    if [ -z "$versions" ]; then
        echo -e "${YELLOW}No Gradle versions specified. Using defaults.${NC}"
        versions="7.6 8.0 8.4"
    fi

    echo "Downloading Gradle versions: $versions"
    echo

    for version in $versions; do
        download_gradle_distribution "$version"
    done

    echo -e "${GREEN}✓ All Gradle distributions downloaded${NC}"
    echo
}

download_gradle_distribution() {
    local version="$1"
    local base_url="https://services.gradle.org/distributions"

    # Download both bin and all distributions
    for dist_type in "bin" "all"; do
        local filename="gradle-${version}-${dist_type}.zip"
        local url="${base_url}/${filename}"
        local output_file="$GRADLE_DIR/distributions/${filename}"

        if [ -f "$output_file" ]; then
            echo -e "${YELLOW}  ✓ $filename already exists, skipping${NC}"
            continue
        fi

        echo "  Downloading: $filename"

        if wget -q -O "$output_file" "$url"; then
            echo -e "${GREEN}    ✓ Downloaded successfully${NC}"
        elif curl -s -L -o "$output_file" "$url"; then
            echo -e "${GREEN}    ✓ Downloaded successfully${NC}"
        else
            echo -e "${RED}    ✗ Failed to download $filename${NC}"
            rm -f "$output_file"
        fi
    done

    # Also download checksum files
    for dist_type in "bin" "all"; do
        local filename="gradle-${version}-${dist_type}.zip"
        local checksum_url="${base_url}/${filename}.sha256"
        local checksum_file="$GRADLE_DIR/distributions/${filename}.sha256"

        wget -q -O "$checksum_file" "$checksum_url" 2>/dev/null || \
        curl -s -L -o "$checksum_file" "$checksum_url" 2>/dev/null || \
        rm -f "$checksum_file"
    done
}

################################################################################
# Download Gradle Wrapper
################################################################################

download_gradle_wrapper() {
    if [ "$DOWNLOAD_GRADLE_WRAPPER" != "true" ]; then
        echo -e "${YELLOW}Skipping Gradle Wrapper download${NC}"
        return
    fi

    echo -e "${GREEN}Downloading Gradle Wrapper files...${NC}"

    # Get the latest wrapper version (use one of the downloaded versions)
    local latest_version=$(echo "$versions" | tr ' ' '\n' | sort -V | tail -1)

    if [ -z "$latest_version" ]; then
        latest_version="8.4"
    fi

    echo "Using Gradle version $latest_version for wrapper files"

    # Download gradle-wrapper.jar
    local wrapper_jar_url="https://raw.githubusercontent.com/gradle/gradle/v${latest_version}/gradle/wrapper/gradle-wrapper.jar"
    local wrapper_jar="$GRADLE_DIR/wrapper/gradle-wrapper.jar"

    wget -O "$wrapper_jar" "$wrapper_jar_url" 2>/dev/null || \
    curl -L -o "$wrapper_jar" "$wrapper_jar_url" 2>/dev/null || \
    echo -e "${YELLOW}  Could not download gradle-wrapper.jar (this is optional)${NC}"

    # Create sample wrapper files
    create_wrapper_template "$latest_version"

    echo -e "${GREEN}✓ Gradle wrapper files ready${NC}"
    echo
}

################################################################################
# Create Gradle Wrapper Template
################################################################################

create_wrapper_template() {
    local version="$1"

    # Create gradle-wrapper.properties template
    cat > "$GRADLE_DIR/wrapper/gradle-wrapper.properties.template" << EOF
distributionBase=GRADLE_USER_HOME
distributionPath=wrapper/dists
distributionUrl=file\:///PATH_TO_OFFLINE_GRADLE/distributions/gradle-${version}-all.zip
zipStoreBase=GRADLE_USER_HOME
zipStorePath=wrapper/dists
EOF

    # Create gradlew script template
    cat > "$GRADLE_DIR/wrapper/README.md" << EOF
# Gradle Wrapper for Offline Use

## Setting up Gradle Wrapper in your project

1. Copy gradle-wrapper.jar to your_project/gradle/wrapper/
2. Copy gradle-wrapper.properties.template to your_project/gradle/wrapper/gradle-wrapper.properties
3. Edit gradle-wrapper.properties and replace PATH_TO_OFFLINE_GRADLE with the actual path
4. Use gradlew as normal: ./gradlew build

## Example gradle-wrapper.properties

\`\`\`properties
distributionBase=GRADLE_USER_HOME
distributionPath=wrapper/dists
distributionUrl=file\:///opt/android-offline/gradle/distributions/gradle-8.4-all.zip
zipStoreBase=GRADLE_USER_HOME
zipStorePath=wrapper/dists
\`\`\`

Note: The distributionUrl must use file:/// protocol for offline access.
EOF

    echo -e "${GREEN}✓ Wrapper template created${NC}"
}

################################################################################
# Create Gradle info file
################################################################################

create_gradle_info() {
    echo -e "${GREEN}Creating Gradle info file...${NC}"

    local info_file="$GRADLE_DIR/gradle-info.txt"

    cat > "$info_file" << EOF
Gradle Downloads Information
=============================

Download Date: $(date)

Downloaded Versions:
--------------------
EOF

    # List downloaded files
    ls -lh "$GRADLE_DIR/distributions/" >> "$info_file"

    echo "" >> "$info_file"
    echo "Total Size:" >> "$info_file"
    du -sh "$GRADLE_DIR" >> "$info_file"

    echo -e "${GREEN}✓ Gradle info saved to: $info_file${NC}"
    echo
}

################################################################################
# Setup offline Gradle configuration
################################################################################

setup_offline_gradle_config() {
    echo -e "${GREEN}Creating offline Gradle configuration...${NC}"

    mkdir -p "$GRADLE_DIR/init.d"

    # Create init script for offline mode
    cat > "$GRADLE_DIR/init.d/offline-repos.gradle" << 'EOF'
// Offline repository configuration
// Copy this file to ~/.gradle/init.d/ on the offline machine

allprojects {
    buildscript {
        repositories {
            // Clear all repositories and use local only
            all { ArtifactRepository repo ->
                if (repo instanceof MavenArtifactRepository) {
                    remove repo
                }
            }

            // Add local Maven repository
            maven {
                url = uri("/opt/android-offline/maven-repo")
            }

            // Add Gradle plugin repository
            maven {
                url = uri("/opt/android-offline/gradle-plugins")
            }
        }
    }

    repositories {
        // Clear all repositories
        all { ArtifactRepository repo ->
            if (repo instanceof MavenArtifactRepository) {
                remove repo
            }
        }

        // Add local Maven repository
        maven {
            url = uri("/opt/android-offline/maven-repo")
        }

        // Add local Google repository
        maven {
            url = uri("/opt/android-offline/google-repo")
        }
    }
}
EOF

    echo -e "${GREEN}✓ Offline Gradle configuration created${NC}"
    echo
}

################################################################################
# Main execution
################################################################################

main() {
    echo "Configuration:"
    echo "  Gradle Directory: $GRADLE_DIR"
    echo "  Versions to download: $GRADLE_VERSIONS"
    echo

    download_gradle_versions
    download_gradle_wrapper
    setup_offline_gradle_config
    create_gradle_info

    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}Step 3 Complete!${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo
    echo "Downloaded files:"
    ls -lh "$GRADLE_DIR/distributions/" 2>/dev/null | head -20
    echo
    echo "Total Gradle size:"
    du -sh "$GRADLE_DIR"
    echo
    echo -e "${YELLOW}Next step: Run ./scripts/04-download-dependencies.sh${NC}"
}

main "$@"
