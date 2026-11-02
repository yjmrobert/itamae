package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/colorprofile"
	"github.com/spf13/cobra"
	"github.com/yjmrobert/itamae/itamae"
)

var logsFlags struct {
	list   bool
	file   string
	follow bool
	grep   string
	clean  bool
}

// Tokyo Night colors for log output
var (
	headerStyle    = lipgloss.NewStyle().Foreground(itamae.TokyoNightMagenta).Bold(true)
	errorStyle     = lipgloss.NewStyle().Foreground(itamae.TokyoNightRed).Bold(true)
	successStyle   = lipgloss.NewStyle().Foreground(itamae.TokyoNightGreen)
	warningStyle   = lipgloss.NewStyle().Foreground(itamae.TokyoNightYellow)
	infoStyle      = lipgloss.NewStyle().Foreground(itamae.TokyoNightBlue)
	timestampStyle = lipgloss.NewStyle().Foreground(itamae.TokyoNightComment)
	pathStyle      = lipgloss.NewStyle().Foreground(itamae.TokyoNightCyan).Underline(true)
	dimStyle       = lipgloss.NewStyle().Foreground(itamae.TokyoNightComment).Italic(true)
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "View installation logs",
	Long: `View installation logs from previous itamae runs.

By default, displays the most recent log file.`,
	Example: `  itamae logs                    # Show most recent log
  itamae logs --list             # List all available logs
  itamae logs --file <filename>  # View specific log file
  itamae logs --follow           # Follow log in real-time (tail -f)
  itamae logs --grep error       # Filter log lines containing "error"
  itamae logs --clean            # Remove all old logs`,
	Args: cobra.NoArgs,
	RunE: runLogs,
}

func init() {
	f := logsCmd.Flags()
	f.BoolVarP(&logsFlags.list, "list", "l", false, "List all available log files")
	f.StringVarP(&logsFlags.file, "file", "f", "", "View specific log file")
	f.BoolVar(&logsFlags.follow, "follow", false, "Follow log output in real-time (like tail -f)")
	f.StringVarP(&logsFlags.grep, "grep", "g", "", "Filter log lines containing this text")
	f.BoolVar(&logsFlags.clean, "clean", false, "Remove all old log files")
	rootCmd.AddCommand(logsCmd)
}

func runLogs(cmd *cobra.Command, _ []string) error {
	out := colorprofile.NewWriter(cmd.OutOrStdout(), os.Environ())

	if logsFlags.clean {
		return cleanupLogs(out, cmd.InOrStdin())
	}
	if logsFlags.list {
		return listAllLogs(out)
	}

	var logPath string
	if logsFlags.file != "" {
		logPath = logsFlags.file
		if !filepath.IsAbs(logPath) {
			logPath = filepath.Join(itamae.GetLogDirectory(), logsFlags.file)
		}
		if _, err := os.Stat(logPath); err != nil {
			return fmt.Errorf("log file not found: %s (run 'itamae logs --list' to see available logs)", logPath)
		}
	} else {
		var err error
		logPath, err = itamae.GetMostRecentLog()
		if err != nil {
			return fmt.Errorf("%w (run 'itamae install' first)", err)
		}
	}

	if logsFlags.follow {
		return followLogFile(out, logPath)
	}
	return displayLogFile(out, logPath)
}

func listAllLogs(out io.Writer) error {
	logs, err := itamae.ListLogFiles()
	if err != nil {
		return fmt.Errorf("failed to list logs: %w", err)
	}

	if len(logs) == 0 {
		fmt.Fprintf(out, "%s No installation logs found\n", warningStyle.Render("⚠"))
		fmt.Fprintf(out, "%s Run 'itamae install' to create logs\n", dimStyle.Render("Tip:"))
		return nil
	}

	fmt.Fprintln(out, headerStyle.Render(fmt.Sprintf("📋 Found %d log file(s):", len(logs))))
	fmt.Fprintln(out)

	for i, log := range logs {
		info, err := os.Stat(log)
		if err != nil {
			continue
		}

		marker := " "
		if i == 0 {
			marker = successStyle.Render("→") // Most recent
		}

		fmt.Fprintf(out, "%s %s %s %s (%d KB)\n",
			marker,
			timestampStyle.Render(info.ModTime().Format("2006-01-02 15:04:05")),
			pathStyle.Render(filepath.Base(log)),
			dimStyle.Render(relativeTime(time.Since(info.ModTime()))),
			info.Size()/1024,
		)
	}

	fmt.Fprintln(out)
	fmt.Fprintf(out, "%s View a log: itamae logs --file <filename>\n", dimStyle.Render("Tip:"))
	return nil
}

// relativeTime formats a duration as a coarse "N units ago" string.
func relativeTime(since time.Duration) string {
	switch {
	case since < time.Minute:
		return "just now"
	case since < time.Hour:
		return fmt.Sprintf("%d minutes ago", int(since.Minutes()))
	case since < 24*time.Hour:
		return fmt.Sprintf("%d hours ago", int(since.Hours()))
	default:
		return fmt.Sprintf("%d days ago", int(since.Hours()/24))
	}
}

func displayLogFile(out io.Writer, logPath string) error {
	file, err := os.Open(logPath)
	if err != nil {
		return fmt.Errorf("failed to open log: %w", err)
	}
	defer file.Close()

	fmt.Fprintln(out, headerStyle.Render("📋 Installation Log"))
	fmt.Fprintf(out, "%s %s\n", dimStyle.Render("File:"), pathStyle.Render(logPath))
	if info, err := os.Stat(logPath); err == nil {
		fmt.Fprintf(out, "%s %s (%d KB)\n",
			dimStyle.Render("Date:"),
			info.ModTime().Format("2006-01-02 15:04:05"),
			info.Size()/1024,
		)
	}
	fmt.Fprintln(out, strings.Repeat("─", 80))
	fmt.Fprintln(out)

	if shouldUsePager(logPath) && logsFlags.grep == "" {
		pager := exec.Command("less", "-R", logPath)
		pager.Stdin = os.Stdin
		pager.Stdout = os.Stdout
		pager.Stderr = os.Stderr
		if err := pager.Run(); err != nil {
			return fmt.Errorf("pager failed: %w", err)
		}
		return nil
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if logsFlags.grep != "" && !strings.Contains(strings.ToLower(line), strings.ToLower(logsFlags.grep)) {
			continue
		}
		fmt.Fprintln(out, colorizeLine(line))
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading log: %w", err)
	}

	fmt.Fprintln(out)
	fmt.Fprintln(out, strings.Repeat("─", 80))
	if logsFlags.grep == "" {
		fmt.Fprintf(out, "%s Filter results: itamae logs --grep <text>\n", dimStyle.Render("Tip:"))
	}
	return nil
}

func followLogFile(out io.Writer, logPath string) error {
	if _, err := os.Stat(logPath); err != nil {
		return fmt.Errorf("log file not found: %s", logPath)
	}

	fmt.Fprintln(out, headerStyle.Render("📋 Following Installation Log"))
	fmt.Fprintf(out, "%s %s\n", dimStyle.Render("File:"), pathStyle.Render(logPath))
	fmt.Fprintf(out, "%s Press Ctrl+C to stop\n", dimStyle.Render("Tip:"))
	fmt.Fprintln(out, strings.Repeat("─", 80))
	fmt.Fprintln(out)

	tail := exec.Command("tail", "-f", logPath)
	tail.Stdout = os.Stdout
	tail.Stderr = os.Stderr
	if err := tail.Run(); err != nil {
		return fmt.Errorf("failed to follow log: %w", err)
	}
	return nil
}

func cleanupLogs(out io.Writer, in io.Reader) error {
	logs, err := itamae.ListLogFiles()
	if err != nil {
		return fmt.Errorf("failed to list logs: %w", err)
	}

	if len(logs) == 0 {
		fmt.Fprintf(out, "%s No logs to clean\n", infoStyle.Render("ℹ"))
		return nil
	}

	fmt.Fprintf(out, "%s About to delete %d log file(s) from %s\n",
		warningStyle.Render("⚠"), len(logs), itamae.GetLogDirectory())
	fmt.Fprint(out, "Continue? (y/N): ")

	var response string
	_, _ = fmt.Fscanln(in, &response)
	if strings.ToLower(strings.TrimSpace(response)) != "y" {
		fmt.Fprintln(out, "Cancelled.")
		return nil
	}

	removed := 0
	for _, log := range logs {
		if err := os.Remove(log); err != nil {
			fmt.Fprintf(out, "%s Failed to remove %s: %s\n", errorStyle.Render("✗"), filepath.Base(log), err.Error())
		} else {
			removed++
		}
	}

	fmt.Fprintf(out, "%s Removed %d log file(s)\n", successStyle.Render("✓"), removed)
	return nil
}

// colorizeLine applies a style to a debug log line based on its content.
func colorizeLine(line string) string {
	switch {
	case strings.Contains(line, "ERROR"):
		return errorStyle.Render(line)
	case strings.Contains(line, "SUCCESS"), strings.Contains(line, "Installation successful"), strings.Contains(line, "✅"):
		return successStyle.Render(line)
	case strings.Contains(line, "WARNING"):
		return warningStyle.Render(line)
	case strings.Contains(line, "Phase"), strings.Contains(line, "Setting up"):
		return infoStyle.Render(line)
	}

	// Default: dim the leading timestamp, e.g. [15:04:05.000]
	if strings.HasPrefix(line, "[") && len(line) > 13 && line[13] == ']' {
		return timestampStyle.Render(line[:14]) + line[14:]
	}
	return line
}

// shouldUsePager reports whether the log is long enough to warrant a pager.
func shouldUsePager(logPath string) bool {
	if !itamae.IsInteractive() {
		return false
	}

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
