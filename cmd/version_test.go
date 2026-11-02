package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestVersionDefaults(t *testing.T) {
	if Version == "" || GitCommit == "" || BuildDate == "" {
		t.Fatal("version variables must have default values")
	}
}

func TestCommandsRegistered(t *testing.T) {
	for _, name := range []string{"install", "check", "list", "logs", "version"} {
		c, _, err := rootCmd.Find([]string{name})
		if err != nil {
			t.Fatalf("command %q not found: %v", name, err)
		}
		if c.Short == "" {
			t.Errorf("command %q should have a short description", name)
		}
	}
}

func TestVersionCommandOutput(t *testing.T) {
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"version"})
	t.Cleanup(func() {
		rootCmd.SetOut(nil)
		rootCmd.SetErr(nil)
		rootCmd.SetArgs(nil)
	})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("version command failed: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"itamae version " + Version, "commit: " + GitCommit, "built:  " + BuildDate} {
		if !strings.Contains(out, want) {
			t.Errorf("version output missing %q:\n%s", want, out)
		}
	}
}
