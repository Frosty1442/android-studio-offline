# Android Studio Offline Installer (Go Edition)

Cross-platform CLI tool for creating and deploying offline Android Studio installations.

## Features

✨ **Cross-Platform**: Single binary for Windows, Linux, and macOS
🚀 **Fast**: Parallel downloads with progress bars
📦 **Complete**: Android Studio, SDK, Gradle, and dependencies
🔧 **Configurable**: YAML configuration for easy customization
💾 **Resume Support**: Interrupted downloads can resume
🎯 **Simple**: Clean CLI interface with colored output

## Quick Start

### Installation

#### Download Pre-built Binary

```bash
# Download for your platform from releases
# Linux
wget https://github.com/Frosty1442/android-studio-offline/releases/download/v1.0.0/android-offline-linux-amd64.tar.gz
tar -xzf android-offline-linux-amd64.tar.gz
sudo mv android-offline /usr/local/bin/

# macOS
wget https://github.com/Frosty1442/android-studio-offline/releases/download/v1.0.0/android-offline-darwin-amd64.tar.gz
tar -xzf android-offline-darwin-amd64.tar.gz
sudo mv android-offline /usr/local/bin/

# Windows
# Download android-offline-windows-amd64.exe
# Add to PATH or run from current directory
```

#### Build from Source

```bash
# Clone repository
git clone https://github.com/Frosty1442/android-studio-offline.git
cd android-studio-offline

# Build
make build

# Install (Linux/macOS)
make install

# Or just run
./bin/android-offline --help
```

### Usage

#### 1. Initialize Configuration

```bash
android-offline init
```

This creates a `config.yaml` file. Edit it to customize:

```yaml
platform: linux  # linux, darwin, darwin-arm, windows
android_studio:
  version: "2024.2.1.10"
sdk:
  api_levels: [33, 34]
  build_tools: ["33.0.2", "34.0.0"]
gradle:
  versions: ["7.6", "8.0", "8.4"]
```

#### 2. Download Components

On an internet-connected machine:

```bash
android-offline download
```

This downloads:
- Android Studio IDE
- Android SDK (platforms, build tools, etc.)
- System images for emulator
- Gradle distributions
- Maven dependencies
- JDK (optional)

**Progress output:**
```
[1/5] Downloading Android Studio 2024.2.1.10
ℹ Downloading from: https://...
android-studio-2024.2.1.10... [████████████████░░░░] 67.32% 1.8GB/2.7GB
```

#### 3. Create Package

```bash
android-offline package
```

Creates `packages/android-studio-offline-YYYYMMDD.tar.gz`

#### 4. Transfer to Offline Machine

Copy the package via USB drive, external HDD, or internal network.

#### 5. Install on Offline Machine

```bash
# Extract package
tar -xzf android-studio-offline-*.tar.gz

# Run installer
android-offline install
```

## Commands

### `android-offline init`
Create default configuration file

```bash
android-offline init
# Creates config.yaml in current directory
```

### `android-offline download`
Download all components

```bash
# Use default config
android-offline download

# Use custom config
android-offline download --config my-config.yaml

# Verbose output
android-offline download -v
```

### `android-offline package`
Create installation package

```bash
android-offline package

# Output: packages/android-studio-offline-YYYYMMDD.tar.gz
```

### `android-offline install`
Install on target machine

```bash
android-offline install

# Custom install directory
android-offline install --install-dir /home/user/android
```

### `android-offline verify`
Verify downloaded components

```bash
android-offline verify
```

## Configuration

### Example `config.yaml`

```yaml
platform: linux
download_dir: downloads
install_dir: /opt/android-offline

android_studio:
  version: "2024.2.1.10"
  download_jdk: true
  jdk_version: "17"

sdk:
  tools_version: "11076708"
  api_levels: [33, 34]
  build_tools: ["33.0.2", "34.0.0"]
  system_images:
    - "system-images;android-33;google_apis;x86_64"
    - "system-images;android-34;google_apis;x86_64"
  download_ndk: true
  ndk_version: "25.2.9519653"

gradle:
  versions: ["7.6", "8.0", "8.4"]

dependencies:
  download_maven: true
  common_libs:
    - "androidx.core:core-ktx:1.12.0"
    - "com.google.android.material:material:1.11.0"
    # ... more libraries

options:
  parallel_downloads: 4
  verify_checksums: true
  resume_downloads: true
  verbose: false
```

### Platform Values

- `linux` - Linux x86_64
- `darwin` - macOS Intel
- `darwin-arm` - macOS Apple Silicon
- `windows` - Windows x86_64

## Building

### Prerequisites

- Go 1.21 or higher
- Make (optional, for using Makefile)

### Build Commands

```bash
# Build for current platform
make build
# Output: bin/android-offline

# Build for all platforms
make build-all
# Output: bin/android-offline-{platform}-{arch}

# Create release packages
make release
# Output: releases/android-offline-VERSION-{platform}-{arch}.tar.gz

# Install locally
make install
# Installs to /usr/local/bin

# Run tests
make test

# Clean build artifacts
make clean
```

### Manual Build

```bash
# Install dependencies
go mod download

# Build
go build -o bin/android-offline ./cmd/android-offline

# Cross-compile for Windows (from Linux/macOS)
GOOS=windows GOARCH=amd64 go build -o bin/android-offline.exe ./cmd/android-offline
```

## Comparison: Go vs Bash Scripts

| Feature | Go Version | Bash Scripts |
|---------|------------|--------------|
| Cross-platform | ✅ Single binary | ⚠️ WSL on Windows |
| Progress bars | ✅ Real-time | ❌ Basic output |
| Parallel downloads | ✅ Concurrent | ⚠️ Limited |
| Resume support | ✅ Built-in | ⚠️ Via wget -c |
| Error handling | ✅ Structured | ⚠️ Basic |
| User experience | ✅ Colored output | ⚠️ Plain text |
| Dependencies | ✅ None (static binary) | ⚠️ wget, curl, bash |
| Performance | ✅ Faster | ⚠️ Slower |

## Advanced Usage

### Custom Configuration per Project

```bash
# Create config for specific project
cat > my-project-config.yaml << EOF
platform: linux
sdk:
  api_levels: [28]  # Only what your project needs
gradle:
  versions: ["8.0"]
EOF

# Download with custom config
android-offline download --config my-project-config.yaml
```

### CI/CD Integration

```yaml
# GitHub Actions example
- name: Create offline installer
  run: |
    ./android-offline init
    ./android-offline download
    ./android-offline package

- name: Upload artifact
  uses: actions/upload-artifact@v3
  with:
    name: android-offline-installer
    path: packages/*.tar.gz
```

### Docker

```dockerfile
FROM golang:1.21 as builder
WORKDIR /app
COPY . .
RUN make build

FROM ubuntu:22.04
COPY --from=builder /app/bin/android-offline /usr/local/bin/
RUN android-offline download
```

## Troubleshooting

### Download Failures

```bash
# Enable verbose logging
android-offline download -v

# Check config file
android-offline verify
```

### Build Issues

```bash
# Update dependencies
go mod tidy

# Clear cache
go clean -cache
```

### Missing Dependencies

```bash
# Install Go dependencies
make deps

# Or manually
go mod download
```

## Development

```bash
# Format code
make fmt

# Run linter
make lint

# Run tests
make test

# Build and run
make run
```

## Contributing

Contributions welcome! See [CONTRIBUTING.md](CONTRIBUTING.md)

## License

MIT License - see [LICENSE](LICENSE)

## Support

- 📖 Documentation: [docs/](docs/)
- 🐛 Issues: [GitHub Issues](https://github.com/Frosty1442/android-studio-offline/issues)
- 💬 Discussions: [GitHub Discussions](https://github.com/Frosty1442/android-studio-offline/discussions)

## Credits

Built with:
- [cobra](https://github.com/spf13/cobra) - CLI framework
- [yaml.v3](https://github.com/go-yaml/yaml) - YAML parsing

---

**Note**: The original bash scripts are still available in the `scripts/` directory for reference or if you prefer shell scripts.
