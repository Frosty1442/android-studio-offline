# Customization Guide

Advanced configuration options for customizing your Android Studio offline installer.

## Custom Component Selection

### Minimal Installation

For development machines with limited storage or specific requirements.

**Edit `config/download.config`**:
```bash
# Minimal configuration
ANDROID_STUDIO_VERSION="2023.3.1.18"
PLATFORM="linux"

# Only latest versions
API_LEVELS="34"
BUILD_TOOLS_VERSIONS="34.0.0"
GRADLE_VERSIONS="8.4"

# Skip optional components
DOWNLOAD_JDK="false"  # Use bundled JDK
DOWNLOAD_NDK="false"  # Skip native development
DOWNLOAD_CMAKE="false"
DOWNLOAD_MAVEN_DEPS="false"  # Manual dependency management
SYSTEM_IMAGES=""  # No emulator

# Minimal SDK
```

**Edit `config/sdk-packages.txt`**:
```
platform-tools
build-tools;34.0.0
platforms;android-34
cmdline-tools;latest
```

**Result**: ~10-15 GB package

### Custom API Level Selection

For projects targeting specific Android versions:

```bash
# Support Android 11 through 14
API_LEVELS="30 31 32 33 34"

# Corresponding build tools
BUILD_TOOLS_VERSIONS="30.0.3 31.0.0 32.0.0 33.0.2 34.0.0"

# System images for testing
SYSTEM_IMAGES="system-images;android-30;google_apis;x86_64 system-images;android-34;google_apis;x86_64"
```

### Game Development Setup

For Unity, Unreal, or native game development:

```bash
# NDK is essential
DOWNLOAD_NDK="true"
NDK_VERSION="25.2.9519653"

# CMake for build scripts
DOWNLOAD_CMAKE="true"

# Multiple architectures for distribution
# Add to sdk-packages.txt:
ndk;25.2.9519653
ndk;21.4.7075529  # For compatibility
cmake;3.22.1

# Additional system images for testing
SYSTEM_IMAGES="system-images;android-33;google_apis;x86_64 system-images;android-33;google_apis;arm64-v8a"
```

### Flutter Development

For Flutter with Android SDK:

```bash
# Flutter needs specific SDK versions
API_LEVELS="33 34"
BUILD_TOOLS_VERSIONS="33.0.2 34.0.0"

# Essential for Flutter
DOWNLOAD_JDK="true"
JDK_VERSION="17"

# Add to common-dependencies.txt:
com.google.android.material:material:1.11.0
androidx.appcompat:appcompat:1.6.1
androidx.annotation:annotation:1.7.0
```

**Additional setup**:
```bash
# Download Flutter SDK separately
# Transfer to offline machine
# Configure Flutter to use offline SDK:
flutter config --android-sdk /opt/android-offline/android-sdk
```

## Custom Dependencies

### Adding Project-Specific Dependencies

#### Method 1: Edit Config File

Add to `config/common-dependencies.txt`:

```
# Your custom libraries
com.squareup.retrofit2:retrofit:2.9.0
com.squareup.okhttp3:okhttp:4.12.0
io.insert-koin:koin-android:3.5.0
com.jakewharton.timber:timber:5.0.1

# Specific versions for your project
androidx.compose.ui:ui:1.5.4
androidx.compose.material3:material3:1.1.2
```

#### Method 2: Use Sample Project

Best for complex dependency trees:

```bash
# Place your project in downloads/sample-project/
cp -r ~/my-android-project downloads/sample-project/

# Configure
SAMPLE_PROJECT_PATH="downloads/sample-project"

# Run dependency download
./scripts/04-download-dependencies.sh

# All project dependencies will be pre-downloaded
```

### Creating Custom Dependency Lists

For organization-wide dependencies:

```bash
# Create custom list
cat > config/company-dependencies.txt << EOF
# Company standard libraries
com.company:core-library:2.0.0
com.company:networking:1.5.0
com.company:analytics:3.1.0

# Third-party requirements
com.google.firebase:firebase-bom:32.7.0
com.google.firebase:firebase-analytics-ktx:21.5.0
EOF
```

**Modify download script**:
```bash
# In scripts/04-download-dependencies.sh
# Add after common dependencies:

if [ -f "$PROJECT_ROOT/config/company-dependencies.txt" ]; then
    while IFS= read -r dep || [ -n "$dep" ]; do
        [[ "$dep" =~ ^#.*$ ]] && continue
        [[ -z "$dep" ]] && continue
        echo "Processing company dependency: $dep"
        download_dependency "$dep" "$maven_central" "$MAVEN_DIR"
    done < "$PROJECT_ROOT/config/company-dependencies.txt"
fi
```

## Custom Gradle Configuration

### Gradle Version Management

Add specific Gradle versions for team projects:

**Edit `config/gradle-versions.txt`**:
```
# Team projects use these versions
7.4.2  # Legacy project
8.0    # Main projects
8.4    # New projects
```

### Custom Gradle Init Scripts

Create team-specific Gradle configurations:

**Create `config/team-gradle-init.gradle`**:
```groovy
// Team-specific Gradle configuration

allprojects {
    // Company artifact repository
    repositories {
        maven {
            url = uri("/opt/android-offline/company-repo")
        }
    }

    // Standard build configurations
    plugins.withId('com.android.application') {
        android {
            compileSdkVersion 34

            defaultConfig {
                minSdkVersion 24
                targetSdkVersion 34
            }

            compileOptions {
                sourceCompatibility JavaVersion.VERSION_17
                targetCompatibility JavaVersion.VERSION_17
            }
        }
    }
}
```

**Install during setup**:
```bash
# In install-offline.sh, add:
if [ -f "$SCRIPT_DIR/config/team-gradle-init.gradle" ]; then
    cp "$SCRIPT_DIR/config/team-gradle-init.gradle" \
       "$USER_HOME/.gradle/init.d/"
fi
```

## Custom Installation Paths

### Multiple Installations

Install different versions side-by-side:

```bash
# Install different versions
./install-offline.sh
# Choose: /opt/android-2023
# Choose: /opt/android-2024

# Switch between them
export ANDROID_HOME=/opt/android-2023/android-sdk
# or
export ANDROID_HOME=/opt/android-2024/android-sdk
```

### Shared Team Installation

For multi-user workstations:

```bash
# Install to shared location
sudo ./install-offline.sh
# Choose: /opt/shared/android-studio

# Set permissions
sudo chmod -R 755 /opt/shared/android-studio
sudo chown -R root:developers /opt/shared/android-studio

# Each user configures
echo "export ANDROID_HOME=/opt/shared/android-studio/android-sdk" >> ~/.bashrc
```

### Portable Installation

For USB drive or portable use:

```bash
# Install to portable drive
./install-offline.sh
# Choose: /media/usb/android-studio-portable

# Create portable launcher
cat > /media/usb/android-studio-portable/launch.sh << 'EOF'
#!/bin/bash
export ANDROID_HOME="$(dirname "$0")/android-sdk"
export ANDROID_SDK_ROOT="$ANDROID_HOME"
export PATH="$PATH:$ANDROID_HOME/platform-tools"
"$(dirname "$0")/android-studio/bin/studio.sh"
EOF

chmod +x /media/usb/android-studio-portable/launch.sh
```

## Custom SDK Packages

### Adding Specialized SDK Components

#### Google Play Licensing Library

```bash
# Add to config/sdk-packages.txt
extras;google;market_licensing

# Or download manually
# Transfer to: android-sdk/extras/google/market_licensing/
```

#### Android TV Components

```bash
# Add to config/sdk-packages.txt
platforms;android-tv-33
platforms;android-tv-34
system-images;android-tv-33;google-tv;x86_64

# System images for TV emulator
```

#### Wear OS

```bash
# Add to config/sdk-packages.txt
platforms;android-wear-33
system-images;android-wear-33;google_apis;x86_64
```

#### Automotive

```bash
# Add to config/sdk-packages.txt
system-images;android-33;google_apis_automotive;x86_64
```

### Custom NDK Configuration

Multiple NDK versions for compatibility:

```bash
# config/sdk-packages.txt
ndk;25.2.9519653  # Latest
ndk;23.1.7779620  # For older projects
ndk;21.4.7075529  # Legacy support

# Each project can specify in build.gradle:
android {
    ndkVersion "25.2.9519653"
}
```

## Custom Build Scripts

### Automated Project Setup

Create script to set up new projects with offline configuration:

**Create `tools/new-offline-project.sh`**:
```bash
#!/bin/bash

PROJECT_NAME="$1"
PACKAGE_NAME="$2"

if [ -z "$PROJECT_NAME" ]; then
    echo "Usage: $0 <project-name> <package-name>"
    exit 1
fi

# Create project directory
mkdir -p ~/AndroidStudioProjects/$PROJECT_NAME

cd ~/AndroidStudioProjects/$PROJECT_NAME

# Create gradle wrapper with offline distribution
mkdir -p gradle/wrapper

cat > gradle/wrapper/gradle-wrapper.properties << EOF
distributionBase=GRADLE_USER_HOME
distributionPath=wrapper/dists
distributionUrl=file\:///opt/android-offline/gradle/distributions/gradle-8.4-all.zip
zipStoreBase=GRADLE_USER_HOME
zipStorePath=wrapper/dists
EOF

cp /opt/android-offline/gradle/wrapper/gradle-wrapper.jar gradle/wrapper/

# Create gradle.properties for offline
cat > gradle.properties << EOF
org.gradle.offline=true
android.useAndroidX=true
android.enableJetifier=true
android.offlineRoot=/opt/android-offline
EOF

# Create local.properties
cat > local.properties << EOF
sdk.dir=/opt/android-offline/android-sdk
EOF

echo "Project $PROJECT_NAME created for offline development"
echo "Open in Android Studio: ~/AndroidStudioProjects/$PROJECT_NAME"
```

### Build All Projects Script

**Create `tools/build-all-offline.sh`**:
```bash
#!/bin/bash

PROJECTS_DIR=~/AndroidStudioProjects

for project in "$PROJECTS_DIR"/*/; do
    echo "Building $(basename "$project")..."
    cd "$project"

    if [ -f "gradlew" ]; then
        ./gradlew clean build --offline || echo "Build failed for $(basename "$project")"
    fi
done

echo "All builds complete"
```

## Custom Emulator Configurations

### Create Standard AVD Profiles

**Create `tools/create-standard-avds.sh`**:
```bash
#!/bin/bash

SDK_ROOT="/opt/android-offline/android-sdk"
AVDMANAGER="$SDK_ROOT/cmdline-tools/latest/bin/avdmanager"

# Pixel 5 for testing
$AVDMANAGER create avd \
    -n "Pixel_5_API_34" \
    -k "system-images;android-34;google_apis;x86_64" \
    -d "pixel_5" \
    --force

# Tablet for large screen testing
$AVDMANAGER create avd \
    -n "Pixel_Tablet_API_34" \
    -k "system-images;android-34;google_apis;x86_64" \
    -d "pixel_tablet" \
    --force

# Older device for compatibility
$AVDMANAGER create avd \
    -n "Pixel_3_API_30" \
    -k "system-images;android-30;google_apis;x86_64" \
    -d "pixel_3" \
    --force

echo "Standard AVDs created"
```

### Custom AVD Settings

Edit AVD config files directly:

```bash
# Edit AVD config
nano ~/.android/avd/Pixel_5_API_34.avd/config.ini

# Customize settings:
hw.ramSize=4096
hw.keyboard=yes
showDeviceFrame=yes
skin.name=pixel_5
hw.gpu.enabled=yes
hw.gpu.mode=host
```

## Custom Repository Mirrors

### Internal Company Repository

If you have internal Maven repository:

**Create `tools/setup-company-repo.sh`**:
```bash
#!/bin/bash

COMPANY_REPO_URL="http://internal.company.com/maven"
COMPANY_REPO_DIR="/opt/android-offline/company-repo"

# Create directory
mkdir -p "$COMPANY_REPO_DIR"

# Mirror company repository (on internet machine)
wget -r -np -nH --cut-dirs=2 \
    -P "$COMPANY_REPO_DIR" \
    "$COMPANY_REPO_URL"

# Create Gradle init script
cat > ~/.gradle/init.d/company-repos.gradle << 'EOF'
allprojects {
    repositories {
        maven {
            url = uri("/opt/android-offline/company-repo")
        }
    }
}
EOF
```

### Custom Plugin Repository

For internal Gradle plugins:

```groovy
// In ~/.gradle/init.d/offline-repos.gradle
// Add before other repositories:

allprojects {
    buildscript {
        repositories {
            maven {
                url = uri("/opt/android-offline/company-plugins")
            }
        }
    }
}
```

## Environment Customization

### Per-Project Environment

**Create `tools/project-env.sh`**:
```bash
#!/bin/bash

# Source this for project-specific environment

export PROJECT_JAVA_HOME="/opt/android-offline/jdk/jdk-17"
export PROJECT_ANDROID_SDK="/opt/android-offline/android-sdk"
export PROJECT_GRADLE_HOME="/opt/android-offline/gradle/gradle-8.4"

export JAVA_HOME="$PROJECT_JAVA_HOME"
export ANDROID_HOME="$PROJECT_ANDROID_SDK"
export PATH="$PROJECT_GRADLE_HOME/bin:$PATH"

echo "Environment configured for offline Android development"
```

**Usage**:
```bash
source tools/project-env.sh
cd ~/AndroidStudioProjects/MyApp
./gradlew build --offline
```

### Custom Android Studio Settings

**Export/Import IDE Settings**:

```bash
# Export settings from configured Studio
# File → Manage IDE Settings → Export Settings
# Transfer android-studio-settings.zip to other machines

# Import on other machines
# File → Manage IDE Settings → Import Settings
```

**Automate settings**:

```bash
# Copy settings to new installation
cp -r ~/.AndroidStudio*/config/* \
      /opt/android-offline/studio-settings/

# Distribute to team
# Each member copies to their ~/.AndroidStudio*/config/
```

## Docker Container (Advanced)

### Create Dockerized Offline Environment

**Create `Dockerfile`**:
```dockerfile
FROM ubuntu:22.04

# Install dependencies
RUN apt-get update && apt-get install -y \
    openjdk-17-jdk \
    unzip \
    git \
    && rm -rf /var/lib/apt/lists/*

# Copy offline installer
COPY downloads /opt/android-offline-installer

# Run installation
RUN cd /opt/android-offline-installer && \
    ./install-offline.sh

# Set environment
ENV ANDROID_HOME=/opt/android-offline/android-sdk
ENV ANDROID_SDK_ROOT=/opt/android-offline/android-sdk
ENV PATH=${PATH}:${ANDROID_HOME}/platform-tools

# Workdir for projects
WORKDIR /workspace

CMD ["/bin/bash"]
```

**Build and use**:
```bash
# Build image
docker build -t android-offline:latest .

# Use for builds
docker run -v $(pwd):/workspace android-offline:latest \
    ./gradlew build --offline
```

## CI/CD Integration

### GitLab CI Example

```yaml
# .gitlab-ci.yml
build:
  image: android-offline:latest
  script:
    - ./gradlew clean build --offline
  artifacts:
    paths:
      - app/build/outputs/
```

### Jenkins Pipeline

```groovy
pipeline {
    agent {
        docker {
            image 'android-offline:latest'
        }
    }
    stages {
        stage('Build') {
            steps {
                sh './gradlew clean build --offline'
            }
        }
    }
}
```

## Maintenance and Updates

### Update Script Template

**Create `tools/update-components.sh`**:
```bash
#!/bin/bash

# Run this on internet-connected machine
# Updates specific components

COMPONENT="$1"

case "$COMPONENT" in
    "studio")
        ./scripts/01-download-android-studio.sh
        ;;
    "sdk")
        ./scripts/02-download-sdk-components.sh
        ;;
    "gradle")
        ./scripts/03-download-gradle.sh
        ;;
    "deps")
        ./scripts/04-download-dependencies.sh
        ;;
    *)
        echo "Usage: $0 {studio|sdk|gradle|deps}"
        exit 1
        ;;
esac

# Create incremental update package
tar -czf "update-${COMPONENT}-$(date +%Y%m%d).tar.gz" downloads/${COMPONENT}/

echo "Update package created: update-${COMPONENT}-$(date +%Y%m%d).tar.gz"
echo "Transfer to offline machine and extract to /opt/android-offline/"
```

## Tips and Best Practices

1. **Version Control**: Keep your configuration files in Git
2. **Documentation**: Document customizations for team
3. **Testing**: Test custom configurations on isolated machine first
4. **Backups**: Keep backups of working configurations
5. **Modularity**: Create modular scripts for different scenarios
6. **Consistency**: Use same configuration across team
7. **Updates**: Plan regular update cycles for components

## Example: Complete Custom Setup

Full example for a specific organization:

```bash
# config/download.config
PLATFORM="linux"
API_LEVELS="32 33 34"
BUILD_TOOLS_VERSIONS="32.0.0 33.0.2 34.0.0"
GRADLE_VERSIONS="8.0 8.2 8.4"
DOWNLOAD_JDK="true"
DOWNLOAD_NDK="true"

# config/company-dependencies.txt
com.company:core:3.0.0
com.company:network:2.1.0
androidx.core:core-ktx:1.12.0
com.google.dagger:hilt-android:2.48.1

# tools/company-setup.sh
./scripts/01-download-android-studio.sh
./scripts/02-download-sdk-components.sh
./scripts/03-download-gradle.sh
./scripts/04-download-dependencies.sh
./tools/add-company-deps.sh
./scripts/05-create-package.sh
```

Customize to your needs!
