# Contributing to Android Studio Offline Installer

Thank you for your interest in contributing! This document provides guidelines for contributing to the project.

## How to Contribute

### Reporting Issues

If you encounter problems:

1. **Search existing issues** first to avoid duplicates
2. **Provide detailed information**:
   - Operating system and version
   - Script/component that failed
   - Complete error message
   - Steps to reproduce
   - Relevant log files

3. **Use issue templates** when available

### Suggesting Enhancements

For new features or improvements:

1. **Check existing suggestions** to avoid duplicates
2. **Explain the use case** clearly
3. **Describe the expected behavior**
4. **Consider backwards compatibility**

### Submitting Pull Requests

1. **Fork the repository**
2. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Make your changes**:
   - Follow existing code style
   - Add comments for complex logic
   - Update documentation
   - Test your changes

4. **Commit your changes**:
   ```bash
   git commit -m "Add feature: description"
   ```
   - Use clear, descriptive commit messages
   - Reference issue numbers if applicable

5. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

6. **Open a Pull Request**:
   - Describe what changes you made
   - Explain why the changes are needed
   - Reference related issues
   - Include testing steps

## Development Guidelines

### Bash Scripts

- Use `#!/bin/bash` shebang
- Enable error checking: `set -e`
- Add error messages for failures
- Include usage/help information
- Add comments for complex sections
- Test on multiple platforms if possible

**Example**:
```bash
#!/bin/bash

set -e

# Description of what this does
function download_component() {
    local url="$1"
    local output="$2"

    # Check if already downloaded
    if [ -f "$output" ]; then
        echo "Already downloaded: $output"
        return 0
    fi

    # Download with error handling
    wget -O "$output" "$url" || {
        echo "Error: Failed to download from $url"
        return 1
    }
}
```

### Documentation

- Use clear, concise language
- Include code examples
- Add troubleshooting tips
- Keep formatting consistent
- Update README if adding features

### Configuration Files

- Use descriptive variable names
- Add comments explaining options
- Provide sensible defaults
- Include examples

### Testing

Before submitting:

1. **Test your changes**:
   - Run modified scripts
   - Check for errors
   - Verify output

2. **Test on target platforms**:
   - Linux (Ubuntu/Debian, Fedora/RHEL)
   - macOS (if applicable)
   - Windows/WSL (if applicable)

3. **Check for regressions**:
   - Ensure existing functionality works
   - Test edge cases

## Code Style

### Shell Scripts

- **Indentation**: 4 spaces
- **Variables**:
  - Uppercase for constants: `DOWNLOAD_DIR`
  - Lowercase for local: `local file_name`
- **Functions**: lowercase with underscores
- **Comments**: Before complex sections
- **Error handling**: Check exit codes

### Markdown

- Use ATX-style headers (`#`)
- Add blank line before/after headers
- Use fenced code blocks with language
- Keep lines under 80-100 characters when possible

## What to Contribute

### High Priority

- Bug fixes
- Documentation improvements
- Platform compatibility fixes
- Error handling improvements
- Performance optimizations

### Welcome Contributions

- New platform support
- Additional component downloads
- Improved configuration options
- Better error messages
- Usage examples

### Examples of Good Contributions

1. **Bug Fix**:
   ```
   Fix: Gradle download fails on macOS

   The wget command wasn't compatible with macOS.
   Changed to use curl as fallback.

   Fixes #123
   ```

2. **Feature Addition**:
   ```
   Add: Support for Kotlin Multiplatform dependencies

   Adds configuration option to download KMP libraries.
   Includes new section in common-dependencies.txt.
   Updates documentation with KMP setup guide.

   Related to #456
   ```

3. **Documentation**:
   ```
   Docs: Add troubleshooting for permission issues

   Many users face permission problems during installation.
   Added detailed troubleshooting section with solutions
   for common scenarios.
   ```

## Community Guidelines

- **Be respectful** and professional
- **Help others** when possible
- **Stay on topic** in discussions
- **Accept feedback** graciously
- **Give constructive feedback** kindly

## Questions?

- **General questions**: Open a Discussion
- **Bug reports**: Open an Issue
- **Feature requests**: Open an Issue
- **Pull requests**: Follow guidelines above

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

## Recognition

Contributors will be recognized in:
- README.md Contributors section
- Release notes for significant contributions

Thank you for contributing! 🎉
