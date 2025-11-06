package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "itamae",
	Short: "Itamae is a tool to set up a developer's Linux workstation.",
	Long:  `A fast and flexible CLI tool to install and manage your development environment on a Linux workstation.`,
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommand is provided, show help
		cmd.Help()
	},
}

func init() {
	rootCmd.PersistentFlags().BoolP("version", "v", false, "Print version information")
}

func Execute() {
	// Check for version flag in os.Args before cobra processes it
	for _, arg := range os.Args[1:] {
		if arg == "--version" || arg == "-v" {
			fmt.Printf("itamae version %s\n", Version)
			fmt.Printf("  commit: %s\n", GitCommit)
			fmt.Printf("  built:  %s\n", BuildDate)
			os.Exit(0)
		}
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your command '%s'", err)
		os.Exit(1)
	}
}
