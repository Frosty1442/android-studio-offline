package download

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/Frosty1442/android-studio-offline/internal/config"
)

const (
	gradleBaseURL    = "https://services.gradle.org/distributions"
	gradleWrapperURL = "https://raw.githubusercontent.com/gradle/gradle/master/gradle/wrapper"
)

// DownloadGradle downloads all specified Gradle versions
func (d *Downloader) DownloadGradle(ctx context.Context, cfg *config.Config) error {
	d.logger.Step(4, 5, "Downloading Gradle distributions")

	if len(cfg.Gradle.Versions) == 0 {
		d.logger.Warning("No Gradle versions specified, skipping")
		return nil
	}

	d.logger.Info("Downloading %d Gradle versions", len(cfg.Gradle.Versions))

	var wg sync.WaitGroup
	errors := make(chan error, len(cfg.Gradle.Versions)*2)
	semaphore := make(chan struct{}, cfg.Options.ParallelDownloads)

	for _, version := range cfg.Gradle.Versions {
		// Download both 'bin' and 'all' distributions
		for _, distType := range []string{"bin", "all"} {
			wg.Add(1)
			go func(ver, dt string) {
				defer wg.Done()

				semaphore <- struct{}{}
				defer func() { <-semaphore }()

				if err := d.downloadGradleVersion(ctx, cfg, ver, dt); err != nil {
					errors <- err
				}
			}(version, distType)
		}
	}

	wg.Wait()
	close(errors)

	// Check for errors
	var errList []error
	for err := range errors {
		errList = append(errList, err)
	}

	if len(errList) > 0 {
		d.logger.Warning("Some Gradle downloads failed: %v", errList)
		// Don't fail completely, some versions might have succeeded
	}

	d.logger.Success("Gradle distributions downloaded")

	// Download Gradle wrapper files
	if err := d.downloadGradleWrapper(ctx, cfg); err != nil {
		d.logger.Warning("Failed to download Gradle wrapper: %v", err)
		// Don't fail completely - wrapper is useful but not critical
	}

	return nil
}

func (d *Downloader) downloadGradleVersion(ctx context.Context, cfg *config.Config, version, distType string) error {
	filename := fmt.Sprintf("gradle-%s-%s.zip", version, distType)
	url := fmt.Sprintf("%s/%s", gradleBaseURL, filename)
	destPath := filepath.Join(cfg.DownloadDir, "gradle", "distributions", filename)

	// Check if already downloaded
	if fileExists(destPath) {
		d.logger.Debug("Gradle %s-%s already downloaded", version, distType)
		return nil
	}

	d.logger.Debug("Downloading Gradle %s-%s", version, distType)

	if err := d.DownloadAndVerify(ctx, url, destPath, cfg.Options.ResumeDownloads, cfg.Options.VerifyChecksums); err != nil {
		return fmt.Errorf("failed to download Gradle %s-%s: %w", version, distType, err)
	}

	return nil
}

// downloadGradleWrapper downloads Gradle wrapper files
func (d *Downloader) downloadGradleWrapper(ctx context.Context, cfg *config.Config) error {
	d.logger.Info("Downloading Gradle wrapper files")

	// Safety check - need at least one Gradle version for wrapper scripts
	if len(cfg.Gradle.Versions) == 0 {
		d.logger.Warning("No Gradle versions specified, skipping wrapper")
		return nil
	}

	wrapperDir := filepath.Join(cfg.DownloadDir, "gradle", "wrapper")

	// Wrapper files to download from GitHub
	wrapperFiles := map[string]string{
		"gradle-wrapper.jar":        filepath.Join(wrapperDir, "gradle-wrapper.jar"),
		"gradle-wrapper.properties": filepath.Join(wrapperDir, "gradle-wrapper.properties"),
	}

	scriptFiles := map[string]string{
		"gradlew":     filepath.Join(cfg.DownloadDir, "gradle", "gradlew"),
		"gradlew.bat": filepath.Join(cfg.DownloadDir, "gradle", "gradlew.bat"),
	}

	// Download wrapper JAR and properties
	for filename, destPath := range wrapperFiles {
		if fileExists(destPath) {
			continue
		}

		url := fmt.Sprintf("%s/%s", gradleWrapperURL, filename)
		if err := d.DownloadFile(ctx, url, destPath, false); err != nil {
			return fmt.Errorf("failed to download %s: %w", filename, err)
		}
	}

	// Download wrapper scripts from latest stable Gradle version
	// Use the highest version specified in config
	latestVersion := cfg.Gradle.Versions[len(cfg.Gradle.Versions)-1]
	gradleScriptURL := fmt.Sprintf("https://raw.githubusercontent.com/gradle/gradle/v%s", latestVersion)

	for filename, destPath := range scriptFiles {
		if fileExists(destPath) {
			continue
		}

		url := fmt.Sprintf("%s/%s", gradleScriptURL, filename)
		if err := d.DownloadFile(ctx, url, destPath, false); err != nil {
			// Try alternate URL (master branch)
			url = fmt.Sprintf("https://raw.githubusercontent.com/gradle/gradle/master/%s", filename)
			if err := d.DownloadFile(ctx, url, destPath, false); err != nil {
				d.logger.Warning("Failed to download %s: %v", filename, err)
				continue
			}
		}

		// Make scripts executable on Unix systems
		if filename == "gradlew" {
			_ = os.Chmod(destPath, 0755)
		}
	}

	d.logger.Success("Gradle wrapper files downloaded")
	return nil
}
