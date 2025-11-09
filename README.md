# Android Studio Offline Installer

A comprehensive solution for creating a fully offline Android Studio development environment that can be deployed on machines without internet access.

## Overview

This repository provides scripts and processes to:
- Download Android Studio and all required components
- Package the Android SDK, build tools, and platform tools
- Collect Gradle distributions and dependencies
- Create a portable offline installer
- Set up a complete Android development environment on air-gapped machines

## What's Included

### Core Components
- **Android Studio IDE** - Latest stable version
- **JDK** - Required Java Development Kit
- **Android SDK** - Software Development Kit
  - SDK Manager
  - SDK Build Tools (multiple versions)
  - SDK Platform Tools (adb, fastboot, etc.)
  - SDK Platforms (Android API levels)
  - SDK Command-line Tools
- **System Images** - For Android emulators
- **Google Play Services**
- **Android Support Libraries**

### Build Tools
- **Gradle** - Multiple versions for compatibility
- **Gradle Wrapper** - For project-specific Gradle versions
- **Maven Dependencies** - Pre-downloaded libraries
- **Google Maven Repository**
- **JCenter/MavenCentral Artifacts**

## Quick Start

### Phase 1: Download Components (On Internet-Connected Machine)

1. Configure your requirements:
```bash
cp config/download.config.template config/download.config
# Edit config/download.config to specify versions and components
```

2. Run the download script:
```bash
./scripts/01-download-android-studio.sh
./scripts/02-download-sdk-components.sh
./scripts/03-download-gradle.sh
./scripts/04-download-dependencies.sh
```

3. Package everything:
```bash
./scripts/05-create-package.sh
```

### Phase 2: Transfer to Offline Machine

Transfer the created package (usually a tar.gz or directory) to your offline machine via:
- USB drive
- External hard drive
- Secure file transfer
- Physical media

### Phase 3: Install on Offline Machine

1. Extract the package
2. Run the installation script:
```bash
./scripts/install-offline.sh
```

3. Configure your environment:
```bash
source ~/.bashrc  # or ~/.zshrc
```

## Detailed Documentation

- [Download Process](docs/01-DOWNLOAD-PROCESS.md) - Detailed guide for downloading components
- [Component List](docs/02-COMPONENTS.md) - Complete list of what's included
- [Installation Guide](docs/03-INSTALLATION-GUIDE.md) - Step-by-step installation instructions
- [Troubleshooting](docs/04-TROUBLESHOOTING.md) - Common issues and solutions
- [Customization](docs/05-CUSTOMIZATION.md) - How to customize for your needs

## Directory Structure

```
android-studio-offline/
├── README.md
├── scripts/
│   ├── 01-download-android-studio.sh
│   ├── 02-download-sdk-components.sh
│   ├── 03-download-gradle.sh
│   ├── 04-download-dependencies.sh
│   ├── 05-create-package.sh
│   └── install-offline.sh
├── config/
│   ├── download.config.template
│   ├── sdk-packages.txt
│   └── gradle-versions.txt
├── docs/
│   ├── 01-DOWNLOAD-PROCESS.md
│   ├── 02-COMPONENTS.md
│   ├── 03-INSTALLATION-GUIDE.md
│   ├── 04-TROUBLESHOOTING.md
│   └── 05-CUSTOMIZATION.md
├── downloads/
│   ├── android-studio/
│   ├── sdk/
│   ├── gradle/
│   └── dependencies/
└── tools/
    └── helper scripts
```

## System Requirements

### Download Machine (Internet-Connected)
- Linux, macOS, or Windows with WSL
- 100+ GB free disk space
- Stable internet connection
- wget or curl
- unzip, tar

### Target Machine (Offline)
- Linux, macOS, or Windows
- 50+ GB free disk space
- 8+ GB RAM (16 GB recommended)
- 64-bit processor

## License

This repository provides scripts and documentation. Individual components (Android Studio, SDK, etc.) are subject to their own licenses:
- Android Studio: [Android Software Development Kit License](https://developer.android.com/studio/terms)
- Gradle: Apache License 2.0

## Contributing

Contributions are welcome! Please see CONTRIBUTING.md for guidelines.

## Support

For issues and questions:
1. Check the [Troubleshooting Guide](docs/04-TROUBLESHOOTING.md)
2. Review existing GitHub issues
3. Create a new issue with detailed information

## Version

Current Version: 1.0.0
Last Updated: 2025-11-09
