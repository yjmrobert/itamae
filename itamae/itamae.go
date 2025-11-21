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

//go:embed scripts/core/* scripts/unverified/*
var scriptsFS embed.FS

type Input struct {
	Name       string
	Prompt     string
	DefaultCmd string // Command to retrieve existing/default value
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
	fmt.Println("\n" + strings.Repeat("═", 60))
	fmt.Println("📊 INSTALLATION SUMMARY")
	fmt.Println(strings.Repeat("═", 60))

	if len(successful) > 0 {
		fmt.Println("\n✅ Successfully installed:")
		for _, name := range successful {
			fmt.Printf("   • %s\n", name)
		}
	}

	if len(failed) > 0 {
		fmt.Println("\n❌ Failed to install:")
		for _, name := range failed {
			fmt.Printf("   • %s\n", name)
		}
	}

	fmt.Println("\n" + strings.Repeat("═", 60))
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

// ensureSudoAccess prompts for sudo password upfront to avoid interruptions during installation.
func ensureSudoAccess() error {
	fmt.Println("\n🔐 Requesting sudo access for installation...")
	cmd := exec.Command("sudo", "-v")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// getDefaultValue executes a command to retrieve an existing/default value.
// Returns empty string on error (silent failure for better UX).
func getDefaultValue(command string) string {
	if command == "" {
		return ""
	}

	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(output))
}

func RunTextInput(question string, defaultValue string) (string, error) {
	value := defaultValue // Pre-populate with default value

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

func SelectCategory() (string, error) {
	var category string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select package category").
				Description("Choose which category of packages to install").
				Options(
					huh.NewOption("Core - Install all essential packages", "core"),
					huh.NewOption("Unverified - Select individual packages", "unverified"),
				).
				Value(&category),
		),
	)

	runner := newFormRunner(form)
	err := runner.Run()
	if err != nil {
		return "", err
	}

	return category, nil
}

func LoadPlugins(category string) ([]ToolPlugin, func(), error) {
	tmpDir, err := os.MkdirTemp("", "itamae-scripts-")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create temp dir: %w", err)
	}

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	var plugins []ToolPlugin

	scriptDir := fmt.Sprintf("scripts/%s", category)
	files, err := scriptsFS.ReadDir(scriptDir)
	if err != nil {
		return nil, cleanup, fmt.Errorf("failed to read embedded scripts dir: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		plugin, err := processPluginFile(file, tmpDir, category)
		if err != nil {
			return nil, cleanup, fmt.Errorf("failed to process plugin %s: %w", file.Name(), err)
		}
		plugins = append(plugins, plugin)
	}

	return plugins, cleanup, nil
}

func processPluginFile(file fs.DirEntry, tmpDir string, category string) (ToolPlugin, error) {
	fileName := file.Name()
	scriptPath := filepath.Join(fmt.Sprintf("scripts/%s", category), fileName)

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
			parts := strings.SplitN(value, "|", 3)
			if len(parts) >= 2 {
				input := Input{
					Name:   parts[0],
					Prompt: parts[1],
				}
				if len(parts) == 3 {
					input.DefaultCmd = parts[2]
				}
				plugin.RequiredInputs = append(plugin.RequiredInputs, input)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return ToolPlugin{}, fmt.Errorf("scanner error while parsing metadata: %w", err)
	}
	return plugin, nil
}

// batchInstallApt installs multiple APT packages in a single command using nala or apt.
// After installation, it runs any post-install tasks defined for each plugin.
func RunInstall(plugins []ToolPlugin, category string) {
	fmt.Println("\n🚀 Starting Itamae setup...")

	// Request sudo access upfront to avoid interruptions during installation
	if err := ensureSudoAccess(); err != nil {
		fmt.Println("\n❌ Failed to obtain sudo access. Installation cancelled.")
		return
	}

	var selectedPlugins []ToolPlugin
	if category == "core" {
		// For core, install everything without prompting
		fmt.Printf("Installing %d core packages\n", len(plugins))
		selectedPlugins = plugins
	} else {
		// For unverified, show multiselect
		selectedPlugins = selectPlugins(plugins)
		if len(selectedPlugins) == 0 {
			fmt.Println("No plugins selected. Exiting.")
			return
		}
	}

	// Gather all required inputs upfront
	requiredInputs := make(map[string]string)
	for _, p := range selectedPlugins {
		for _, input := range p.RequiredInputs {
			if _, ok := requiredInputs[input.Name]; !ok {
				defaultValue := getDefaultValue(input.DefaultCmd)
				value, err := RunTextInput(input.Prompt, defaultValue)
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
		fmt.Println("\nInstallation cancelled.")
		return
	}

	processInstall(selectedPlugins, requiredInputs)
}

func processInstall(selectedPlugins []ToolPlugin, requiredInputs map[string]string) {
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
		fmt.Printf("\n▶️  Installing %s...\n", p.Name)
		if err := executeScript(p, "install", requiredInputs); err != nil {
			Logger.Errorf("❌ Error installing %s: %v\n", p.Name, err)
			failed = append(failed, p.Name)
		} else {
			successful = append(successful, p.Name)
		}
	}

	// Display summary
	displayInstallationSummary(successful, failed)

	fmt.Println("\n✅ Itamae setup complete!")
}

func selectPlugins(plugins []ToolPlugin) []ToolPlugin {
	if len(plugins) == 0 {
		return []ToolPlugin{}
	}

	fmt.Println("\n📦 Select the tools you'd like to install:")

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
		fmt.Printf("Error: %v\n", err)
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

	fmt.Printf("\n📦 Installing %d APT package(s)\n", len(plugins))

	// Check if nala is available
	_, err := exec.LookPath("nala")
	useNala := err == nil

	// Collect package names
	packages := []string{}
	for _, p := range plugins {
		if p.PackageName != "" {
			packages = append(packages, p.PackageName)
			fmt.Printf("   • %s\n", p.Name)
		}
	}

	if len(packages) == 0 {
		fmt.Println("No APT packages to install.")
		return nil
	}

	// Build install command
	var cmd *exec.Cmd
	if useNala {
		args := append([]string{"nala", "install", "-y"}, packages...)
		cmd = exec.Command("sudo", args...)
		fmt.Println()
	} else {
		args := append([]string{"apt", "install", "-y"}, packages...)
		cmd = exec.Command("sudo", args...)
		fmt.Println()
	}

	// Execute with live output
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("batch APT installation failed: %w", err)
	}

	fmt.Printf("\n✅ Successfully installed %d APT package(s)\n", len(plugins))

	// Run post-install tasks with progress
	if hasPostInstall := countPostInstalls(plugins); hasPostInstall > 0 {
		fmt.Println("\n⚙️  Running post-installation tasks...")
		for _, p := range plugins {
			if p.PostInstall != "" {
				fmt.Printf("   • %s... ", p.Name)
				if err := executeScript(p, "post_install", env); err != nil {
					fmt.Println("❌")
				} else {
					fmt.Println("✅")
				}
			}
		}
	}

	return nil
}
