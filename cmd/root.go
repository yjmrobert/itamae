package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/yjmrobert/itamae/itamae"
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
			itamae.Logger.Infof("itamae version %s", Version)
			itamae.Logger.Infof("  commit: %s", GitCommit)
			itamae.Logger.Infof("  built:  %s", BuildDate)
			os.Exit(0)
		}
	}

	if err := rootCmd.Execute(); err != nil {
		itamae.Logger.Errorf("Whoops. There was an error while executing your command '%s'", err)
		os.Exit(1)
	}
}
