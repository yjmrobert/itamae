package itamae

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var (
	mockDir     string
	logFilePath string
	plugins     []ToolPlugin
)

// MOCK_COMMANDS is a list of commands that should be mocked during testing.
var MOCK_COMMANDS = []string{
	"sudo", "nala", "apt-get", "ln", "mkdir", "rm", "chsh", "git", "curl", "unzip", "stow", "chmod", "bash", "sh", "pipx", "rustup", "wget", "tar", "fc-cache",
}

// pluginAssertion defines the expected commands for a plugin's install and remove actions.
type pluginAssertion struct {
	install string
	remove  string
}

// pluginAssertions is a map of plugin IDs to their expected command assertions.
var pluginAssertions = map[string]pluginAssertion{
	// Core plugins (OMAKASE: true)
	"alacritty":           {install: "sudo nala install -y alacritty", remove: "sudo apt-get purge -y alacritty"},
	"ansible":             {install: "pipx install", remove: "pipx uninstall ansible"},
	"apt-transport-https": {install: "sudo nala install -y apt-transport-https", remove: "sudo apt-get purge -y apt-transport-https"},
	"atuin":               {install: "bash", remove: "bash -s -- --uninstall"},
	"bat":                 {install: "sudo nala install -y bat", remove: "sudo apt-get purge -y bat"},
	"bin":                 {install: "curl -sL", remove: "rm -f"},
	"btop":                {install: "sudo nala install -y btop", remove: "sudo apt-get purge -y btop"},
	"ca-certificates":     {install: "sudo nala install -y ca-certificates", remove: "sudo apt-get purge -y ca-certificates"},
	"curl":                {install: "sudo nala install -y curl", remove: "sudo apt-get purge -y curl"},
	"dotnet-sdk-8.0":      {install: "wget", remove: "sudo apt-get purge -y dotnet-sdk-8.0"},
	"fd":                  {install: "sudo nala install -y fd-find", remove: "sudo apt-get purge -y fd-find"},
	"fzf":                 {install: "sudo nala install -y fzf", remove: "sudo apt-get purge -y fzf"},
	"gh":                  {install: "sudo mkdir -p /etc/apt/keyrings", remove: "sudo apt-get purge -y gh"},
	"git":                 {install: "sudo nala install -y git", remove: "sudo apt-get purge -y git"},
	"gnupg":               {install: "sudo nala install -y gnupg", remove: "sudo apt-get purge -y gnupg"},
	"httpie":              {install: "sudo nala install -y httpie", remove: "sudo apt-get purge -y httpie"},
	"java":                {install: "sudo mkdir -p /etc/apt/keyrings", remove: "sudo apt-get purge -y temurin-21-jdk"},
	"jq":                  {install: "sudo nala install -y jq", remove: "sudo apt-get purge -y jq"},
	"kubectl":             {install: "curl -LO", remove: "sudo rm -f /usr/local/bin/kubectl"},
	"lsd":                 {install: "sudo nala install -y lsd", remove: "sudo apt-get purge -y lsd"},
	"maven":               {install: "wget", remove: "sudo rm -rf /opt/maven"},
	"nala":                {install: "sudo apt-get install -y nala", remove: "sudo apt-get purge -y nala"},
	"nodejs":              {install: "curl -fsSL", remove: "sudo apt-get purge -y nodejs"},
	"npm":                 {install: "sudo nala install -y npm", remove: "sudo apt-get purge -y npm"},
	"pass":                {install: "sudo nala install -y pass", remove: "sudo apt-get purge -y pass"},
	"pipx":                {install: "sudo nala install -y pipx", remove: "sudo apt-get purge -y pipx"},
	"python3-full":        {install: "sudo nala install -y python3-full", remove: "sudo apt-get purge -y python3-full"},
	"ripgrep":             {install: "sudo nala install -y ripgrep", remove: "sudo apt-get purge -y ripgrep"},
	"ruby":                {install: "sudo nala install -y ruby-full", remove: "sudo apt-get purge -y ruby-full"},
	"rust":                {install: "curl --proto", remove: "rustup self uninstall -y"},
	"sdkman":              {install: "curl -s", remove: "rm -rf"},
	"semgrep":             {install: "pipx install semgrep", remove: "pipx uninstall semgrep"},
	"starship":            {install: "curl -sS https://starship.rs/install.sh", remove: "sh -c rm \"$(command -v starship)\""},
	"stow":                {install: "sudo nala install -y stow", remove: "sudo apt-get purge -y stow"},
	"task":                {install: "sh -c", remove: "rm -f"},
	"tldr":                {install: "curl -L https://github.com/tealdeer-rs/tealdeer/releases/latest/download/tealdeer-linux-x86_64-musl", remove: "rm"},
	"vscode":              {install: "sudo apt-get install -y /tmp/vscode-itamae.deb", remove: "sudo apt-get purge -y code"},
	"wget":                {install: "sudo nala install -y wget", remove: "sudo apt-get purge -y wget"},
	"wireguard":           {install: "sudo nala install -y wireguard", remove: "sudo apt-get purge -y wireguard"},
	"yq":                  {install: "sudo curl -L https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64 -o /usr/local/bin/yq", remove: "sudo rm /usr/local/bin/yq"},
	"helm":                {install: "curl", remove: "sudo rm -f /usr/local/bin/helm"},

	// À la carte plugins (OMAKASE: false)
	"btop-desktop":  {install: "sudo nala install -y btop", remove: "sudo apt-get purge -y btop"},
	"cascadia-code": {install: "mkdir -p", remove: "rm -f"},
	"chezmoi":       {install: "sh -c", remove: "rm"},
	"dunst":         {install: "sudo nala install -y dunst", remove: "sudo apt-get purge -y dunst"},
	"flameshot":     {install: "sudo nala install -y flameshot", remove: "sudo apt-get purge -y flameshot"},
	"ghostty":       {install: "mkdir -p", remove: "rm -f"},
	"kubecolor":     {install: "curl -s", remove: "rm -f"},
	"meld":          {install: "sudo nala install -y meld", remove: "sudo apt-get purge -y meld"},
	"ncdu":          {install: "sudo nala install -y ncdu", remove: "sudo apt-get purge -y ncdu"},
	"polybar":       {install: "sudo nala install -y polybar", remove: "sudo apt-get purge -y polybar"},
	"rofi":          {install: "sudo nala install -y rofi", remove: "sudo apt-get purge -y rofi"},
	"zellij":        {install: "curl -L https://github.com/zellij-project/zellij/releases/latest/download/zellij-x86_64-unknown-linux-musl.tar.gz", remove: "sudo rm /usr/local/bin/zellij"},
	"zoxide":        {install: "curl -sS https://raw.githubusercontent.com/ajeetdsouza/zoxide/main/install.sh", remove: "rm"},
	"zsh":           {install: "sudo nala install -y zsh", remove: "sudo apt-get purge -y zsh"},
} // TestMain sets up the test environment for the entire package.
func TestMain(m *testing.M) {
	var cleanupPlugins func()
	var err error
	plugins, cleanupPlugins, err = LoadPlugins("core")
	if err != nil {
		fmt.Printf("Failed to load plugins in TestMain: %v\n", err)
		os.Exit(1)
	}

	var cleanupMocks func()
	mockDir, logFilePath, cleanupMocks = setupTestEnvironment()

	exitCode := m.Run()

	cleanupPlugins()
	cleanupMocks()
	os.Exit(exitCode)
}

// setupTestEnvironment creates a temporary directory for mock commands.
func setupTestEnvironment() (string, string, func()) {
	dir, err := os.MkdirTemp("", "itamae-test-mocks-")
	if err != nil {
		panic(fmt.Sprintf("Failed to create mock dir: %v", err))
	}

	logFile, err := os.CreateTemp("", "itamae-test-log-")
	if err != nil {
		panic(fmt.Sprintf("Failed to create log file: %v", err))
	}
	path := logFile.Name()
	logFile.Close()

	for _, cmdName := range MOCK_COMMANDS {
		mockScriptPath := filepath.Join(dir, cmdName)
		scriptContent := fmt.Sprintf(`#!/bin/bash
# Mock for '%s'
echo "%s $@" >> %s
`, cmdName, cmdName, path)
		if err := os.WriteFile(mockScriptPath, []byte(scriptContent), 0755); err != nil {
			panic(fmt.Sprintf("Failed to write mock script for %s: %v", cmdName, err))
		}
	}

	cleanup := func() {
		os.RemoveAll(dir)
		os.Remove(path)
	}
	return dir, path, cleanup
}

// runPluginTest executes a plugin's command and verifies the output.
func runPluginTest(t *testing.T, plugin ToolPlugin, command, expectedLog string) {
	t.Helper()

	if err := os.WriteFile(logFilePath, []byte{}, 0644); err != nil {
		t.Fatalf("Failed to clear log file: %v", err)
	}

	cmd := exec.Command("bash", plugin.ScriptPath, command)
	originalPath := os.Getenv("PATH")
	newPath := fmt.Sprintf("%s:%s", mockDir, originalPath)
	cmd.Env = append(os.Environ(), fmt.Sprintf("PATH=%s", newPath))

	homeDir, err := os.MkdirTemp("", "itamae-test-home-")
	if err != nil {
		t.Fatalf("Failed to create temp home dir: %v", err)
	}
	defer os.RemoveAll(homeDir)
	cmd.Env = append(cmd.Env, fmt.Sprintf("HOME=%s", homeDir))

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Script execution for plugin '%s' failed. Output:\n%s\nError: %v", plugin.ID, string(output), err)
	}

	logBytes, err := os.ReadFile(logFilePath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}
	logContent := strings.TrimSpace(string(logBytes))

	if expectedLog != "" && !strings.Contains(logContent, expectedLog) {
		t.Errorf("Expected log to contain '%s', but it was:\n%s", expectedLog, logContent)
	} else if expectedLog == "" && logContent != "" {
		// If we don't expect a command, the log should be empty.
		// This handles scripts that are manual.
		t.Errorf("Expected log to be empty, but it was:\n%s", logContent)
	}
}

// TestInstall runs the install command for every plugin.
func TestInstall(t *testing.T) {
	for _, plugin := range plugins {
		plugin := plugin
		t.Run(fmt.Sprintf("install_%s", plugin.ID), func(t *testing.T) {
			assertion, ok := pluginAssertions[plugin.ID]
			if !ok {
				t.Skip("No assertion defined for this plugin.")
			}
			runPluginTest(t, plugin, "install", assertion.install)
		})
	}
}

// TestUninstall runs the uninstall command for every plugin.
func TestUninstall(t *testing.T) {
	for _, plugin := range plugins {
		plugin := plugin
		t.Run(fmt.Sprintf("uninstall_%s", plugin.ID), func(t *testing.T) {
			assertion, ok := pluginAssertions[plugin.ID]
			if !ok {
				t.Skip("No assertion defined for this plugin.")
			}
			runPluginTest(t, plugin, "remove", assertion.remove)
		})
	}
}

// TestMetadataParsing verifies that all plugins have the new metadata fields set correctly.
func TestMetadataParsing(t *testing.T) {
	expectedInstallMethods := map[string]string{
		// Core plugins (OMAKASE: true)
		"alacritty":           "apt",
		"ansible":             "binary",
		"apt-transport-https": "apt",
		"atuin":               "binary",
		"bat":                 "apt",
		"bin":                 "binary",
		"btop":                "apt",
		"ca-certificates":     "apt",
		"curl":                "apt",
		"dotnet-sdk-8.0":      "binary",
		"fd":                  "apt",
		"fzf":                 "apt",
		"gh":                  "binary",
		"git":                 "apt",
		"gnupg":               "apt",
		"helm":                "binary",
		"httpie":              "apt",
		"java":                "binary",
		"jq":                  "apt",
		"kubectl":             "binary",
		"lsd":                 "apt",
		"maven":               "binary",
		"nala":                "apt",
		"nodejs":              "binary",
		"npm":                 "apt",
		"pass":                "apt",
		"pipx":                "apt",
		"python3-full":        "apt",
		"ripgrep":             "apt",
		"ruby":                "apt",
		"rust":                "binary",
		"sdkman":              "binary",
		"semgrep":             "binary",
		"starship":            "binary",
		"stow":                "apt",
		"task":                "binary",
		"tldr":                "binary",
		"vscode":              "binary",
		"wget":                "apt",
		"wireguard":           "apt",
		"yq":                  "binary",

		// À la carte plugins (OMAKASE: false)
		"btop-desktop":  "apt",
		"cascadia-code": "binary",
		"chezmoi":       "binary",
		"dunst":         "apt",
		"flameshot":     "apt",
		"ghostty":       "manual",
		"kubecolor":     "binary",
		"meld":          "apt",
		"ncdu":          "apt",
		"polybar":       "apt",
		"rofi":          "apt",
		"zellij":        "binary",
		"zoxide":        "binary",
		"zsh":           "apt",
	}

	expectedPackageNames := map[string]string{
		// Core plugins
		"alacritty":           "alacritty",
		"apt-transport-https": "apt-transport-https",
		"bat":                 "bat",
		"btop":                "btop",
		"ca-certificates":     "ca-certificates",
		"curl":                "curl",
		"fd":                  "fd-find",
		"fzf":                 "fzf",
		"git":                 "git",
		"gnupg":               "gnupg",
		"httpie":              "httpie",
		"jq":                  "jq",
		"lsd":                 "lsd",
		"nala":                "nala",
		"npm":                 "npm",
		"pass":                "pass",
		"pipx":                "pipx",
		"python3-full":        "python3-full",
		"ripgrep":             "ripgrep",
		"ruby":                "ruby-full",
		"stow":                "stow",
		"wget":                "wget",
		"wireguard":           "wireguard",

		// À la carte plugins
		"btop-desktop": "btop",
		"dunst":        "dunst",
		"flameshot":    "flameshot",
		"meld":         "meld",
		"ncdu":         "ncdu",
		"polybar":      "polybar",
		"rofi":         "rofi",
		"zsh":          "zsh",
	}

	expectedPostInstall := map[string]string{
		"bat": "post_install",
		"fd":  "post_install",
	}

	for _, plugin := range plugins {
		plugin := plugin
		t.Run(fmt.Sprintf("metadata_%s", plugin.ID), func(t *testing.T) {
			// Check InstallMethod
			expectedMethod, ok := expectedInstallMethods[plugin.ID]
			if !ok {
				t.Errorf("Plugin '%s' has no expected install method defined in test", plugin.ID)
				return
			}
			if plugin.InstallMethod != expectedMethod {
				t.Errorf("Plugin '%s': expected InstallMethod='%s', got '%s'", plugin.ID, expectedMethod, plugin.InstallMethod)
			}

			// Check PackageName for APT plugins
			if plugin.InstallMethod == "apt" {
				expectedPkg, ok := expectedPackageNames[plugin.ID]
				if !ok {
					t.Errorf("Plugin '%s' is APT-based but has no expected package name defined in test", plugin.ID)
				} else if plugin.PackageName != expectedPkg {
					t.Errorf("Plugin '%s': expected PackageName='%s', got '%s'", plugin.ID, expectedPkg, plugin.PackageName)
				}
			}

			// Check PostInstall for plugins that need it
			if expectedPost, ok := expectedPostInstall[plugin.ID]; ok {
				if plugin.PostInstall != expectedPost {
					t.Errorf("Plugin '%s': expected PostInstall='%s', got '%s'", plugin.ID, expectedPost, plugin.PostInstall)
				}
			}
		})
	}
}

// TestBatchInstallSeparation verifies that plugins are correctly separated by install method.
func TestBatchInstallSeparation(t *testing.T) {
	aptPlugins := []ToolPlugin{}
	binaryPlugins := []ToolPlugin{}
	manualPlugins := []ToolPlugin{}

	// Core plugins don't use OMAKASE flag - they're all installed automatically
	for _, p := range plugins {
		switch p.InstallMethod {
		case "apt":
			aptPlugins = append(aptPlugins, p)
		case "binary":
			binaryPlugins = append(binaryPlugins, p)
		case "manual":
			manualPlugins = append(manualPlugins, p)
		}
	}

	// Expected counts for Core plugins
	// APT: git = 1
	expectedAptCount := 1
	// Binary: helm = 1
	expectedBinaryCount := 1

	if len(aptPlugins) != expectedAptCount {
		t.Errorf("Expected %d Core APT plugins, got %d: %v", expectedAptCount, len(aptPlugins), getPluginNames(aptPlugins))
	}

	if len(binaryPlugins) != expectedBinaryCount {
		t.Errorf("Expected %d Core binary plugins, got %d: %v", expectedBinaryCount, len(binaryPlugins), getPluginNames(binaryPlugins))
	}

	// Verify all APT plugins have package names
	for _, p := range aptPlugins {
		if p.PackageName == "" {
			t.Errorf("APT plugin '%s' is missing PackageName", p.ID)
		}
	}
}

// Helper function to get plugin names for debugging
func getPluginNames(plugins []ToolPlugin) []string {
	names := make([]string, len(plugins))
	for i, p := range plugins {
		names[i] = p.ID
	}
	return names
}
