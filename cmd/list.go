package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/colorprofile"
	"github.com/spf13/cobra"
	"github.com/yjmrobert/itamae/itamae"
)

var listFlags struct {
	category string
	idsOnly  bool
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the available tools",
	Long: `List every tool itamae knows about, grouped by category, with its install
method and description. Use --ids to print only plugin IDs, one per line, which
is handy for scripting "itamae install --plugins".`,
	Example: `  itamae list
  itamae list --category essentials
  itamae list -c unverified --ids`,
	Args: cobra.NoArgs,
	RunE: runList,
}

func init() {
	f := listCmd.Flags()
	f.StringVarP(&listFlags.category, "category", "c", "", "Only list this category (default: all)")
	f.BoolVar(&listFlags.idsOnly, "ids", false, "Print plugin IDs only, one per line")
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, _ []string) error {
	plugins, cleanup, err := loadForCategory(listFlags.category)
	if err != nil {
		return err
	}
	defer cleanup()

	if listFlags.idsOnly {
		for _, p := range plugins {
			fmt.Fprintln(cmd.OutOrStdout(), p.ID)
		}
		return nil
	}

	out := colorprofile.NewWriter(cmd.OutOrStdout(), os.Environ())
	header := lipgloss.NewStyle().Bold(true).Foreground(itamae.TokyoNightMagenta)
	dim := lipgloss.NewStyle().Foreground(itamae.TokyoNightComment)

	current := ""
	var tw *tabwriter.Writer
	flush := func() {
		if tw != nil {
			_ = tw.Flush()
			fmt.Fprintln(out)
		}
	}
	for _, p := range plugins {
		if p.Category != current {
			flush()
			current = p.Category
			fmt.Fprintln(out, header.Render(fmt.Sprintf("%s (%d)", current, countCategory(plugins, current))))
			tw = tabwriter.NewWriter(out, 0, 4, 2, ' ', 0)
		}
		fmt.Fprintf(tw, "  %s\t%s\t%s\n", p.ID, dim.Render(p.InstallMethod), p.Description)
	}
	flush()
	return nil
}

func countCategory(plugins []itamae.ToolPlugin, category string) int {
	n := 0
	for _, p := range plugins {
		if p.Category == category {
			n++
		}
	}
	return n
}
