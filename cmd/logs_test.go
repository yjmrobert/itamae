package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/yjmrobert/itamae/itamae"
)

func TestColorizeLinePreservesContent(t *testing.T) {
	testCases := []string{
		"[15:04:05.123] Simple message",
		"[15:04:05.123] ERROR: Something failed",
		"[15:04:05.123] ✅ Installation successful",
		"[15:04:05.123] Phase 1 starting",
		"[15:04:05.123] WARNING: careful",
		"No timestamp message",
		"",
	}

	for _, tc := range testCases {
		result := colorizeLine(tc)
		// Styles may add escape codes but must never drop the text itself.
		for _, word := range strings.Fields(tc) {
			if !strings.Contains(result, word) {
				t.Errorf("colorizeLine(%q) lost %q: %q", tc, word, result)
			}
		}
	}
}

func TestRelativeTime(t *testing.T) {
	cases := map[time.Duration]string{
		10 * time.Second: "just now",
		5 * time.Minute:  "5 minutes ago",
		3 * time.Hour:    "3 hours ago",
		49 * time.Hour:   "2 days ago",
	}
	for d, want := range cases {
		if got := relativeTime(d); got != want {
			t.Errorf("relativeTime(%v) = %q, want %q", d, got, want)
		}
	}
}

func TestShouldUsePager(t *testing.T) {
	tmpDir := t.TempDir()

	small := filepath.Join(tmpDir, "small.log")
	if err := os.WriteFile(small, []byte("[15:04:05.123] Line 1\n[15:04:05.124] Line 2\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if shouldUsePager(small) {
		t.Error("shouldUsePager should return false for a small file")
	}

	large := filepath.Join(tmpDir, "large.log")
	var sb strings.Builder
	for i := 0; i < 60; i++ {
		sb.WriteString("[15:04:05.123] Log line\n")
	}
	if err := os.WriteFile(large, []byte(sb.String()), 0o644); err != nil {
		t.Fatal(err)
	}
	// Tests never run on an interactive terminal, so the pager must stay off
	// regardless of the file size.
	if shouldUsePager(large) {
		t.Error("shouldUsePager must return false when not attached to a terminal")
	}
}

func TestLogDirectoryHelpers(t *testing.T) {
	logDir := itamae.GetLogDirectory()
	if filepath.Base(logDir) != "itamae-logs" {
		t.Errorf("GetLogDirectory should end with 'itamae-logs', got: %s", logDir)
	}

	logs, err := itamae.ListLogFiles()
	if err != nil {
		t.Fatalf("ListLogFiles failed: %v", err)
	}
	for i := 0; i+1 < len(logs); i++ {
		a, errA := os.Stat(logs[i])
		b, errB := os.Stat(logs[i+1])
		if errA != nil || errB != nil {
			continue
		}
		if a.ModTime().Before(b.ModTime()) {
			t.Error("ListLogFiles must sort newest first")
		}
	}
}

func BenchmarkColorizeLine(b *testing.B) {
	line := "[15:04:05.123] This is a test log message"
	for i := 0; i < b.N; i++ {
		colorizeLine(line)
	}
}
