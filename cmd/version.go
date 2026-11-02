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
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		out := cmd.OutOrStdout()
		fmt.Fprintf(out, "itamae version %s\n", Version)
		fmt.Fprintf(out, "  commit: %s\n", GitCommit)
		fmt.Fprintf(out, "  built:  %s\n", BuildDate)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
