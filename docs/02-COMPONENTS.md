# Components Guide

This document provides detailed information about all components included in the Android Studio offline installer.

## Core Components

### 1. Android Studio IDE

**Description**: The official IDE for Android development, based on IntelliJ IDEA.

**Version**: Configurable (default: latest stable)

**Size**: 2-3 GB

**What's Included**:
- Complete IDE with code editor
- Visual layout editor
- Built-in emulator
- Profiling tools
- APK analyzer
- Lint tools
- Bundled JDK

**Platform Variants**:
- Linux: `.tar.gz`
- macOS (Intel): `.dmg`
- macOS (Apple Silicon): `.dmg`
- Windows: `.exe` or `.zip`

**Location After Install**: `/opt/android-offline/android-studio/`

### 2. Java Development Kit (JDK)

**Description**: Required for running Android Studio and building Android apps.

**Version**: JDK 17 (recommended for Android Studio)

**Size**: 300-400 MB

**Source**: Eclipse Temurin (AdoptOpenJDK)

**What's Included**:
- Java Runtime Environment (JRE)
- Java compiler
- Development tools

**Note**: Android Studio includes a bundled JDK, so this is optional.

**Location After Install**: `/opt/android-offline/jdk/`

## Android SDK Components

### 3. SDK Command Line Tools

**Description**: Command-line tools for managing the Android SDK.

**Size**: ~100 MB

**What's Included**:
- `sdkmanager` - SDK package manager
- `avdmanager` - AVD (emulator) manager
- `apkanalyzer` - APK analysis tool
- Package management utilities

**Location**: `/opt/android-offline/android-sdk/cmdline-tools/`

### 4. Platform Tools

**Description**: Essential tools for Android development and debugging.

**Size**: 10-20 MB

**What's Included**:
- `adb` - Android Debug Bridge
- `fastboot` - Flashing tool
- `systrace` - Performance analysis
- `sqlite3` - Database tool
- `dmtracedump`, `etc1tool`, `hprof-conv`

**Location**: `/opt/android-offline/android-sdk/platform-tools/`

### 5. SDK Build Tools

**Description**: Tools for building Android apps.

**Versions**: Multiple versions for compatibility (e.g., 33.0.2, 34.0.0)

**Size per version**: 50-70 MB

**What's Included**:
- `aapt2` - Android Asset Packaging Tool
- `apksigner` - APK signing tool
- `zipalign` - APK optimization
- `d8`, `dx` - Dex compilers
- `aidl` - Android Interface Definition Language compiler
- Split APK tools

**Location**: `/opt/android-offline/android-sdk/build-tools/`

### 6. SDK Platforms

**Description**: Android framework libraries for specific API levels.

**Versions**: Configurable (e.g., Android 13/API 33, Android 14/API 34)

**Size per platform**: 50-80 MB

**What's Included**:
- `android.jar` - Framework API
- Platform-specific tools
- Data files
- Optional sources

**Location**: `/opt/android-offline/android-sdk/platforms/`

### 7. System Images

**Description**: OS images for the Android Emulator.

**Size per image**: 400 MB - 2 GB

**Types**:
- **Google APIs**: Includes Google services
- **Google Play**: Includes Play Store
- **Default**: Vanilla AOSP
- **Android TV**, **Wear OS**, **Automotive**

**Architectures**:
- `x86_64` - For Intel/AMD processors (fastest on most dev machines)
- `arm64-v8a` - ARM 64-bit
- `armeabi-v7a` - ARM 32-bit (legacy)

**Example**:
```
system-images;android-34;google_apis;x86_64
```

**Location**: `/opt/android-offline/android-sdk/system-images/`

### 8. NDK (Native Development Kit)

**Description**: Toolchain for native C/C++ development.

**Version**: 25.2.9519653 (configurable)

**Size**: 1-2 GB

**What's Included**:
- GCC and Clang compilers
- STL implementations
- Cross-compilation toolchains
- CMake integration
- Native libraries

**Use Cases**:
- High-performance code
- Game engines
- Existing C/C++ libraries
- Hardware-specific operations

**Location**: `/opt/android-offline/android-sdk/ndk/`

### 9. CMake

**Description**: Build system for NDK projects.

**Version**: 3.22.1

**Size**: 50-100 MB

**Location**: `/opt/android-offline/android-sdk/cmake/`

### 10. SDK Extras

**Description**: Additional Google components.

**What's Included**:
- Google Play Services
- Google Repository (Maven)
- Android Support Repository
- Google USB Driver (Windows only)
- Intel HAXM (optional, for emulator acceleration)

**Size**: 500 MB - 2 GB combined

## Build System Components

### 11. Gradle Distributions

**Description**: Build automation tool for Android projects.

**Versions**: Multiple versions (7.x, 8.x series)

**Size per version**: 100-150 MB

**Distribution Types**:
- **bin**: Runtime only (~100 MB)
- **all**: Runtime + sources + docs (~150 MB)

**Recommended downloads**:
```
gradle-7.6-all.zip
gradle-8.0-all.zip
gradle-8.4-all.zip
```

**Location**: `/opt/android-offline/gradle/distributions/`

### 12. Gradle Wrapper

**Description**: Project-specific Gradle launcher.

**What's Included**:
- `gradle-wrapper.jar`
- Configuration templates
- Offline setup instructions

**Location**: `/opt/android-offline/gradle/wrapper/`

## Dependency Components

### 13. Maven Central Repository (Partial)

**Description**: Java/Kotlin libraries from Maven Central.

**Size**: 5-10 GB (for common dependencies)

**What's Included**:
- Kotlin standard library
- Testing frameworks (JUnit, Mockito)
- Common utilities (Gson, OkHttp, Retrofit)
- Third-party libraries

**Location**: `/opt/android-offline/maven-repo/`

### 14. Google Maven Repository

**Description**: Google's Android libraries.

**Size**: 5-15 GB

**What's Included**:
- **AndroidX**: Jetpack libraries
  - Core, AppCompat, ConstraintLayout
  - Lifecycle, ViewModel, LiveData
  - Navigation, Room, WorkManager
  - Compose (if included)

- **Material Design**: Material components

- **Google Play Services**:
  - Maps, Location
  - Firebase
  - AdMob

- **Android Tools**:
  - Android Gradle Plugin (AGP)
  - Lint tools
  - Data binding

**Location**: `/opt/android-offline/google-repo/`

### 15. Gradle Plugins

**Description**: Gradle build plugins.

**What's Included**:
- Android Gradle Plugin (AGP)
- Kotlin Gradle Plugin
- Dependency management plugins

**Versions**: Multiple versions for compatibility

**Size**: 1-2 GB

**Location**: `/opt/android-offline/gradle-plugins/`

## Common AndroidX Libraries

### Core Libraries
- `androidx.core:core-ktx` - Kotlin extensions
- `androidx.appcompat:appcompat` - Backward compatibility
- `androidx.activity:activity-ktx` - Activity APIs
- `androidx.fragment:fragment-ktx` - Fragment APIs

### UI Components
- `androidx.constraintlayout:constraintlayout` - Layout manager
- `com.google.android.material:material` - Material Design
- `androidx.recyclerview:recyclerview` - List views
- `androidx.cardview:cardview` - Card UI

### Architecture Components
- `androidx.lifecycle:lifecycle-*` - Lifecycle management
- `androidx.navigation:navigation-*` - Navigation framework
- `androidx.room:room-*` - Database ORM
- `androidx.paging:paging-*` - Pagination

### Jetpack Compose (Optional)
- `androidx.compose.ui:ui`
- `androidx.compose.material3:material3`
- `androidx.compose.foundation:foundation`

### Networking
- `com.squareup.retrofit2:retrofit` - HTTP client
- `com.squareup.okhttp3:okhttp` - Network library
- `com.google.code.gson:gson` - JSON parser

### Image Loading
- `com.github.bumptech.glide:glide` - Image loading
- `io.coil-kt:coil` - Kotlin image loader

### Dependency Injection
- `com.google.dagger:hilt-android` - DI framework
- `javax.inject:javax.inject` - Injection annotations

### Testing
- `junit:junit` - Unit testing
- `androidx.test.ext:junit` - Android testing
- `androidx.test.espresso:espresso-core` - UI testing
- `org.mockito:mockito-core` - Mocking framework

## Component Selection Guide

### Minimal Setup (Beginner)
```
✓ Android Studio
✓ Platform Tools
✓ Latest Build Tools
✓ Latest SDK Platform
✓ Latest Gradle
✓ Basic AndroidX libraries
Size: ~15-20 GB
```

### Standard Setup (Most Users)
```
✓ Android Studio
✓ Platform Tools
✓ Build Tools (2-3 versions)
✓ SDK Platforms (2-3 versions)
✓ System Images (1-2)
✓ Gradle (3-4 versions)
✓ Common dependencies
Size: ~30-40 GB
```

### Complete Setup (Professional)
```
✓ All above
✓ NDK
✓ CMake
✓ Multiple system images
✓ Many SDK versions
✓ Full dependency cache
✓ All Gradle versions
Size: 60-100 GB
```

### Specialized Setups

**Game Development**:
- NDK (required)
- CMake (required)
- ARM system images
- Native libraries

**IoT/Embedded**:
- Older SDK platforms
- ARM toolchains
- Minimal dependencies

**Enterprise**:
- Multiple API levels for testing
- All build tool versions
- Complete dependency mirror
- All Gradle versions

## Version Compatibility

### Android Studio vs SDK
- Android Studio 2023.x → SDK 33-34
- Android Studio 2022.x → SDK 31-33
- Always use matching or newer SDK

### Gradle vs Android Gradle Plugin
| Gradle Version | AGP Version | Android Studio |
|----------------|-------------|----------------|
| 7.5+ | 7.4+ | 2022.2+ |
| 8.0+ | 8.0+ | 2023.1+ |
| 8.2+ | 8.2+ | 2023.2+ |

### SDK Platform vs Min/Target SDK
- `compileSdk` should be latest available
- `targetSdk` should match `compileSdk`
- `minSdk` based on app requirements

## Updating Components

Components can be updated individually:

1. **Android Studio**: Download new version
2. **SDK Components**: Use SDK Manager
3. **Gradle**: Download new distributions
4. **Dependencies**: Re-run dependency scripts

For offline updates, repeat the download process on an internet-connected machine.

## Component Dependencies

Some components require others:

- **Emulator** requires: System Images, Platform Tools
- **NDK development** requires: CMake, NDK, Build Tools
- **Compose** requires: Kotlin 1.5+, specific AndroidX versions
- **Gradle builds** require: Matching AGP version, SDK platforms

## Storage Optimization

To reduce storage:

1. **Skip sources**: Don't download source JARs
2. **Limit architectures**: x86_64 only for emulator
3. **Fewer versions**: Latest only
4. **No NDK**: If not doing native dev
5. **Selective dependencies**: Only what you use

## Advanced: Custom Components

You can add custom components by:

1. Downloading to appropriate directory
2. Updating configuration files
3. Repackaging

Example:
```bash
# Custom Maven dependency
downloads/dependencies/maven-repo/com/example/library/...
```
