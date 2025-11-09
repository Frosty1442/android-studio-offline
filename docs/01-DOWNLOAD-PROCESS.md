# Download Process Guide

This guide walks you through downloading all components needed for the offline Android Studio installer.

## Prerequisites

### System Requirements (Download Machine)
- **Operating System**: Linux, macOS, or Windows with WSL
- **Disk Space**: 100+ GB free (actual size varies by selected components)
- **Internet**: Stable broadband connection
- **Tools Required**:
  - `wget` or `curl`
  - `unzip`
  - `tar`
  - `bash`

### Verify Tools

```bash
# Check if tools are installed
which wget curl unzip tar bash

# Install missing tools (Ubuntu/Debian)
sudo apt-get install wget curl unzip tar

# Install missing tools (macOS)
brew install wget
```

## Configuration

### Step 1: Copy Configuration Template

```bash
cd android-studio-offline
cp config/download.config.template config/download.config
```

### Step 2: Edit Configuration

Open `config/download.config` and customize:

```bash
# Essential settings
PLATFORM="linux"  # Options: linux, mac, mac_arm, windows
ANDROID_STUDIO_VERSION="2023.3.1.18"  # Check for latest version
API_LEVELS="33 34"  # Android versions to support
BUILD_TOOLS_VERSIONS="33.0.2 34.0.0"
GRADLE_VERSIONS="7.6 8.0 8.4"

# Optional components
DOWNLOAD_JDK="true"
DOWNLOAD_NDK="true"
DOWNLOAD_MAVEN_DEPS="true"
```

### Step 3: Review Package Lists

Edit these files to customize what gets downloaded:

- `config/sdk-packages.txt` - SDK components
- `config/gradle-versions.txt` - Gradle versions
- `config/common-dependencies.txt` - Maven dependencies

## Download Process

### Step 1: Download Android Studio and JDK

```bash
./scripts/01-download-android-studio.sh
```

**What this downloads:**
- Android Studio IDE (2-3 GB)
- JDK 17 (300-400 MB)
- SDK Command Line Tools (100 MB)

**Time estimate**: 15-30 minutes depending on connection

**Verification:**
```bash
ls -lh downloads/android-studio/
ls -lh downloads/jdk/
ls -lh downloads/sdk/
```

### Step 2: Download SDK Components

```bash
./scripts/02-download-sdk-components.sh
```

**What this downloads:**
- Platform Tools (adb, fastboot)
- Build Tools (multiple versions)
- Android Platforms (API levels)
- System Images (for emulator)
- NDK (Native Development Kit)
- Google Play Services
- SDK extras

**Size**: 10-30 GB depending on selections

**Time estimate**: 30-90 minutes

**Note**: You'll be prompted to accept SDK licenses. Type 'y' for each.

**Verification:**
```bash
cat downloads/sdk/sdk-info.txt
du -sh downloads/sdk/
```

### Step 3: Download Gradle Distributions

```bash
./scripts/03-download-gradle.sh
```

**What this downloads:**
- Multiple Gradle versions (bin and all distributions)
- Gradle wrapper files
- Offline configuration templates

**Size**: 1-3 GB

**Time estimate**: 10-20 minutes

**Verification:**
```bash
ls -lh downloads/gradle/distributions/
cat downloads/gradle/gradle-info.txt
```

### Step 4: Download Maven Dependencies

```bash
./scripts/04-download-dependencies.sh
```

**What this downloads:**
- Common AndroidX libraries
- Material Design components
- Kotlin libraries
- Popular third-party libraries
- Android Gradle Plugin
- Kotlin Gradle Plugin

**Size**: 5-20 GB depending on configuration

**Time estimate**: 30-120 minutes

**Optional**: If you have a sample Android project, you can pre-download its dependencies:

```bash
# In download.config, set:
SAMPLE_PROJECT_PATH="/path/to/your/android/project"
```

**Verification:**
```bash
du -sh downloads/dependencies/
cat downloads/dependencies/dependencies-info.txt
```

### Step 5: Create Package

```bash
./scripts/05-create-package.sh
```

**What this does:**
- Verifies all downloads
- Creates manifest and checksums
- Compresses everything into a portable package
- Creates split archives if needed (>4GB)

**Time estimate**: 20-60 minutes depending on size

**Output:**
- `packages/android-studio-offline-YYYYMMDD.tar.gz`

**Verification:**
```bash
ls -lh packages/
cat packages/PACKAGE-INFO.txt
```

## Download Sizes (Approximate)

| Component | Size Range | Required |
|-----------|-----------|----------|
| Android Studio | 2-3 GB | Yes |
| JDK | 300-400 MB | Optional |
| SDK Command Tools | 100 MB | Yes |
| SDK Components | 10-30 GB | Yes |
| Gradle | 1-3 GB | Yes |
| Maven Dependencies | 5-20 GB | Recommended |
| **Total** | **20-60 GB** | - |

**Compressed package**: Usually 40-60% of total size

## Customization Tips

### Minimal Installation

For a minimal installation (only essential components):

```bash
# In download.config
API_LEVELS="34"  # Latest only
BUILD_TOOLS_VERSIONS="34.0.0"  # Latest only
GRADLE_VERSIONS="8.4"  # Latest only
DOWNLOAD_NDK="false"
DOWNLOAD_MAVEN_DEPS="false"
SYSTEM_IMAGES=""  # Skip emulator images
```

**Size**: ~15-20 GB

### Full Installation

For complete Android development:

```bash
# In download.config
API_LEVELS="28 29 30 31 32 33 34"  # Multiple versions
BUILD_TOOLS_VERSIONS="30.0.3 31.0.0 32.0.0 33.0.2 34.0.0"
GRADLE_VERSIONS="7.0 7.2 7.4 7.6 8.0 8.2 8.4"
DOWNLOAD_NDK="true"
DOWNLOAD_MAVEN_DEPS="true"
SYSTEM_IMAGES="system-images;android-33;google_apis;x86_64 system-images;android-34;google_apis;x86_64"
```

**Size**: 60-100 GB

### For Specific Project

If you have an existing project:

1. Note the versions from your project:
   - `build.gradle` - Check `compileSdk`, `buildToolsVersion`
   - `gradle/wrapper/gradle-wrapper.properties` - Check Gradle version
   - `build.gradle` dependencies

2. Download only those versions:
```bash
API_LEVELS="33"  # From compileSdk
BUILD_TOOLS_VERSIONS="33.0.2"  # From buildToolsVersion
GRADLE_VERSIONS="8.0"  # From gradle-wrapper.properties
```

3. Set sample project path to pre-download exact dependencies:
```bash
SAMPLE_PROJECT_PATH="/path/to/project"
```

## Resuming Interrupted Downloads

All scripts support resuming interrupted downloads:

```bash
# In download.config
RESUME_DOWNLOADS="true"
```

Simply re-run the script that was interrupted. Already downloaded files will be skipped.

## Troubleshooting

### Download Failures

**Issue**: wget/curl fails with SSL errors

**Solution**:
```bash
# For wget
wget --no-check-certificate ...

# For curl
curl -k ...
```

**Issue**: Slow downloads

**Solution**:
```bash
# Enable parallel downloads
PARALLEL_DOWNLOADS="8"  # In download.config
```

### Disk Space

**Issue**: Running out of disk space

**Solution**:
1. Download to external drive
2. Use minimal configuration
3. Download components separately

```bash
# Download to external drive
DOWNLOAD_DIR="/mnt/external/downloads"
```

### SDK Manager Issues

**Issue**: sdkmanager fails to accept licenses

**Solution**:
```bash
# Manually accept licenses
cd downloads/sdk/android-sdk
yes | cmdline-tools/latest/bin/sdkmanager --licenses
```

### Checksum Verification

If you want to verify downloads:

```bash
# Check SHA256 checksums
cd downloads
sha256sum -c SHA256SUMS.txt
```

## Next Steps

After downloading:

1. Verify the package was created successfully
2. Check `packages/PACKAGE-INFO.txt` for details
3. Transfer to offline machine (see [Installation Guide](03-INSTALLATION-GUIDE.md))
4. For split packages, transfer all parts

## Updating Components

To update an existing download:

```bash
# Delete specific component
rm -rf downloads/android-studio/

# Re-run specific script
./scripts/01-download-android-studio.sh

# Recreate package
./scripts/05-create-package.sh
```

## Tips for Large Downloads

1. **Use a download manager**: For very large downloads, consider using a download manager
2. **Schedule during off-peak**: Downloads are faster during off-peak hours
3. **Check mirrors**: Some components have mirror servers
4. **Incremental approach**: Download components over multiple days if needed
