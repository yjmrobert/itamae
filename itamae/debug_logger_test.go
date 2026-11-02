package itamae

import (
	"os"
	"strings"
	"testing"
	"time"
)

// TestDebugLogLifecycle guards against the logger deadlocking: Init and Close
// used to call DebugLog while holding the mutex.
func TestDebugLogLifecycle(t *testing.T) {
	t.Setenv("TMPDIR", t.TempDir())

	done := make(chan error, 1)
	go func() {
		if err := InitDebugLog(); err != nil {
			done <- err
			return
		}
		if err := InitDebugLog(); err != nil { // second call is a no-op
			done <- err
			return
		}
		DebugLog("hello %s", "world")
		var sb strings.Builder
		CloseDebugLog(&sb)
		CloseDebugLog(nil) // closing twice is harmless
		if !strings.Contains(sb.String(), "Debug log saved to") {
			t.Errorf("CloseDebugLog should report the log path, got %q", sb.String())
		}
		done <- nil
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("debug logger deadlocked")
	}

	path, err := GetMostRecentLog()
	if err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"Itamae Installation Log", "hello world", "Installation Complete"} {
		if !strings.Contains(string(content), want) {
			t.Errorf("log missing %q:\n%s", want, content)
		}
	}

	logs, err := ListLogFiles()
	if err != nil || len(logs) != 1 {
		t.Fatalf("expected exactly one log file, got %v (%v)", logs, err)
	}

	// Logging after close must not panic or block.
	DebugLog("after close")
}
