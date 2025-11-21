package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/yjmrobert/itamae/itamae"
)

var (
	listLogs   bool
	logFile    string
	followLog  bool
	grepFilter string
	cleanLogs  bool
)

// Tokyo Night colors for log output
var (
	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#bb9af7")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f7768e")).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9ece6a"))

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#e0af68"))

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7aa2f7"))

	timestampStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#565f89"))

	pathStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7dcfff")).
			Underline(true)

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#565f89")).
			Italic(true)
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "View installation logs",
	Long: `View installation logs from previous itamae runs.

By default, displays the most recent log file.

Examples:
  itamae logs                    # Show most recent log
  itamae logs --list             # List all available logs
  itamae logs --file <filename>  # View specific log file
  itamae logs --follow           # Follow log in real-time (tail -f)
  itamae logs --grep error       # Filter log lines containing "error"
  itamae logs --clean            # Remove all old logs`,
	Run: runLogs,
}

func init() {
	logsCmd.Flags().BoolVarP(&listLogs, "list", "l", false, "List all available log files")
	logsCmd.Flags().StringVarP(&logFile, "file", "f", "", "View specific log file")
	logsCmd.Flags().BoolVar(&followLog, "follow", false, "Follow log output in real-time (like tail -f)")
	logsCmd.Flags().StringVarP(&grepFilter, "grep", "g", "", "Filter log lines containing this text")
	logsCmd.Flags().BoolVar(&cleanLogs, "clean", false, "Remove all old log files")
	rootCmd.AddCommand(logsCmd)
}

func runLogs(cmd *cobra.Command, args []string) {
	// Handle cleanup
	if cleanLogs {
		cleanupLogs()
		return
	}

	// Handle list
	if listLogs {
		listAllLogs()
		return
	}

	// Determine which log file to show
	var logPath string
	var err error

	if logFile != "" {
		// User specified a file
		logPath = logFile
		if !filepath.IsAbs(logPath) {
			// If not absolute, try to find it in the log directory
			logPath = filepath.Join(itamae.GetLogDirectory(), logFile)
		}

		if _, err := os.Stat(logPath); os.IsNotExist(err) {
			fmt.Printf("%s Log file not found: %s\n", errorStyle.Render("✗"), logPath)
			fmt.Printf("%s Run 'itamae logs --list' to see available logs\n", dimStyle.Render("Tip:"))
			os.Exit(1)
		}
	} else {
		// Get most recent log
		logPath, err = itamae.GetMostRecentLog()
		if err != nil {
			fmt.Printf("%s %s\n", errorStyle.Render("✗"), err.Error())
			fmt.Printf("%s No installation logs found. Run 'itamae install' first.\n", dimStyle.Render("Tip:"))
			os.Exit(1)
		}
	}

	// Display the log
	if followLog {
		followLogFile(logPath)
	} else {
		displayLogFile(logPath)
	}
}

func listAllLogs() {
	logs, err := itamae.ListLogFiles()
	if err != nil {
		fmt.Printf("%s Failed to list logs: %s\n", errorStyle.Render("✗"), err.Error())
		os.Exit(1)
	}

	if len(logs) == 0 {
		fmt.Printf("%s No installation logs found\n", warningStyle.Render("⚠"))
		fmt.Printf("%s Run 'itamae install' to create logs\n", dimStyle.Render("Tip:"))
		return
	}

	fmt.Println(headerStyle.Render(fmt.Sprintf("📋 Found %d log file(s):", len(logs))))
	fmt.Println()

	for i, log := range logs {
		info, err := os.Stat(log)
		if err != nil {
			continue
		}

		// Format relative time
		timeSince := time.Since(info.ModTime())
		var timeStr string
		if timeSince < time.Minute {
			timeStr = "just now"
		} else if timeSince < time.Hour {
			timeStr = fmt.Sprintf("%d minutes ago", int(timeSince.Minutes()))
		} else if timeSince < 24*time.Hour {
			timeStr = fmt.Sprintf("%d hours ago", int(timeSince.Hours()))
		} else {
			timeStr = fmt.Sprintf("%d days ago", int(timeSince.Hours()/24))
		}

		marker := " "
		if i == 0 {
			marker = successStyle.Render("→") // Most recent
		}

		filename := filepath.Base(log)
		sizeKB := info.Size() / 1024

		fmt.Printf("%s %s %s %s (%d KB)\n",
			marker,
			timestampStyle.Render(info.ModTime().Format("2006-01-02 15:04:05")),
			pathStyle.Render(filename),
			dimStyle.Render(timeStr),
			sizeKB,
		)
	}

	fmt.Println()
	fmt.Printf("%s View a log: itamae logs --file <filename>\n", dimStyle.Render("Tip:"))
}

func displayLogFile(logPath string) {
	file, err := os.Open(logPath)
	if err != nil {
		fmt.Printf("%s Failed to open log: %s\n", errorStyle.Render("✗"), err.Error())
		os.Exit(1)
	}
	defer file.Close()

	// Print header
	info, _ := os.Stat(logPath)
	fmt.Println(headerStyle.Render("📋 Installation Log"))
	fmt.Printf("%s %s\n", dimStyle.Render("File:"), pathStyle.Render(logPath))
	if info != nil {
		fmt.Printf("%s %s (%d KB)\n",
			dimStyle.Render("Date:"),
			info.ModTime().Format("2006-01-02 15:04:05"),
			info.Size()/1024,
		)
	}
	fmt.Println(strings.Repeat("─", 80))
	fmt.Println()

	// Check if we should use a pager
	usePager := shouldUsePager(logPath)

	if usePager && grepFilter == "" {
		// Use less for viewing
		cmd := exec.Command("less", "-R", logPath)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		return
	}

	// Read and display line by line
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		// Apply grep filter if specified
		if grepFilter != "" && !strings.Contains(strings.ToLower(line), strings.ToLower(grepFilter)) {
			continue
		}

		// Colorize the line
		coloredLine := colorizeLine(line)
		fmt.Println(coloredLine)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("\n%s Error reading log: %s\n", errorStyle.Render("✗"), err.Error())
	}

	// Print footer
	fmt.Println()
	fmt.Println(strings.Repeat("─", 80))
	fmt.Printf("%s Use ↑/↓ to scroll, 'q' to quit (if using pager)\n", dimStyle.Render("Tip:"))
	if grepFilter == "" {
		fmt.Printf("%s Filter results: itamae logs --grep <text>\n", dimStyle.Render("Tip:"))
	}
}

func followLogFile(logPath string) {
	// Check if file exists
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		fmt.Printf("%s Log file not found: %s\n", errorStyle.Render("✗"), logPath)
		os.Exit(1)
	}

	// Print header
	fmt.Println(headerStyle.Render("📋 Following Installation Log"))
	fmt.Printf("%s %s\n", dimStyle.Render("File:"), pathStyle.Render(logPath))
	fmt.Printf("%s Press Ctrl+C to stop\n", dimStyle.Render("Tip:"))
	fmt.Println(strings.Repeat("─", 80))
	fmt.Println()

	// Use tail -f
	cmd := exec.Command("tail", "-f", logPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("\n%s Failed to follow log: %s\n", errorStyle.Render("✗"), err.Error())
		os.Exit(1)
	}
}

func cleanupLogs() {
	logDir := itamae.GetLogDirectory()

	logs, err := itamae.ListLogFiles()
	if err != nil {
		fmt.Printf("%s Failed to list logs: %s\n", errorStyle.Render("✗"), err.Error())
		os.Exit(1)
	}

	if len(logs) == 0 {
		fmt.Printf("%s No logs to clean\n", infoStyle.Render("ℹ"))
		return
	}

	fmt.Printf("%s About to delete %d log file(s) from %s\n",
		warningStyle.Render("⚠"),
		len(logs),
		logDir,
	)
	fmt.Print("Continue? (y/N): ")

	var response string
	fmt.Scanln(&response)

	if strings.ToLower(response) != "y" {
		fmt.Println("Cancelled.")
		return
	}

	// Delete all logs
	removed := 0
	for _, log := range logs {
		if err := os.Remove(log); err != nil {
			fmt.Printf("%s Failed to remove %s: %s\n", errorStyle.Render("✗"), filepath.Base(log), err.Error())
		} else {
			removed++
		}
	}

	fmt.Printf("%s Removed %d log file(s)\n", successStyle.Render("✓"), removed)
}

func colorizeLine(line string) string {
	// Check for log level indicators
	if strings.Contains(line, "ERROR:") {
		return errorStyle.Render(line)
	}
	if strings.Contains(line, "Installation successful") || strings.Contains(line, "✅") {
		return successStyle.Render(line)
	}
	if strings.Contains(line, "Phase") {
		return infoStyle.Render(line)
	}
	if strings.Contains(line, "Setting up") {
		return infoStyle.Render(line)
	}

	// Default: check for timestamp and colorize it
	if strings.HasPrefix(line, "[") && len(line) > 13 {
		// Extract timestamp (e.g., [15:04:05.000])
		if line[12] == ']' {
			timestamp := line[:13]
			rest := line[13:]
			return timestampStyle.Render(timestamp) + rest
		}
	}

	return line
}

func shouldUsePager(logPath string) bool {
	// Don't use pager if output is not a terminal
	if !isTerminal() {
		return false
	}

	// Check file size - use pager if > 50 lines
	file, err := os.Open(logPath)
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := 0
	for scanner.Scan() {
		lines++
		if lines > 50 {
			return true
		}
	}

	return false
}

func isTerminal() bool {
	fileInfo, _ := os.Stdout.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}
