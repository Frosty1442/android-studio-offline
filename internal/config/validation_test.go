package config

import (
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: &Config{
				Platform:    "linux",
				DownloadDir: "downloads",
				InstallDir:  "/opt/android",
				AndroidStudio: Studio{
					Version:     "2024.1.1.1",
					DownloadJDK: true,
					JDKVersion:  "17",
				},
				SDK: SDK{
					ToolsVersion: "11076708",
					APILevels:    []int{33, 34},
					BuildTools:   []string{"34.0.0"},
				},
				Gradle: Gradle{
					Versions: []string{"8.0"},
				},
				Options: Options{
					ParallelDownloads: 4,
				},
			},
			wantErr: false,
		},
		{
			name: "invalid platform",
			cfg: &Config{
				Platform: "invalid",
				AndroidStudio: Studio{
					Version: "2024.1.1.1",
				},
				SDK: SDK{
					ToolsVersion: "11076708",
					APILevels:    []int{33},
					BuildTools:   []string{"34.0.0"},
				},
				Gradle: Gradle{
					Versions: []string{"8.0"},
				},
				Options: Options{
					ParallelDownloads: 4,
				},
			},
			wantErr: true,
		},
		{
			name: "missing API levels",
			cfg: &Config{
				Platform: "linux",
				AndroidStudio: Studio{
					Version: "2024.1.1.1",
				},
				SDK: SDK{
					ToolsVersion: "11076708",
					APILevels:    []int{},
					BuildTools:   []string{"34.0.0"},
				},
				Gradle: Gradle{
					Versions: []string{"8.0"},
				},
				Options: Options{
					ParallelDownloads: 4,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid API level",
			cfg: &Config{
				Platform: "linux",
				AndroidStudio: Studio{
					Version: "2024.1.1.1",
				},
				SDK: SDK{
					ToolsVersion: "11076708",
					APILevels:    []int{50}, // Too high
					BuildTools:   []string{"34.0.0"},
				},
				Gradle: Gradle{
					Versions: []string{"8.0"},
				},
				Options: Options{
					ParallelDownloads: 4,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid parallel downloads",
			cfg: &Config{
				Platform: "linux",
				AndroidStudio: Studio{
					Version: "2024.1.1.1",
				},
				SDK: SDK{
					ToolsVersion: "11076708",
					APILevels:    []int{33},
					BuildTools:   []string{"34.0.0"},
				},
				Gradle: Gradle{
					Versions: []string{"8.0"},
				},
				Options: Options{
					ParallelDownloads: 20, // Too high
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
