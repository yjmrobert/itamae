package cmd

import (
	"context"
	"os"
	"syscall"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "itamae",
	Short: "Set up a developer's Linux workstation",
	Long: `Itamae installs and manages a curated set of development tools on
Debian and Ubuntu based workstations.

Tools are grouped into three categories:
  core        Essential tools, installed together
  essentials  Common developer extras, installed together
  unverified  Optional tools you pick individually

Run "itamae install" for the interactive installer, or pass --category and
--yes to install without any prompts (for example in a provisioning script).`,
}

// Execute runs the CLI. It is the entry point used by main.
func Execute() {
	if err := fang.Execute(
		context.Background(),
		rootCmd,
		fang.WithVersion(Version),
		fang.WithCommit(GitCommit),
		fang.WithNotifySignal(os.Interrupt, syscall.SIGTERM),
	); err != nil {
		os.Exit(1)
	}
}
