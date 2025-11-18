package install

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/Frosty1442/android-studio-offline/internal/config"
	"github.com/Frosty1442/android-studio-offline/internal/ui"
)

// Packager handles creating installation packages
type Packager struct {
	logger *ui.Logger
}

// NewPackager creates a new packager
func NewPackager(logger *ui.Logger) *Packager {
	return &Packager{logger: logger}
}

// CreatePackage creates a compressed tar.gz package of all downloads
func (p *Packager) CreatePackage(cfg *config.Config) (string, error) {
	timestamp := time.Now().Format("20060102")
	packageName := fmt.Sprintf("android-studio-offline-%s.tar.gz", timestamp)
	packagePath := filepath.Join("packages", packageName)

	// Create packages directory
	if err := os.MkdirAll("packages", 0755); err != nil {
		return "", fmt.Errorf("failed to create packages directory: %w", err)
	}

	p.logger.Info("Creating package: %s", packageName)
	p.logger.Info("This may take several minutes...")

	// Create package file
	outFile, err := os.Create(packagePath)
	if err != nil {
		return "", fmt.Errorf("failed to create package file: %w", err)
	}
	defer outFile.Close()

	// Create gzip writer
	gzWriter := gzip.NewWriter(outFile)
	defer gzWriter.Close()

	// Create tar writer
	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	// Add downloads directory
	if err := p.addDirectoryToTar(tarWriter, cfg.DownloadDir, ""); err != nil {
		return "", fmt.Errorf("failed to add downloads to package: %w", err)
	}

	// Add install script
	if err := p.addFileToTar(tarWriter, "scripts/install-offline.sh", "install-offline.sh"); err != nil {
		p.logger.Warning("Failed to add install script: %v", err)
	}

	// Add README
	if err := p.addFileToTar(tarWriter, "README.md", "README.md"); err != nil {
		p.logger.Warning("Failed to add README: %v", err)
	}

	p.logger.Success("Package created successfully")
	return packagePath, nil
}

func (p *Packager) addDirectoryToTar(tw *tar.Writer, sourceDir, targetDir string) error {
	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create tar header
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}

		// Update header name
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		if targetDir != "" {
			header.Name = filepath.Join(targetDir, relPath)
		} else {
			header.Name = relPath
		}

		// Write header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// Write file content if it's a regular file
		if info.Mode().IsRegular() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			if _, err := io.Copy(tw, file); err != nil {
				return err
			}

			p.logger.Debug("Added to package: %s", header.Name)
		}

		return nil
	})
}

func (p *Packager) addFileToTar(tw *tar.Writer, sourcePath, targetPath string) error {
	file, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return err
	}

	header.Name = targetPath

	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	if _, err := io.Copy(tw, file); err != nil {
		return err
	}

	return nil
}
