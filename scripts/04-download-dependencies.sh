#!/bin/bash

################################################################################
# Maven Dependencies Download Script
# Downloads Maven dependencies for offline Android development
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
MAVEN_DIR="$DOWNLOAD_DIR/dependencies/maven-repo"
GOOGLE_DIR="$DOWNLOAD_DIR/dependencies/google-repo"
GRADLE_PLUGINS_DIR="$DOWNLOAD_DIR/dependencies/gradle-plugins"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Android Studio Offline Installer${NC}"
echo -e "${BLUE}Step 4: Download Dependencies${NC}"
echo -e "${BLUE}========================================${NC}"
echo

mkdir -p "$MAVEN_DIR"
mkdir -p "$GOOGLE_DIR"
mkdir -p "$GRADLE_PLUGINS_DIR"

################################################################################
# Download specific dependency
################################################################################

download_dependency() {
    local dep="$1"
    local repo_url="$2"
    local output_dir="$3"

    # Parse dependency (format: group:artifact:version)
    IFS=':' read -ra PARTS <<< "$dep"
    local group="${PARTS[0]}"
    local artifact="${PARTS[1]}"
    local version="${PARTS[2]}"

    if [ -z "$group" ] || [ -z "$artifact" ] || [ -z "$version" ]; then
        echo -e "${RED}Invalid dependency format: $dep${NC}"
        return 1
    fi

    # Convert group to path
    local group_path="${group//./\/}"
    local base_path="$group_path/$artifact/$version"
    local artifact_file="$artifact-$version"

    # Create directory structure
    mkdir -p "$output_dir/$base_path"

    # Download various artifact types
    local downloaded=false

    # Try .aar (Android library)
    if download_file "$repo_url/$base_path/$artifact_file.aar" "$output_dir/$base_path/$artifact_file.aar"; then
        downloaded=true
    fi

    # Try .jar
    if download_file "$repo_url/$base_path/$artifact_file.jar" "$output_dir/$base_path/$artifact_file.jar"; then
        downloaded=true
    fi

    # Download POM file (always needed)
    download_file "$repo_url/$base_path/$artifact_file.pom" "$output_dir/$base_path/$artifact_file.pom"

    # Download sources if available
    download_file "$repo_url/$base_path/$artifact_file-sources.jar" "$output_dir/$base_path/$artifact_file-sources.jar" || true

    # Download module metadata
    download_file "$repo_url/$base_path/$artifact_file.module" "$output_dir/$base_path/$artifact_file.module" || true

    if [ "$downloaded" = true ]; then
        echo -e "${GREEN}  ✓ $dep${NC}"
    else
        echo -e "${YELLOW}  ⚠ $dep (might not be available)${NC}"
    fi
}

download_file() {
    local url="$1"
    local output="$2"

    if [ -f "$output" ]; then
        return 0
    fi

    wget -q -O "$output" "$url" 2>/dev/null && return 0
    curl -s -L -o "$output" "$url" 2>/dev/null && return 0
    rm -f "$output"
    return 1
}

################################################################################
# Download common dependencies
################################################################################

download_common_dependencies() {
    echo -e "${GREEN}Downloading common Android dependencies...${NC}"
    echo "This will take a while..."
    echo

    local deps_file="$PROJECT_ROOT/config/common-dependencies.txt"

    if [ ! -f "$deps_file" ]; then
        echo -e "${YELLOW}Warning: $deps_file not found${NC}"
        return
    fi

    # Maven Central
    local maven_central="https://repo1.maven.org/maven2"

    # Google Maven
    local google_maven="https://dl.google.com/dl/android/maven2"

    # Read dependencies from file
    while IFS= read -r dep || [ -n "$dep" ]; do
        # Skip comments and empty lines
        [[ "$dep" =~ ^#.*$ ]] && continue
        [[ -z "$dep" ]] && continue

        echo "Processing: $dep"

        # Try Google repo first for androidx and google dependencies
        if [[ "$dep" =~ ^(androidx|com\.google|com\.android) ]]; then
            download_dependency "$dep" "$google_maven" "$GOOGLE_DIR"
        else
            download_dependency "$dep" "$maven_central" "$MAVEN_DIR"
        fi

    done < "$deps_file"

    echo -e "${GREEN}✓ Common dependencies downloaded${NC}"
    echo
}

################################################################################
# Download Android Gradle Plugin
################################################################################

download_android_gradle_plugin() {
    echo -e "${GREEN}Downloading Android Gradle Plugin...${NC}"

    local agp_versions="8.1.0 8.1.1 8.1.2 8.2.0 8.3.0"
    local google_maven="https://dl.google.com/dl/android/maven2"

    for version in $agp_versions; do
        echo "  Downloading AGP $version..."
        download_dependency "com.android.tools.build:gradle:$version" "$google_maven" "$GRADLE_PLUGINS_DIR"
    done

    # Download Kotlin Gradle Plugin
    echo "  Downloading Kotlin Gradle Plugin..."
    local kotlin_versions="1.9.0 1.9.10 1.9.20 1.9.21"
    local maven_central="https://repo1.maven.org/maven2"

    for version in $kotlin_versions; do
        download_dependency "org.jetbrains.kotlin:kotlin-gradle-plugin:$version" "$maven_central" "$GRADLE_PLUGINS_DIR"
    done

    echo -e "${GREEN}✓ Gradle plugins downloaded${NC}"
    echo
}

################################################################################
# Download from sample project
################################################################################

download_from_sample_project() {
    if [ -z "$SAMPLE_PROJECT_PATH" ] || [ ! -d "$SAMPLE_PROJECT_PATH" ]; then
        echo -e "${YELLOW}No sample project specified or found${NC}"
        return
    fi

    echo -e "${GREEN}Downloading dependencies from sample project...${NC}"

    local sample_project="$SAMPLE_PROJECT_PATH"
    local temp_gradle_home="$DOWNLOAD_DIR/dependencies/.gradle-temp"

    mkdir -p "$temp_gradle_home"

    # Build the project to download all dependencies
    cd "$sample_project"

    export GRADLE_USER_HOME="$temp_gradle_home"

    if [ -f "gradlew" ]; then
        ./gradlew dependencies --refresh-dependencies || echo "Build completed with warnings"
    elif command -v gradle &> /dev/null; then
        gradle dependencies --refresh-dependencies || echo "Build completed with warnings"
    else
        echo -e "${YELLOW}Gradle not found, skipping sample project${NC}"
        return
    fi

    # Copy downloaded dependencies
    if [ -d "$temp_gradle_home/caches/modules-2/files-2.1" ]; then
        echo "Copying downloaded dependencies..."
        cp -r "$temp_gradle_home/caches/modules-2/files-2.1"/* "$MAVEN_DIR/" 2>/dev/null || true
    fi

    cd "$PROJECT_ROOT"

    echo -e "${GREEN}✓ Sample project dependencies downloaded${NC}"
    echo
}

################################################################################
# Mirror entire Google Maven Repository (optional)
################################################################################

mirror_google_maven() {
    echo -e "${YELLOW}Mirroring Google Maven Repository...${NC}"
    echo "This will download a LOT of data. Press Ctrl+C to skip or wait 5 seconds..."
    sleep 5

    # Use wget to mirror the repository
    wget -r -np -nH --cut-dirs=3 \
        -R "index.html*" \
        -P "$GOOGLE_DIR" \
        https://dl.google.com/dl/android/maven2/ \
        2>&1 | grep -v "^--" || true

    echo -e "${GREEN}✓ Google Maven mirrored (partial)${NC}"
    echo
}

################################################################################
# Download Gradle build cache dependencies
################################################################################

download_gradle_dependencies() {
    echo -e "${GREEN}Downloading Gradle build dependencies...${NC}"

    # Create a temporary Gradle project to download build dependencies
    local temp_project="$DOWNLOAD_DIR/dependencies/.temp-project"
    mkdir -p "$temp_project"

    cat > "$temp_project/build.gradle" << 'EOF'
plugins {
    id 'com.android.application' version '8.2.0'
    id 'org.jetbrains.kotlin.android' version '1.9.21'
}

android {
    compileSdk 34
    namespace 'com.example.offline'

    defaultConfig {
        applicationId "com.example.offline"
        minSdk 24
        targetSdk 34
        versionCode 1
        versionName "1.0"
    }
}

dependencies {
    implementation 'androidx.core:core-ktx:1.12.0'
    implementation 'androidx.appcompat:appcompat:1.6.1'
    implementation 'com.google.android.material:material:1.11.0'
}
EOF

    cat > "$temp_project/settings.gradle" << 'EOF'
pluginManagement {
    repositories {
        google()
        mavenCentral()
        gradlePluginPortal()
    }
}

dependencyResolutionManagement {
    repositoriesMode.set(RepositoriesMode.FAIL_ON_PROJECT_REPOS)
    repositories {
        google()
        mavenCentral()
    }
}

rootProject.name = "offline-temp"
EOF

    mkdir -p "$temp_project/src/main/java"

    cd "$temp_project"

    # Use a specific Gradle version
    local gradle_zip="$DOWNLOAD_DIR/gradle/distributions/gradle-8.4-all.zip"
    if [ -f "$gradle_zip" ]; then
        export GRADLE_USER_HOME="$DOWNLOAD_DIR/dependencies/.gradle-home"
        unzip -q "$gradle_zip" -d "$DOWNLOAD_DIR/dependencies/"
        "$DOWNLOAD_DIR/dependencies/gradle-8.4/bin/gradle" dependencies || true
    fi

    cd "$PROJECT_ROOT"

    echo -e "${GREEN}✓ Gradle dependencies cached${NC}"
    echo
}

################################################################################
# Create dependency info file
################################################################################

create_dependency_info() {
    echo -e "${GREEN}Creating dependency info file...${NC}"

    local info_file="$DOWNLOAD_DIR/dependencies/dependencies-info.txt"

    cat > "$info_file" << EOF
Dependencies Download Information
==================================

Download Date: $(date)

Directories:
-----------
Maven Repository: $MAVEN_DIR
Google Repository: $GOOGLE_DIR
Gradle Plugins: $GRADLE_PLUGINS_DIR

Size Information:
-----------------
EOF

    du -sh "$MAVEN_DIR" >> "$info_file" 2>/dev/null || echo "Maven: N/A" >> "$info_file"
    du -sh "$GOOGLE_DIR" >> "$info_file" 2>/dev/null || echo "Google: N/A" >> "$info_file"
    du -sh "$GRADLE_PLUGINS_DIR" >> "$info_file" 2>/dev/null || echo "Gradle Plugins: N/A" >> "$info_file"

    echo "" >> "$info_file"
    echo "Total Dependencies Size:" >> "$info_file"
    du -sh "$DOWNLOAD_DIR/dependencies" >> "$info_file"

    echo -e "${GREEN}✓ Dependency info saved${NC}"
    echo
}

################################################################################
# Main execution
################################################################################

main() {
    if [ "$DOWNLOAD_MAVEN_DEPS" != "true" ]; then
        echo -e "${YELLOW}Maven dependency download disabled in config${NC}"
        exit 0
    fi

    echo "Starting dependency downloads..."
    echo

    download_common_dependencies
    download_android_gradle_plugin
    download_from_sample_project
    # mirror_google_maven  # Uncomment if you want full mirror
    create_dependency_info

    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}Step 4 Complete!${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo
    echo "Total dependencies size:"
    du -sh "$DOWNLOAD_DIR/dependencies"
    echo
    echo -e "${YELLOW}Next step: Run ./scripts/05-create-package.sh${NC}"
}

main "$@"
