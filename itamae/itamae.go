package itamae

import (
	"bufio"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/charmbracelet/huh"
)

//go:embed scripts/unverified/*
var scriptsFS embed.FS

type Input struct {
	Name   string
	Prompt string
}

type ToolPlugin struct {
	ID             string // "vscode", "ripgrep"
	Name           string // "Visual Studio Code"
	Description    string
	Omakase        bool
	ScriptPath     string // The path to the executable in the /tmp/ directory
	InstallMethod  string // "apt", "binary", "manual"
	PackageName    string // For apt packages, the actual package name
	PostInstall    string // Function name for post-install tasks (optional)
	RequiredInputs []Input
}

func RunInstall(plugins []ToolPlugin) {
	Logger.Info("Starting Itamae setup...")

	// Separate Omakase and Optional plugins
	omakasePlugins := []ToolPlugin{}
	optionalPlugins := []ToolPlugin{}

	for _, p := range plugins {
		if p.Omakase {
			omakasePlugins = append(omakasePlugins, p)
		} else {
			optionalPlugins = append(optionalPlugins, p)
		}
	}

	// Get user selections for optional plugins
	selectedOptional := selectOptionalPlugins(optionalPlugins)

	// Combine Omakase + User-selected
	allSelectedPlugins := append(omakasePlugins, selectedOptional...)

	// Group by installation method
	aptPlugins := []ToolPlugin{}
	binaryPlugins := []ToolPlugin{}
	manualPlugins := []ToolPlugin{}

	for _, p := range allSelectedPlugins {
		switch p.InstallMethod {
		case "apt":
			aptPlugins = append(aptPlugins, p)
		case "binary":
			binaryPlugins = append(binaryPlugins, p)
		case "manual":
			manualPlugins = append(manualPlugins, p)
		}
	}

	// Display installation plan
	displayInstallationPlan(aptPlugins, binaryPlugins, manualPlugins)

	// Confirm before proceeding
	if !confirmInstallation() {
		Logger.Info("Installation cancelled.")
		return
	}

	// Track success/failure
	successful := []string{}
	failed := []string{}

	// Phase 1: Batch install all APT packages
	if len(aptPlugins) > 0 {
		Logger.Info("\n" + strings.Repeat("=", 60))
		Logger.Info("=== Phase 1: Installing APT packages ===")
		Logger.Info(strings.Repeat("=", 60))
		if err := batchInstallApt(aptPlugins, nil); err != nil {
			Logger.Errorf("❌ Error in batch APT installation: %v\n", err)
			for _, p := range aptPlugins {
				failed = append(failed, p.Name)
			}
		} else {
			for _, p := range aptPlugins {
				successful = append(successful, p.Name)
			}
		}
	}

	// Configure Git
	if err := configureGit(); err != nil {
		Logger.Warnf("⚠️  Error configuring Git: %v\n", err)
	}

	// Phase 2: Install binary plugins individually
	if len(binaryPlugins) > 0 {
		Logger.Info("\n" + strings.Repeat("=", 60))
		Logger.Info("=== Phase 2: Installing binary packages ===")
		Logger.Info(strings.Repeat("=", 60))
		for _, p := range binaryPlugins {
			Logger.Infof("\n--- Installing %s ---\n", p.Name)
			if err := executeScript(p, "install", nil); err != nil {
				Logger.Errorf("❌ Error installing %s: %v\n", p.Name, err)
				failed = append(failed, p.Name)
			} else {
				successful = append(successful, p.Name)
			}
		}
	}

	// Phase 3: Install manual plugins individually
	if len(manualPlugins) > 0 {
		Logger.Info("\n" + strings.Repeat("=", 60))
		Logger.Info("=== Phase 3: Manual installation required ===")
		Logger.Info(strings.Repeat("=", 60))
		for _, p := range manualPlugins {
			Logger.Infof("\n--- %s ---\n", p.Name)
			if err := executeScript(p, "install", nil); err != nil {
				Logger.Errorf("❌ Error installing %s: %v\n", p.Name, err)
				failed = append(failed, p.Name)
			} else {
				successful = append(successful, p.Name)
			}
		}
	}

	// Display summary
	displayInstallationSummary(successful, failed)

	Logger.Info("\n✅ Itamae setup complete!")
}

func selectOptionalPlugins(plugins []ToolPlugin) []ToolPlugin {
	if len(plugins) == 0 {
		return []ToolPlugin{}
	}

	Logger.Info("\n🍱 Core tools will be installed automatically (Omakase).")
	Logger.Info("📦 Select additional tools you'd like to install:")

	// Create options for multi-select
	options := []huh.Option[string]{}
	for _, p := range plugins {
		// Format: "Tool Name - Description"
		label := fmt.Sprintf("%s - %s", p.Name, p.Description)
		options = append(options, huh.NewOption(label, p.ID))
	}

	var selectedIDs []string

	// Create multi-select form
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Additional Tools").
				Description("Use space to select, enter to confirm").
				Options(options...).
				Value(&selectedIDs).
				Height(15),
		),
	)

	// Run the form
	err := form.Run()
	if err != nil {
		Logger.Errorf("Error: %v\n", err)
		return []ToolPlugin{}
	}

	// Convert selected IDs back to plugins
	selectedPlugins := []ToolPlugin{}
	selectedMap := make(map[string]bool)
	for _, id := range selectedIDs {
		selectedMap[id] = true
	}

	for _, p := range plugins {
		if selectedMap[p.ID] {
			selectedPlugins = append(selectedPlugins, p)
		}
	}

	return selectedPlugins
}

func displayInstallationPlan(aptPlugins, binaryPlugins, manualPlugins []ToolPlugin) {
	Logger.Info("\n" + strings.Repeat("=", 60))
	Logger.Info("📋 INSTALLATION PLAN")
	Logger.Info(strings.Repeat("=", 60))

	if len(aptPlugins) > 0 {
		Logger.Info("\n📦 APT Packages (batch installation):")
		for _, p := range aptPlugins {
			marker := "🍱"
			if !p.Omakase {
				marker = "📌"
			}
			Logger.Infof("  %s %s (%s)\n", marker, p.Name, p.PackageName)
		}
	}

	if len(binaryPlugins) > 0 {
		Logger.Info("\n🔧 Binary Installations (individual):")
		for _, p := range binaryPlugins {
			marker := "🍱"
			if !p.Omakase {
				marker = "📌"
			}
			Logger.Infof("  %s %s\n", marker, p.Name)
		}
	}

	if len(manualPlugins) > 0 {
		Logger.Info("\n⚠️  Manual Installations (requires attention):")
		for _, p := range manualPlugins {
			Logger.Infof("  ⚙️  %s\n", p.Name)
		}
	}

	Logger.Info("\n" + strings.Repeat("=", 60))
	Logger.Infof("Total: %d tools\n", len(aptPlugins)+len(binaryPlugins)+len(manualPlugins))
	Logger.Info(strings.Repeat("=", 60))
}

func confirmInstallation() bool {
	var confirm bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Proceed with installation?").
				Description("This will install the selected tools on your system.").
				Value(&confirm),
		),
	)

	err := form.Run()
	if err != nil {
		return false
	}

	return confirm
}

func displayInstallationSummary(successful, failed []string) {
	Logger.Info("\n" + strings.Repeat("=", 60))
	Logger.Info("📊 INSTALLATION SUMMARY")
	Logger.Info(strings.Repeat("=", 60))

	if len(successful) > 0 {
		Logger.Info("\n✅ Successfully installed:")
		for _, name := range successful {
			Logger.Infof("  • %s\n", name)
		}
	}

	if len(failed) > 0 {
		Logger.Info("\n❌ Failed to install:")
		for _, name := range failed {
			Logger.Infof("  • %s\n", name)
		}
	}

	Logger.Info("\n" + strings.Repeat("=", 60))
}

func countPostInstalls(plugins []ToolPlugin) int {
	count := 0
	for _, p := range plugins {
		if p.PostInstall != "" {
			count++
		}
	}
	return count
}

type formRunner interface {
	Run() error
}

type realFormRunner struct {
	form *huh.Form
}

func (r *realFormRunner) Run() error {
	return r.form.Run()
}

var newFormRunner = func(form *huh.Form) formRunner {
	return &realFormRunner{form: form}
}

func RunTextInput(question string) (string, error) {
	var value string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title(question).
				Value(&value).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("input cannot be empty")
					}
					return nil
				}),
		),
	)

	runner := newFormRunner(form)
	err := runner.Run()
	if err != nil {
		return "", err
	}

	return value, nil
}

func RunUninstall(plugins []ToolPlugin) {
	Logger.Info("Checking for installed plugins...")
	for _, p := range plugins {
		if err := checkScript(p); err == nil {
			Logger.Infof("--- Uninstalling %s ---\n", p.Name)
			if err := executeScript(p, "remove", nil); err != nil {
				Logger.Errorf("Error uninstalling %s: %v\n", p.Name, err)
			}
		}
	}
	Logger.Info("Uninstallation complete.")
}

func checkScript(plugin ToolPlugin) error {
	cmd := exec.Command("bash", plugin.ScriptPath, "check")
	return cmd.Run()
}

func executeScript(plugin ToolPlugin, command string, env map[string]string) error {
	cmd := exec.Command("bash", plugin.ScriptPath, command)

	// Set environment variables
	for key, value := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		io.Copy(os.Stdout, stdout)
	}()

	go func() {
		defer wg.Done()
		io.Copy(os.Stderr, stderr)
	}()

	wg.Wait()

	return cmd.Wait()
}

func LoadPlugins() ([]ToolPlugin, func(), error) {
	tmpDir, err := os.MkdirTemp("", "itamae-scripts-")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create temp dir: %w", err)
	}

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	var plugins []ToolPlugin

	files, err := scriptsFS.ReadDir("scripts/unverified")
	if err != nil {
		return nil, cleanup, fmt.Errorf("failed to read embedded scripts dir: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		plugin, err := processPluginFile(file, tmpDir)
		if err != nil {
			return nil, cleanup, fmt.Errorf("failed to process plugin %s: %w", file.Name(), err)
		}
		plugins = append(plugins, plugin)
	}

	return plugins, cleanup, nil
}

func processPluginFile(file fs.DirEntry, tmpDir string) (ToolPlugin, error) {
	fileName := file.Name()
	scriptPath := filepath.Join("scripts/unverified", fileName)

	// Read content for parsing
	content, err := scriptsFS.ReadFile(scriptPath)
	if err != nil {
		return ToolPlugin{}, fmt.Errorf("failed to read embedded file %s: %w", scriptPath, err)
	}

	// Parse metadata
	plugin, err := parseMetadata(string(content))
	if err != nil {
		return ToolPlugin{}, fmt.Errorf("failed to parse metadata for %s: %w", fileName, err)
	}
	plugin.ID = strings.TrimSuffix(fileName, ".sh")

	// Unpack script to temp directory
	destPath := filepath.Join(tmpDir, fileName)
	if err := os.WriteFile(destPath, content, 0755); err != nil {
		return ToolPlugin{}, fmt.Errorf("failed to write script to temp dir: %w", err)
	}
	plugin.ScriptPath = destPath

	return plugin, nil
}

func parseMetadata(content string) (ToolPlugin, error) {
	plugin := ToolPlugin{}
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "#") {
			continue // Stop parsing metadata once we hit code
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(strings.TrimPrefix(parts[0], "#"))
		value := strings.TrimSpace(parts[1])

		switch key {
		case "NAME":
			plugin.Name = value
		case "OMAKASE":
			plugin.Omakase = (value == "true")
		case "DESCRIPTION":
			plugin.Description = value
		case "INSTALL_METHOD":
			plugin.InstallMethod = value
		case "PACKAGE_NAME":
			plugin.PackageName = value
		case "POST_INSTALL":
			plugin.PostInstall = value
		case "REQUIRES":
			parts := strings.SplitN(value, "|", 2)
			if len(parts) == 2 {
				plugin.RequiredInputs = append(plugin.RequiredInputs, Input{
					Name:   parts[0],
					Prompt: parts[1],
				})
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return ToolPlugin{}, fmt.Errorf("scanner error while parsing metadata: %w", err)
	}
	return plugin, nil
}

func configureGit() error {
	Logger.Info("\n--- Configuring Git ---")
	name, err := RunTextInput("Enter your Git user name")
	if err != nil {
		return err
	}
	email, err := RunTextInput("Enter your Git user email")
	if err != nil {
		return err
	}

	if err := exec.Command("git", "config", "--global", "user.name", name).Run(); err != nil {
		return fmt.Errorf("failed to set git user.name: %w", err)
	}
	if err := exec.Command("git", "config", "--global", "user.email", email).Run(); err != nil {
		return fmt.Errorf("failed to set git user.email: %w", err)
	}

	Logger.Info("✅ Git configured.")
	return nil
}

// batchInstallApt installs multiple APT packages in a single command using nala or apt.
// After installation, it runs any post-install tasks defined for each plugin.
func RunCustom(plugins []ToolPlugin) {
	Logger.Info("Starting Itamae custom setup...")
	selectedPlugins := selectCustomPlugins(plugins)
	if len(selectedPlugins) == 0 {
		Logger.Info("No plugins selected. Exiting.")
		return
	}

	// Gather all required inputs upfront
	requiredInputs := make(map[string]string)
	for _, p := range selectedPlugins {
		for _, input := range p.RequiredInputs {
			if _, ok := requiredInputs[input.Name]; !ok {
				value, err := RunTextInput(input.Prompt)
				if err != nil {
					Logger.Errorf("Error getting input for %s: %v", input.Name, err)
					return
				}
				requiredInputs[input.Name] = value
			}
		}
	}

	// Confirm before proceeding
	if !confirmInstallation() {
		Logger.Info("Installation cancelled.")
		return
	}

	processCustomInstall(selectedPlugins, requiredInputs)
}

func processCustomInstall(selectedPlugins []ToolPlugin, requiredInputs map[string]string) {
	// Separate plugins by install method
	aptPlugins := []ToolPlugin{}
	otherPlugins := []ToolPlugin{}
	for _, p := range selectedPlugins {
		if p.InstallMethod == "apt" {
			aptPlugins = append(aptPlugins, p)
		} else {
			otherPlugins = append(otherPlugins, p)
		}
	}

	// Track success/failure
	successful := []string{}
	failed := []string{}

	// Phase 1: Batch install all APT packages
	if len(aptPlugins) > 0 {
		if err := batchInstallApt(aptPlugins, requiredInputs); err != nil {
			Logger.Errorf("❌ Error in batch APT installation: %v\n", err)
			for _, p := range aptPlugins {
				failed = append(failed, p.Name)
			}
		} else {
			for _, p := range aptPlugins {
				successful = append(successful, p.Name)
			}
		}
	}

	// Phase 2: Install other plugins individually
	for _, p := range otherPlugins {
		Logger.Infof("\n--- Installing %s ---\n", p.Name)
		if err := executeScript(p, "install", requiredInputs); err != nil {
			Logger.Errorf("❌ Error installing %s: %v\n", p.Name, err)
			failed = append(failed, p.Name)
		} else {
			successful = append(successful, p.Name)
		}
	}

	// Display summary
	displayInstallationSummary(successful, failed)

	Logger.Info("\n✅ Itamae custom setup complete!")
}

func selectCustomPlugins(plugins []ToolPlugin) []ToolPlugin {
	if len(plugins) == 0 {
		return []ToolPlugin{}
	}

	Logger.Info("📦 Select the tools you'd like to install:")

	options := []huh.Option[string]{}
	for _, p := range plugins {
		label := fmt.Sprintf("%s - %s", p.Name, p.Description)
		options = append(options, huh.NewOption(label, p.ID))
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

	runner := newFormRunner(form)
	err := runner.Run()
	if err != nil {
		Logger.Errorf("Error: %v\n", err)
		return []ToolPlugin{}
	}

	selectedPlugins := []ToolPlugin{}
	selectedMap := make(map[string]bool)
	for _, id := range selectedIDs {
		selectedMap[id] = true
	}

	for _, p := range plugins {
		if selectedMap[p.ID] {
			selectedPlugins = append(selectedPlugins, p)
		}
	}

	return selectedPlugins
}

func batchInstallApt(plugins []ToolPlugin, env map[string]string) error {
	if len(plugins) == 0 {
		return nil
	}

	Logger.Infof("\n⏳ Installing %d APT packages...\n\n", len(plugins))

	// Check if nala is available
	_, err := exec.LookPath("nala")
	useNala := err == nil

	// Collect package names
	packages := []string{}
	for _, p := range plugins {
		if p.PackageName != "" {
			packages = append(packages, p.PackageName)
			Logger.Infof("  • %s (%s)\n", p.Name, p.PackageName)
		}
	}

	if len(packages) == 0 {
		Logger.Info("No APT packages to install.")
		return nil
	}

	// Build install command
	var cmd *exec.Cmd
	if useNala {
		args := append([]string{"nala", "install", "-y"}, packages...)
		cmd = exec.Command("sudo", args...)
		Logger.Infof("\n▶️  Running: sudo nala install -y %s\n\n", strings.Join(packages, " "))
	} else {
		args := append([]string{"apt", "install", "-y"}, packages...)
		cmd = exec.Command("sudo", args...)
		Logger.Infof("\n▶️  Running: sudo apt install -y %s\n\n", strings.Join(packages, " "))
	}

	// Execute with live output
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("batch APT installation failed: %w", err)
	}

	Logger.Infof("\n✅ Successfully installed %d APT packages\n", len(plugins))

	// Run post-install tasks with progress
	if hasPostInstall := countPostInstalls(plugins); hasPostInstall > 0 {
		Logger.Info("\n⚙️  Running post-installation tasks...\n")
		for _, p := range plugins {
			if p.PostInstall != "" {
				Logger.Infof("  • %s... ", p.Name)
				if err := executeScript(p, "post_install", env); err != nil {
					Logger.Error("❌ failed\n")
				} else {
					Logger.Info("✅\n")
				}
			}
		}
	}

	return nil
}
