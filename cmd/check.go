package cmd

import (
	"fmt"
	"os"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/colorprofile"
	"github.com/spf13/cobra"
	"github.com/yjmrobert/itamae/itamae"
)

var checkFlags struct {
	category string
	plugins  []string
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Report which tools are installed",
	Long: `Run every plugin's check() function and report whether the tool is present.

The command exits with a non-zero status when at least one of the selected tools
is missing, so it can be used to verify a provisioning run.`,
	Example: `  itamae check
  itamae check --category core
  itamae check --plugins git,ripgrep`,
	Args: cobra.NoArgs,
	RunE: runCheck,
}

func init() {
	f := checkCmd.Flags()
	f.StringVarP(&checkFlags.category, "category", "c", "", "Only check this category (default: all)")
	f.StringSliceVarP(&checkFlags.plugins, "plugins", "p", nil, "Comma-separated plugin IDs to check")
	rootCmd.AddCommand(checkCmd)
}

func runCheck(cmd *cobra.Command, _ []string) error {
	plugins, cleanup, err := loadForCategory(checkFlags.category)
	if err != nil {
		return err
	}
	defer cleanup()

	if len(checkFlags.plugins) > 0 {
		plugins, err = itamae.FilterPlugins(plugins, checkFlags.plugins)
		if err != nil {
			return err
		}
	}

	out := colorprofile.NewWriter(cmd.OutOrStdout(), os.Environ())
	ok := lipgloss.NewStyle().Foreground(itamae.StatusSuccess)
	bad := lipgloss.NewStyle().Foreground(itamae.StatusError)
	dim := lipgloss.NewStyle().Foreground(itamae.TokyoNightComment)

	missing := 0
	for _, p := range plugins {
		if itamae.CheckInstalled(cmd.Context(), p) {
			fmt.Fprintf(out, "%s %-20s %s\n", ok.Render("✓"), p.ID, dim.Render(p.Name))
		} else {
			missing++
			fmt.Fprintf(out, "%s %-20s %s\n", bad.Render("✗"), p.ID, dim.Render(p.Name+" (not installed)"))
		}
	}

	fmt.Fprintf(out, "\n%d of %d tool(s) installed\n", len(plugins)-missing, len(plugins))
	if missing > 0 {
		return fmt.Errorf("%d tool(s) not installed", missing)
	}
	return nil
}

// loadForCategory loads one category, or every category when name is empty.
func loadForCategory(name string) ([]itamae.ToolPlugin, func(), error) {
	if name == "" {
		return itamae.LoadAllPlugins()
	}
	return itamae.LoadPlugins(name)
}
