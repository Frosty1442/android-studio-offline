package download

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Frosty1442/android-studio-offline/internal/ui"
)

// Downloader handles file downloads with progress tracking
type Downloader struct {
	client     *http.Client
	logger     *ui.Logger
	maxRetries int
	parallel   int
}

// NewDownloader creates a new downloader
func NewDownloader(logger *ui.Logger, parallel int) *Downloader {
	return &Downloader{
		client: &http.Client{
			Timeout: 30 * time.Minute,
			Transport: &http.Transport{
				MaxIdleConns:        parallel,
				MaxIdleConnsPerHost: parallel,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		logger:     logger,
		maxRetries: 3,
		parallel:   parallel,
	}
}

// DownloadFile downloads a file with progress tracking and resume support
func (d *Downloader) DownloadFile(ctx context.Context, url, destPath string, resume bool) error {
	// Create destination directory
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Check if file already exists
	existingSize := int64(0)
	if resume {
		if info, err := os.Stat(destPath); err == nil {
			existingSize = info.Size()
		}
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add range header for resume
	if resume && existingSize > 0 {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", existingSize))
	}

	// Execute request
	resp, err := d.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	// Check status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Get total size
	totalSize := resp.ContentLength
	if resume && resp.StatusCode == http.StatusPartialContent {
		totalSize += existingSize
	}

	// Open file for writing
	flags := os.O_CREATE | os.O_WRONLY
	if resume && resp.StatusCode == http.StatusPartialContent {
		flags |= os.O_APPEND
	} else {
		flags |= os.O_TRUNC
		existingSize = 0
	}

	out, err := os.OpenFile(destPath, flags, 0644)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Create progress bar
	filename := filepath.Base(destPath)
	progress := ui.NewProgressBar(filename, totalSize)
	if existingSize > 0 {
		progress.Set(existingSize)
	}

	// Download with progress
	writer := &progressWriter{
		writer:   out,
		progress: progress,
	}

	_, err = io.Copy(writer, resp.Body)
	progress.Finish()

	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}

type progressWriter struct {
	writer   io.Writer
	progress *ui.ProgressBar
	written  int64
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	n, err := pw.writer.Write(p)
	pw.written += int64(n)
	pw.progress.Add(int64(n))
	return n, err
}

// DownloadFiles downloads multiple files in parallel
func (d *Downloader) DownloadFiles(ctx context.Context, files map[string]string, resume bool) error {
	type downloadJob struct {
		url  string
		dest string
	}

	jobs := make(chan downloadJob, len(files))
	for url, dest := range files {
		jobs <- downloadJob{url: url, dest: dest}
	}
	close(jobs)

	var wg sync.WaitGroup
	errors := make(chan error, len(files))
	semaphore := make(chan struct{}, d.parallel)

	for job := range jobs {
		wg.Add(1)
		go func(j downloadJob) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if err := d.DownloadFile(ctx, j.url, j.dest, resume); err != nil {
				errors <- fmt.Errorf("failed to download %s: %w", j.url, err)
			}
		}(job)
	}

	wg.Wait()
	close(errors)

	// Collect errors
	var downloadErrors []error
	for err := range errors {
		downloadErrors = append(downloadErrors, err)
	}

	if len(downloadErrors) > 0 {
		return fmt.Errorf("download failures: %v", downloadErrors)
	}

	return nil
}

// VerifyChecksum verifies a file's SHA256 checksum
func VerifyChecksum(filePath, expectedHash string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return false, err
	}

	actualHash := fmt.Sprintf("%x", hash.Sum(nil))
	return actualHash == expectedHash, nil
}

// GetFileSize gets the size of a remote file without downloading
func (d *Downloader) GetFileSize(ctx context.Context, url string) (int64, error) {
	req, err := http.NewRequestWithContext(ctx, "HEAD", url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return resp.ContentLength, nil
}
