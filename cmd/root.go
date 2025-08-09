package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/turahe/master-data-rest-api/configs"
	"github.com/turahe/master-data-rest-api/pkg/logger"
)

var (
	cfgFile string
	config  *configs.Config
	log     *logger.Logger
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "master-data-api",
	Short: "Master Data REST API - A comprehensive API for managing master data",
	Long: `Master Data REST API is a comprehensive REST API for managing master data including:
- Geographical data (countries, provinces, cities, districts, villages) using nested set model
- Banking information
- Currency data with status management
- Language information with localization support
- API key management with authentication

The API provides full CRUD operations, search capabilities, and robust authentication.`,
	Version: "1.0.0",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .env)")
	rootCmd.PersistentFlags().String("log-level", "info", "log level (trace, debug, info, warn, error)")
	rootCmd.PersistentFlags().String("log-format", "text", "log format (text, json)")
	rootCmd.PersistentFlags().String("log-output", "stdout", "log output (stdout, stderr, or file path)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Load configuration
	config = configs.Load()

	// Override config with flags if provided
	if level, _ := rootCmd.PersistentFlags().GetString("log-level"); level != "" {
		config.Logging.Level = level
	}
	if format, _ := rootCmd.PersistentFlags().GetString("log-format"); format != "" {
		config.Logging.Format = format
	}
	if output, _ := rootCmd.PersistentFlags().GetString("log-output"); output != "" {
		config.Logging.Output = output
	}

	// Initialize logger
	log = logger.New(logger.Config{
		Level:  config.Logging.Level,
		Format: config.Logging.Format,
		Output: config.Logging.Output,
	})
}

// GetConfig returns the loaded configuration
func GetConfig() *configs.Config {
	return config
}

// GetLogger returns the initialized logger
func GetLogger() *logger.Logger {
	return log
}

// CheckConfig ensures configuration is loaded
func CheckConfig(cmd *cobra.Command) error {
	if config == nil {
		return fmt.Errorf("configuration not loaded")
	}
	if log == nil {
		return fmt.Errorf("logger not initialized")
	}
	return nil
}
