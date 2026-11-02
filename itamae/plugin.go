package itamae

import (
	"bufio"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

//go:embed scripts/core/* scripts/essentials/* scripts/unverified/*
var scriptsFS embed.FS

// Categories lists the plugin categories in display order.
var Categories = []string{"core", "essentials", "unverified"}

// IsValidCategory reports whether name is one of the known plugin categories.
func IsValidCategory(name string) bool {
	for _, c := range Categories {
		if c == name {
			return true
		}
	}
	return false
}

// Input describes a value a plugin needs from the user before it can run.
type Input struct {
	Name       string // Environment variable name passed to the script, e.g. GIT_USER_NAME
	Prompt     string // Question shown to the user
	DefaultCmd string // Optional shell command whose output pre-fills the prompt
}

// ToolPlugin is a single installable tool backed by an embedded shell script.
type ToolPlugin struct {
	ID             string // Script file name without extension, e.g. "ripgrep"
	Name           string // Human readable name, e.g. "ripgrep"
	Description    string
	Category       string // "core", "essentials" or "unverified"
	ScriptPath     string // Path of the script unpacked on disk
	InstallMethod  string // "apt", "binary" or "manual"
	PackageName    string // APT package name (apt plugins only)
	RepoSetup      string // Script function that configures an APT repository (optional)
	PostInstall    string // Script function run after the batch APT install (optional)
	RequiredInputs []Input
}

// LoadPlugins unpacks the scripts of one category to a temporary directory and
// parses their metadata. The returned cleanup function removes that directory.
func LoadPlugins(category string) ([]ToolPlugin, func(), error) {
	if !IsValidCategory(category) {
		return nil, func() {}, fmt.Errorf("unknown category %q (expected one of: %s)", category, strings.Join(Categories, ", "))
	}

	tmpDir, err := os.MkdirTemp("", "itamae-scripts-")
	if err != nil {
		return nil, func() {}, fmt.Errorf("failed to create temp dir: %w", err)
	}
	cleanup := func() { _ = os.RemoveAll(tmpDir) }

	plugins, err := loadCategory(category, tmpDir)
	if err != nil {
		cleanup()
		return nil, func() {}, err
	}
	return plugins, cleanup, nil
}

// LoadAllPlugins loads the plugins of every category, in category order.
func LoadAllPlugins() ([]ToolPlugin, func(), error) {
	tmpDir, err := os.MkdirTemp("", "itamae-scripts-")
	if err != nil {
		return nil, func() {}, fmt.Errorf("failed to create temp dir: %w", err)
	}
	cleanup := func() { _ = os.RemoveAll(tmpDir) }

	var all []ToolPlugin
	for _, category := range Categories {
		plugins, err := loadCategory(category, filepath.Join(tmpDir, category))
		if err != nil {
			cleanup()
			return nil, func() {}, err
		}
		all = append(all, plugins...)
	}
	return all, cleanup, nil
}

func loadCategory(category, tmpDir string) ([]ToolPlugin, error) {
	if err := os.MkdirAll(tmpDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create script dir: %w", err)
	}

	scriptDir := "scripts/" + category
	files, err := scriptsFS.ReadDir(scriptDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded scripts dir %s: %w", scriptDir, err)
	}

	var plugins []ToolPlugin
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sh") {
			continue
		}
		plugin, err := processPluginFile(file, tmpDir, category)
		if err != nil {
			return nil, fmt.Errorf("failed to process plugin %s: %w", file.Name(), err)
		}
		plugins = append(plugins, plugin)
	}

	sort.Slice(plugins, func(i, j int) bool { return plugins[i].ID < plugins[j].ID })
	return plugins, nil
}

func processPluginFile(file fs.DirEntry, tmpDir string, category string) (ToolPlugin, error) {
	fileName := file.Name()
	scriptPath := "scripts/" + category + "/" + fileName

	content, err := scriptsFS.ReadFile(scriptPath)
	if err != nil {
		return ToolPlugin{}, fmt.Errorf("failed to read embedded file %s: %w", scriptPath, err)
	}

	plugin, err := parseMetadata(string(content))
	if err != nil {
		return ToolPlugin{}, fmt.Errorf("failed to parse metadata for %s: %w", fileName, err)
	}
	plugin.ID = strings.TrimSuffix(fileName, ".sh")
	plugin.Category = category
	if plugin.Name == "" {
		return ToolPlugin{}, fmt.Errorf("%s is missing the NAME metadata field", fileName)
	}
	if plugin.InstallMethod == "" {
		return ToolPlugin{}, fmt.Errorf("%s is missing the INSTALL_METHOD metadata field", fileName)
	}
	if plugin.InstallMethod == "apt" && plugin.PackageName == "" {
		return ToolPlugin{}, fmt.Errorf("%s uses INSTALL_METHOD apt but has no PACKAGE_NAME", fileName)
	}

	destPath := filepath.Join(tmpDir, fileName)
	if err := os.WriteFile(destPath, content, 0o755); err != nil { //nolint:gosec // scripts must be executable
		return ToolPlugin{}, fmt.Errorf("failed to write script to temp dir: %w", err)
	}
	plugin.ScriptPath = destPath

	return plugin, nil
}

// parseMetadata reads the "# KEY: value" header block at the top of a plugin script.
func parseMetadata(content string) (ToolPlugin, error) {
	plugin := ToolPlugin{}
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "#") {
			continue
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
		case "DESCRIPTION":
			plugin.Description = value
		case "INSTALL_METHOD":
			plugin.InstallMethod = value
		case "PACKAGE_NAME":
			plugin.PackageName = value
		case "REPO_SETUP":
			plugin.RepoSetup = value
		case "POST_INSTALL":
			plugin.PostInstall = value
		case "REQUIRES":
			fields := strings.SplitN(value, "|", 3)
			if len(fields) >= 2 {
				input := Input{Name: strings.TrimSpace(fields[0]), Prompt: strings.TrimSpace(fields[1])}
				if len(fields) == 3 {
					input.DefaultCmd = strings.TrimSpace(fields[2])
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

// FilterPlugins returns the plugins whose IDs are listed, preserving the order of
// plugins. Unknown IDs produce an error.
func FilterPlugins(plugins []ToolPlugin, ids []string) ([]ToolPlugin, error) {
	wanted := make(map[string]bool, len(ids))
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id != "" {
			wanted[id] = true
		}
	}

	var selected []ToolPlugin
	for _, p := range plugins {
		if wanted[p.ID] {
			selected = append(selected, p)
			delete(wanted, p.ID)
		}
	}

	if len(wanted) > 0 {
		unknown := make([]string, 0, len(wanted))
		for id := range wanted {
			unknown = append(unknown, id)
		}
		sort.Strings(unknown)
		return nil, fmt.Errorf("unknown plugin(s): %s", strings.Join(unknown, ", "))
	}
	return selected, nil
}

// RequiredInputs returns the distinct inputs required by the given plugins, in order.
func RequiredInputs(plugins []ToolPlugin) []Input {
	seen := make(map[string]bool)
	var inputs []Input
	for _, p := range plugins {
		for _, in := range p.RequiredInputs {
			if !seen[in.Name] {
				seen[in.Name] = true
				inputs = append(inputs, in)
			}
		}
	}
	return inputs
}
