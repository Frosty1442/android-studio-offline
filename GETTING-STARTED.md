# Getting Started with Android Studio Offline Installer

Quick start guide to get you up and running quickly.

## Overview

This tool helps you create a complete offline Android development environment by:
1. Downloading all necessary components on an internet-connected machine
2. Packaging everything into a portable installer
3. Installing on an air-gapped/offline machine

## Prerequisites

### On Internet-Connected Machine (Download)
- Linux, macOS, or Windows with WSL
- 100+ GB free disk space
- Broadband internet connection
- Basic command-line knowledge

### On Offline Machine (Installation)
- Linux, macOS, or Windows
- 50+ GB free disk space
- 8+ GB RAM (16 GB recommended)

## Quick Start (5 Steps)

### Step 1: Clone Repository

```bash
git clone https://github.com/yourusername/android-studio-offline.git
cd android-studio-offline
```

### Step 2: Configure

```bash
# Copy configuration template
cp config/download.config.template config/download.config

# Edit for your needs (or use defaults)
nano config/download.config
```

**Minimal changes needed**:
```bash
PLATFORM="linux"  # Change to: linux, mac, mac_arm, or windows
```

### Step 3: Download Components

Run all download scripts sequentially:

```bash
# This will take 1-3 hours depending on internet speed
./scripts/01-download-android-studio.sh
./scripts/02-download-sdk-components.sh
./scripts/03-download-gradle.sh
./scripts/04-download-dependencies.sh
./scripts/05-create-package.sh
```

**Or use one command** (if all scripts have no issues):
```bash
# Run all at once
for script in scripts/0*.sh; do bash "$script"; done
```

Package will be created in `packages/` directory.

### Step 4: Transfer to Offline Machine

Copy the package to your offline machine:

```bash
# Using USB drive
cp packages/android-studio-offline-*.tar.gz /media/usb/

# Or network (if available)
scp packages/android-studio-offline-*.tar.gz user@offline-machine:~/
```

### Step 5: Install on Offline Machine

```bash
# Extract package
tar -xzf android-studio-offline-*.tar.gz
cd downloads

# Run installer
./install-offline.sh

# Follow prompts (press Enter for defaults)
```

After installation:
```bash
# Reload environment
source ~/.bashrc

# Start Android Studio
/opt/android-offline/android-studio/bin/studio.sh
```

## First Project

### Create New Project

1. **Start Android Studio**
2. Click **"New Project"**
3. Select **"Empty Activity"**
4. Configure:
   - Name: `HelloOffline`
   - Package: `com.example.hellooffline`
   - Language: `Kotlin`
   - Minimum SDK: `API 24`
5. Click **Finish**
6. Wait for Gradle sync
7. Click **▶ Run**

### Build from Command Line

```bash
cd ~/AndroidStudioProjects/HelloOffline
./gradlew build --offline
```

Success! You're building Android apps offline! 🎉

## Common Configurations

### For Beginners (Minimal)

Smallest download size, just basics:

```bash
# In config/download.config
API_LEVELS="34"
BUILD_TOOLS_VERSIONS="34.0.0"
GRADLE_VERSIONS="8.4"
DOWNLOAD_NDK="false"
DOWNLOAD_MAVEN_DEPS="true"  # Still recommended
```

**Size**: ~20 GB

### For Most Users (Recommended)

Balanced configuration:

```bash
API_LEVELS="33 34"
BUILD_TOOLS_VERSIONS="33.0.2 34.0.0"
GRADLE_VERSIONS="8.0 8.4"
DOWNLOAD_NDK="true"
DOWNLOAD_MAVEN_DEPS="true"
```

**Size**: ~40 GB

### For Professional/Team (Complete)

Everything you might need:

```bash
API_LEVELS="30 31 32 33 34"
BUILD_TOOLS_VERSIONS="30.0.3 31.0.0 32.0.0 33.0.2 34.0.0"
GRADLE_VERSIONS="7.4 7.6 8.0 8.2 8.4"
DOWNLOAD_NDK="true"
DOWNLOAD_MAVEN_DEPS="true"
```

**Size**: ~60-80 GB

## Troubleshooting Quick Fixes

### Download Issues

**Problem**: Download fails

```bash
# Resume download
RESUME_DOWNLOADS="true"  # In config/download.config

# Re-run failed script
./scripts/02-download-sdk-components.sh
```

### Installation Issues

**Problem**: Permission denied

```bash
# Install to home directory instead
# When prompted: ~/android-studio-offline
```

### Build Issues

**Problem**: Gradle can't find dependencies

```bash
# Check offline mode
cat ~/.gradle/init.d/offline-repos.gradle

# Should exist with repository paths

# Force offline mode
./gradlew build --offline
```

### Emulator Issues

**Problem**: Emulator won't start

```bash
# Linux: Enable KVM
sudo usermod -aG kvm $USER
# Log out and back in

# Or use software rendering
/opt/android-offline/android-sdk/emulator/emulator \
  -avd Pixel_5_API_34 -gpu swiftshader_indirect
```

## Next Steps

After getting started:

1. **Read Full Documentation**:
   - [Download Process](docs/01-DOWNLOAD-PROCESS.md) - Detailed download guide
   - [Components Guide](docs/02-COMPONENTS.md) - What's included
   - [Installation Guide](docs/03-INSTALLATION-GUIDE.md) - Complete installation
   - [Troubleshooting](docs/04-TROUBLESHOOTING.md) - Solve common issues
   - [Customization](docs/05-CUSTOMIZATION.md) - Advanced configuration

2. **Optimize Your Setup**:
   - Customize dependency list for your projects
   - Add project-specific libraries
   - Configure team standards

3. **Create Projects**:
   - Start building Android apps
   - Test on emulator
   - Deploy to devices (if available)

## Tips for Success

### Before Downloading

✅ **Check disk space**: Ensure you have enough
✅ **Stable connection**: Use wired connection if possible
✅ **Review configuration**: Match to your needs
✅ **Plan download time**: Can take several hours

### During Download

✅ **Monitor progress**: Watch for errors
✅ **Don't interrupt**: Let scripts complete
✅ **Check logs**: Review output for issues
✅ **Save bandwidth**: Download during off-peak if metered

### After Installation

✅ **Disable updates**: Prevent Android Studio from trying to update
✅ **Test build**: Create simple project to verify
✅ **Document setup**: Note any custom configurations
✅ **Backup**: Keep package for future installations

## Time Estimates

| Phase | Time | Notes |
|-------|------|-------|
| Configuration | 10-15 min | One-time setup |
| Download | 1-3 hours | Depends on internet speed |
| Package Creation | 20-60 min | Depends on size |
| Transfer | 10-60 min | Depends on method |
| Installation | 10-20 min | Mostly automated |
| **Total** | **2-5 hours** | First time; faster for updates |

## Size Estimates

| Component | Size | Required |
|-----------|------|----------|
| Android Studio | 2-3 GB | Yes |
| SDK | 10-30 GB | Yes |
| Gradle | 1-3 GB | Yes |
| Dependencies | 5-20 GB | Recommended |
| **Total Download** | **20-60 GB** | - |
| **Compressed Package** | **10-30 GB** | - |
| **Installed Size** | **25-70 GB** | - |

## Support

### Getting Help

1. **Check Documentation**: Most issues are documented
2. **Review Logs**: Error messages often explain the issue
3. **Search Issues**: Someone may have had the same problem
4. **Ask Questions**: Open a GitHub issue if stuck

### Providing Feedback

- **Found a bug?** Open an issue
- **Have a suggestion?** Open a feature request
- **Want to contribute?** See [CONTRIBUTING.md](CONTRIBUTING.md)

## Summary

**That's it!** You now have a complete offline Android development environment.

The basic workflow is:
1. **Download** components (internet machine)
2. **Package** everything (internet machine)
3. **Transfer** package (USB/network)
4. **Install** on offline machine
5. **Develop** Android apps!

For more details, see the full documentation in the `docs/` directory.

Happy coding! 🚀
