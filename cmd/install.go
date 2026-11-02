package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yjmrobert/itamae/itamae"
)

var installFlags struct {
	category   string
	plugins    []string
	all        bool
	yes        bool
	plain      bool
	inputs     []string
	skipVerify bool
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install a category of tools, or a hand-picked set",
	Long: `Install development tools.

Without flags the installer is fully interactive: pick a category, pick the
tools (for the unverified category), answer any questions the tools need and
confirm. Every prompt can be answered up front with flags, which makes the
command usable from scripts and CI:

  --category   which category to install
  --plugins    install only these plugin IDs (see "itamae list")
  --all        install every plugin of the category (required for unverified
               when not running interactively)
  --input      answer a plugin question, e.g. --input GIT_USER_NAME=Jane.
               Environment variables with the same name are used as well.
  --yes        skip the confirmation prompt
  --plain      line-oriented output instead of the TUI (automatic when the
               output is not a terminal)

The command exits with a non-zero status if any tool fails to install or does
not pass its post-install check.`,
	Example: `  itamae install
  itamae install --category core --yes
  itamae install --category unverified --plugins zsh,btop --yes --plain
  GIT_USER_NAME=Jane GIT_USER_EMAIL=jane@example.com itamae install -c core -y`,
	Args: cobra.NoArgs,
	RunE: runInstall,
}

func init() {
	f := installCmd.Flags()
	f.StringVarP(&installFlags.category, "category", "c", "", "Plugin category to install (core, essentials, unverified)")
	f.StringSliceVarP(&installFlags.plugins, "plugins", "p", nil, "Comma-separated plugin IDs to install (default: whole category)")
	f.BoolVar(&installFlags.all, "all", false, "Install every plugin of the category without prompting")
	f.BoolVarP(&installFlags.yes, "yes", "y", false, "Skip the confirmation prompt")
	f.BoolVar(&installFlags.plain, "plain", false, "Disable the TUI and print plain progress lines")
	f.StringArrayVarP(&installFlags.inputs, "input", "i", nil, "Answer a plugin question as NAME=VALUE (repeatable)")
	f.BoolVar(&installFlags.skipVerify, "skip-verify", false, "Do not run each tool's check() after installing it")
	rootCmd.AddCommand(installCmd)
}

func runInstall(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()
	out := cmd.OutOrStdout()
	interactive := itamae.IsInteractive() && !installFlags.plain

	category := installFlags.category
	if category == "" {
		if !interactive {
			return errors.New("a category is required when not running interactively (use --category)")
		}
		selected, err := itamae.SelectCategory()
		if err != nil {
			return fmt.Errorf("selecting category: %w", err)
		}
		category = selected
	}

	plugins, cleanup, err := itamae.LoadPlugins(category)
	if err != nil {
		return err
	}
	defer cleanup()

	selected, err := selectPlugins(plugins, category, interactive)
	if err != nil {
		return err
	}
	if len(selected) == 0 {
		fmt.Fprintln(out, "No plugins selected. Nothing to do.")
		return nil
	}

	inputs, err := resolveInputs(selected, installFlags.inputs, interactive)
	if err != nil {
		return err
	}

	if !installFlags.yes {
		if !interactive {
			return errors.New("confirmation is required: pass --yes when not running interactively")
		}
		ok, err := itamae.ConfirmInstallation(len(selected))
		if err != nil {
			return err
		}
		if !ok {
			fmt.Fprintln(out, "Installation cancelled.")
			return nil
		}
	}

	if err := itamae.EnsureSudoAccess(); err != nil {
		return err
	}

	if err := itamae.InitDebugLog(); err != nil {
		itamae.Logger.Warnf("Could not initialise the debug log: %v", err)
	}
	defer itamae.CloseDebugLog(out)

	var result itamae.InstallResult
	if interactive {
		result, err = itamae.RunInstallTUI(ctx, selected, inputs)
		if err != nil {
			return err
		}
	} else {
		fmt.Fprintf(out, "Installing %d %s tool(s)\n", len(selected), category)
		result = itamae.RunInstallPlain(ctx, selected, inputs, out)
	}

	itamae.PrintSummary(out, result)
	if !result.OK() {
		return fmt.Errorf("%d of %d tool(s) failed to install", len(result.Failed), len(selected))
	}
	return nil
}

// selectPlugins decides which plugins of the category to install.
func selectPlugins(plugins []itamae.ToolPlugin, category string, interactive bool) ([]itamae.ToolPlugin, error) {
	if len(installFlags.plugins) > 0 {
		return itamae.FilterPlugins(plugins, installFlags.plugins)
	}
	if category != "unverified" || installFlags.all {
		return plugins, nil
	}
	if !interactive {
		return nil, errors.New("the unverified category needs --all or --plugins when not running interactively")
	}
	return itamae.SelectPlugins(plugins)
}

// resolveInputs collects every value the selected plugins need, from --input flags,
// the environment, or (interactively) a prompt.
func resolveInputs(plugins []itamae.ToolPlugin, flagInputs []string, interactive bool) (map[string]string, error) {
	values := make(map[string]string)
	for _, kv := range flagInputs {
		name, value, ok := strings.Cut(kv, "=")
		if !ok || strings.TrimSpace(name) == "" {
			return nil, fmt.Errorf("invalid --input %q, expected NAME=VALUE", kv)
		}
		values[strings.TrimSpace(name)] = value
	}

	var missing []string
	for _, in := range itamae.RequiredInputs(plugins) {
		if _, ok := values[in.Name]; ok {
			continue
		}
		if v, ok := os.LookupEnv(in.Name); ok && v != "" {
			values[in.Name] = v
			continue
		}
		if !interactive {
			missing = append(missing, in.Name)
			continue
		}
		value, err := itamae.PromptInput(in.Prompt, itamae.DefaultInputValue(in.DefaultCmd))
		if err != nil {
			return nil, fmt.Errorf("reading input %s: %w", in.Name, err)
		}
		values[in.Name] = value
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required input(s): %s (pass --input NAME=VALUE or set them in the environment)", strings.Join(missing, ", "))
	}
	return values, nil
}
