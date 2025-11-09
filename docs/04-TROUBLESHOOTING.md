# Troubleshooting Guide

Common issues and solutions for Android Studio offline installation and usage.

## Installation Issues

### Extraction Errors

**Problem**: `tar: Unexpected EOF` or `tar: Error is not recoverable`

**Causes**:
- Corrupted download
- Incomplete file transfer
- Disk space during extraction

**Solutions**:
```bash
# Verify file integrity
sha256sum android-studio-offline-*.tar.gz
# Compare with SHA256SUMS.txt from download machine

# Check disk space
df -h

# Try extraction with verbose output
tar -xzvf android-studio-offline-*.tar.gz

# If split package, verify all parts are present
ls android-studio-offline-*.tar.gz.part*
```

### Permission Issues

**Problem**: `Permission denied` during installation

**Solution 1**: Change installation directory
```bash
# Install to home directory instead
./install-offline.sh
# When prompted, enter: ~/android-studio-offline
```

**Solution 2**: Fix permissions
```bash
# Create directory with proper permissions
sudo mkdir -p /opt/android-offline
sudo chown -R $USER:$(id -gn) /opt/android-offline

# Run installer again
./install-offline.sh
```

### Insufficient Disk Space

**Problem**: Installation fails due to disk space

**Solutions**:
```bash
# Check available space
df -h

# Clean up space
sudo apt-get clean        # Linux
brew cleanup              # macOS

# Install to different location with more space
./install-offline.sh
# Choose location with adequate space

# Or mount external drive and install there
```

## Android Studio Launch Issues

### Won't Start

**Problem**: Android Studio doesn't launch

**Solution 1**: Check log files
```bash
# View startup logs
tail -f ~/.android/studio*/system/log/idea.log

# Look for errors
```

**Solution 2**: Increase memory
```bash
# Edit studio.vmoptions
nano /opt/android-offline/android-studio/bin/studio64.vmoptions

# Adjust heap size
-Xms256m
-Xmx4096m
```

**Solution 3**: Reset configuration
```bash
# Backup and remove config
mv ~/.android ~/.android.backup
mv ~/.AndroidStudio* ~/.AndroidStudio.backup

# Restart Android Studio
```

### UI Issues on Linux

**Problem**: Blank window or rendering issues

**Solutions**:
```bash
# Try different rendering mode
export _JAVA_OPTIONS='-Dsun.java2d.opengl=true'
/opt/android-offline/android-studio/bin/studio.sh

# Or disable hardware acceleration
export STUDIO_JDK=/opt/android-offline/jdk/...
/opt/android-offline/android-studio/bin/studio.sh
```

### macOS "Cannot Open App"

**Problem**: macOS blocks Android Studio

**Solutions**:
```bash
# Remove quarantine attribute
xattr -dr com.apple.quarantine /opt/android-offline/android-studio

# Allow from System Preferences
# System Preferences → Security & Privacy → Allow
```

## Gradle Issues

### Gradle Sync Fails

**Problem**: "Could not resolve all dependencies"

**Diagnosis**:
```bash
# Check if offline mode is enabled
cat ~/.gradle/init.d/offline-repos.gradle

# Verify repository paths
ls /opt/android-offline/maven-repo/
ls /opt/android-offline/google-repo/
```

**Solutions**:

**Solution 1**: Verify offline repositories configuration
```bash
# Check init script exists
cat ~/.gradle/init.d/offline-repos.gradle

# Should contain repository redirects
# If missing, copy from installation:
cp /opt/android-offline/gradle/init.d/offline-repos.gradle \
   ~/.gradle/init.d/
```

**Solution 2**: Update repository paths
```bash
# Edit init script
nano ~/.gradle/init.d/offline-repos.gradle

# Update paths to match your installation
# Change /opt/android-offline to your actual path
```

**Solution 3**: Check project gradle files
```groovy
// In build.gradle (project level)
// Ensure repositories are defined but will be overridden by init script

buildscript {
    repositories {
        google()
        mavenCentral()
    }
}

allprojects {
    repositories {
        google()
        mavenCentral()
    }
}
```

### Specific Dependency Not Found

**Problem**: "Could not find com.example:library:1.0.0"

**Causes**:
- Dependency not included in offline package
- Version mismatch

**Solutions**:

**Solution 1**: Check if dependency exists
```bash
# Search in local repos
find /opt/android-offline -name "*library*"

# Check specific dependency
ls /opt/android-offline/maven-repo/com/example/library/
```

**Solution 2**: Add missing dependency manually
```bash
# On internet-connected machine, download dependency
cd /tmp
wget https://repo1.maven.org/maven2/com/example/library/1.0.0/library-1.0.0.jar
wget https://repo1.maven.org/maven2/com/example/library/1.0.0/library-1.0.0.pom

# Transfer to offline machine
# Copy to correct location
mkdir -p /opt/android-offline/maven-repo/com/example/library/1.0.0/
cp library-1.0.0.* /opt/android-offline/maven-repo/com/example/library/1.0.0/
```

**Solution 3**: Use alternative version
```kotlin
// In build.gradle, try a different version that is available
implementation("com.example:library:1.1.0")  // instead of 1.0.0
```

### Gradle Wrapper Issues

**Problem**: Gradle wrapper downloads Gradle distribution

**Solutions**:

**Solution 1**: Update gradle-wrapper.properties
```properties
# In gradle/wrapper/gradle-wrapper.properties
distributionBase=GRADLE_USER_HOME
distributionPath=wrapper/dists
distributionUrl=file\:///opt/android-offline/gradle/distributions/gradle-8.4-all.zip
zipStoreBase=GRADLE_USER_HOME
zipStorePath=wrapper/dists
```

**Solution 2**: Use local Gradle
```bash
# Instead of ./gradlew, use system Gradle
/opt/android-offline/gradle/gradle-8.4/bin/gradle build --offline
```

## SDK Issues

### SDK Path Not Set

**Problem**: "SDK location not found"

**Solutions**:
```bash
# Set in Android Studio
# File → Project Structure → SDK Location
# Set to: /opt/android-offline/android-sdk

# Or create/edit local.properties in project root
echo "sdk.dir=/opt/android-offline/android-sdk" > local.properties

# Or set environment variable
export ANDROID_HOME=/opt/android-offline/android-sdk
```

### Missing Build Tools

**Problem**: "Build-Tools 34.0.0 is missing"

**Check what's installed**:
```bash
ls /opt/android-offline/android-sdk/build-tools/
```

**Solutions**:

**Solution 1**: Use available version
```groovy
// In app/build.gradle
android {
    buildToolsVersion "33.0.2"  // Use installed version
}
```

**Solution 2**: Download missing version
```bash
# On internet machine, download SDK package
# Extract and transfer to offline machine
# Copy to /opt/android-offline/android-sdk/build-tools/34.0.0/
```

### Missing Platform

**Problem**: "Failed to find target with hash string 'android-34'"

**Check installed platforms**:
```bash
ls /opt/android-offline/android-sdk/platforms/
```

**Solutions**:

**Solution 1**: Use available platform
```groovy
// In app/build.gradle
android {
    compileSdk 33  // Use installed version
    targetSdk 33
}
```

**Solution 2**: Add missing platform
```bash
# Download on internet machine
# Transfer platform-34 directory to:
# /opt/android-offline/android-sdk/platforms/android-34/
```

## Emulator Issues

### Emulator Won't Start

**Problem**: Emulator fails to launch

**Check logs**:
```bash
# Run emulator from command line to see errors
/opt/android-offline/android-sdk/emulator/emulator -avd Pixel_5_API_34 -verbose
```

**Solutions**:

**Solution 1**: Enable KVM (Linux)
```bash
# Check KVM support
egrep -c '(vmx|svm)' /proc/cpuinfo  # Should be > 0

# Install KVM
sudo apt-get install qemu-kvm libvirt-daemon-system

# Add user to kvm group
sudo usermod -aG kvm $USER

# Log out and log back in
```

**Solution 2**: Use software rendering
```bash
# Start emulator with software graphics
/opt/android-offline/android-sdk/emulator/emulator \
  -avd Pixel_5_API_34 \
  -gpu swiftshader_indirect
```

**Solution 3**: Increase RAM allocation
```bash
# Edit AVD config
nano ~/.android/avd/Pixel_5_API_34.avd/config.ini

# Adjust memory
hw.ramSize=2048
vm.heapSize=512
```

### System Image Not Found

**Problem**: "No system images installed"

**Check installed images**:
```bash
ls /opt/android-offline/android-sdk/system-images/
```

**Solutions**:

**Solution 1**: Create AVD with available image
```bash
# List available images
/opt/android-offline/android-sdk/cmdline-tools/latest/bin/avdmanager list

# Create AVD with available image
/opt/android-offline/android-sdk/cmdline-tools/latest/bin/avdmanager create avd \
  -n Pixel_5_API_33 \
  -k "system-images;android-33;google_apis;x86_64" \
  -d pixel_5
```

**Solution 2**: Add missing system image
```bash
# Download on internet machine
# Transfer to: /opt/android-offline/android-sdk/system-images/
```

### Emulator Performance Issues

**Problem**: Emulator is very slow

**Solutions**:

**Solution 1**: Enable hardware acceleration
```bash
# Linux: Use KVM (see above)

# macOS: Hypervisor.framework (automatic)

# Check acceleration
/opt/android-offline/android-sdk/emulator/emulator -accel-check
```

**Solution 2**: Reduce graphics quality
```bash
# Start with specific graphics mode
/opt/android-offline/android-sdk/emulator/emulator \
  -avd Pixel_5_API_34 \
  -gpu host
```

**Solution 3**: Adjust emulator settings
- AVD Manager → Edit AVD
- Show Advanced Settings
- Reduce RAM: 2048 MB
- Reduce Internal Storage: 2048 MB
- Graphics: Hardware or Automatic

## Build Issues

### Build Fails with "Task failed"

**Problem**: Build fails during compilation

**Diagnosis**:
```bash
# Build with verbose output
./gradlew build --offline --stacktrace --debug

# Check specific issue in output
```

**Common Causes**:

**Cause 1**: Dependency version conflict
```bash
# View dependency tree
./gradlew dependencies --offline

# Look for version conflicts
```

**Solution**:
```kotlin
// Force specific version in build.gradle
configurations.all {
    resolutionStrategy {
        force("com.example:library:1.0.0")
    }
}
```

**Cause 2**: Java version mismatch
```bash
# Check Java version
java -version

# Set JAVA_HOME to bundled JDK
export JAVA_HOME=/opt/android-offline/android-studio/jbr
```

**Cause 3**: Corrupted build cache
```bash
# Clean and rebuild
./gradlew clean build --offline

# Or delete build cache
rm -rf ~/.gradle/caches/
rm -rf app/build/
```

### Kotlin Compilation Errors

**Problem**: "Could not initialize Kotlin compiler"

**Solutions**:

```kotlin
// In build.gradle, ensure Kotlin plugin version matches
plugins {
    id 'org.jetbrains.kotlin.android' version '1.9.21'
}

// Ensure kotlin-stdlib is available
dependencies {
    implementation "org.jetbrains.kotlin:kotlin-stdlib:1.9.21"
}
```

### AAPT2 Errors

**Problem**: "AAPT: error: resource not found"

**Solutions**:

**Solution 1**: Clean project
```bash
./gradlew clean --offline
# Build → Clean Project in Android Studio
```

**Solution 2**: Check resource files
- Verify all referenced resources exist
- Check for typos in resource names
- Ensure XML files are well-formed

**Solution 3**: Invalidate caches
- File → Invalidate Caches → Invalidate and Restart

## ADB Issues

### ADB Devices Not Showing

**Problem**: `adb devices` shows no devices

**Solutions**:

**Solution 1**: Restart ADB
```bash
adb kill-server
adb start-server
adb devices
```

**Solution 2**: Check USB connection (physical device)
```bash
# Linux: Check udev rules
sudo usermod -aG plugdev $USER

# Create udev rules
sudo nano /etc/udev/rules.d/51-android.rules
# Add: SUBSYSTEM=="usb", ATTR{idVendor}=="XXXX", MODE="0666"
# Replace XXXX with device vendor ID

sudo udevadm control --reload-rules
```

**Solution 3**: For emulator
```bash
# Verify emulator is running
/opt/android-offline/android-sdk/emulator/emulator -list-avds

# Check emulator processes
ps aux | grep emulator
```

### ADB Connection Errors

**Problem**: "device offline" or "unauthorized"

**Solutions**:

**Solution 1**: Revoke USB debugging authorizations
```bash
# On device: Settings → Developer Options
# → Revoke USB debugging authorizations

# Reconnect and accept prompt
```

**Solution 2**: Reset ADB keys
```bash
rm ~/.android/adbkey*
adb kill-server
adb start-server
```

## Environment Issues

### Environment Variables Not Set

**Problem**: Commands like `adb` not found

**Solutions**:

**Solution 1**: Reload shell config
```bash
source ~/.bashrc  # or ~/.zshrc
```

**Solution 2**: Manually set for current session
```bash
export ANDROID_HOME=/opt/android-offline/android-sdk
export PATH=$PATH:$ANDROID_HOME/platform-tools:$ANDROID_HOME/tools
```

**Solution 3**: Verify shell config file
```bash
# Check if variables are in correct file
cat ~/.bashrc | grep ANDROID_HOME

# If using zsh, should be in ~/.zshrc instead
```

## Network/Proxy Issues

### Gradle Tries to Access Internet

**Problem**: Gradle attempts to download despite offline mode

**Solutions**:

**Solution 1**: Force offline mode
```bash
./gradlew build --offline
```

**Solution 2**: Enable in Android Studio
- File → Settings → Build, Execution, Deployment → Gradle
- Check "Offline mode"

**Solution 3**: Gradle properties
```properties
# In ~/.gradle/gradle.properties
systemProp.offline=true
```

**Solution 4**: Block network (extreme)
```bash
# Temporary firewall rule (Linux)
sudo iptables -A OUTPUT -p tcp --dport 443 -j REJECT
sudo iptables -A OUTPUT -p tcp --dport 80 -j REJECT

# Remember to remove after testing!
```

## Performance Issues

### Slow Build Times

**Solutions**:

**Solution 1**: Optimize Gradle settings
```properties
# In ~/.gradle/gradle.properties
org.gradle.daemon=true
org.gradle.parallel=true
org.gradle.caching=true
org.gradle.configureondemand=true
org.gradle.jvmargs=-Xmx6g -XX:MaxMetaspaceSize=512m
```

**Solution 2**: Exclude unnecessary directories from indexing
- File → Settings → Project Settings → Directories
- Mark build folders as Excluded

**Solution 3**: Increase Android Studio memory
```bash
# Edit studio.vmoptions
nano /opt/android-offline/android-studio/bin/studio64.vmoptions

# Increase heap
-Xms1024m
-Xmx8192m
```

### Slow IDE

**Solutions**:

**Solution 1**: Disable unnecessary plugins
- File → Settings → Plugins
- Disable unused plugins

**Solution 2**: Reduce inspection scope
- File → Settings → Editor → Inspections
- Disable heavy inspections

**Solution 3**: Power save mode
- File → Power Save Mode (toggle on)

## Logs and Diagnostics

### Where to Find Logs

**Android Studio Logs**:
```bash
~/.android/studio*/system/log/idea.log
```

**Gradle Logs**:
```bash
# Build output
./gradlew build --offline --debug > build.log 2>&1

# Daemon logs
~/.gradle/daemon/*/daemon-*.out.log
```

**Emulator Logs**:
```bash
# Run with verbose output
/opt/android-offline/android-sdk/emulator/emulator -avd <name> -verbose

# Logcat from running emulator
adb logcat
```

**ADB Logs**:
```bash
# Kill server with logging
adb kill-server
ADB_TRACE=all adb start-server
```

### Collecting Diagnostic Information

```bash
# System info
uname -a
lsb_release -a  # Linux
sw_vers  # macOS

# Java version
java -version

# Environment
echo $ANDROID_HOME
echo $ANDROID_SDK_ROOT
echo $PATH

# Installed components
ls /opt/android-offline/

# Gradle info
./gradlew --version

# SDK info
/opt/android-offline/android-sdk/cmdline-tools/latest/bin/sdkmanager --list
```

## Getting More Help

1. **Check logs** first (see above)
2. **Search** error messages in documentation
3. **Try minimal** project to isolate issue
4. **Compare** with working setup if available
5. **Document** exact steps to reproduce

## Emergency Recovery

### Complete Reset

If all else fails:

```bash
# Backup projects
cp -r ~/AndroidStudioProjects ~/AndroidStudioProjects.backup

# Remove Android Studio config
rm -rf ~/.android
rm -rf ~/.AndroidStudio*
rm -rf ~/.gradle

# Remove installation (optional)
sudo rm -rf /opt/android-offline

# Re-run installation
cd downloads
./install-offline.sh

# Restore projects
cp -r ~/AndroidStudioProjects.backup/* ~/AndroidStudioProjects/
```

Remember: Always keep backups of your projects!
