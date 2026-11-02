package itamae

import (
	"strings"
	"testing"
)

func TestParseMetadata(t *testing.T) {
	content := `#!/bin/bash
#
# METADATA
# NAME: Git
# DESCRIPTION: A version control system.
# INSTALL_METHOD: apt
# PACKAGE_NAME: git
# REPO_SETUP: setup_repo
# POST_INSTALL: post_install
# REQUIRES: GIT_USER_NAME|Enter your Git user name|git config --global user.name 2>/dev/null || echo ''
# REQUIRES: GIT_USER_EMAIL|Enter your Git user email

install() { :; }
`
	p, err := parseMetadata(content)
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != "Git" || p.Description != "A version control system." || p.InstallMethod != "apt" || p.PackageName != "git" {
		t.Errorf("basic fields not parsed: %+v", p)
	}
	if p.RepoSetup != "setup_repo" || p.PostInstall != "post_install" {
		t.Errorf("hook fields not parsed: %+v", p)
	}
	if len(p.RequiredInputs) != 2 {
		t.Fatalf("expected 2 inputs, got %d", len(p.RequiredInputs))
	}
	if p.RequiredInputs[0].DefaultCmd != "git config --global user.name 2>/dev/null || echo ''" {
		t.Errorf("default command not parsed: %q", p.RequiredInputs[0].DefaultCmd)
	}
	if p.RequiredInputs[1].DefaultCmd != "" || p.RequiredInputs[1].Prompt != "Enter your Git user email" {
		t.Errorf("second input not parsed: %+v", p.RequiredInputs[1])
	}
}

func TestLoadAllPluginsCoversEveryCategory(t *testing.T) {
	all, cleanup, err := LoadAllPlugins()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	counts := make(map[string]int)
	ids := make(map[string]bool)
	for _, p := range all {
		counts[p.Category]++
		if ids[p.ID] {
			t.Errorf("duplicate plugin ID across categories: %s", p.ID)
		}
		ids[p.ID] = true
		if p.Name == "" || p.Description == "" || p.InstallMethod == "" || p.ScriptPath == "" {
			t.Errorf("plugin %s has incomplete metadata: %+v", p.ID, p)
		}
		switch p.InstallMethod {
		case "apt", "binary", "manual":
		default:
			t.Errorf("plugin %s has unknown install method %q", p.ID, p.InstallMethod)
		}
	}
	for _, c := range Categories {
		if counts[c] == 0 {
			t.Errorf("category %s has no plugins", c)
		}
	}
}

func TestLoadPluginsRejectsUnknownCategory(t *testing.T) {
	_, cleanup, err := LoadPlugins("nope")
	cleanup()
	if err == nil || !strings.Contains(err.Error(), "unknown category") {
		t.Fatalf("expected unknown category error, got %v", err)
	}
}

func TestFilterPlugins(t *testing.T) {
	plugins := []ToolPlugin{{ID: "a"}, {ID: "b"}, {ID: "c"}}

	got, err := FilterPlugins(plugins, []string{"c", " a "})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0].ID != "a" || got[1].ID != "c" {
		t.Errorf("FilterPlugins should keep plugin order, got %v", got)
	}

	_, err = FilterPlugins(plugins, []string{"a", "zzz", "yyy"})
	if err == nil || !strings.Contains(err.Error(), "yyy, zzz") {
		t.Errorf("expected sorted unknown IDs in error, got %v", err)
	}
}

func TestRequiredInputsDeduplicates(t *testing.T) {
	plugins := []ToolPlugin{
		{ID: "a", RequiredInputs: []Input{{Name: "X"}, {Name: "Y"}}},
		{ID: "b", RequiredInputs: []Input{{Name: "Y"}, {Name: "Z"}}},
	}
	got := RequiredInputs(plugins)
	if len(got) != 3 || got[0].Name != "X" || got[1].Name != "Y" || got[2].Name != "Z" {
		t.Errorf("unexpected inputs: %+v", got)
	}
}

func TestOutputTail(t *testing.T) {
	if got := outputTail("", 3); got != "" {
		t.Errorf("empty output should give empty tail, got %q", got)
	}
	got := outputTail("a\n\nb\nc\nd\n", 2)
	if !strings.Contains(got, "c") || !strings.Contains(got, "d") || strings.Contains(got, "a") {
		t.Errorf("tail should keep the last 2 lines, got %q", got)
	}
}
