package itamae

import (
	"os"
	"strings"
	"testing"
)

func TestProcessCustomInstall(t *testing.T) {
	mockDir, logPath, cleanup := setupTestEnvironment()
	defer cleanup()

	originalPath := os.Getenv("PATH")
	newPath := mockDir + ":" + originalPath
	os.Setenv("PATH", newPath)
	defer os.Setenv("PATH", originalPath)

	if err := os.WriteFile(logPath, []byte{}, 0644); err != nil {
		t.Fatalf("Failed to clear log file: %v", err)
	}

	mockPlugins := []ToolPlugin{
		{ID: "git", Name: "Git", InstallMethod: "apt", PackageName: "git", ScriptPath: "/tmp/itamae-test-git.sh"},
		{ID: "my-script", Name: "My Script", InstallMethod: "binary", ScriptPath: "/tmp/itamae-test-script.sh"},
	}
	requiredInputs := map[string]string{
		"GIT_USER_NAME":  "test-user",
		"GIT_USER_EMAIL": "test@example.com",
	}

	os.WriteFile("/tmp/itamae-test-git.sh", []byte("#!/bin/bash\necho 'git script'"), 0755)
	os.WriteFile("/tmp/itamae-test-script.sh", []byte("#!/bin/bash\necho 'my script'"), 0755)
	defer os.Remove("/tmp/itamae-test-git.sh")
	defer os.Remove("/tmp/itamae-test-script.sh")

	processCustomInstall(mockPlugins, requiredInputs)

	logBytes, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}
	logContent := string(logBytes)

	expectedAptLog := "sudo nala install -y git"
	if !strings.Contains(logContent, expectedAptLog) {
		t.Errorf("Expected log to contain '%s', but got:\n%s", expectedAptLog, logContent)
	}

	expectedBinaryLog := "bash /tmp/itamae-test-script.sh install"
	if !strings.Contains(logContent, expectedBinaryLog) {
		t.Errorf("Expected log to contain '%s', but got:\n%s", expectedBinaryLog, logContent)
	}
}
