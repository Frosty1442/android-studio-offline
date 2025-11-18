package download

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Frosty1442/android-studio-offline/internal/config"
)

// DownloadSDKPackages downloads SDK packages using sdkmanager
func (d *Downloader) DownloadSDKPackages(ctx context.Context, cfg *config.Config) error {
	d.logger.Step(5, 5, "Downloading SDK packages")

	// First, ensure we have SDK tools
	sdkRoot := filepath.Join(cfg.DownloadDir, "sdk", "android-sdk")
	sdkManager := filepath.Join(sdkRoot, "cmdline-tools", "latest", "bin", "sdkmanager")

	// Extract SDK command-line tools if not already done
	if err := d.extractSDKTools(cfg); err != nil {
		return fmt.Errorf("failed to setup SDK tools: %w", err)
	}

	// Make sdkmanager executable
	if err := os.Chmod(sdkManager, 0755); err != nil {
		d.logger.Warning("Failed to make sdkmanager executable: %v", err)
	}

	// Accept licenses
	d.logger.Info("Accepting SDK licenses...")
	if err := d.acceptLicenses(sdkManager, sdkRoot); err != nil {
		d.logger.Warning("Failed to accept licenses automatically: %v", err)
	}

	// Download platforms
	for _, apiLevel := range cfg.SDK.APILevels {
		pkgName := fmt.Sprintf("platforms;android-%d", apiLevel)
		d.logger.Info("Downloading API level %d...", apiLevel)
		if err := d.installSDKPackage(sdkManager, sdkRoot, pkgName); err != nil {
			d.logger.Warning("Failed to download %s: %v", pkgName, err)
		}
	}

	// Download build tools
	for _, version := range cfg.SDK.BuildTools {
		pkgName := fmt.Sprintf("build-tools;%s", version)
		d.logger.Info("Downloading Build Tools %s...", version)
		if err := d.installSDKPackage(sdkManager, sdkRoot, pkgName); err != nil {
			d.logger.Warning("Failed to download %s: %v", pkgName, err)
		}
	}

	// Download system images
	for _, image := range cfg.SDK.SystemImages {
		d.logger.Info("Downloading system image: %s", image)
		if err := d.installSDKPackage(sdkManager, sdkRoot, image); err != nil {
			d.logger.Warning("Failed to download %s: %v", image, err)
		}
	}

	// Download NDK
	if cfg.SDK.DownloadNDK {
		pkgName := fmt.Sprintf("ndk;%s", cfg.SDK.NDKVersion)
		d.logger.Info("Downloading NDK %s...", cfg.SDK.NDKVersion)
		if err := d.installSDKPackage(sdkManager, sdkRoot, pkgName); err != nil {
			d.logger.Warning("Failed to download NDK: %v", err)
		}
	}

	// Download CMake
	if cfg.SDK.DownloadCMake {
		d.logger.Info("Downloading CMake...")
		if err := d.installSDKPackage(sdkManager, sdkRoot, "cmake;3.22.1"); err != nil {
			d.logger.Warning("Failed to download CMake: %v", err)
		}
	}

	// Download extra packages
	for _, pkg := range cfg.SDK.ExtraPackages {
		d.logger.Info("Downloading: %s", pkg)
		if err := d.installSDKPackage(sdkManager, sdkRoot, pkg); err != nil {
			d.logger.Warning("Failed to download %s: %v", pkg, err)
		}
	}

	d.logger.Success("SDK packages downloaded")
	return nil
}

func (d *Downloader) extractSDKTools(cfg *config.Config) error {
	// Find SDK tools zip
	sdkDir := filepath.Join(cfg.DownloadDir, "sdk")
	sdkRoot := filepath.Join(sdkDir, "android-sdk")

	// Check if already extracted
	sdkManager := filepath.Join(sdkRoot, "cmdline-tools", "latest", "bin", "sdkmanager")
	if fileExists(sdkManager) {
		d.logger.Debug("SDK tools already extracted")
		return nil
	}

	d.logger.Info("Extracting SDK command-line tools...")

	// Find the zip file
	matches, err := filepath.Glob(filepath.Join(sdkDir, "commandlinetools-*.zip"))
	if err != nil || len(matches) == 0 {
		return fmt.Errorf("SDK tools zip not found")
	}

	zipFile := matches[0]

	// Create cmdline-tools directory
	cmdlineDir := filepath.Join(sdkRoot, "cmdline-tools")
	if err := os.MkdirAll(cmdlineDir, 0755); err != nil {
		return err
	}

	// Extract to temporary location
	tempDir := filepath.Join(cmdlineDir, "temp")
	if err := extractZip(zipFile, tempDir); err != nil {
		return fmt.Errorf("failed to extract SDK tools: %w", err)
	}

	// Move cmdline-tools to latest
	srcDir := filepath.Join(tempDir, "cmdline-tools")
	dstDir := filepath.Join(cmdlineDir, "latest")

	if err := os.Rename(srcDir, dstDir); err != nil {
		return fmt.Errorf("failed to move SDK tools: %w", err)
	}

	// Clean up temp
	os.RemoveAll(tempDir)

	d.logger.Success("SDK tools extracted")
	return nil
}

func (d *Downloader) acceptLicenses(sdkManager, sdkRoot string) error {
	cmd := exec.Command(sdkManager, "--licenses", fmt.Sprintf("--sdk_root=%s", sdkRoot))

	// Pipe "y" to accept all licenses
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	// Write "y" multiple times to accept all licenses
	for i := 0; i < 10; i++ {
		stdin.Write([]byte("y\n"))
	}
	stdin.Close()

	return cmd.Wait()
}

func (d *Downloader) installSDKPackage(sdkManager, sdkRoot, packageName string) error {
	cmd := exec.Command(sdkManager, packageName, fmt.Sprintf("--sdk_root=%s", sdkRoot))

	// Suppress verbose output but capture errors
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, string(output))
	}

	return nil
}

// extractZip extracts a zip file to destination
func extractZip(zipPath, destPath string) error {
	if err := os.MkdirAll(destPath, 0755); err != nil {
		return err
	}

	cmd := exec.Command("unzip", "-q", "-o", zipPath, "-d", destPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("unzip failed: %w", err)
	}

	return nil
}

// ListInstalledPackages lists all installed SDK packages
func (d *Downloader) ListInstalledPackages(cfg *config.Config) ([]string, error) {
	sdkRoot := filepath.Join(cfg.DownloadDir, "sdk", "android-sdk")
	sdkManager := filepath.Join(sdkRoot, "cmdline-tools", "latest", "bin", "sdkmanager")

	if !fileExists(sdkManager) {
		return nil, fmt.Errorf("sdkmanager not found")
	}

	cmd := exec.Command(sdkManager, "--list_installed", fmt.Sprintf("--sdk_root=%s", sdkRoot))
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	var packages []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "---") && !strings.HasPrefix(line, "Installed packages:") {
			packages = append(packages, line)
		}
	}

	return packages, nil
}
