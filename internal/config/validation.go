package config

import (
	"fmt"
	"os"
)

// Validate validates the configuration
func (c *Config) Validate() error {
	// Validate platform
	validPlatforms := map[string]bool{
		"linux":      true,
		"darwin":     true,
		"darwin-arm": true,
		"windows":    true,
	}

	if !validPlatforms[c.Platform] {
		return fmt.Errorf("invalid platform: %s (must be: linux, darwin, darwin-arm, windows)", c.Platform)
	}

	// Validate Android Studio version
	if c.AndroidStudio.Version == "" {
		return fmt.Errorf("android_studio.version is required")
	}

	// Validate JDK version if downloading
	if c.AndroidStudio.DownloadJDK {
		if c.AndroidStudio.JDKVersion == "" {
			return fmt.Errorf("jdk_version is required when download_jdk is true")
		}
	}

	// Validate SDK
	if c.SDK.ToolsVersion == "" {
		return fmt.Errorf("sdk.tools_version is required")
	}

	if len(c.SDK.APILevels) == 0 {
		return fmt.Errorf("at least one API level is required in sdk.api_levels")
	}

	// Validate API levels are reasonable
	for _, api := range c.SDK.APILevels {
		if api < 21 || api > 40 {
			return fmt.Errorf("invalid API level: %d (expected 21-40)", api)
		}
	}

	if len(c.SDK.BuildTools) == 0 {
		return fmt.Errorf("at least one build tools version is required in sdk.build_tools")
	}

	// Validate Gradle
	if len(c.Gradle.Versions) == 0 {
		return fmt.Errorf("at least one Gradle version is required")
	}

	// Validate directories
	if c.DownloadDir == "" {
		c.DownloadDir = "downloads"
	}

	if c.InstallDir == "" {
		c.InstallDir = "/opt/android-offline"
	}

	// Validate download options
	if c.Options.ParallelDownloads < 1 || c.Options.ParallelDownloads > 16 {
		return fmt.Errorf("parallel_downloads must be between 1 and 16")
	}

	return nil
}

// ValidatePaths validates that required paths are accessible
func (c *Config) ValidatePaths() error {
	// Check download directory is writable
	if err := os.MkdirAll(c.DownloadDir, 0755); err != nil {
		return fmt.Errorf("cannot create download directory %s: %w", c.DownloadDir, err)
	}

	// Test write access
	testFile := fmt.Sprintf("%s/.test-write", c.DownloadDir)
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		return fmt.Errorf("download directory %s is not writable: %w", c.DownloadDir, err)
	}
	os.Remove(testFile)

	return nil
}

// ValidateTools checks that required system tools are available
func ValidateTools() error {
	// No longer need unzip - using pure Go
	// Could check for other tools if needed in the future
	return nil
}
