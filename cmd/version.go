package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version information - set via ldflags during build
var (
	Version   = "dev"     // Semantic version (e.g., "1.0.0")
	GitCommit = "unknown" // Git commit hash
	BuildDate = "unknown" // Build timestamp
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Display the version, git commit, and build date of itamae.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("itamae version %s\n", Version)
		fmt.Printf("  commit: %s\n", GitCommit)
		fmt.Printf("  built:  %s\n", BuildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
