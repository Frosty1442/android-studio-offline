package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the main configuration
type Config struct {
	Platform        string   `yaml:"platform"`         // linux, darwin, windows
	DownloadDir     string   `yaml:"download_dir"`     // Where to download files
	InstallDir      string   `yaml:"install_dir"`      // Where to install on target machine
	AndroidStudio   Studio   `yaml:"android_studio"`
	SDK             SDK      `yaml:"sdk"`
	Gradle          Gradle   `yaml:"gradle"`
	Dependencies    Deps     `yaml:"dependencies"`
	Options         Options  `yaml:"options"`
}

type Studio struct {
	Version     string `yaml:"version"`
	DownloadJDK bool   `yaml:"download_jdk"`
	JDKVersion  string `yaml:"jdk_version"`
}

type SDK struct {
	ToolsVersion    string   `yaml:"tools_version"`
	APILevels       []int    `yaml:"api_levels"`
	BuildTools      []string `yaml:"build_tools"`
	SystemImages    []string `yaml:"system_images"`
	DownloadNDK     bool     `yaml:"download_ndk"`
	NDKVersion      string   `yaml:"ndk_version"`
	DownloadCMake   bool     `yaml:"download_cmake"`
	ExtraPackages   []string `yaml:"extra_packages"`
}

type Gradle struct {
	Versions []string `yaml:"versions"`
}

type Deps struct {
	DownloadMaven  bool     `yaml:"download_maven"`
	CommonLibs     []string `yaml:"common_libs"`
	SampleProject  string   `yaml:"sample_project"`
}

type Options struct {
	ParallelDownloads int  `yaml:"parallel_downloads"`
	VerifyChecksums   bool `yaml:"verify_checksums"`
	ResumeDownloads   bool `yaml:"resume_downloads"`
	Verbose           bool `yaml:"verbose"`
}

// LoadConfig loads configuration from a YAML file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Set defaults
	if cfg.DownloadDir == "" {
		cfg.DownloadDir = "downloads"
	}
	if cfg.InstallDir == "" {
		cfg.InstallDir = "/opt/android-offline"
	}
	if cfg.Options.ParallelDownloads == 0 {
		cfg.Options.ParallelDownloads = 4
	}

	return &cfg, nil
}

// SaveConfig saves configuration to a YAML file
func SaveConfig(cfg *Config, path string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		Platform:    "linux",
		DownloadDir: "downloads",
		InstallDir:  "/opt/android-offline",
		AndroidStudio: Studio{
			Version:     "2024.2.1.10",
			DownloadJDK: true,
			JDKVersion:  "17",
		},
		SDK: SDK{
			ToolsVersion: "11076708",
			APILevels:    []int{33, 34},
			BuildTools:   []string{"33.0.2", "34.0.0"},
			SystemImages: []string{
				"system-images;android-33;google_apis;x86_64",
				"system-images;android-34;google_apis;x86_64",
			},
			DownloadNDK:   true,
			NDKVersion:    "25.2.9519653",
			DownloadCMake: true,
			ExtraPackages: []string{
				"platform-tools",
				"emulator",
				"extras;google;google_play_services",
			},
		},
		Gradle: Gradle{
			Versions: []string{"7.6", "8.0", "8.4"},
		},
		Dependencies: Deps{
			DownloadMaven: true,
			CommonLibs: []string{
				"androidx.core:core-ktx:1.12.0",
				"androidx.appcompat:appcompat:1.6.1",
				"com.google.android.material:material:1.11.0",
				"androidx.constraintlayout:constraintlayout:2.1.4",
			},
		},
		Options: Options{
			ParallelDownloads: 4,
			VerifyChecksums:   true,
			ResumeDownloads:   true,
			Verbose:           false,
		},
	}
}
