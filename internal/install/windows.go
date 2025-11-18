package install

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// configureWindowsEnvironment sets up Windows environment variables
func (fi *FullInstaller) configureWindowsEnvironment() error {
	fi.logger.Info("Configuring Windows environment variables...")

	androidHome := filepath.Join(fi.cfg.InstallDir, "android-sdk")
	studioPath := filepath.Join(fi.cfg.InstallDir, "android-studio", "bin")
	platformTools := filepath.Join(androidHome, "platform-tools")

	// Use PowerShell to set system environment variables
	commands := []struct {
		name  string
		value string
	}{
		{"ANDROID_HOME", androidHome},
		{"ANDROID_SDK_ROOT", androidHome},
	}

	for _, cmd := range commands {
		// Set user environment variable (doesn't require admin)
		psCmd := fmt.Sprintf("[Environment]::SetEnvironmentVariable('%s', '%s', 'User')",
			cmd.name, cmd.value)

		if err := fi.runPowerShell(psCmd); err != nil {
			fi.logger.Warning("Failed to set %s: %v", cmd.name, err)
			fi.logger.Info("Please set manually: setx %s \"%s\"", cmd.name, cmd.value)
		} else {
			fi.logger.Success("Set %s", cmd.name)
		}
	}

	// Update PATH
	fi.logger.Info("Updating PATH...")

	// Get current user PATH
	getPathCmd := "[Environment]::GetEnvironmentVariable('Path', 'User')"
	currentPath, err := fi.runPowerShellOutput(getPathCmd)
	if err != nil {
		fi.logger.Warning("Failed to get current PATH: %v", err)
		fi.logger.Info("Please add to PATH manually:")
		fi.logger.Info("  %s", platformTools)
		fi.logger.Info("  %s", studioPath)
		return nil
	}

	// Add new paths if not already present
	newPaths := []string{platformTools, studioPath}
	pathToSet := strings.TrimSpace(currentPath)

	for _, newPath := range newPaths {
		if !strings.Contains(strings.ToLower(pathToSet), strings.ToLower(newPath)) {
			if pathToSet != "" {
				pathToSet += ";"
			}
			pathToSet += newPath
		}
	}

	// Set updated PATH
	setPathCmd := fmt.Sprintf("[Environment]::SetEnvironmentVariable('Path', '%s', 'User')",
		pathToSet)

	if err := fi.runPowerShell(setPathCmd); err != nil {
		fi.logger.Warning("Failed to update PATH: %v", err)
		fi.logger.Info("Please add to PATH manually:")
		for _, path := range newPaths {
			fi.logger.Info("  %s", path)
		}
	} else {
		fi.logger.Success("PATH updated")
	}

	// Create batch file for easy environment reload
	batchFile := filepath.Join(fi.cfg.InstallDir, "setup-environment.bat")
	batchContent := fmt.Sprintf(`@echo off
REM Android Studio Offline Environment Setup
REM Run this file to set environment variables for current session

set ANDROID_HOME=%s
set ANDROID_SDK_ROOT=%s
set PATH=%s;%s;%%PATH%%

echo Environment configured for current session!
echo.
echo To make permanent, run PowerShell as Administrator:
echo   setx ANDROID_HOME "%s"
echo   setx ANDROID_SDK_ROOT "%s"
echo.
`, androidHome, androidHome, platformTools, studioPath, androidHome, androidHome)

	if err := os.WriteFile(batchFile, []byte(batchContent), 0644); err != nil {
		fi.logger.Warning("Failed to create setup batch file: %v", err)
	} else {
		fi.logger.Info("Created setup script: %s", batchFile)
		fi.logger.Info("Run this in new command prompts to set environment")
	}

	fi.logger.Success("Windows environment configured")
	fi.logger.Warning("You may need to log out and back in for changes to take effect")
	fi.logger.Info("Or run: %s", batchFile)

	return nil
}

func (fi *FullInstaller) runPowerShell(command string) error {
	cmd := exec.Command("powershell", "-Command", command)
	return cmd.Run()
}

func (fi *FullInstaller) runPowerShellOutput(command string) (string, error) {
	cmd := exec.Command("powershell", "-Command", command)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
