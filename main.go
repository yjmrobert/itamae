package main

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
)

//go:embed scripts/*
var scriptsFS embed.FS

type ToolPlugin struct {
	ID          string // "vscode", "ripgrep"
	Name        string // "Visual Studio Code"
	Description string
	Omakase     bool
	ScriptPath  string // The path to the executable in the /tmp/ directory
}

func main() {
	plugins, cleanup, err := loadPlugins()
	if err != nil {
		fmt.Printf("Error loading plugins: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	args := os.Args[1:]
	if len(args) == 0 {
		runOmakase(plugins)
		return
	}

	command := args[0]
	switch command {
	case "customize":
		runCustomize(plugins)
	case "remove":
		runRemove(plugins)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Usage: itamae [customize|remove]")
		os.Exit(1)
	}
}

func runOmakase(plugins []ToolPlugin) {
	fmt.Println("Starting Itamae (Chef's Choice) setup...")
	for _, p := range plugins {
		if p.Omakase {
			fmt.Printf("--- Installing %s ---\n", p.Name)
			if err := executeScript(p, "install"); err != nil {
				fmt.Printf("Error installing %s: %v\n", p.Name, err)
			}
		}
	}
	fmt.Println("Setup complete.")
}

func runCustomize(plugins []ToolPlugin) {
	selected, err := runTUI(plugins, "Itamae - À La Carte Setup")
	if err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
		os.Exit(1)
	}

	if len(selected) == 0 {
		fmt.Println("No plugins selected. Exiting.")
		return
	}

	fmt.Println("Installing selected plugins...")
	for _, p := range selected {
		fmt.Printf("--- Installing %s ---\n", p.Name)
		if err := executeScript(p, "install"); err != nil {
			fmt.Printf("Error installing %s: %v\n", p.Name, err)
		}
	}
	fmt.Println("Installation complete.")
}

func runRemove(plugins []ToolPlugin) {
	selected, err := runTUI(plugins, "Itamae - Remove Software")
	if err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
		os.Exit(1)
	}

	if len(selected) == 0 {
		fmt.Println("No plugins selected. Exiting.")
		return
	}

	fmt.Println("Removing selected plugins...")
	for _, p := range selected {
		fmt.Printf("--- Removing %s ---\n", p.Name)
		if err := executeScript(p, "remove"); err != nil {
			fmt.Printf("Error removing %s: %v\n", p.Name, err)
		}
	}
	fmt.Println("Removal complete.")
}

func executeScript(plugin ToolPlugin, command string) error {
	cmd := exec.Command("bash", plugin.ScriptPath, command)

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

func loadPlugins() ([]ToolPlugin, func(), error) {
	tmpDir, err := os.MkdirTemp("", "itamae-scripts-")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create temp dir: %w", err)
	}

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	var plugins []ToolPlugin

	files, err := scriptsFS.ReadDir("scripts")
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
	scriptPath := filepath.Join("scripts", fileName)

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
		}
	}
	if err := scanner.Err(); err != nil {
		return ToolPlugin{}, fmt.Errorf("scanner error while parsing metadata: %w", err)
	}
	return plugin, nil
}
