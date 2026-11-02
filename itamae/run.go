package itamae

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/colorprofile"
)

// ErrCancelled is returned when the user aborts an installation from the TUI.
var ErrCancelled = errors.New("installation cancelled")

// EnsureSudoAccess prompts for sudo credentials up front so the installation
// never blocks on a password prompt half way through.
func EnsureSudoAccess() error {
	if os.Geteuid() == 0 {
		return nil
	}
	if !commandExists("sudo") {
		return errors.New("sudo is required to install packages but was not found on PATH")
	}

	fmt.Fprintln(os.Stderr, "🔐 Requesting sudo access for installation...")
	cmd := exec.Command("sudo", "-v")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to obtain sudo access: %w", err)
	}
	return nil
}

// RunInstallTUI installs the plugins while showing the multi-pane TUI.
func RunInstallTUI(ctx context.Context, plugins []ToolPlugin, inputs map[string]string) (InstallResult, error) {
	DebugLog("RunInstallTUI started with %d plugins", len(plugins))

	model := NewInstallModel(plugins)
	p := tea.NewProgram(model, tea.WithContext(ctx))

	results := make(chan InstallResult, 1)
	go func() {
		installer := &Installer{Plugins: plugins, Inputs: inputs, Sender: p}
		results <- installer.Run(ctx)
	}()

	finalModel, err := p.Run()
	if err != nil {
		return InstallResult{}, fmt.Errorf("TUI failed: %w", err)
	}
	if m, ok := finalModel.(InstallModel); ok && !m.complete {
		DebugLog("TUI closed before the installation finished")
		return InstallResult{}, ErrCancelled
	}

	DebugLog("TUI program exited normally")
	return <-results, nil
}

// RunInstallPlain installs the plugins with line-oriented output suitable for
// logs, CI and non-interactive terminals.
func RunInstallPlain(ctx context.Context, plugins []ToolPlugin, inputs map[string]string, out io.Writer) InstallResult {
	DebugLog("RunInstallPlain started with %d plugins", len(plugins))
	sender := NewPlainSender(out, plugins)
	installer := &Installer{Plugins: plugins, Inputs: inputs, Sender: sender, Output: sender.indented()}
	return installer.Run(ctx)
}

// PlainSender renders installer messages as plain, colour-aware lines.
type PlainSender struct {
	w     *colorprofile.Writer
	names map[string]string
}

// NewPlainSender creates a PlainSender writing to w. Colours degrade to what the
// destination supports and disappear entirely when w is not a terminal.
func NewPlainSender(w io.Writer, plugins []ToolPlugin) *PlainSender {
	names := make(map[string]string, len(plugins))
	for _, p := range plugins {
		names[p.ID] = p.Name
	}
	return &PlainSender{w: colorprofile.NewWriter(w, os.Environ()), names: names}
}

// indented returns a writer that indents raw script output under the status lines.
func (s *PlainSender) indented() io.Writer {
	return indentWriter{w: s.w, prefix: "      "}
}

func (s *PlainSender) name(id string) string {
	if n, ok := s.names[id]; ok {
		return n
	}
	return id
}

// Send implements Sender.
func (s *PlainSender) Send(msg tea.Msg) {
	dim := lipgloss.NewStyle().Foreground(TokyoNightComment)
	bold := lipgloss.NewStyle().Bold(true).Foreground(TokyoNightMagenta)
	pkg := lipgloss.NewStyle().Foreground(TokyoNightBlue)

	switch msg := msg.(type) {
	case PhaseStartMsg:
		fmt.Fprintf(s.w, "\n%s %s\n", bold.Render("▶ Phase:"), bold.Render(fmt.Sprintf("%s (%d items)", msg.Phase, msg.Count)))
	case PhaseCompleteMsg:
		fmt.Fprintf(s.w, "%s\n", dim.Render(fmt.Sprintf("  ✓ phase %s complete", msg.Phase)))
	case PackageStartMsg:
		fmt.Fprintf(s.w, "  %s %s %s\n", dim.Render("→"), pkg.Render(s.name(msg.PackageID)), dim.Render(phaseLabel(msg.Phase)))
	case ProgressMsg:
		fmt.Fprintf(s.w, "  %s %s %s\n", dim.Render("·"), pkg.Render(s.name(msg.PackageID)), dim.Render(msg.Progress))
	case PackageCompleteMsg:
		switch {
		case msg.Skipped:
			fmt.Fprintf(s.w, "  %s %s %s\n", lipgloss.NewStyle().Foreground(StatusSkipped).Render("⊘"), pkg.Render(s.name(msg.PackageID)), dim.Render(msg.Error))
		case msg.Success:
			fmt.Fprintf(s.w, "  %s %s\n", lipgloss.NewStyle().Foreground(StatusSuccess).Render("✓"), pkg.Render(s.name(msg.PackageID)))
		default:
			fmt.Fprintf(s.w, "  %s %s %s\n", lipgloss.NewStyle().Foreground(StatusError).Render("✗"), pkg.Render(s.name(msg.PackageID)), lipgloss.NewStyle().Foreground(StatusError).Render(msg.Error))
		}
	case LogMsg:
		line := msg.Message
		if msg.Package != "" {
			line = pkg.Render(s.name(msg.Package)) + " → " + line
		}
		fmt.Fprintf(s.w, "  %s %s %s\n",
			dim.Render(time.Now().Format("15:04:05")),
			lipgloss.NewStyle().Foreground(LevelColor(msg.Level)).Render(LevelIcon(msg.Level)),
			lipgloss.NewStyle().Foreground(LevelColor(msg.Level)).Render(line))
	case ErrorMsg, SummaryMsg:
		// Errors are already logged; the summary is printed by PrintSummary.
	}
}

// PrintSummary writes the final installation summary to w.
func PrintSummary(w io.Writer, result InstallResult) {
	cw := colorprofile.NewWriter(w, os.Environ())
	success := lipgloss.NewStyle().Foreground(StatusSuccess).Bold(true)
	failure := lipgloss.NewStyle().Foreground(StatusError).Bold(true)
	skipped := lipgloss.NewStyle().Foreground(StatusSkipped).Bold(true)
	rule := lipgloss.NewStyle().Foreground(TokyoNightComment).Render(strings.Repeat("─", 60))

	fmt.Fprintln(cw)
	fmt.Fprintln(cw, rule)
	fmt.Fprintln(cw, lipgloss.NewStyle().Bold(true).Foreground(TokyoNightMagenta).Render("📊 Installation summary"))
	fmt.Fprintln(cw, rule)
	if len(result.Successful) > 0 {
		fmt.Fprintln(cw, success.Render(fmt.Sprintf("✓ Installed (%d):", len(result.Successful))))
		for _, name := range result.Successful {
			fmt.Fprintf(cw, "   • %s\n", name)
		}
	}
	if len(result.Skipped) > 0 {
		fmt.Fprintln(cw, skipped.Render(fmt.Sprintf("⊘ Manual steps required (%d):", len(result.Skipped))))
		for _, name := range result.Skipped {
			fmt.Fprintf(cw, "   • %s\n", name)
		}
	}
	if len(result.Failed) > 0 {
		fmt.Fprintln(cw, failure.Render(fmt.Sprintf("✗ Failed (%d):", len(result.Failed))))
		for _, name := range result.Failed {
			fmt.Fprintf(cw, "   • %s\n", name)
		}
	} else {
		fmt.Fprintln(cw, success.Render("🎉 All packages installed successfully!"))
	}
	fmt.Fprintln(cw, rule)
}

// indentWriter prefixes every line written through it.
type indentWriter struct {
	w      io.Writer
	prefix string
}

func (iw indentWriter) Write(p []byte) (int, error) {
	text := strings.TrimRight(string(p), "\n")
	if strings.TrimSpace(text) == "" {
		return len(p), nil
	}
	for _, line := range strings.Split(text, "\n") {
		if _, err := fmt.Fprintf(iw.w, "%s%s\n", iw.prefix, strings.TrimRight(line, "\r")); err != nil {
			return 0, err
		}
	}
	return len(p), nil
}
