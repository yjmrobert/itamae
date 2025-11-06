package cmd

import (
	"testing"
)

func TestVersionVariables(t *testing.T) {
	// Default values when not set via ldflags
	if Version == "" {
		t.Error("Version should have a default value")
	}

	if GitCommit == "" {
		t.Error("GitCommit should have a default value")
	}

	if BuildDate == "" {
		t.Error("BuildDate should have a default value")
	}

	// Verify default values
	if Version != "dev" {
		t.Logf("Version is set to: %s (expected 'dev' without ldflags)", Version)
	}

	if GitCommit != "unknown" {
		t.Logf("GitCommit is set to: %s (expected 'unknown' without ldflags)", GitCommit)
	}

	if BuildDate != "unknown" {
		t.Logf("BuildDate is set to: %s (expected 'unknown' without ldflags)", BuildDate)
	}
}

func TestVersionCommand(t *testing.T) {
	// Test that version command exists
	cmd := rootCmd
	versionCmd, _, err := cmd.Find([]string{"version"})
	if err != nil {
		t.Fatalf("version command not found: %v", err)
	}

	if versionCmd.Use != "version" {
		t.Errorf("Expected 'version', got '%s'", versionCmd.Use)
	}

	if versionCmd.Short == "" {
		t.Error("Version command should have a short description")
	}

	if versionCmd.Long == "" {
		t.Error("Version command should have a long description")
	}
}

func TestVersionFlag(t *testing.T) {
	// Test that --version flag exists
	flag := rootCmd.PersistentFlags().Lookup("version")
	if flag == nil {
		t.Fatal("--version flag not found")
	}

	if flag.Shorthand != "v" {
		t.Errorf("Expected shorthand 'v', got '%s'", flag.Shorthand)
	}

	if flag.Usage == "" {
		t.Error("Version flag should have usage text")
	}
}
