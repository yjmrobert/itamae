package cmd

import (
	"fmt"
	"github.com/yjmrobert/itamae/itamae"

	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall all installed software.",
	Run: func(cmd *cobra.Command, args []string) {
		plugins, cleanup, err := itamae.LoadPlugins()
		if err != nil {
			fmt.Printf("Error loading plugins: %v\n", err)
			return
		}
		defer cleanup()
		itamae.RunUninstall(plugins)
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
