package download

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Frosty1442/android-studio-offline/internal/config"
)

const (
	mavenCentral  = "https://repo1.maven.org/maven2"
	googleMaven   = "https://dl.google.com/dl/android/maven2"
	gradlePlugins = "https://plugins.gradle.org/m2"
)

// DownloadMavenDependencies downloads common Maven dependencies
func (d *Downloader) DownloadMavenDependencies(ctx context.Context, cfg *config.Config) error {
	if !cfg.Dependencies.DownloadMaven {
		d.logger.Info("Maven dependency download disabled")
		return nil
	}

	d.logger.Info("Downloading Maven dependencies...")

	var wg sync.WaitGroup
	errors := make(chan error, len(cfg.Dependencies.CommonLibs))
	semaphore := make(chan struct{}, cfg.Options.ParallelDownloads)

	for _, dep := range cfg.Dependencies.CommonLibs {
		wg.Add(1)
		go func(dependency string) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if err := d.downloadDependency(ctx, cfg, dependency); err != nil {
				errors <- fmt.Errorf("failed to download %s: %w", dependency, err)
			}
		}(dep)
	}

	wg.Wait()
	close(errors)

	// Log errors but don't fail completely
	errorCount := 0
	for err := range errors {
		d.logger.Warning("%v", err)
		errorCount++
	}

	if errorCount > 0 {
		d.logger.Warning("%d dependencies failed to download", errorCount)
	} else {
		d.logger.Success("All Maven dependencies downloaded")
	}

	return nil
}

func (d *Downloader) downloadDependency(ctx context.Context, cfg *config.Config, dependency string) error {
	parts := strings.Split(dependency, ":")
	if len(parts) != 3 {
		return fmt.Errorf("invalid dependency format: %s", dependency)
	}

	group := parts[0]
	artifact := parts[1]
	version := parts[2]

	// Determine repository
	repo := mavenCentral
	if strings.HasPrefix(group, "androidx.") || strings.HasPrefix(group, "com.google.android") ||
	   strings.HasPrefix(group, "com.android.") {
		repo = googleMaven
	}

	// Convert group to path
	groupPath := strings.ReplaceAll(group, ".", "/")
	basePath := fmt.Sprintf("%s/%s/%s", groupPath, artifact, version)

	// Base filename
	baseFilename := fmt.Sprintf("%s-%s", artifact, version)

	// Determine output directory
	var outputBase string
	if repo == googleMaven {
		outputBase = filepath.Join(cfg.DownloadDir, "dependencies", "google-repo")
	} else {
		outputBase = filepath.Join(cfg.DownloadDir, "dependencies", "maven-repo")
	}

	outputDir := filepath.Join(outputBase, basePath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	// Download artifacts
	artifacts := []string{
		baseFilename + ".pom",     // Always download POM
		baseFilename + ".jar",     // JAR if exists
		baseFilename + ".aar",     // AAR for Android libraries
		baseFilename + "-sources.jar", // Sources
	}

	downloaded := false
	for _, artifactFile := range artifacts {
		url := fmt.Sprintf("%s/%s/%s", repo, basePath, artifactFile)
		destPath := filepath.Join(outputDir, artifactFile)

		// Skip if already exists
		if fileExists(destPath) {
			downloaded = true
			continue
		}

		// Try to download (some may not exist, which is OK)
		err := d.DownloadFile(ctx, url, destPath, false)
		if err == nil {
			downloaded = true
			d.logger.Debug("Downloaded: %s", artifactFile)
		} else {
			// Remove failed download
			os.Remove(destPath)
		}
	}

	if !downloaded {
		return fmt.Errorf("no artifacts downloaded for %s", dependency)
	}

	return nil
}

// DownloadAndroidGradlePlugin downloads Android Gradle Plugin versions
func (d *Downloader) DownloadAndroidGradlePlugin(ctx context.Context, cfg *config.Config) error {
	d.logger.Info("Downloading Android Gradle Plugin...")

	agpVersions := []string{"8.1.0", "8.1.1", "8.1.2", "8.2.0", "8.3.0"}

	for _, version := range agpVersions {
		dep := fmt.Sprintf("com.android.tools.build:gradle:%s", version)
		if err := d.downloadDependency(ctx, cfg, dep); err != nil {
			d.logger.Warning("Failed to download AGP %s: %v", version, err)
		}
	}

	d.logger.Success("Android Gradle Plugin downloaded")
	return nil
}

// DownloadKotlinGradlePlugin downloads Kotlin Gradle Plugin
func (d *Downloader) DownloadKotlinGradlePlugin(ctx context.Context, cfg *config.Config) error {
	d.logger.Info("Downloading Kotlin Gradle Plugin...")

	kotlinVersions := []string{"1.9.0", "1.9.10", "1.9.20", "1.9.21", "1.9.22"}

	for _, version := range kotlinVersions {
		dep := fmt.Sprintf("org.jetbrains.kotlin:kotlin-gradle-plugin:%s", version)
		if err := d.downloadDependency(ctx, cfg, dep); err != nil {
			d.logger.Warning("Failed to download Kotlin %s: %v", version, err)
		}
	}

	d.logger.Success("Kotlin Gradle Plugin downloaded")
	return nil
}
