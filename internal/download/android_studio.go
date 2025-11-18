package download

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Frosty1442/android-studio-offline/internal/config"
)

const (
	studioBaseURL = "https://redirector.gvt1.com/edgedl/android/studio/ide-zips"
	jdkBaseURL    = "https://api.adoptium.net/v3/binary/latest"
)

// DownloadAndroidStudio downloads Android Studio for the specified platform
func (d *Downloader) DownloadAndroidStudio(ctx context.Context, cfg *config.Config) error {
	d.logger.Step(1, 5, "Downloading Android Studio %s", cfg.AndroidStudio.Version)

	// Determine filename based on platform
	var filename string
	switch cfg.Platform {
	case "linux":
		filename = fmt.Sprintf("android-studio-%s-linux.tar.gz", cfg.AndroidStudio.Version)
	case "darwin":
		filename = fmt.Sprintf("android-studio-%s-mac.dmg", cfg.AndroidStudio.Version)
	case "darwin-arm":
		filename = fmt.Sprintf("android-studio-%s-mac_arm.dmg", cfg.AndroidStudio.Version)
	case "windows":
		filename = fmt.Sprintf("android-studio-%s-windows.zip", cfg.AndroidStudio.Version)
	default:
		return fmt.Errorf("unsupported platform: %s", cfg.Platform)
	}

	url := fmt.Sprintf("%s/%s/%s", studioBaseURL, cfg.AndroidStudio.Version, filename)
	destPath := filepath.Join(cfg.DownloadDir, "android-studio", filename)

	// Check if already downloaded
	if fileExists(destPath) {
		d.logger.Info("Android Studio already downloaded, skipping")
		return nil
	}

	d.logger.Info("Downloading from: %s", url)
	if err := d.DownloadAndVerify(ctx, url, destPath, cfg.Options.ResumeDownloads, cfg.Options.VerifyChecksums); err != nil {
		return fmt.Errorf("failed to download Android Studio: %w", err)
	}

	d.logger.Success("Android Studio downloaded successfully")
	return nil
}

// DownloadJDK downloads the Java Development Kit
func (d *Downloader) DownloadJDK(ctx context.Context, cfg *config.Config) error {
	if !cfg.AndroidStudio.DownloadJDK {
		d.logger.Info("JDK download disabled, skipping")
		return nil
	}

	d.logger.Step(2, 5, "Downloading JDK %s", cfg.AndroidStudio.JDKVersion)

	// Determine architecture
	var arch, ext string
	switch cfg.Platform {
	case "linux":
		arch = "linux-x64"
		ext = "tar.gz"
	case "darwin":
		arch = "mac-x64"
		ext = "tar.gz"
	case "darwin-arm":
		arch = "mac-aarch64"
		ext = "tar.gz"
	case "windows":
		arch = "windows-x64"
		ext = "zip"
	default:
		return fmt.Errorf("unsupported platform: %s", cfg.Platform)
	}

	filename := fmt.Sprintf("jdk-%s-%s.%s", cfg.AndroidStudio.JDKVersion, arch, ext)
	url := fmt.Sprintf("%s/%s/ga/%s/jdk/hotspot/normal/eclipse",
		jdkBaseURL, cfg.AndroidStudio.JDKVersion, arch)
	destPath := filepath.Join(cfg.DownloadDir, "jdk", filename)

	// Check if already downloaded
	if fileExists(destPath) {
		d.logger.Info("JDK already downloaded, skipping")
		return nil
	}

	d.logger.Info("Downloading from: %s", url)
	if err := d.DownloadAndVerify(ctx, url, destPath, cfg.Options.ResumeDownloads, cfg.Options.VerifyChecksums); err != nil {
		return fmt.Errorf("failed to download JDK: %w", err)
	}

	d.logger.Success("JDK downloaded successfully")
	return nil
}

// DownloadSDKCommandLineTools downloads the Android SDK command-line tools
func (d *Downloader) DownloadSDKCommandLineTools(ctx context.Context, cfg *config.Config) error {
	d.logger.Step(3, 5, "Downloading SDK Command Line Tools")

	// Determine platform suffix
	var platformSuffix string
	switch cfg.Platform {
	case "linux":
		platformSuffix = "linux"
	case "darwin", "darwin-arm":
		platformSuffix = "mac"
	case "windows":
		platformSuffix = "win"
	default:
		return fmt.Errorf("unsupported platform: %s", cfg.Platform)
	}

	filename := fmt.Sprintf("commandlinetools-%s-%s_latest.zip", platformSuffix, cfg.SDK.ToolsVersion)
	url := fmt.Sprintf("https://dl.google.com/android/repository/%s", filename)
	destPath := filepath.Join(cfg.DownloadDir, "sdk", filename)

	// Check if already downloaded
	if fileExists(destPath) {
		d.logger.Info("SDK command-line tools already downloaded, skipping")
		return nil
	}

	d.logger.Info("Downloading from: %s", url)
	if err := d.DownloadAndVerify(ctx, url, destPath, cfg.Options.ResumeDownloads, cfg.Options.VerifyChecksums); err != nil {
		return fmt.Errorf("failed to download SDK tools: %w", err)
	}

	d.logger.Success("SDK command-line tools downloaded successfully")
	return nil
}

// DownloadPlatformTools downloads standalone platform tools
func (d *Downloader) DownloadPlatformTools(ctx context.Context, cfg *config.Config) error {
	d.logger.Info("Downloading Platform Tools (adb, fastboot, etc.)")

	var platformSuffix string
	switch cfg.Platform {
	case "linux":
		platformSuffix = "linux"
	case "darwin", "darwin-arm":
		platformSuffix = "darwin"
	case "windows":
		platformSuffix = "windows"
	default:
		return fmt.Errorf("unsupported platform: %s", cfg.Platform)
	}

	filename := fmt.Sprintf("platform-tools-latest-%s.zip", platformSuffix)
	url := fmt.Sprintf("https://dl.google.com/android/repository/%s", filename)
	destPath := filepath.Join(cfg.DownloadDir, "sdk", filename)

	// Check if already downloaded
	if fileExists(destPath) {
		d.logger.Info("Platform tools already downloaded, skipping")
		return nil
	}

	if err := d.DownloadAndVerify(ctx, url, destPath, cfg.Options.ResumeDownloads, cfg.Options.VerifyChecksums); err != nil {
		return fmt.Errorf("failed to download platform tools: %w", err)
	}

	d.logger.Success("Platform tools downloaded successfully")
	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
