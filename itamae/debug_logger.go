package itamae

import (
	"fmt"
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

// InitDebugLog initializes the debug log file
func InitDebugLog() error {
	debugLogMux.Lock()
	defer debugLogMux.Unlock()

	if debugLog != nil {
		return nil // Already initialized
	}

	// Create log directory
	logDir := filepath.Join(os.TempDir(), "itamae-logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Create log file with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	debugLogPath = filepath.Join(logDir, fmt.Sprintf("itamae-install-%s.log", timestamp))

	var err error
	debugLog, err = os.OpenFile(debugLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	// Write header
	DebugLog("=== Itamae Installation Log ===")
	DebugLog("Started at: %s", time.Now().Format(time.RFC3339))
	DebugLog("Log file: %s", debugLogPath)
	DebugLog("================================\n")

	return nil
}

// DebugLog writes a formatted message to the debug log file
func DebugLog(format string, args ...interface{}) {
	debugLogMux.Lock()
	defer debugLogMux.Unlock()

	if debugLog == nil {
		return // Not initialized
	}

	timestamp := time.Now().Format("15:04:05.000")
	message := fmt.Sprintf(format, args...)
	line := fmt.Sprintf("[%s] %s\n", timestamp, message)

	debugLog.WriteString(line)
	debugLog.Sync() // Flush to disk immediately
}

// CloseDebugLog closes the debug log file and prints its location
func CloseDebugLog() {
	debugLogMux.Lock()
	defer debugLogMux.Unlock()

	if debugLog != nil {
		DebugLog("\n=== Installation Complete ===")
		DebugLog("Ended at: %s", time.Now().Format(time.RFC3339))
		DebugLog("=============================")

		debugLog.Close()

		// Print location to stdout
		fmt.Printf("\n📝 Debug log saved to: %s\n", debugLogPath)

		debugLog = nil
	}
}

// GetDebugLogPath returns the current debug log file path
func GetDebugLogPath() string {
	debugLogMux.Lock()
	defer debugLogMux.Unlock()
	return debugLogPath
}

// GetLogDirectory returns the path to the log directory
func GetLogDirectory() string {
	return filepath.Join(os.TempDir(), "itamae-logs")
}

// ListLogFiles returns a sorted list of log files (newest first)
func ListLogFiles() ([]string, error) {
	logDir := GetLogDirectory()

	// Check if directory exists
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		return []string{}, nil // No logs yet
	}

	// Read directory
	files, err := os.ReadDir(logDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read log directory: %w", err)
	}

	// Collect log files with their modification times
	type logFile struct {
		path    string
		modTime time.Time
	}
	var logFiles []logFile

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".log") {
			continue
		}

		fullPath := filepath.Join(logDir, file.Name())
		info, err := os.Stat(fullPath)
		if err != nil {
			continue
		}

		logFiles = append(logFiles, logFile{
			path:    fullPath,
			modTime: info.ModTime(),
		})
	}

	// Sort by modification time (newest first)
	sort.Slice(logFiles, func(i, j int) bool {
		return logFiles[i].modTime.After(logFiles[j].modTime)
	})

	// Extract paths
	result := make([]string, len(logFiles))
	for i, lf := range logFiles {
		result[i] = lf.path
	}

	return result, nil
}

// GetMostRecentLog returns the path to the most recent log file
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
