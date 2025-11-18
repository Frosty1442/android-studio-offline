package download

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/Frosty1442/android-studio-offline/internal/config"
)

const gradleBaseURL = "https://services.gradle.org/distributions"

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

	if err := d.DownloadFile(ctx, url, destPath, cfg.Options.ResumeDownloads); err != nil {
		return fmt.Errorf("failed to download Gradle %s-%s: %w", version, distType, err)
	}

	// Also download checksum
	checksumURL := url + ".sha256"
	checksumPath := destPath + ".sha256"
	_ = d.DownloadFile(ctx, checksumURL, checksumPath, false) // Ignore checksum download errors

	return nil
}
