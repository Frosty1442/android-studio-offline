package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/Frosty1442/android-studio-offline/internal/config"
	"github.com/Frosty1442/android-studio-offline/internal/download"
	"github.com/Frosty1442/android-studio-offline/internal/install"
	"github.com/Frosty1442/android-studio-offline/internal/ui"
	"github.com/Frosty1442/android-studio-offline/internal/util"
	"github.com/spf13/cobra"
)

var (
	version   = "1.0.0"
	cfgFile   string
	verbose   bool
	logger    *ui.Logger
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "android-offline",
	Short: "Android Studio Offline Installer",
	Long: `A cross-platform tool to create and deploy offline Android Studio installations.

This tool helps you download all necessary Android development components
on an internet-connected machine and install them on air-gapped systems.`,
	Version: version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logger = &ui.Logger{Verbose: verbose}
	},
}

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download all Android Studio components",
	Long:  `Downloads Android Studio, SDK, Gradle, and dependencies for offline installation.`,
	RunE:  runDownload,
}

var packageCmd = &cobra.Command{
	Use:   "package",
	Short: "Create portable installation package",
	Long:  `Packages all downloaded components into a compressed archive for transfer.`,
	RunE:  runPackage,
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install on offline machine",
	Long:  `Installs Android Studio and all components on the target machine.`,
	RunE:  runInstall,
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create default configuration file",
	Long:  `Generates a default configuration file that you can customize.`,
	RunE:  runInit,
}

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify downloaded components",
	Long:  `Checks that all configured components have been downloaded successfully.`,
	RunE:  runVerify,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.yaml", "config file")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	rootCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(packageCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(verifyCmd)
}

func runDownload(cmd *cobra.Command, args []string) error {
	ctx, cancel := setupSignalHandler()
	defer cancel()

	// Load configuration
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	// Validate configuration before downloading
	if err := cfg.Validate(); err != nil {
		logger.Error("Configuration validation failed")
		return fmt.Errorf("invalid configuration: %w", err)
	}

	logger.Info("Starting download process")
	logger.Info("Platform: %s", cfg.Platform)
	logger.Info("Download directory: %s", cfg.DownloadDir)
	fmt.Println()

	// Create downloader
	dl := download.NewDownloader(logger, cfg.Options.ParallelDownloads)

	// Download Android Studio
	if err := dl.DownloadAndroidStudio(ctx, cfg); err != nil {
		return fmt.Errorf("failed to download Android Studio: %w", err)
	}

	// Download JDK
	if err := dl.DownloadJDK(ctx, cfg); err != nil {
		return fmt.Errorf("failed to download JDK: %w", err)
	}

	// Download SDK tools
	if err := dl.DownloadSDKCommandLineTools(ctx, cfg); err != nil {
		return fmt.Errorf("failed to download SDK tools: %w", err)
	}

	// Download platform tools
	if err := dl.DownloadPlatformTools(ctx, cfg); err != nil {
		return fmt.Errorf("failed to download platform tools: %w", err)
	}

	// Download Gradle
	if err := dl.DownloadGradle(ctx, cfg); err != nil {
		return fmt.Errorf("failed to download Gradle: %w", err)
	}

	// Download SDK packages
	logger.Info("")
	if err := dl.DownloadSDKPackages(ctx, cfg); err != nil {
		logger.Warning("SDK packages download had errors: %v", err)
		logger.Info("Some SDK packages may be missing, but continuing...")
	}

	// Download Maven dependencies
	logger.Info("")
	if err := dl.DownloadMavenDependencies(ctx, cfg); err != nil {
		logger.Warning("Maven dependencies download had errors: %v", err)
	}

	// Download Android Gradle Plugin
	if err := dl.DownloadAndroidGradlePlugin(ctx, cfg); err != nil {
		logger.Warning("AGP download had errors: %v", err)
	}

	// Download Kotlin Gradle Plugin
	if err := dl.DownloadKotlinGradlePlugin(ctx, cfg); err != nil {
		logger.Warning("Kotlin plugin download had errors: %v", err)
	}

	fmt.Println()
	logger.Success("All downloads completed successfully!")
	logger.Info("Next step: Run 'android-offline package' to create installation package")

	return nil
}

func runPackage(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	logger.Info("Creating installation package...")

	pkgr := install.NewPackager(logger)
	packagePath, err := pkgr.CreatePackage(cfg)
	if err != nil {
		return fmt.Errorf("failed to create package: %w", err)
	}

	logger.Success("Package created: %s", packagePath)
	logger.Info("Transfer this file to your offline machine and run 'android-offline install'")

	return nil
}

func runInstall(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	// Validate configuration before installing
	if err := cfg.Validate(); err != nil {
		logger.Error("Configuration validation failed")
		return fmt.Errorf("invalid configuration: %w", err)
	}

	logger.Info("Starting installation...")
	fmt.Println()

	fullInstaller := install.NewFullInstaller(logger, cfg)
	if err := fullInstaller.Install(); err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}

	fmt.Println()
	logger.Success("Installation completed successfully!")
	logger.Info("Android Studio installed to: %s", cfg.InstallDir)

	startCmd := filepath.Join(cfg.InstallDir, "android-studio", "bin", "studio.sh")
	logger.Info("Start Android Studio: %s", startCmd)
	logger.Info("Or use the desktop launcher (Linux only)")

	return nil
}

func runInit(cmd *cobra.Command, args []string) error {
	if _, err := os.Stat(cfgFile); err == nil {
		logger.Warning("Config file already exists: %s", cfgFile)
		fmt.Print("Overwrite? (y/N): ")
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			logger.Info("Aborted")
			return nil
		}
	}

	cfg := config.DefaultConfig()
	if err := config.SaveConfig(cfg, cfgFile); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	logger.Success("Configuration file created: %s", cfgFile)
	logger.Info("Edit this file to customize your installation, then run 'android-offline download'")

	return nil
}

func runVerify(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	logger.Info("Verifying downloaded components...")
	fmt.Println()

	verified := 0
	missing := 0

	// Check Android Studio
	studioDir := filepath.Join(cfg.DownloadDir, "android-studio")
	if util.FileExists(studioDir) && hasFiles(studioDir) {
		logger.Success("Android Studio: Found")
		verified++
	} else {
		logger.Error("Android Studio: Missing")
		missing++
	}

	// Check JDK
	if cfg.AndroidStudio.DownloadJDK {
		jdkDir := filepath.Join(cfg.DownloadDir, "jdk")
		if util.FileExists(jdkDir) && hasFiles(jdkDir) {
			logger.Success("JDK: Found")
			verified++
		} else {
			logger.Error("JDK: Missing")
			missing++
		}
	}

	// Check SDK
	sdkDir := filepath.Join(cfg.DownloadDir, "sdk")
	if util.FileExists(sdkDir) && hasFiles(sdkDir) {
		logger.Success("SDK Tools: Found")
		verified++
	} else {
		logger.Error("SDK Tools: Missing")
		missing++
	}

	// Check Gradle
	gradleDir := filepath.Join(cfg.DownloadDir, "gradle", "distributions")
	if util.FileExists(gradleDir) && hasFiles(gradleDir) {
		count := countFiles(gradleDir)
		logger.Success("Gradle: Found (%d files)", count)
		verified++
	} else {
		logger.Error("Gradle: Missing")
		missing++
	}

	// Check dependencies
	if cfg.Dependencies.DownloadMaven {
		depsDir := filepath.Join(cfg.DownloadDir, "dependencies")
		if util.FileExists(depsDir) && hasFiles(depsDir) {
			logger.Success("Dependencies: Found")
			verified++
		} else {
			logger.Warning("Dependencies: Missing (optional)")
		}
	}

	// Calculate total size
	totalSize, err := util.DirSize(cfg.DownloadDir)
	if err == nil {
		logger.Info("Total download size: %s", formatBytes(totalSize))
	}

	fmt.Println()
	if missing > 0 {
		logger.Warning("Verification: %d verified, %d missing", verified, missing)
		logger.Info("Run 'android-offline download' to download missing components")
		return fmt.Errorf("missing components")
	}

	logger.Success("Verification complete: All components present!")
	logger.Info("Ready to create package with 'android-offline package'")

	return nil
}

func hasFiles(dir string) bool {
	entries, err := os.ReadDir(dir)
	return err == nil && len(entries) > 0
}

func countFiles(dir string) int {
	entries, _ := os.ReadDir(dir)
	return len(entries)
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

func loadConfig() (*config.Config, error) {
	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		logger.Warning("Config file not found: %s", cfgFile)
		logger.Info("Run 'android-offline init' to create a default configuration")
		return nil, fmt.Errorf("config file not found")
	}

	cfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Override verbose from flag
	if verbose {
		cfg.Options.Verbose = true
	}

	return cfg, nil
}

func setupSignalHandler() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Println("\n\nReceived interrupt signal, cancelling...")
		cancel()
	}()

	return ctx, cancel
}
