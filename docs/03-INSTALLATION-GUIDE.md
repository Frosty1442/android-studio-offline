# Installation Guide

Complete guide for installing Android Studio offline installer on an air-gapped machine.

## Pre-Installation

### Transfer Package to Offline Machine

Choose your transfer method:

#### Option 1: USB Drive
```bash
# On internet-connected machine
cp packages/android-studio-offline-*.tar.gz /media/usb/

# On offline machine
cp /media/usb/android-studio-offline-*.tar.gz ~/
```

#### Option 2: External Hard Drive
```bash
# Large packages (>32GB) may need external drive
# Format: exFAT for cross-platform or ext4 for Linux
```

#### Option 3: Network Transfer (Internal Network)
```bash
# Using scp (if internal network available)
scp android-studio-offline-*.tar.gz user@offline-machine:~/

# Or use rsync
rsync -avh --progress android-studio-offline-*.tar.gz user@offline-machine:~/
```

### For Split Packages

If package was split (>4GB):

```bash
# Transfer all parts
android-studio-offline-*.tar.gz.partaa
android-studio-offline-*.tar.gz.partab
android-studio-offline-*.tar.gz.partac
...

# On offline machine, join them:
cat android-studio-offline-*.tar.gz.part* > android-studio-offline.tar.gz

# Verify
tar -tzf android-studio-offline.tar.gz | head
```

## System Requirements

### Minimum Requirements
- **OS**: Linux (Ubuntu 18.04+), macOS (10.14+), Windows 10/11
- **Disk Space**: 50 GB free
- **RAM**: 8 GB
- **Processor**: 64-bit x86 CPU

### Recommended Requirements
- **Disk Space**: 100 GB free (SSD preferred)
- **RAM**: 16 GB
- **Processor**: Multi-core 64-bit CPU
- **Display**: 1920x1080 or higher

### Check System Requirements

```bash
# Check disk space
df -h ~

# Check RAM
free -h

# Check CPU
lscpu | grep -E "Architecture|CPU\(s\)"

# Check OS version
cat /etc/os-release  # Linux
sw_vers              # macOS
```

## Installation Steps

### Step 1: Extract Package

```bash
# Navigate to where you transferred the package
cd ~

# Extract the package
tar -xzf android-studio-offline-*.tar.gz

# This creates a 'downloads' directory
ls -l downloads/
```

### Step 2: Run Installation Script

```bash
# Navigate to downloads directory
cd downloads

# Make script executable (if not already)
chmod +x install-offline.sh

# Run installation
./install-offline.sh
```

### Step 3: Follow Installation Prompts

The script will ask you:

1. **Installation directory**
   ```
   Default: /opt/android-offline

   Options:
   - System-wide: /opt/android-offline (requires sudo)
   - User-local: ~/android-studio-offline
   - Custom: /path/of/your/choice
   ```

2. **Confirmation**
   ```
   Review the summary and confirm (y/n)
   ```

The installation will then proceed automatically.

### Step 4: Reload Environment

After installation completes:

```bash
# Reload shell configuration
source ~/.bashrc    # For Bash
source ~/.zshrc     # For Zsh

# Or start a new terminal session
```

### Step 5: Verify Installation

```bash
# Check environment variables
echo $ANDROID_HOME
echo $ANDROID_SDK_ROOT

# Check if Android Studio is accessible
ls -l /opt/android-offline/android-studio/bin/studio.sh

# Check SDK
ls -l /opt/android-offline/android-sdk/
```

## First Run

### Starting Android Studio

#### From Command Line
```bash
/opt/android-offline/android-studio/bin/studio.sh
```

#### From Desktop (Linux)
- Look for "Android Studio (Offline)" in your application menu
- Or double-click the desktop entry

#### Create Alias (Optional)
```bash
# Add to ~/.bashrc or ~/.zshrc
alias android-studio='/opt/android-offline/android-studio/bin/studio.sh'

# Now you can start with:
android-studio
```

### Initial Setup Wizard

1. **Import Settings**
   - Choose "Do not import settings" (first install)
   - Or import from previous installation

2. **Data Sharing**
   - Choose "Don't send" (recommended for offline)

3. **Install Type**
   - Select **"Custom"** installation

4. **UI Theme**
   - Choose your preferred theme (Darcula/Light)

5. **SDK Components**
   - **IMPORTANT**: Click "..." next to Android SDK Location
   - Set to: `/opt/android-offline/android-sdk`
   - Uncheck "Download SDK components" (already have them)

6. **Emulator Settings**
   - Keep default settings
   - Emulator is already installed

7. **Verify Settings**
   - Review and click "Finish"

### Disable Automatic Updates

**IMPORTANT**: Disable updates to prevent errors in offline mode.

1. Go to **Settings/Preferences** (File → Settings)

2. **Appearance & Behavior → System Settings → Updates**
   - Uncheck "Automatically check updates for:"
   - Set to "Never check"

3. **Appearance & Behavior → System Settings → Android SDK**
   - Verify SDK location: `/opt/android-offline/android-sdk`

4. **Build, Execution, Deployment → Gradle**
   - Gradle user home: `~/.gradle`
   - Enable "Offline mode" checkbox

5. Click **Apply** and **OK**

## Configuring for Offline Use

### Gradle Configuration

The installer creates `~/.gradle/gradle.properties` with offline settings.

Verify it contains:

```properties
org.gradle.daemon=true
org.gradle.parallel=true
org.gradle.caching=true
android.useAndroidX=true
android.offlineRoot=/opt/android-offline
```

### Offline Repository Configuration

Check `~/.gradle/init.d/offline-repos.gradle` exists:

```bash
cat ~/.gradle/init.d/offline-repos.gradle
```

This file redirects Maven repositories to local copies.

### Project-Level Configuration

For each Android project, add to `gradle.properties`:

```properties
# Enable offline mode
org.gradle.offline=true
```

Or use the Gradle wrapper with offline flag:

```bash
./gradlew build --offline
```

## Creating Your First Project

### 1. Start Android Studio

```bash
/opt/android-offline/android-studio/bin/studio.sh
```

### 2. New Project

1. Click **"New Project"**

2. Choose a template (e.g., "Empty Activity")

3. Configure project:
   - **Name**: MyOfflineApp
   - **Package**: com.example.myofflineapp
   - **Save location**: ~/AndroidStudioProjects/MyOfflineApp
   - **Language**: Kotlin (or Java)
   - **Minimum SDK**: API 24 (Android 7.0) or as needed
   - Click **Finish**

### 3. Wait for Gradle Sync

First sync may take a few minutes as Gradle:
- Sets up build cache
- Indexes dependencies
- Configures project

Watch the bottom status bar for "Gradle sync finished"

### 4. Test Build

```bash
# In Android Studio terminal or external terminal
cd ~/AndroidStudioProjects/MyOfflineApp
./gradlew build --offline
```

Should complete without errors!

## Setting Up Emulator

### 1. Open AVD Manager

- **Tools → Device Manager** (or **AVD Manager**)

### 2. Create Virtual Device

1. Click **"Create Device"**

2. Select hardware:
   - Choose **"Pixel 5"** or similar
   - Click **Next**

3. Select system image:
   - Choose from **downloaded images** (e.g., API 34)
   - If image is missing, you need to add it to the package
   - Click **Next**

4. Verify configuration:
   - Name: Pixel_5_API_34
   - Startup orientation: Portrait
   - Click **Finish**

### 3. Launch Emulator

- Click **▶ Play** button next to your AVD
- Wait for emulator to boot (first boot is slower)

### 4. Deploy App

- With emulator running, click **▶ Run** in Android Studio
- Select the running emulator
- App will install and launch

## Troubleshooting Installation

### Issue: Permission Denied

```bash
# If installation directory requires sudo
sudo mkdir -p /opt/android-offline
sudo chown -R $USER:$USER /opt/android-offline

# Then re-run installer
./install-offline.sh
```

### Issue: ANDROID_HOME Not Set

```bash
# Manually add to ~/.bashrc or ~/.zshrc
export ANDROID_HOME="/opt/android-offline/android-sdk"
export ANDROID_SDK_ROOT="/opt/android-offline/android-sdk"
export PATH="$PATH:$ANDROID_HOME/platform-tools:$ANDROID_HOME/tools"

# Reload
source ~/.bashrc
```

### Issue: Gradle Can't Find Dependencies

Check offline repositories:

```bash
# Verify init script
cat ~/.gradle/init.d/offline-repos.gradle

# Check repository paths are correct
ls /opt/android-offline/maven-repo/
ls /opt/android-offline/google-repo/

# Test with specific project
cd your-project
./gradlew dependencies --offline
```

### Issue: SDK Manager Shows Errors

This is normal in offline mode. The SDK Manager expects internet access.

To manage SDK:
```bash
# Use command line
cd /opt/android-offline/android-sdk
./cmdline-tools/latest/bin/sdkmanager --list

# SDK is already installed, so this is informational only
```

### Issue: Emulator Won't Start

```bash
# Check emulator
/opt/android-offline/android-sdk/emulator/emulator -list-avds

# Check system images
ls /opt/android-offline/android-sdk/system-images/

# Try from command line
/opt/android-offline/android-sdk/emulator/emulator -avd Pixel_5_API_34
```

## Platform-Specific Notes

### Linux

#### Desktop Entry
Located at: `~/.local/share/applications/android-studio-offline.desktop`

If not working:
```bash
chmod +x ~/.local/share/applications/android-studio-offline.desktop
```

#### KVM Acceleration (for Emulator)
```bash
# Check if KVM is available
egrep -c '(vmx|svm)' /proc/cpuinfo  # Should be > 0

# Check KVM modules
lsmod | grep kvm

# Add user to kvm group
sudo usermod -aG kvm $USER

# Log out and back in
```

### macOS

#### Allow Apps from Unidentified Developers
```bash
# If macOS blocks Android Studio
sudo spctl --master-disable

# Or individually
xattr -d com.apple.quarantine /opt/android-offline/android-studio
```

#### Emulator Acceleration
macOS uses Hypervisor.framework (built-in), no additional setup needed.

### Windows (via WSL)

#### Using WSL2
```bash
# Install WSL2 first
# Then follow Linux instructions

# Note: GUI apps require WSLg (Windows 11) or X server
```

#### Native Windows Installation
Extract Android Studio `.zip` or run `.exe` installer, then:
1. Point to extracted SDK
2. Configure Gradle offline repos manually
3. Edit environment variables via System Properties

## Verification Checklist

After installation, verify:

- [ ] Android Studio launches successfully
- [ ] SDK location is set to offline SDK
- [ ] Automatic updates are disabled
- [ ] Can create a new project
- [ ] Gradle sync completes without internet
- [ ] Can build project with `--offline` flag
- [ ] AVD (emulator) can be created
- [ ] Emulator starts successfully
- [ ] Can deploy app to emulator
- [ ] ADB works: `adb devices`

## Quick Reference

### Important Paths
```
Android Studio: /opt/android-offline/android-studio
Android SDK:    /opt/android-offline/android-sdk
Gradle:         /opt/android-offline/gradle
JDK:            /opt/android-offline/jdk
Maven Repo:     /opt/android-offline/maven-repo
Google Repo:    /opt/android-offline/google-repo
```

### Common Commands
```bash
# Start Android Studio
/opt/android-offline/android-studio/bin/studio.sh

# ADB
adb devices
adb logcat

# Gradle (in project directory)
./gradlew build --offline
./gradlew clean
./gradlew assembleDebug --offline

# SDK Manager
cd /opt/android-offline/android-sdk
./cmdline-tools/latest/bin/sdkmanager --list

# Emulator
/opt/android-offline/android-sdk/emulator/emulator -list-avds
/opt/android-offline/android-sdk/emulator/emulator -avd <name>
```

## Next Steps

1. Read [QUICKSTART.txt](file:///opt/android-offline/QUICKSTART.txt) for quick tips
2. Review [Troubleshooting Guide](04-TROUBLESHOOTING.md) for common issues
3. Check [Customization Guide](05-CUSTOMIZATION.md) for advanced configuration
4. Start building your Android apps offline!

## Getting Help

For issues not covered here:
1. Check logs: `~/.android/` and Android Studio logs
2. Review installation output
3. Consult offline installer repository documentation
4. Check Android Studio offline mode documentation
