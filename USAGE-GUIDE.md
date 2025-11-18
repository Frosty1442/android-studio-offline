# Complete Usage Guide

End-to-end guide for using the Android Studio Offline Installer.

## Overview

This tool downloads Android Studio and all dependencies on an internet-connected machine, packages them, and installs on an offline/air-gapped system.

**Workflow:**
```
Internet Machine → Download → Package → Transfer → Offline Machine → Install
```

## Prerequisites

### Internet-Connected Machine

- **OS**: Linux, macOS, or Windows
- **Disk Space**: 100+ GB free
- **Internet**: Broadband connection
- **Tools**:
  - Go 1.21+ (to build from source)
  - OR use pre-built binary

### Offline Machine

- **OS**: Linux, macOS, or Windows
- **Disk Space**: 50+ GB free
- **RAM**: 8+ GB (16 GB recommended)
- **CPU**: 64-bit processor

## Installation

### Option 1: Download Pre-Built Binary

```bash
# Linux
wget https://github.com/YOUR_REPO/releases/download/v1.0.0/android-offline-linux-amd64
chmod +x android-offline-linux-amd64
sudo mv android-offline-linux-amd64 /usr/local/bin/android-offline

# macOS Intel
wget https://github.com/YOUR_REPO/releases/download/v1.0.0/android-offline-darwin-amd64
chmod +x android-offline-darwin-amd64
sudo mv android-offline-darwin-amd64 /usr/local/bin/android-offline

# macOS Apple Silicon
wget https://github.com/YOUR_REPO/releases/download/v1.0.0/android-offline-darwin-arm64
chmod +x android-offline-darwin-arm64
sudo mv android-offline-darwin-arm64 /usr/local/bin/android-offline

# Windows
# Download android-offline-windows-amd64.exe
# Add to PATH or run from current directory
```

### Option 2: Build from Source

```bash
git clone https://github.com/YOUR_REPO/android-studio-offline.git
cd android-studio-offline
make build
sudo make install
```

## Complete Workflow

### Phase 1: Configuration (Internet Machine)

#### Step 1.1: Initialize Configuration

```bash
android-offline init
```

This creates `config.yaml` with defaults.

#### Step 1.2: Customize Configuration

Edit `config.yaml`:

```yaml
# Essential settings
platform: linux  # Change based on target: linux, darwin, darwin-arm, windows

android_studio:
  version: "2024.2.1.10"  # Check for latest version
  download_jdk: true
  jdk_version: "17"

sdk:
  api_levels: [33, 34]  # Android versions you need
  build_tools: ["33.0.2", "34.0.0"]
  system_images:
    - "system-images;android-34;google_apis;x86_64"
  download_ndk: true
  download_cmake: true

gradle:
  versions: ["8.0", "8.4"]  # Versions your projects use

dependencies:
  download_maven: true
  common_libs:
    - "androidx.core:core-ktx:1.12.0"
    - "com.google.android.material:material:1.11.0"
    # Add your project dependencies

options:
  parallel_downloads: 4
  verify_checksums: true
```

**Customization Tips:**

**Minimal Setup** (20-25 GB):
```yaml
api_levels: [34]
build_tools: ["34.0.0"]
gradle_versions: ["8.4"]
download_ndk: false
system_images: []  # Skip emulator
```

**Full Setup** (60-80 GB):
```yaml
api_levels: [30, 31, 32, 33, 34]
build_tools: ["30.0.3", "31.0.0", "32.0.0", "33.0.2", "34.0.0"]
gradle_versions: ["7.4", "7.6", "8.0", "8.2", "8.4"]
download_ndk: true
system_images:
  - "system-images;android-33;google_apis;x86_64"
  - "system-images;android-34;google_apis;x86_64"
```

### Phase 2: Download (Internet Machine)

#### Step 2.1: Download All Components

```bash
android-offline download
```

**What this downloads:**

- Android Studio IDE (~2-3 GB)
- JDK 17 (~300 MB)
- SDK Command Line Tools (~100 MB)
- Platform Tools (~10 MB)
- Gradle Distributions (~500 MB - 2 GB)
- SDK Packages:
  - Android Platforms (50-80 MB each)
  - Build Tools (50-70 MB each)
  - System Images (400 MB - 2 GB each)
  - NDK (~1-2 GB, if enabled)
  - CMake (~50 MB, if enabled)
- Maven Dependencies (~2-10 GB)
- Android Gradle Plugin
- Kotlin Gradle Plugin

**Expected Output:**

```
[1/5] Downloading Android Studio 2024.2.1.10
ℹ Downloading from: https://...
android-studio-2024.2.1.10... [████████████████░░░░] 67.32% 1.8GB/2.7GB
✓ Android Studio downloaded successfully

[2/5] Downloading JDK 17
...

[5/5] Downloading SDK packages
ℹ Downloading API level 33...
ℹ Downloading API level 34...
✓ SDK packages downloaded

ℹ Downloading Maven dependencies...
✓ All Maven dependencies downloaded

✓ All downloads completed successfully!
ℹ Next step: Run 'android-offline package' to create installation package
```

**Time Estimate**: 1-3 hours depending on:
- Internet speed
- Components selected
- Number of parallel downloads

**Download Progress:**
- Real-time progress bars
- Download speed shown
- Can resume if interrupted (Ctrl+C then re-run)

#### Step 2.2: Verify Downloads

```bash
android-offline verify
```

**Output:**
```
ℹ Verifying downloaded components...

✓ Android Studio: Found
✓ JDK: Found
✓ SDK Tools: Found
✓ Gradle: Found (18 files)
✓ Dependencies: Found
ℹ Total download size: 45.3 GB

✓ Verification complete: All components present!
ℹ Ready to create package with 'android-offline package'
```

### Phase 3: Package (Internet Machine)

#### Step 3.1: Create Installation Package

```bash
android-offline package
```

**What this does:**

- Creates `packages/android-studio-offline-YYYYMMDD.tar.gz`
- Includes all downloads
- Adds installation scripts
- Compresses everything

**Expected Output:**

```
ℹ Creating package: android-studio-offline-20241118.tar.gz
ℹ This may take several minutes...
✓ Package created successfully
✓ Package created: packages/android-studio-offline-20241118.tar.gz
ℹ Transfer this file to your offline machine and run 'android-offline install'
```

**Package Size**: 15-50 GB compressed (depending on components)

**Time Estimate**: 20-60 minutes depending on:
- Total download size
- CPU speed
- Disk I/O

### Phase 4: Transfer

Transfer the package to your offline machine using:

**USB Drive:**
```bash
cp packages/android-studio-offline-*.tar.gz /media/usb/
```

**External Hard Drive:**
```bash
rsync -avh --progress packages/android-studio-offline-*.tar.gz /mnt/external/
```

**Internal Network (if available):**
```bash
scp packages/android-studio-offline-*.tar.gz user@offline-machine:~/
```

**For Large Packages (>4GB):**

If package is split automatically:
```bash
# Transfer all parts
scp packages/android-studio-offline-*.tar.gz.part* user@offline-machine:~/

# On offline machine, join:
cat android-studio-offline-*.tar.gz.part* > android-studio-offline.tar.gz
```

### Phase 5: Installation (Offline Machine)

#### Step 5.1: Extract Package

```bash
cd ~
tar -xzf android-studio-offline-20241118.tar.gz
cd downloads
```

#### Step 5.2: Verify android-offline Binary

```bash
ls -lh
# Should see: bin/android-offline (or android-offline.exe on Windows)

# Make executable (Linux/macOS)
chmod +x bin/android-offline

# Or use the transferred binary
```

#### Step 5.3: Initialize Configuration

If you didn't transfer config with package:

```bash
./bin/android-offline init
```

Edit `config.yaml` to set install directory:

```yaml
install_dir: /opt/android-offline  # Or ~/android-studio for user install
```

#### Step 5.4: Install

```bash
./bin/android-offline install
```

**What this does:**

1. Creates installation directories
2. Extracts Android Studio
3. Extracts JDK (if downloaded)
4. Copies Android SDK
5. Copies Gradle distributions
6. Copies Maven dependencies
7. Configures environment variables
8. Creates desktop launcher (Linux)

**Expected Output:**

```
ℹ Starting Android Studio offline installation
ℹ Install directory: /opt/android-offline

[1/8] Create directories
✓ Directories created

[2/8] Install Android Studio
ℹ Extracting: android-studio-2024.2.1.10-linux.tar.gz
✓ Android Studio installed

[3/8] Install JDK
✓ JDK installed

[4/8] Install Android SDK
ℹ Copying Android SDK...
✓ Android SDK installed

[5/8] Install Gradle
✓ Gradle installed

[6/8] Install dependencies
✓ Dependencies installed

[7/8] Configure environment
✓ Environment configured
ℹ Reload your shell: source ~/.bashrc

[8/8] Create launcher
✓ Desktop launcher created

✓ Installation completed successfully!
ℹ Android Studio installed to: /opt/android-offline
ℹ Start Android Studio: /opt/android-offline/android-studio/bin/studio.sh
ℹ Or use the desktop launcher (Linux only)
```

**Time Estimate**: 10-20 minutes

#### Step 5.5: Reload Environment

```bash
# Bash
source ~/.bashrc

# Zsh
source ~/.zshrc

# Or start a new terminal
```

#### Step 5.6: Verify Installation

```bash
# Check environment variables
echo $ANDROID_HOME
# Should show: /opt/android-offline/android-sdk

# Check adb
which adb
adb --version

# Check Android Studio
ls /opt/android-offline/android-studio/bin/studio.sh
```

### Phase 6: First Run

#### Step 6.1: Launch Android Studio

**Linux/macOS:**
```bash
/opt/android-offline/android-studio/bin/studio.sh

# Or from desktop launcher
# Applications → Development → Android Studio (Offline)
```

**Windows:**
```cmd
C:\android-offline\android-studio\bin\studio64.exe
```

#### Step 6.2: Initial Setup Wizard

1. **Import Settings**: Choose "Do not import settings"

2. **Data Sharing**: Choose "Don't send"

3. **Install Type**: Select **"Custom"**

4. **SDK Location**:
   - Click "..." next to Android SDK Location
   - Set to: `/opt/android-offline/android-sdk`
   - ✅ **DO NOT** check "Download SDK components"

5. **Verify Settings**: Click "Finish"

#### Step 6.3: Disable Updates

**IMPORTANT**: Prevent Android Studio from trying to update in offline mode.

1. **File → Settings** (or **Android Studio → Preferences** on macOS)

2. **Appearance & Behavior → System Settings → Updates**
   - Uncheck "Automatically check updates"
   - Set to "Stable Channel" and "Never check"

3. **Appearance & Behavior → System Settings → Android SDK**
   - Verify: `/opt/android-offline/android-sdk`

4. **Build, Execution, Deployment → Gradle**
   - ✅ Enable "Offline mode"

5. Click **Apply** and **OK**

### Phase 7: Create First Project

#### Step 7.1: New Project

1. Click **"New Project"**

2. Select **"Empty Activity"**

3. Configure:
   ```
   Name: HelloOffline
   Package: com.example.hellooffline
   Save location: ~/AndroidStudioProjects/HelloOffline
   Language: Kotlin
   Minimum SDK: API 24 (or your choice)
   ```

4. Click **Finish**

#### Step 7.2: Wait for Gradle Sync

First sync takes 2-5 minutes as Gradle:
- Sets up project
- Indexes dependencies
- Configures build

Watch bottom status bar for "Gradle sync finished"

#### Step 7.3: Build Project

```bash
# From Android Studio: Build → Build Bundle(s)/APK(s) → Build APK

# Or command line:
cd ~/AndroidStudioProjects/HelloOffline
./gradlew build --offline
```

**Expected Output:**
```
BUILD SUCCESSFUL in 45s
```

✅ **Success!** You're building Android apps offline!

### Phase 8: Run on Emulator

#### Step 8.1: Create AVD

1. **Tools → Device Manager**

2. Click **"Create Device"**

3. **Hardware**: Select "Pixel 5"

4. **System Image**:
   - Select from **downloaded images** (e.g., API 34)
   - Should see "System Image: android-34"

5. **Verify**: Click **Finish**

#### Step 8.2: Launch Emulator

1. Click **▶ Play** button next to AVD

2. Wait for boot (first boot takes 2-3 minutes)

#### Step 8.3: Run App

1. With emulator running, click **▶ Run** in Android Studio

2. Select running emulator

3. App will build, install, and launch

## Troubleshooting

### Download Issues

**Problem**: Download fails

**Solutions:**
```bash
# Downloads are resumable, just re-run:
android-offline download

# Use verbose mode for debugging:
android-offline download -v

# Check config:
cat config.yaml
```

**Problem**: "Permission denied" when saving

**Solution:**
```bash
# Run from directory with write permissions
cd ~/android-downloads
android-offline download
```

### Installation Issues

**Problem**: "Permission denied" during install

**Solution 1** - User installation:
```yaml
# In config.yaml
install_dir: ~/android-studio-offline
```

**Solution 2** - System installation:
```bash
# Create directory with sudo
sudo mkdir -p /opt/android-offline
sudo chown $USER:$USER /opt/android-offline

# Then install
android-offline install
```

**Problem**: Environment variables not set

**Solution:**
```bash
# Manually add to ~/.bashrc:
export ANDROID_HOME=/opt/android-offline/android-sdk
export PATH=$PATH:$ANDROID_HOME/platform-tools

# Reload:
source ~/.bashrc
```

### Build Issues

**Problem**: "Could not find dependency"

**Solution:**
```bash
# Check offline repos are configured:
cat ~/.gradle/init.d/offline-repos.gradle

# Force offline mode in project:
# Add to gradle.properties:
org.gradle.offline=true

# Use --offline flag:
./gradlew build --offline
```

**Problem**: Gradle wrapper tries to download

**Solution:**
```groovy
// In gradle/wrapper/gradle-wrapper.properties:
distributionUrl=file\:///opt/android-offline/gradle/distributions/gradle-8.4-all.zip
```

### Emulator Issues

**Problem**: Emulator won't start

**Solution (Linux):**
```bash
# Enable KVM:
sudo usermod -aG kvm $USER
# Log out and back in

# Or use software rendering:
/opt/android-offline/android-sdk/emulator/emulator -avd Pixel_5_API_34 -gpu swiftshader_indirect
```

**Solution (macOS):**
```bash
# Uses Hypervisor.framework automatically
# Check permissions: System Preferences → Security
```

## Tips & Best Practices

### 1. Keep Package Updated

Recreate package every few months for latest versions:

```bash
# Update config.yaml with new versions
android-offline download
android-offline package
```

### 2. Multiple Configurations

Create different configs for different use cases:

```bash
# Minimal config
android-offline init --config minimal.yaml
# Edit minimal.yaml for bare minimum

# Full config
android-offline init --config full.yaml
# Edit full.yaml for everything

# Download with specific config
android-offline download --config minimal.yaml
```

### 3. Project-Specific Dependencies

Add your project dependencies to config:

```yaml
dependencies:
  common_libs:
    - "com.your-company:library:1.0.0"
    - "io.ktor:ktor-client:2.3.0"
    # etc.
```

### 4. Verify Before Transfer

Always verify before creating package:

```bash
android-offline verify
```

### 5. Document Your Setup

Keep notes about:
- Android Studio version
- SDK versions
- Which projects work with this package
- Custom dependencies

## Command Reference

### Global Flags

```bash
-c, --config string   # Config file (default "config.yaml")
-v, --verbose         # Verbose output
--version             # Show version
```

### Commands

```bash
android-offline init              # Create config file
android-offline download          # Download components
android-offline package           # Create package
android-offline verify            # Verify downloads
android-offline install           # Install on offline machine
android-offline help [command]    # Help for command
```

### Examples

```bash
# Custom config location
android-offline download --config /path/to/config.yaml

# Verbose mode
android-offline download -v

# Multiple commands
android-offline download && android-offline verify && android-offline package
```

## Support

- **Documentation**: See [docs/](docs/) directory
- **Bash Version**: Original scripts in [scripts/](scripts/)
- **GitHub**: Report issues and contribute

## Success Checklist

- [ ] Downloaded all components
- [ ] Verified downloads
- [ ] Created package
- [ ] Transferred to offline machine
- [ ] Installed successfully
- [ ] Environment configured
- [ ] Android Studio launches
- [ ] Created test project
- [ ] Built project offline
- [ ] Emulator works (optional)

**You're ready to develop Android apps offline!** 🎉
