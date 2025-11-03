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
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your command '%s'", err)
		os.Exit(1)
	}
}
