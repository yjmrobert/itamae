package cmd

import (
	"github.com/yjmrobert/itamae/itamae"

	"github.com/spf13/cobra"
)

var customCmd = &cobra.Command{
	Use:   "custom",
	Short: "Install a custom set of software.",
	Run: func(cmd *cobra.Command, args []string) {
		plugins, cleanup, err := itamae.LoadPlugins()
		if err != nil {
			itamae.Logger.Errorf("Error loading plugins: %v\n", err)
			return
		}
		defer cleanup()
		itamae.RunCustom(plugins)
	},
}

func init() {
	rootCmd.AddCommand(customCmd)
}
