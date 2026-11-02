package itamae

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
)

// Sender receives progress messages from an Installer. *tea.Program satisfies it,
// as does PlainSender for non-interactive runs.
type Sender interface {
	Send(msg tea.Msg)
}

// InstallResult summarises an installation run.
type InstallResult struct {
	Successful []string // Names of tools that installed and passed their check
	Failed     []string // Names of tools that failed to install or verify
	Skipped    []string // Names of tools that need manual steps
}

// OK reports whether every tool installed successfully.
func (r InstallResult) OK() bool {
	return len(r.Failed) == 0
}

// prerequisites are the packages the repository setup scripts rely on.
var prerequisites = []struct {
	pkg   string
	check func() bool
}{
	{"curl", func() bool { return commandExists("curl") }},
	{"gnupg", func() bool { return commandExists("gpg") }},
	{"ca-certificates", func() bool { return fileExists("/etc/ssl/certs/ca-certificates.crt") }},
}

// Installer drives a full installation: repository setup, batch APT install,
// individual installs and verification. Progress is reported to Sender.
type Installer struct {
	Plugins []ToolPlugin
	Inputs  map[string]string
	Sender  Sender
	// Output receives the raw output of scripts and package manager commands.
	// Leave nil to keep it in the debug log only (TUI mode).
	Output io.Writer
	// SkipVerify disables running each plugin's check() after installation.
	SkipVerify bool

	ctx context.Context

	// lookPath is a hook for tests; nil means exec.LookPath.
	lookPath func(string) (string, error)
	// prereqs is a hook for tests; nil means the real prerequisite probes.
	prereqs func() []string
}

func (in *Installer) send(msg tea.Msg) {
	if in.Sender != nil {
		in.Sender.Send(msg)
	}
}

func (in *Installer) log(level, pkg, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	DebugLog("%s: %s%s", strings.ToUpper(level), prefixPkg(pkg), msg)
	in.send(LogMsg{Level: level, Package: pkg, Message: msg})
}

func prefixPkg(pkg string) string {
	if pkg == "" {
		return ""
	}
	return pkg + ": "
}

// packageManager returns the APT front end to use: nala when present, apt-get otherwise.
func (in *Installer) packageManager() string {
	look := in.lookPath
	if look == nil {
		look = exec.LookPath
	}
	if _, err := look("nala"); err == nil {
		return "nala"
	}
	return "apt-get"
}

func (in *Installer) cancelled() bool {
	return in.ctx.Err() != nil
}

// fail records a failed plugin.
func (in *Installer) fail(plugin ToolPlugin, phase, reason string, result *InstallResult) {
	in.log("error", plugin.ID, "%s", reason)
	in.send(ErrorMsg{Package: plugin.ID, Phase: phase, Message: reason})
	in.send(PackageCompleteMsg{PackageID: plugin.ID, Success: false, Error: firstLine(reason)})
	result.Failed = append(result.Failed, plugin.Name)
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}

// Run performs the installation and returns the result. Script failures never
// abort the run; every failure is reported through Sender and the result.
// Cancelling ctx stops the run after the current command.
func (in *Installer) Run(ctx context.Context) InstallResult {
	if ctx == nil {
		ctx = context.Background()
	}
	in.ctx = ctx

	var result InstallResult
	stopKeepalive := in.startSudoKeepalive()
	defer stopKeepalive()

	var aptPlugins, otherPlugins []ToolPlugin
	for _, plugin := range in.Plugins {
		if plugin.InstallMethod == "apt" {
			aptPlugins = append(aptPlugins, plugin)
		} else {
			otherPlugins = append(otherPlugins, plugin)
		}
	}
	DebugLog("Installer started: %d plugins (%d apt, %d other), package manager: %s",
		len(in.Plugins), len(aptPlugins), len(otherPlugins), in.packageManager())

	if len(aptPlugins) > 0 {
		if in.prepareApt(aptPlugins, &result) {
			in.installApt(aptPlugins, &result)
		}
	}

	if len(otherPlugins) > 0 {
		in.installIndividually(otherPlugins, &result)
	}

	DebugLog("Installer finished: %d succeeded, %d failed, %d skipped",
		len(result.Successful), len(result.Failed), len(result.Skipped))
	in.send(SummaryMsg(result))
	return result
}

// prepareApt runs Phase 0: prerequisites, custom repositories and the package list update.
// It returns false when the APT phase cannot continue; all APT plugins are then marked failed.
func (in *Installer) prepareApt(aptPlugins []ToolPlugin, result *InstallResult) bool {
	pm := in.packageManager()

	var repoPlugins []ToolPlugin
	for _, plugin := range aptPlugins {
		if plugin.RepoSetup != "" {
			repoPlugins = append(repoPlugins, plugin)
		}
	}

	failAll := func(reason string) bool {
		in.log("error", "", "%s", reason)
		in.send(ErrorMsg{Phase: "repo_setup", Message: reason})
		for _, plugin := range aptPlugins {
			in.send(PackageCompleteMsg{PackageID: plugin.ID, Success: false, Error: firstLine(reason)})
			result.Failed = append(result.Failed, plugin.Name)
		}
		return false
	}

	in.send(PhaseStartMsg{Phase: "repo_setup", Count: len(repoPlugins)})

	if len(repoPlugins) > 0 {
		if missing := in.missingPrerequisites(); len(missing) > 0 {
			in.log("info", "", "Installing prerequisites: %s", strings.Join(missing, ", "))
			if err := runSudo(in.ctx, in.Output, "apt-get", "update"); err != nil {
				return failAll(fmt.Sprintf("Package list update failed: %v", err))
			}
			args := append([]string{"apt-get", "install", "-y"}, missing...)
			if err := runSudo(in.ctx, in.Output, args...); err != nil {
				return failAll(fmt.Sprintf("Prerequisite installation failed: %v", err))
			}
		}

		for _, plugin := range repoPlugins {
			if in.cancelled() {
				return failAll("Installation cancelled")
			}
			in.send(PackageStartMsg{PackageID: plugin.ID, Phase: "repo_setup"})
			in.log("info", plugin.ID, "Setting up custom repository...")
			if err := runScript(in.ctx, plugin, plugin.RepoSetup, in.Inputs, in.Output); err != nil {
				return failAll(fmt.Sprintf("Repository setup failed for %s: %v", plugin.Name, err))
			}
			in.send(ProgressMsg{PackageID: plugin.ID, Progress: "Repository configured"})
		}
	}

	if in.cancelled() {
		return failAll("Installation cancelled")
	}
	in.log("info", "", "Updating package lists with %s...", pm)
	if err := runSudo(in.ctx, in.Output, pm, "update"); err != nil {
		return failAll(fmt.Sprintf("Package list update failed: %v", err))
	}
	in.log("success", "", "Package lists updated")
	in.send(PhaseCompleteMsg{Phase: "repo_setup"})
	return true
}

// installApt runs Phase 1: one batched install, falling back to per-package installs
// when the batch fails, followed by post-install hooks and verification.
func (in *Installer) installApt(aptPlugins []ToolPlugin, result *InstallResult) {
	pm := in.packageManager()
	in.send(PhaseStartMsg{Phase: "apt_batch", Count: len(aptPlugins)})

	packages := make([]string, 0, len(aptPlugins))
	for _, plugin := range aptPlugins {
		in.send(PackageStartMsg{PackageID: plugin.ID, Phase: "install"})
		packages = append(packages, plugin.PackageName)
	}

	installed := make(map[string]bool, len(aptPlugins))
	in.log("info", "", "Installing %d APT packages with %s...", len(packages), pm)
	args := append([]string{pm, "install", "-y"}, packages...)
	if err := runSudo(in.ctx, in.Output, args...); err == nil {
		in.log("success", "", "Installed %d APT packages", len(packages))
		for _, plugin := range aptPlugins {
			installed[plugin.ID] = true
		}
	} else {
		in.log("warning", "", "Batch install failed, retrying packages one at a time: %v", err)
		for _, plugin := range aptPlugins {
			if in.cancelled() {
				in.fail(plugin, "install", "Installation cancelled", result)
				continue
			}
			in.send(ProgressMsg{PackageID: plugin.ID, Progress: "Installing individually..."})
			if err := runSudo(in.ctx, in.Output, pm, "install", "-y", plugin.PackageName); err != nil {
				in.fail(plugin, "install", fmt.Sprintf("Install failed: %v", err), result)
				continue
			}
			installed[plugin.ID] = true
		}
	}

	for _, plugin := range aptPlugins {
		if !installed[plugin.ID] {
			continue
		}
		if plugin.PostInstall != "" {
			in.send(PackageStartMsg{PackageID: plugin.ID, Phase: "post_install"})
			in.log("info", plugin.ID, "Running post-install tasks...")
			if err := runScript(in.ctx, plugin, plugin.PostInstall, in.Inputs, in.Output); err != nil {
				in.fail(plugin, "post_install", fmt.Sprintf("Post-install failed: %v", err), result)
				continue
			}
		}
		in.finish(plugin, result)
	}

	in.send(PhaseCompleteMsg{Phase: "apt_batch"})
}

// installIndividually runs Phase 2 for binary and manual plugins.
func (in *Installer) installIndividually(plugins []ToolPlugin, result *InstallResult) {
	in.send(PhaseStartMsg{Phase: "individual", Count: len(plugins)})

	for _, plugin := range plugins {
		in.send(PackageStartMsg{PackageID: plugin.ID, Phase: "install"})
		if in.cancelled() {
			in.fail(plugin, "install", "Installation cancelled", result)
			continue
		}
		in.log("info", plugin.ID, "Installing...")

		if err := runScript(in.ctx, plugin, "install", in.Inputs, in.Output); err != nil {
			in.fail(plugin, "install", fmt.Sprintf("Install failed: %v", err), result)
			continue
		}

		if plugin.InstallMethod == "manual" {
			in.log("warning", plugin.ID, "Requires manual installation steps, see the log")
			in.send(PackageCompleteMsg{PackageID: plugin.ID, Success: true, Skipped: true, Error: "manual steps required"})
			result.Skipped = append(result.Skipped, plugin.Name)
			continue
		}

		in.finish(plugin, result)
	}

	in.send(PhaseCompleteMsg{Phase: "individual"})
}

// finish verifies an installed plugin and records the outcome.
func (in *Installer) finish(plugin ToolPlugin, result *InstallResult) {
	if !in.SkipVerify {
		in.send(PackageStartMsg{PackageID: plugin.ID, Phase: "verify"})
		if !checkInstalled(in.ctx, plugin, in.Inputs) {
			in.fail(plugin, "verify", "Installed, but the post-install check failed", result)
			return
		}
	}
	in.log("success", plugin.ID, "Installed")
	in.send(PackageCompleteMsg{PackageID: plugin.ID, Success: true})
	result.Successful = append(result.Successful, plugin.Name)
}

func (in *Installer) missingPrerequisites() []string {
	if in.prereqs != nil {
		return in.prereqs()
	}
	var missing []string
	for _, p := range prerequisites {
		if !p.check() {
			missing = append(missing, p.pkg)
		}
	}
	return missing
}

// startSudoKeepalive refreshes the sudo timestamp while a long install runs so the
// scripts never block on a password prompt. It returns a stop function.
func (in *Installer) startSudoKeepalive() func() {
	if os.Geteuid() == 0 || !commandExists("sudo") {
		return func() {}
	}
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				_ = exec.Command("sudo", "-n", "-v").Run()
			}
		}
	}()
	return func() { close(done) }
}

// CheckInstalled runs a plugin's check() function and reports whether the tool is present.
func CheckInstalled(ctx context.Context, plugin ToolPlugin) bool {
	return checkInstalled(ctx, plugin, nil)
}

func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
