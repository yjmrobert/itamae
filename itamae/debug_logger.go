package itamae

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

var (
	debugLog     *os.File
	debugLogMux  sync.Mutex
	debugLogPath string
)

// InitDebugLog creates a new timestamped debug log file in GetLogDirectory().
// Calling it again while a log is open is a no-op.
func InitDebugLog() error {
	debugLogMux.Lock()
	defer debugLogMux.Unlock()

	if debugLog != nil {
		return nil
	}

	logDir := GetLogDirectory()
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	path := filepath.Join(logDir, fmt.Sprintf("itamae-install-%s.log", timestamp))

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	debugLog = f
	debugLogPath = path

	writeDebugLine("=== Itamae Installation Log ===")
	writeDebugLine("Started at: " + time.Now().Format(time.RFC3339))
	writeDebugLine("Log file: " + path)
	writeDebugLine("================================")
	return nil
}

// DebugLog writes a formatted, timestamped line to the debug log. It is safe to
// call from any goroutine and does nothing when no log is open.
func DebugLog(format string, args ...any) {
	debugLogMux.Lock()
	defer debugLogMux.Unlock()
	writeDebugLine(fmt.Sprintf(format, args...))
}

// writeDebugLine writes one line; the caller must hold debugLogMux.
func writeDebugLine(message string) {
	if debugLog == nil {
		return
	}
	timestamp := time.Now().Format("15:04:05.000")
	_, _ = fmt.Fprintf(debugLog, "[%s] %s\n", timestamp, message)
}

// CloseDebugLog flushes and closes the debug log and, when w is not nil, prints
// where it was saved.
func CloseDebugLog(w io.Writer) {
	debugLogMux.Lock()
	defer debugLogMux.Unlock()

	if debugLog == nil {
		return
	}

	writeDebugLine("=== Installation Complete ===")
	writeDebugLine("Ended at: " + time.Now().Format(time.RFC3339))
	_ = debugLog.Sync()
	_ = debugLog.Close()
	debugLog = nil

	if w != nil {
		fmt.Fprintf(w, "\n📝 Debug log saved to: %s\n", debugLogPath)
	}
}

// GetDebugLogPath returns the path of the current (or most recently closed) debug log.
func GetDebugLogPath() string {
	debugLogMux.Lock()
	defer debugLogMux.Unlock()
	return debugLogPath
}

// GetLogDirectory returns the directory that holds the debug logs.
func GetLogDirectory() string {
	return filepath.Join(os.TempDir(), "itamae-logs")
}

// ListLogFiles returns the debug log files, newest first.
func ListLogFiles() ([]string, error) {
	logDir := GetLogDirectory()

	entries, err := os.ReadDir(logDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read log directory: %w", err)
	}

	type logFile struct {
		path    string
		modTime time.Time
	}
	var logFiles []logFile
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".log") {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		logFiles = append(logFiles, logFile{path: filepath.Join(logDir, entry.Name()), modTime: info.ModTime()})
	}

	sort.Slice(logFiles, func(i, j int) bool {
		return logFiles[i].modTime.After(logFiles[j].modTime)
	})

	result := make([]string, len(logFiles))
	for i, lf := range logFiles {
		result[i] = lf.path
	}
	return result, nil
}

// GetMostRecentLog returns the path to the newest debug log.
func GetMostRecentLog() (string, error) {
	logs, err := ListLogFiles()
	if err != nil {
		return "", err
	}
	if len(logs) == 0 {
		return "", fmt.Errorf("no log files found in %s", GetLogDirectory())
	}
	return logs[0], nil
}
