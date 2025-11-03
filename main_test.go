package main

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
	"sudo", "nala", "apt-get", "ln", "mkdir", "rm", "chsh", "git", "curl", "unzip", "stow", "chmod", "bash", "sh",
}

// pluginAssertion defines the expected commands for a plugin's install and remove actions.
type pluginAssertion struct {
	install string
	remove  string
}

// pluginAssertions is a map of plugin IDs to their expected command assertions.
var pluginAssertions = map[string]pluginAssertion{
	"atuin":     {install: "bash", remove: "bash -s -- --uninstall"},
	"bat":       {install: "sudo nala install -y batcat", remove: "sudo apt-get purge -y batcat"},
	"btop":      {install: "sudo nala install -y btop", remove: "sudo apt-get purge -y btop"},
	"chezmoi":   {install: "sh -c", remove: "rm"},
	"dunst":     {install: "sudo nala install -y dunst", remove: "sudo apt-get purge -y dunst"},
	"fd":        {install: "sudo nala install -y fd-find", remove: "sudo apt-get purge -y fd-find"},
	"flameshot": {install: "sudo nala install -y flameshot", remove: "sudo apt-get purge -y flameshot"},
	"fzf":       {install: "sudo nala install -y fzf", remove: "sudo apt-get purge -y fzf"},
	"ghostty":   {install: "mkdir -p", remove: "rm -f"},
	"httpie":    {install: "sudo nala install -y httpie", remove: "sudo apt-get purge -y httpie"},
	"jq":        {install: "sudo nala install -y jq", remove: "sudo apt-get purge -y jq"},
	"lsd":       {install: "sudo nala install -y lsd", remove: "sudo apt-get purge -y lsd"},
	"pass":      {install: "sudo nala install -y pass", remove: "sudo apt-get purge -y pass"},
	"polybar":   {install: "sudo nala install -y polybar", remove: "sudo apt-get purge -y polybar"},
	"ripgrep":   {install: "sudo nala install -y ripgrep", remove: "sudo apt-get purge -y ripgrep"},
	"rofi":      {install: "sudo nala install -y rofi", remove: "sudo apt-get purge -y rofi"},
	"starship":  {install: "curl -sS https://starship.rs/install.sh", remove: "sh -c rm \"$(command -v starship)\""},
	"stow":      {install: "sudo nala install -y stow", remove: "sudo apt-get purge -y stow"},
	"tldr":      {install: "curl -L https://github.com/tealdeer-rs/tealdeer/releases/latest/download/tealdeer-linux-x86_64-musl", remove: "rm"},
	"vscode":    {install: "sudo apt-get install -y /tmp/vscode-itamae.deb", remove: "sudo apt-get purge -y code"},
	"yq":        {install: "sudo curl -L https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64 -o /usr/local/bin/yq", remove: "sudo rm /usr/local/bin/yq"},
	"zellij":    {install: "curl -L https://github.com/zellij-project/zellij/releases/latest/download/zellij-x86_64-unknown-linux-musl.tar.gz", remove: "sudo rm /usr/local/bin/zellij"},
	"zoxide":    {install: "curl -sS https://raw.githubusercontent.com/ajeetdsouza/zoxide/main/install.sh", remove: "rm"},
	"zsh":       {install: "sudo nala install -y zsh", remove: "sudo apt-get purge -y zsh"},
}

// TestMain sets up the test environment for the entire package.
func TestMain(m *testing.M) {
	var cleanupPlugins func()
	var err error
	plugins, cleanupPlugins, err = loadPlugins()
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

// TestRemove runs the remove command for every plugin.
func TestRemove(t *testing.T) {
	for _, plugin := range plugins {
		plugin := plugin
		t.Run(fmt.Sprintf("remove_%s", plugin.ID), func(t *testing.T) {
			assertion, ok := pluginAssertions[plugin.ID]
			if !ok {
				t.Skip("No assertion defined for this plugin.")
			}
			runPluginTest(t, plugin, "remove", assertion.remove)
		})
	}
}
