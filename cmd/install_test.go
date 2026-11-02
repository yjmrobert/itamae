package cmd

import (
	"strings"
	"testing"

	"github.com/yjmrobert/itamae/itamae"
)

func TestResolveInputsFromFlagsAndEnv(t *testing.T) {
	plugins := []itamae.ToolPlugin{{
		ID: "git",
		RequiredInputs: []itamae.Input{
			{Name: "GIT_USER_NAME", Prompt: "name"},
			{Name: "GIT_USER_EMAIL", Prompt: "email"},
		},
	}}

	t.Setenv("GIT_USER_EMAIL", "env@example.com")
	values, err := resolveInputs(plugins, []string{"GIT_USER_NAME=Flag Name"}, false)
	if err != nil {
		t.Fatalf("resolveInputs failed: %v", err)
	}
	if values["GIT_USER_NAME"] != "Flag Name" {
		t.Errorf("flag input not applied: %v", values)
	}
	if values["GIT_USER_EMAIL"] != "env@example.com" {
		t.Errorf("environment input not applied: %v", values)
	}
}

func TestResolveInputsMissingNonInteractive(t *testing.T) {
	plugins := []itamae.ToolPlugin{{
		ID:             "git",
		RequiredInputs: []itamae.Input{{Name: "ITAMAE_TEST_MISSING_INPUT", Prompt: "x"}},
	}}
	t.Setenv("ITAMAE_TEST_MISSING_INPUT", "")

	_, err := resolveInputs(plugins, nil, false)
	if err == nil || !strings.Contains(err.Error(), "ITAMAE_TEST_MISSING_INPUT") {
		t.Fatalf("expected a missing input error naming the input, got %v", err)
	}
}

func TestResolveInputsRejectsMalformedFlag(t *testing.T) {
	_, err := resolveInputs(nil, []string{"NOEQUALS"}, false)
	if err == nil {
		t.Fatal("expected an error for a malformed --input value")
	}
}

func TestSelectPluginsNonInteractive(t *testing.T) {
	plugins := []itamae.ToolPlugin{{ID: "a", Name: "A"}, {ID: "b", Name: "B"}}
	t.Cleanup(func() {
		installFlags.plugins = nil
		installFlags.all = false
	})

	// Core and essentials install everything by default.
	got, err := selectPlugins(plugins, "core", false)
	if err != nil || len(got) != 2 {
		t.Fatalf("core should select all plugins, got %d (%v)", len(got), err)
	}

	// Unverified refuses to guess without --all or --plugins.
	if _, err := selectPlugins(plugins, "unverified", false); err == nil {
		t.Fatal("unverified without --all/--plugins should fail when not interactive")
	}

	installFlags.all = true
	got, err = selectPlugins(plugins, "unverified", false)
	if err != nil || len(got) != 2 {
		t.Fatalf("--all should select all plugins, got %d (%v)", len(got), err)
	}

	installFlags.all = false
	installFlags.plugins = []string{"b"}
	got, err = selectPlugins(plugins, "unverified", false)
	if err != nil || len(got) != 1 || got[0].ID != "b" {
		t.Fatalf("--plugins should select only b, got %v (%v)", got, err)
	}

	installFlags.plugins = []string{"nope"}
	if _, err := selectPlugins(plugins, "core", false); err == nil {
		t.Fatal("unknown plugin ID should fail")
	}
}
