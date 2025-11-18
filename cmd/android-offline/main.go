package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Frosty1442/android-studio-offline/internal/config"
	"github.com/Frosty1442/android-studio-offline/internal/download"
	"github.com/Frosty1442/android-studio-offline/internal/install"
	"github.com/Frosty1442/android-studio-offline/internal/ui"
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

	logger.Info("Starting installation...")

	installer := install.NewInstaller(logger)
	if err := installer.Install(cfg); err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}

	logger.Success("Installation completed successfully!")
	logger.Info("Android Studio installed to: %s", cfg.InstallDir)
	logger.Info("Start Android Studio: %s/android-studio/bin/studio.sh", cfg.InstallDir)

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
	_, err := loadConfig()
	if err != nil {
		return err
	}

	logger.Info("Verifying downloaded components...")

	// TODO: Implement verification logic
	logger.Info("Verification not yet implemented")

	return nil
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
