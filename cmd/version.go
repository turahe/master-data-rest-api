package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long: `Display detailed version information about the Master Data REST API including:
- Application version
- Go version used to build the binary
- Build information
- Runtime information`,
	Run: func(cmd *cobra.Command, args []string) {
		showVersion()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func showVersion() {
	config := GetConfig()

	fmt.Println("Master Data REST API")
	fmt.Println("====================")
	fmt.Printf("Version:     %s\n", config.App.Version)
	fmt.Printf("Environment: %s\n", config.App.Env)
	fmt.Printf("Go Version:  %s\n", runtime.Version())
	fmt.Printf("OS/Arch:     %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("Build Date:  %s\n", "2025-08-09") // This could be set at build time
	fmt.Println()
	fmt.Println("Features:")
	fmt.Println("- ✅ Geographical data management with nested set model")
	fmt.Println("- ✅ Banks, currencies, and languages management")
	fmt.Println("- ✅ API key authentication and management")
	fmt.Println("- ✅ PostgreSQL with pgx driver")
	fmt.Println("- ✅ Comprehensive logging and monitoring")
	fmt.Println("- ✅ Auto-generated Swagger documentation")
	fmt.Println("- ✅ Fiber v2 web framework")
	fmt.Println("- ✅ Cobra CLI framework")
}
