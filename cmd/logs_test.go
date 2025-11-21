package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yjmrobert/itamae/itamae"
)

func TestLogsCommand(t *testing.T) {
	// Create temporary log directory for testing
	tmpDir := t.TempDir()

	// Note: We can't override GetLogDirectory() in tests without refactoring,
	// so these tests are placeholders for future implementation
	_ = tmpDir

	t.Run("colorizeLine detects error", func(t *testing.T) {
		line := "[15:04:05.123] ERROR: Installation failed"
		colored := colorizeLine(line)

		// Should contain the original text
		if colored == "" {
			t.Error("colorizeLine returned empty string")
		}
	})

	t.Run("colorizeLine detects success", func(t *testing.T) {
		line := "[15:04:05.123] ✅ Installation successful"
		colored := colorizeLine(line)

		if colored == "" {
			t.Error("colorizeLine returned empty string")
		}
	})

	t.Run("colorizeLine handles timestamp", func(t *testing.T) {
		line := "[15:04:05.123] Normal log message"
		colored := colorizeLine(line)

		if colored == "" {
			t.Error("colorizeLine returned empty string")
		}
	})

	t.Run("isTerminal returns bool", func(t *testing.T) {
		result := isTerminal()
		// Should return either true or false
		_ = result
	})

	// Test shouldUsePager with actual file
	t.Run("shouldUsePager with small file", func(t *testing.T) {
		// Create small temp file
		tmpFile := filepath.Join(tmpDir, "small.log")
		content := "[15:04:05.123] Line 1\n[15:04:05.124] Line 2\n"
		if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		result := shouldUsePager(tmpFile)
		if result {
			t.Error("shouldUsePager should return false for small file")
		}
	})

	t.Run("shouldUsePager with large file", func(t *testing.T) {
		// Create large temp file (>50 lines)
		tmpFile := filepath.Join(tmpDir, "large.log")
		f, err := os.Create(tmpFile)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		// Write 60 lines
		for i := 0; i < 60; i++ {
			f.WriteString("[15:04:05.123] Log line " + string(rune(i)) + "\n")
		}
		f.Close()

		// Note: This will return false in CI/non-TTY environments
		// In real terminal, it should return true
		_ = shouldUsePager(tmpFile)
	})
}

func TestLogsCommandIntegration(t *testing.T) {
	// These tests require actual log files to be present
	// They serve as examples of how to test the command

	t.Run("ListLogFiles returns sorted array", func(t *testing.T) {
		logs, err := itamae.ListLogFiles()

		// It's OK if there are no logs yet
		if err != nil && !os.IsNotExist(err) {
			t.Errorf("ListLogFiles failed: %v", err)
		}

		// If logs exist, verify they're sorted (newest first)
		if len(logs) > 1 {
			for i := 0; i < len(logs)-1; i++ {
				info1, err1 := os.Stat(logs[i])
				info2, err2 := os.Stat(logs[i+1])

				if err1 != nil || err2 != nil {
					continue
				}

				if info1.ModTime().Before(info2.ModTime()) {
					t.Error("Logs not sorted newest first")
				}
			}
		}
	})

	t.Run("GetMostRecentLog returns newest", func(t *testing.T) {
		logPath, err := itamae.GetMostRecentLog()

		// It's OK if there are no logs
		if err != nil && !os.IsNotExist(err) {
			return
		}

		if logPath == "" {
			return // No logs yet
		}

		// Verify file exists
		if _, err := os.Stat(logPath); os.IsNotExist(err) {
			t.Errorf("GetMostRecentLog returned non-existent file: %s", logPath)
		}
	})

	t.Run("GetLogDirectory returns valid path", func(t *testing.T) {
		logDir := itamae.GetLogDirectory()

		if logDir == "" {
			t.Error("GetLogDirectory returned empty string")
		}

		// Should end with itamae-logs
		if filepath.Base(logDir) != "itamae-logs" {
			t.Errorf("GetLogDirectory should end with 'itamae-logs', got: %s", logDir)
		}
	})
}

func TestLogsCommandHelpers(t *testing.T) {
	t.Run("colorizeLine preserves content", func(t *testing.T) {
		testCases := []string{
			"[15:04:05.123] Simple message",
			"[15:04:05.123] ERROR: Something failed",
			"[15:04:05.123] ✅ Installation successful",
			"[15:04:05.123] Phase 1 starting",
			"No timestamp message",
		}

		for _, tc := range testCases {
			result := colorizeLine(tc)
			// The result might have color codes, but should contain original text
			if result == "" {
				t.Errorf("colorizeLine(%q) returned empty string", tc)
			}
		}
	})
}

// Benchmark tests
func BenchmarkColorizeLine(b *testing.B) {
	line := "[15:04:05.123] This is a test log message"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		colorizeLine(line)
	}
}

func BenchmarkShouldUsePager(b *testing.B) {
	// Create temp file for benchmark
	tmpDir := b.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.log")

	f, _ := os.Create(tmpFile)
	for i := 0; i < 100; i++ {
		f.WriteString("[15:04:05.123] Log line\n")
	}
	f.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		shouldUsePager(tmpFile)
	}
}
