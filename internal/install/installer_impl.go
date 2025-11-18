package install

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Frosty1442/android-studio-offline/internal/config"
	"github.com/Frosty1442/android-studio-offline/internal/ui"
	"github.com/Frosty1442/android-studio-offline/internal/util"
)

// FullInstaller implements complete installation
type FullInstaller struct {
	logger *ui.Logger
	cfg    *config.Config
}

// NewFullInstaller creates a new full installer
func NewFullInstaller(logger *ui.Logger, cfg *config.Config) *FullInstaller {
	return &FullInstaller{
		logger: logger,
		cfg:    cfg,
	}
}

// Install performs complete installation
func (fi *FullInstaller) Install() error {
	fi.logger.Info("Starting Android Studio offline installation")
	fi.logger.Info("Install directory: %s", fi.cfg.InstallDir)
	fmt.Println()

	steps := []struct {
		name string
		fn   func() error
	}{
		{"Create directories", fi.createDirectories},
		{"Install Android Studio", fi.installAndroidStudio},
		{"Install JDK", fi.installJDK},
		{"Install Android SDK", fi.installAndroidSDK},
		{"Install Gradle", fi.installGradle},
		{"Install dependencies", fi.installDependencies},
		{"Configure environment", fi.configureEnvironment},
		{"Create launcher", fi.createLauncher},
	}

	for i, step := range steps {
		fi.logger.Step(i+1, len(steps), "%s", step.name)
		if err := step.fn(); err != nil {
			return fmt.Errorf("failed at step '%s': %w", step.name, err)
		}
	}

	fi.logger.Success("Installation completed successfully!")
	return nil
}

func (fi *FullInstaller) createDirectories() error {
	dirs := []string{
		filepath.Join(fi.cfg.InstallDir, "android-studio"),
		filepath.Join(fi.cfg.InstallDir, "android-sdk"),
		filepath.Join(fi.cfg.InstallDir, "gradle"),
		filepath.Join(fi.cfg.InstallDir, "jdk"),
		filepath.Join(fi.cfg.InstallDir, "maven-repo"),
		filepath.Join(fi.cfg.InstallDir, "google-repo"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create %s: %w", dir, err)
		}
	}

	return nil
}

func (fi *FullInstaller) installAndroidStudio() error {
	downloadDir := filepath.Join(fi.cfg.DownloadDir, "android-studio")

	// Find Android Studio archive
	var archivePath string
	err := filepath.Walk(downloadDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.Contains(info.Name(), "android-studio") {
			archivePath = path
			return filepath.SkipDir
		}
		return nil
	})

	if err != nil || archivePath == "" {
		fi.logger.Warning("Android Studio archive not found, skipping")
		return nil
	}

	fi.logger.Info("Extracting: %s", filepath.Base(archivePath))

	installDir := fi.cfg.InstallDir
	if err := util.ExtractArchive(archivePath, installDir); err != nil {
		return err
	}

	fi.logger.Success("Android Studio installed")
	return nil
}

func (fi *FullInstaller) installJDK() error {
	downloadDir := filepath.Join(fi.cfg.DownloadDir, "jdk")

	// Find JDK archive
	var archivePath string
	err := filepath.Walk(downloadDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasPrefix(info.Name(), "jdk-") {
			archivePath = path
			return filepath.SkipDir
		}
		return nil
	})

	if err != nil || archivePath == "" {
		fi.logger.Info("JDK archive not found, skipping (Android Studio includes JDK)")
		return nil
	}

	fi.logger.Info("Extracting: %s", filepath.Base(archivePath))

	jdkDir := filepath.Join(fi.cfg.InstallDir, "jdk")
	if err := util.ExtractArchive(archivePath, jdkDir); err != nil {
		return err
	}

	fi.logger.Success("JDK installed")
	return nil
}

func (fi *FullInstaller) installAndroidSDK() error {
	downloadSDK := filepath.Join(fi.cfg.DownloadDir, "sdk", "android-sdk")
	installSDK := filepath.Join(fi.cfg.InstallDir, "android-sdk")

	if !util.FileExists(downloadSDK) {
		fi.logger.Warning("Android SDK not found in downloads, skipping")
		return nil
	}

	fi.logger.Info("Copying Android SDK...")

	if err := util.CopyDir(downloadSDK, installSDK); err != nil {
		return fmt.Errorf("failed to copy SDK: %w", err)
	}

	fi.logger.Success("Android SDK installed")
	return nil
}

func (fi *FullInstaller) installGradle() error {
	downloadGradle := filepath.Join(fi.cfg.DownloadDir, "gradle", "distributions")
	installGradle := filepath.Join(fi.cfg.InstallDir, "gradle", "distributions")

	if !util.FileExists(downloadGradle) {
		fi.logger.Warning("Gradle distributions not found, skipping")
		return nil
	}

	fi.logger.Info("Copying Gradle distributions...")

	if err := util.CopyDir(downloadGradle, installGradle); err != nil {
		return fmt.Errorf("failed to copy Gradle: %w", err)
	}

	fi.logger.Success("Gradle installed")
	return nil
}

func (fi *FullInstaller) installDependencies() error {
	// Copy Maven repositories
	for _, repo := range []string{"maven-repo", "google-repo"} {
		src := filepath.Join(fi.cfg.DownloadDir, "dependencies", repo)
		dst := filepath.Join(fi.cfg.InstallDir, repo)

		if !util.FileExists(src) {
			continue
		}

		fi.logger.Info("Copying %s...", repo)
		if err := util.CopyDir(src, dst); err != nil {
			fi.logger.Warning("Failed to copy %s: %v", repo, err)
		}
	}

	fi.logger.Success("Dependencies installed")
	return nil
}

func (fi *FullInstaller) configureEnvironment() error {
	fi.logger.Info("Configuring environment...")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// Handle Windows separately
	if runtime.GOOS == "windows" {
		return fi.configureWindowsEnvironment()
	}

	// Determine shell config file for Unix-like systems
	var rcFile string
	if util.FileExists(filepath.Join(homeDir, ".zshrc")) {
		rcFile = filepath.Join(homeDir, ".zshrc")
	} else {
		rcFile = filepath.Join(homeDir, ".bashrc")
	}

	// Backup existing file
	if util.FileExists(rcFile) {
		backupFile := rcFile + ".backup-android-offline"
		util.CopyFile(rcFile, backupFile)
	}

	// Append environment variables
	envVars := fmt.Sprintf(`
# Android Studio Offline Environment
# Added by android-offline installer on %s
export ANDROID_HOME="%s/android-sdk"
export ANDROID_SDK_ROOT="%s/android-sdk"
export PATH="$PATH:$ANDROID_HOME/platform-tools:$ANDROID_HOME/tools"
export PATH="$PATH:%s/android-studio/bin"
`,
		os.Getenv("USER"),
		fi.cfg.InstallDir,
		fi.cfg.InstallDir,
		fi.cfg.InstallDir,
	)

	file, err := os.OpenFile(rcFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(envVars); err != nil {
		return err
	}

	// Create Gradle config
	gradleDir := filepath.Join(homeDir, ".gradle")
	os.MkdirAll(gradleDir, 0755)

	gradleProps := filepath.Join(gradleDir, "gradle.properties")
	propsContent := fmt.Sprintf(`# Android Offline Configuration
org.gradle.daemon=true
org.gradle.parallel=true
org.gradle.caching=true
android.useAndroidX=true
android.offlineRoot=%s
`, fi.cfg.InstallDir)

	os.WriteFile(gradleProps, []byte(propsContent), 0644)

	// Create Gradle init script for offline repos
	initDir := filepath.Join(gradleDir, "init.d")
	os.MkdirAll(initDir, 0755)

	initScript := filepath.Join(initDir, "offline-repos.gradle")
	initContent := fmt.Sprintf(`// Offline repository configuration
allprojects {
    repositories {
        maven { url = uri("%s/google-repo") }
        maven { url = uri("%s/maven-repo") }
    }
    buildscript {
        repositories {
            maven { url = uri("%s/google-repo") }
            maven { url = uri("%s/maven-repo") }
        }
    }
}
`, fi.cfg.InstallDir, fi.cfg.InstallDir, fi.cfg.InstallDir, fi.cfg.InstallDir)

	os.WriteFile(initScript, []byte(initContent), 0644)

	fi.logger.Success("Environment configured")
	fi.logger.Info("Reload your shell: source %s", rcFile)
	return nil
}

func (fi *FullInstaller) createLauncher() error {
	if runtime.GOOS != "linux" {
		return nil
	}

	fi.logger.Info("Creating desktop launcher...")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	desktopDir := filepath.Join(homeDir, ".local", "share", "applications")
	os.MkdirAll(desktopDir, 0755)

	desktopFile := filepath.Join(desktopDir, "android-studio-offline.desktop")
	iconPath := filepath.Join(fi.cfg.InstallDir, "android-studio", "bin", "studio.png")
	execPath := filepath.Join(fi.cfg.InstallDir, "android-studio", "bin", "studio.sh")

	content := fmt.Sprintf(`[Desktop Entry]
Version=1.0
Type=Application
Name=Android Studio (Offline)
Icon=%s
Exec=%s
Comment=Android Studio IDE (Offline Installation)
Categories=Development;IDE;
Terminal=false
StartupWMClass=jetbrains-studio
`, iconPath, execPath)

	if err := os.WriteFile(desktopFile, []byte(content), 0755); err != nil {
		fi.logger.Warning("Failed to create desktop entry: %v", err)
		return nil
	}

	fi.logger.Success("Desktop launcher created")
	return nil
}
