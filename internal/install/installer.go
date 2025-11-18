package install

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Frosty1442/android-studio-offline/internal/config"
	"github.com/Frosty1442/android-studio-offline/internal/ui"
)

// Installer handles installation on the target machine
type Installer struct {
	logger *ui.Logger
}

// NewInstaller creates a new installer
func NewInstaller(logger *ui.Logger) *Installer {
	return &Installer{logger: logger}
}

// Install performs the installation
func (i *Installer) Install(cfg *config.Config) error {
	i.logger.Step(1, 5, "Creating installation directories")
	if err := i.createDirectories(cfg); err != nil {
		return err
	}

	i.logger.Step(2, 5, "Installing Android Studio")
	if err := i.installAndroidStudio(cfg); err != nil {
		return err
	}

	i.logger.Step(3, 5, "Installing Android SDK")
	if err := i.installSDK(cfg); err != nil {
		return err
	}

	i.logger.Step(4, 5, "Installing Gradle")
	if err := i.installGradle(cfg); err != nil {
		return err
	}

	i.logger.Step(5, 5, "Configuring environment")
	if err := i.configureEnvironment(cfg); err != nil {
		return err
	}

	return nil
}

func (i *Installer) createDirectories(cfg *config.Config) error {
	dirs := []string{
		filepath.Join(cfg.InstallDir, "android-studio"),
		filepath.Join(cfg.InstallDir, "android-sdk"),
		filepath.Join(cfg.InstallDir, "gradle"),
		filepath.Join(cfg.InstallDir, "jdk"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	i.logger.Success("Installation directories created")
	return nil
}

func (i *Installer) installAndroidStudio(cfg *config.Config) error {
	// TODO: Extract Android Studio archive
	i.logger.Info("Android Studio installation not yet implemented")
	return nil
}

func (i *Installer) installSDK(cfg *config.Config) error {
	// TODO: Extract SDK components
	i.logger.Info("SDK installation not yet implemented")
	return nil
}

func (i *Installer) installGradle(cfg *config.Config) error {
	// TODO: Copy Gradle distributions
	i.logger.Info("Gradle installation not yet implemented")
	return nil
}

func (i *Installer) configureEnvironment(cfg *config.Config) error {
	// TODO: Configure shell environment
	i.logger.Info("Environment configuration not yet implemented")
	return nil
}
