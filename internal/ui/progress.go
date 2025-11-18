package ui

import (
	"fmt"
	"sync"
	"time"
)

// ProgressWriter wraps an io.Writer and reports progress
type ProgressWriter struct {
	Total      int64
	Downloaded int64
	Name       string
	mu         sync.Mutex
	onProgress func(downloaded, total int64)
}

func NewProgressWriter(name string, total int64, onProgress func(downloaded, total int64)) *ProgressWriter {
	return &ProgressWriter{
		Name:       name,
		Total:      total,
		onProgress: onProgress,
	}
}

func (pw *ProgressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.mu.Lock()
	pw.Downloaded += int64(n)
	downloaded := pw.Downloaded
	pw.mu.Unlock()

	if pw.onProgress != nil {
		pw.onProgress(downloaded, pw.Total)
	}

	return n, nil
}

// ProgressBar represents a simple progress bar
type ProgressBar struct {
	Total       int64
	Current     int64
	Description string
	mu          sync.Mutex
	lastUpdate  time.Time
}

func NewProgressBar(description string, total int64) *ProgressBar {
	return &ProgressBar{
		Description: description,
		Total:       total,
		lastUpdate:  time.Now(),
	}
}

func (pb *ProgressBar) Add(n int64) {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	pb.Current += n

	// Update display every 100ms to avoid flickering
	if time.Since(pb.lastUpdate) > 100*time.Millisecond {
		pb.render()
		pb.lastUpdate = time.Now()
	}
}

func (pb *ProgressBar) Set(current int64) {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	pb.Current = current

	if time.Since(pb.lastUpdate) > 100*time.Millisecond {
		pb.render()
		pb.lastUpdate = time.Now()
	}
}

func (pb *ProgressBar) Finish() {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	pb.Current = pb.Total
	pb.render()
	fmt.Println()
}

func (pb *ProgressBar) render() {
	percent := float64(0)
	if pb.Total > 0 {
		percent = float64(pb.Current) / float64(pb.Total) * 100
	}

	barWidth := 40
	filled := int(float64(barWidth) * percent / 100)
	if filled > barWidth {
		filled = barWidth
	}

	bar := ""
	for i := 0; i < barWidth; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}

	// Format size
	current := formatBytes(pb.Current)
	total := formatBytes(pb.Total)

	fmt.Printf("\r%-30s [%s] %6.2f%% %s/%s",
		truncate(pb.Description, 30),
		bar,
		percent,
		current,
		total,
	)
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

// Logger provides colored logging
type Logger struct {
	Verbose bool
}

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
)

func (l *Logger) Info(format string, args ...interface{}) {
	fmt.Printf(colorBlue+"ℹ "+colorReset+format+"\n", args...)
}

func (l *Logger) Success(format string, args ...interface{}) {
	fmt.Printf(colorGreen+"✓ "+colorReset+format+"\n", args...)
}

func (l *Logger) Warning(format string, args ...interface{}) {
	fmt.Printf(colorYellow+"⚠ "+colorReset+format+"\n", args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	fmt.Printf(colorRed+"✗ "+colorReset+format+"\n", args...)
}

func (l *Logger) Debug(format string, args ...interface{}) {
	if l.Verbose {
		fmt.Printf(colorCyan+"🔍 "+colorReset+format+"\n", args...)
	}
}

func (l *Logger) Step(step int, total int, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf(colorCyan+"[%d/%d]"+colorReset+" %s\n", step, total, msg)
}
