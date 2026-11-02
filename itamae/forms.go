package itamae

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"charm.land/huh/v2"
)

// IsTerminal reports whether stdout is attached to a terminal.
func IsTerminal() bool {
	info, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

// SelectCategory asks the user which plugin category to work with.
func SelectCategory() (string, error) {
	var category string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select package category").
				Description("Choose which category of packages to install").
				Options(
					huh.NewOption("Core - Install all essential packages", "core"),
					huh.NewOption("Essentials - Install common developer extras", "essentials"),
					huh.NewOption("Unverified - Select individual packages", "unverified"),
				).
				Value(&category),
		),
	)

	if err := form.Run(); err != nil {
		return "", err
	}
	return category, nil
}

// SelectPlugins shows a multi-select of the given plugins and returns the chosen ones.
func SelectPlugins(plugins []ToolPlugin) ([]ToolPlugin, error) {
	if len(plugins) == 0 {
		return nil, nil
	}

	options := make([]huh.Option[string], 0, len(plugins))
	for _, p := range plugins {
		options = append(options, huh.NewOption(fmt.Sprintf("%s - %s", p.Name, p.Description), p.ID))
	}

	var selectedIDs []string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Available Tools").
				Description("Use space to select, enter to confirm").
				Options(options...).
				Value(&selectedIDs).
				Height(30),
		),
	)

	if err := form.Run(); err != nil {
		return nil, err
	}
	if len(selectedIDs) == 0 {
		return nil, nil
	}
	return FilterPlugins(plugins, selectedIDs)
}

// ConfirmInstallation asks the user to confirm before anything is installed.
func ConfirmInstallation(count int) (bool, error) {
	var confirm bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Proceed with installation?").
				Description(fmt.Sprintf("This will install %d tool(s) on your system.", count)).
				Value(&confirm),
		),
	)

	if err := form.Run(); err != nil {
		return false, err
	}
	return confirm, nil
}

// PromptInput asks the user for a single text value, pre-filled with defaultValue.
func PromptInput(question string, defaultValue string) (string, error) {
	value := defaultValue

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title(question).
				Value(&value).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("input cannot be empty")
					}
					return nil
				}),
		),
	)

	if err := form.Run(); err != nil {
		return "", err
	}
	return value, nil
}

// DefaultInputValue runs an input's DefaultCmd and returns its trimmed output.
// Failures yield an empty string so the user simply gets an empty prompt.
func DefaultInputValue(command string) string {
	if command == "" {
		return ""
	}
	output, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// IsInteractive reports whether both stdin and stdout are terminals, which the
// interactive forms and the TUI require.
func IsInteractive() bool {
	if !IsTerminal() {
		return false
	}
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}
