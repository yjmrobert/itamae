package cmd

import (
	"github.com/yjmrobert/itamae/itamae"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install a custom set of software.",
	Run: func(cmd *cobra.Command, args []string) {
		// First, prompt user to select category
		category, err := itamae.SelectCategory()
		if err != nil {
			itamae.Logger.Errorf("Error selecting category: %v\n", err)
			return
		}

		plugins, cleanup, err := itamae.LoadPlugins(category)
		if err != nil {
			itamae.Logger.Errorf("Error loading plugins: %v\n", err)
			return
		}
		defer cleanup()
		itamae.RunInstall(plugins, category)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
