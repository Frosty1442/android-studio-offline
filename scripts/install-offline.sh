#!/bin/bash

################################################################################
# Offline Installation Script
# Installs Android Studio and all components on an offline machine
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

# Default installation directory
DEFAULT_INSTALL_DIR="/opt/android-offline"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Android Studio Offline Installer${NC}"
echo -e "${BLUE}Installation Script${NC}"
echo -e "${BLUE}========================================${NC}"
echo

################################################################################
# Get installation preferences
################################################################################

get_install_preferences() {
    echo -e "${GREEN}Installation Configuration${NC}"
    echo

    # Installation directory
    read -p "Installation directory [$DEFAULT_INSTALL_DIR]: " INSTALL_DIR
    INSTALL_DIR="${INSTALL_DIR:-$DEFAULT_INSTALL_DIR}"

    # User home directory for configs
    USER_HOME="${HOME}"

    echo
    echo "Installation Summary:"
    echo "  Installation Directory: $INSTALL_DIR"
    echo "  User Home: $USER_HOME"
    echo "  Android Studio: $INSTALL_DIR/android-studio"
    echo "  Android SDK: $INSTALL_DIR/android-sdk"
    echo "  Gradle: $INSTALL_DIR/gradle"
    echo

    read -p "Continue with installation? (y/n): " CONFIRM
    if [ "$CONFIRM" != "y" ] && [ "$CONFIRM" != "Y" ]; then
        echo "Installation cancelled."
        exit 0
    fi

    echo
}

################################################################################
# Check system requirements
################################################################################

check_requirements() {
    echo -e "${GREEN}Checking system requirements...${NC}"

    # Check disk space
    local available_space=$(df -BG "$HOME" | awk 'NR==2 {print $4}' | sed 's/G//')
    if [ "$available_space" -lt 50 ]; then
        echo -e "${RED}Warning: Less than 50GB available. Installation may fail.${NC}"
        read -p "Continue anyway? (y/n): " CONFIRM
        [ "$CONFIRM" != "y" ] && [ "$CONFIRM" != "Y" ] && exit 1
    else
        echo -e "${GREEN}  ✓ Sufficient disk space${NC}"
    fi

    # Check for required tools
    for cmd in tar unzip; do
        if ! command -v $cmd &> /dev/null; then
            echo -e "${RED}  ✗ $cmd not found. Please install it first.${NC}"
            exit 1
        else
            echo -e "${GREEN}  ✓ $cmd found${NC}"
        fi
    done

    echo
}

################################################################################
# Create installation directories
################################################################################

create_directories() {
    echo -e "${GREEN}Creating installation directories...${NC}"

    sudo mkdir -p "$INSTALL_DIR"
    sudo chown -R $USER:$(id -gn) "$INSTALL_DIR"

    mkdir -p "$INSTALL_DIR/android-studio"
    mkdir -p "$INSTALL_DIR/android-sdk"
    mkdir -p "$INSTALL_DIR/gradle"
    mkdir -p "$INSTALL_DIR/jdk"
    mkdir -p "$INSTALL_DIR/maven-repo"
    mkdir -p "$INSTALL_DIR/google-repo"
    mkdir -p "$INSTALL_DIR/gradle-plugins"

    echo -e "${GREEN}✓ Directories created${NC}"
    echo
}

################################################################################
# Install Android Studio
################################################################################

install_android_studio() {
    echo -e "${GREEN}Installing Android Studio...${NC}"

    local studio_archive=$(find "$SCRIPT_DIR/android-studio" -name "android-studio-*" -type f | head -1)

    if [ -z "$studio_archive" ]; then
        echo -e "${YELLOW}Android Studio archive not found, skipping...${NC}"
        return
    fi

    echo "Found: $(basename "$studio_archive")"

    # Extract based on file type
    case "$studio_archive" in
        *.tar.gz)
            tar -xzf "$studio_archive" -C "$INSTALL_DIR/"
            ;;
        *.zip)
            unzip -q "$studio_archive" -d "$INSTALL_DIR/"
            ;;
        *)
            echo -e "${YELLOW}Unsupported archive format. Please install manually.${NC}"
            echo "File: $studio_archive"
            return
            ;;
    esac

    echo -e "${GREEN}✓ Android Studio installed${NC}"
    echo
}

################################################################################
# Install JDK
################################################################################

install_jdk() {
    echo -e "${GREEN}Installing JDK...${NC}"

    local jdk_archive=$(find "$SCRIPT_DIR/jdk" -name "jdk-*" -type f 2>/dev/null | head -1)

    if [ -z "$jdk_archive" ]; then
        echo -e "${YELLOW}JDK archive not found, skipping...${NC}"
        echo -e "${YELLOW}Android Studio includes a bundled JDK${NC}"
        return
    fi

    echo "Found: $(basename "$jdk_archive")"

    case "$jdk_archive" in
        *.tar.gz)
            tar -xzf "$jdk_archive" -C "$INSTALL_DIR/jdk/"
            ;;
        *.zip)
            unzip -q "$jdk_archive" -d "$INSTALL_DIR/jdk/"
            ;;
    esac

    # Find the extracted JDK directory
    JDK_HOME=$(find "$INSTALL_DIR/jdk" -maxdepth 2 -name "bin" -type d | head -1 | xargs dirname)

    echo -e "${GREEN}✓ JDK installed${NC}"
    echo
}

################################################################################
# Install Android SDK
################################################################################

install_android_sdk() {
    echo -e "${GREEN}Installing Android SDK...${NC}"

    if [ -d "$SCRIPT_DIR/sdk/android-sdk" ]; then
        echo "Copying SDK files..."
        cp -r "$SCRIPT_DIR/sdk/android-sdk"/* "$INSTALL_DIR/android-sdk/"
        echo -e "${GREEN}✓ Android SDK installed${NC}"
    else
        echo -e "${YELLOW}Android SDK not found in package${NC}"
        echo "You can install it later using Android Studio's SDK Manager"
    fi

    echo
}

################################################################################
# Install Gradle
################################################################################

install_gradle() {
    echo -e "${GREEN}Installing Gradle distributions...${NC}"

    if [ -d "$SCRIPT_DIR/gradle/distributions" ]; then
        echo "Copying Gradle distributions..."
        cp -r "$SCRIPT_DIR/gradle/distributions" "$INSTALL_DIR/gradle/"

        # Also copy wrapper files
        if [ -d "$SCRIPT_DIR/gradle/wrapper" ]; then
            cp -r "$SCRIPT_DIR/gradle/wrapper" "$INSTALL_DIR/gradle/"
        fi

        # Copy init scripts
        if [ -d "$SCRIPT_DIR/gradle/init.d" ]; then
            mkdir -p "$USER_HOME/.gradle/init.d"
            cp "$SCRIPT_DIR/gradle/init.d"/* "$USER_HOME/.gradle/init.d/" 2>/dev/null || true
        fi

        echo -e "${GREEN}✓ Gradle installed${NC}"
    else
        echo -e "${YELLOW}Gradle distributions not found${NC}"
    fi

    echo
}

################################################################################
# Install Maven Dependencies
################################################################################

install_dependencies() {
    echo -e "${GREEN}Installing Maven dependencies...${NC}"

    # Install Maven repository
    if [ -d "$SCRIPT_DIR/dependencies/maven-repo" ]; then
        echo "Copying Maven repository..."
        cp -r "$SCRIPT_DIR/dependencies/maven-repo"/* "$INSTALL_DIR/maven-repo/" 2>/dev/null || true
        echo -e "${GREEN}✓ Maven repository installed${NC}"
    fi

    # Install Google repository
    if [ -d "$SCRIPT_DIR/dependencies/google-repo" ]; then
        echo "Copying Google repository..."
        cp -r "$SCRIPT_DIR/dependencies/google-repo"/* "$INSTALL_DIR/google-repo/" 2>/dev/null || true
        echo -e "${GREEN}✓ Google repository installed${NC}"
    fi

    # Install Gradle plugins
    if [ -d "$SCRIPT_DIR/dependencies/gradle-plugins" ]; then
        echo "Copying Gradle plugins..."
        cp -r "$SCRIPT_DIR/dependencies/gradle-plugins"/* "$INSTALL_DIR/gradle-plugins/" 2>/dev/null || true
        echo -e "${GREEN}✓ Gradle plugins installed${NC}"
    fi

    echo
}

################################################################################
# Configure environment variables
################################################################################

configure_environment() {
    echo -e "${GREEN}Configuring environment variables...${NC}"

    local shell_rc=""

    # Detect shell
    if [ -n "$BASH_VERSION" ]; then
        shell_rc="$USER_HOME/.bashrc"
    elif [ -n "$ZSH_VERSION" ]; then
        shell_rc="$USER_HOME/.zshrc"
    else
        shell_rc="$USER_HOME/.profile"
    fi

    # Backup existing rc file
    if [ -f "$shell_rc" ]; then
        cp "$shell_rc" "$shell_rc.backup-$(date +%Y%m%d)"
    fi

    # Add environment variables
    cat >> "$shell_rc" << EOF

# Android Studio Offline Environment
# Added by Android Studio Offline Installer on $(date)
export ANDROID_HOME="$INSTALL_DIR/android-sdk"
export ANDROID_SDK_ROOT="$INSTALL_DIR/android-sdk"
export PATH="\$PATH:\$ANDROID_HOME/platform-tools:\$ANDROID_HOME/tools:\$ANDROID_HOME/cmdline-tools/latest/bin"

# Android Studio
export STUDIO_HOME="$INSTALL_DIR/android-studio"
export PATH="\$PATH:\$STUDIO_HOME/bin"

# JDK (if installed separately)
EOF

    if [ -n "$JDK_HOME" ]; then
        cat >> "$shell_rc" << EOF
export JAVA_HOME="$JDK_HOME"
export PATH="\$PATH:\$JAVA_HOME/bin"
EOF
    fi

    cat >> "$shell_rc" << EOF

# Gradle
export GRADLE_USER_HOME="$USER_HOME/.gradle"

# Offline Maven repositories
export ANDROID_OFFLINE_ROOT="$INSTALL_DIR"
EOF

    echo -e "${GREEN}✓ Environment configured${NC}"
    echo -e "${YELLOW}Shell configuration file: $shell_rc${NC}"
    echo
}

################################################################################
# Configure Gradle for offline use
################################################################################

configure_gradle_offline() {
    echo -e "${GREEN}Configuring Gradle for offline use...${NC}"

    mkdir -p "$USER_HOME/.gradle"

    # Create gradle.properties
    cat > "$USER_HOME/.gradle/gradle.properties" << EOF
# Gradle Offline Configuration
# Created by Android Studio Offline Installer

# Performance optimizations
org.gradle.daemon=true
org.gradle.parallel=true
org.gradle.caching=true
org.gradle.configureondemand=true

# Memory settings
org.gradle.jvmargs=-Xmx4096m -XX:MaxMetaspaceSize=512m -XX:+HeapDumpOnOutOfMemoryError

# Android specific
android.useAndroidX=true
android.enableJetifier=true

# Offline mode repositories
android.offlineRoot=$INSTALL_DIR
EOF

    # Update init script with correct paths
    if [ -f "$USER_HOME/.gradle/init.d/offline-repos.gradle" ]; then
        sed -i.bak "s|/opt/android-offline|$INSTALL_DIR|g" "$USER_HOME/.gradle/init.d/offline-repos.gradle"
    fi

    echo -e "${GREEN}✓ Gradle configured for offline use${NC}"
    echo
}

################################################################################
# Create desktop entry (Linux)
################################################################################

create_desktop_entry() {
    if [ "$(uname)" != "Linux" ]; then
        return
    fi

    echo -e "${GREEN}Creating desktop entry...${NC}"

    mkdir -p "$USER_HOME/.local/share/applications"

    cat > "$USER_HOME/.local/share/applications/android-studio-offline.desktop" << EOF
[Desktop Entry]
Version=1.0
Type=Application
Name=Android Studio (Offline)
Icon=$INSTALL_DIR/android-studio/bin/studio.png
Exec=$INSTALL_DIR/android-studio/bin/studio.sh
Comment=Android Studio IDE (Offline Installation)
Categories=Development;IDE;
Terminal=false
StartupWMClass=jetbrains-studio
EOF

    chmod +x "$USER_HOME/.local/share/applications/android-studio-offline.desktop"

    echo -e "${GREEN}✓ Desktop entry created${NC}"
    echo
}

################################################################################
# Create quick start guide
################################################################################

create_quick_start_guide() {
    echo -e "${GREEN}Creating quick start guide...${NC}"

    cat > "$INSTALL_DIR/QUICKSTART.txt" << EOF
Android Studio Offline Installation
====================================

Installation completed successfully!

Installed Components:
--------------------
- Android Studio: $INSTALL_DIR/android-studio
- Android SDK: $INSTALL_DIR/android-sdk
- Gradle: $INSTALL_DIR/gradle
- Maven Repositories: $INSTALL_DIR/maven-repo, $INSTALL_DIR/google-repo

Starting Android Studio:
-----------------------
Run: $INSTALL_DIR/android-studio/bin/studio.sh

Or use the desktop launcher (Linux only)

Environment Variables:
---------------------
ANDROID_HOME=$INSTALL_DIR/android-sdk
ANDROID_SDK_ROOT=$INSTALL_DIR/android-sdk
STUDIO_HOME=$INSTALL_DIR/android-studio

First Run:
----------
1. Start Android Studio
2. Choose "Do not import settings" (if first time)
3. Select "Custom" setup
4. Point SDK location to: $INSTALL_DIR/android-sdk
5. Disable automatic updates and telemetry for offline use

Creating New Projects:
---------------------
1. When creating a project, Android Studio will use the offline SDK
2. Gradle will use the pre-downloaded distributions
3. Dependencies will be resolved from local Maven repositories

Using Gradle Offline:
--------------------
Add to your project's gradle.properties:

offline.mode=true

Or run with: ./gradlew build --offline

Troubleshooting:
---------------
If dependencies are not found:
1. Check that ~/.gradle/init.d/offline-repos.gradle exists
2. Verify repository paths in the init script
3. Ensure ANDROID_OFFLINE_ROOT is set to: $INSTALL_DIR

For Gradle wrapper in new projects:
1. Use the templates in: $INSTALL_DIR/gradle/wrapper/
2. Update distributionUrl to point to local Gradle distribution

Need to reload environment variables?
Run: source ~/.bashrc  (or ~/.zshrc)

Support:
--------
For issues, refer to the offline installer documentation
or check the original repository.

EOF

    echo -e "${GREEN}✓ Quick start guide created: $INSTALL_DIR/QUICKSTART.txt${NC}"
    echo
}

################################################################################
# Verify installation
################################################################################

verify_installation() {
    echo -e "${GREEN}Verifying installation...${NC}"

    local errors=0

    # Check Android Studio
    if [ -f "$INSTALL_DIR/android-studio/bin/studio.sh" ]; then
        echo -e "${GREEN}  ✓ Android Studio executable found${NC}"
    else
        echo -e "${RED}  ✗ Android Studio not found${NC}"
        ((errors++))
    fi

    # Check SDK
    if [ -d "$INSTALL_DIR/android-sdk" ]; then
        echo -e "${GREEN}  ✓ Android SDK directory exists${NC}"
    else
        echo -e "${YELLOW}  ⚠ Android SDK directory not found${NC}"
    fi

    # Check environment
    if grep -q "ANDROID_HOME" "$USER_HOME/.bashrc" 2>/dev/null || \
       grep -q "ANDROID_HOME" "$USER_HOME/.zshrc" 2>/dev/null; then
        echo -e "${GREEN}  ✓ Environment variables configured${NC}"
    else
        echo -e "${YELLOW}  ⚠ Environment variables may not be configured${NC}"
    fi

    if [ $errors -eq 0 ]; then
        echo -e "${GREEN}✓ Installation verified successfully${NC}"
    else
        echo -e "${YELLOW}⚠ Installation completed with warnings${NC}"
    fi

    echo
}

################################################################################
# Main installation
################################################################################

main() {
    get_install_preferences
    check_requirements
    create_directories
    install_android_studio
    install_jdk
    install_android_sdk
    install_gradle
    install_dependencies
    configure_environment
    configure_gradle_offline
    create_desktop_entry
    create_quick_start_guide
    verify_installation

    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}Installation Complete!${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo
    echo -e "${BLUE}Next Steps:${NC}"
    echo "1. Reload your shell configuration:"
    echo "   ${YELLOW}source ~/.bashrc${NC}  (or source ~/.zshrc)"
    echo
    echo "2. Start Android Studio:"
    echo "   ${YELLOW}$INSTALL_DIR/android-studio/bin/studio.sh${NC}"
    echo
    echo "3. Read the quick start guide:"
    echo "   ${YELLOW}cat $INSTALL_DIR/QUICKSTART.txt${NC}"
    echo
    echo -e "${GREEN}Happy coding!${NC}"
}

main "$@"
