package itamae

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

// sudoEnv is how every package manager call shows up in the mock sudo log.
const sudoEnv = "sudo env DEBIAN_FRONTEND=noninteractive "

// recordingSender collects every message the installer sends.
type recordingSender struct {
	msgs []tea.Msg
}

func (r *recordingSender) Send(msg tea.Msg) { r.msgs = append(r.msgs, msg) }

func (r *recordingSender) completions() map[string]PackageCompleteMsg {
	out := make(map[string]PackageCompleteMsg)
	for _, m := range r.msgs {
		if c, ok := m.(PackageCompleteMsg); ok {
			out[c.PackageID] = c
		}
	}
	return out
}

func (r *recordingSender) summary(t *testing.T) SummaryMsg {
	t.Helper()
	last := r.msgs[len(r.msgs)-1]
	s, ok := last.(SummaryMsg)
	if !ok {
		t.Fatalf("last message should be SummaryMsg, got %T", last)
	}
	return s
}

// engineEnv is a sandbox with a mocked sudo/package manager and a command log.
type engineEnv struct {
	t       *testing.T
	dir     string
	logFile string
}

// newEngineEnv creates mock commands on PATH. The mock "sudo" logs its arguments and
// fails whenever they mention "badpkg", which lets tests exercise the fallback path.
func newEngineEnv(t *testing.T, withNala bool) *engineEnv {
	t.Helper()
	dir := t.TempDir()
	logFile := filepath.Join(dir, "commands.log")

	sudo := fmt.Sprintf(`#!/bin/bash
echo "sudo $*" >> %q
case "$*" in *badpkg*) echo "E: Unable to locate package badpkg" >&2; exit 100;; esac
exit 0
`, logFile)
	if err := os.WriteFile(filepath.Join(dir, "sudo"), []byte(sudo), 0o755); err != nil {
		t.Fatal(err)
	}
	if withNala {
		if err := os.WriteFile(filepath.Join(dir, "nala"), []byte("#!/bin/bash\nexit 0\n"), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("HOME", dir)
	return &engineEnv{t: t, dir: dir, logFile: logFile}
}

// plugin writes a fake plugin script whose functions log their invocation.
// checkExit is the exit status of check(); failing is the exit status of install().
func (e *engineEnv) plugin(id, method string, checkExit, installExit int, extra ...string) ToolPlugin {
	e.t.Helper()
	script := fmt.Sprintf(`#!/bin/bash
log() { echo "%s $1 ${EXTRA_INPUT:-}" >> %q; }
case "$1" in
  install) log install; echo "installing %s"; exit %d ;;
  setup_repo) log setup_repo; exit 0 ;;
  post_install) log post_install; exit 0 ;;
  check) log check; exit %d ;;
  *) echo "usage" >&2; exit 1 ;;
esac
`, id, e.logFile, id, installExit, checkExit)
	path := filepath.Join(e.dir, id+".sh")
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		e.t.Fatal(err)
	}
	p := ToolPlugin{ID: id, Name: strings.ToUpper(id), InstallMethod: method, ScriptPath: path}
	if method == "apt" {
		p.PackageName = id
	}
	for _, opt := range extra {
		switch opt {
		case "repo":
			p.RepoSetup = "setup_repo"
		case "post":
			p.PostInstall = "post_install"
		}
	}
	return p
}

func (e *engineEnv) commands() []string {
	data, _ := os.ReadFile(e.logFile)
	var lines []string
	for _, l := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		if l != "" {
			lines = append(lines, strings.TrimSpace(l))
		}
	}
	return lines
}

// look resolves commands against the mock directory only, so a nala installed on
// the host cannot leak into the tests.
func (e *engineEnv) look(name string) (string, error) {
	p := filepath.Join(e.dir, name)
	if _, err := os.Stat(p); err != nil {
		return "", err
	}
	return p, nil
}

func (e *engineEnv) run(plugins []ToolPlugin, inputs map[string]string) (*recordingSender, InstallResult) {
	sender := &recordingSender{}
	in := &Installer{Plugins: plugins, Inputs: inputs, Sender: sender, lookPath: e.look, prereqs: func() []string { return nil }}
	return sender, in.Run(context.Background())
}

func indexOf(lines []string, prefix string) int {
	for i, l := range lines {
		if strings.HasPrefix(l, prefix) {
			return i
		}
	}
	return -1
}

func TestInstallerHappyPath(t *testing.T) {
	env := newEngineEnv(t, true)
	plugins := []ToolPlugin{
		env.plugin("gitx", "apt", 0, 0, "repo", "post"),
		env.plugin("jqx", "apt", 0, 0),
		env.plugin("binx", "binary", 0, 0),
		env.plugin("manx", "manual", 1, 0),
	}

	sender, result := env.run(plugins, map[string]string{"EXTRA_INPUT": "hello"})
	cmds := env.commands()

	if got, want := strings.Join(result.Successful, ","), "GITX,JQX,BINX"; got != want {
		t.Errorf("successful = %q, want %q\ncommands:\n%s", got, want, strings.Join(cmds, "\n"))
	}
	if len(result.Failed) != 0 {
		t.Errorf("unexpected failures: %v", result.Failed)
	}
	if got, want := strings.Join(result.Skipped, ","), "MANX"; got != want {
		t.Errorf("skipped = %q, want %q", got, want)
	}

	// Phase ordering: repo setup -> nala update -> batch install -> post install -> checks -> binary.
	order := []string{
		"gitx setup_repo hello",
		sudoEnv + "nala update",
		sudoEnv + "nala install -y gitx jqx",
		"gitx post_install hello",
		"gitx check hello",
		"jqx check hello",
		"binx install hello",
		"binx check hello",
		"manx install hello",
	}
	last := -1
	for _, step := range order {
		idx := indexOf(cmds, step)
		if idx < 0 {
			t.Fatalf("expected command %q in log:\n%s", step, strings.Join(cmds, "\n"))
		}
		if idx <= last {
			t.Errorf("command %q ran out of order:\n%s", step, strings.Join(cmds, "\n"))
		}
		last = idx
	}
	if indexOf(cmds, "manx check") >= 0 {
		t.Error("manual plugins must not be verified")
	}

	done := sender.completions()
	if !done["gitx"].Success || !done["binx"].Success || !done["manx"].Skipped {
		t.Errorf("unexpected completion messages: %+v", done)
	}
	s := sender.summary(t)
	if len(s.Successful) != 3 || len(s.Skipped) != 1 {
		t.Errorf("summary mismatch: %+v", s)
	}
}

func TestInstallerFallsBackToAptGetWithoutNala(t *testing.T) {
	env := newEngineEnv(t, false)
	plugins := []ToolPlugin{env.plugin("jqx", "apt", 0, 0)}

	_, result := env.run(plugins, nil)
	cmds := env.commands()

	if !result.OK() {
		t.Fatalf("expected success, got %+v\n%s", result, strings.Join(cmds, "\n"))
	}
	if indexOf(cmds, sudoEnv+"apt-get update") < 0 || indexOf(cmds, sudoEnv+"apt-get install -y jqx") < 0 {
		t.Errorf("expected apt-get commands, got:\n%s", strings.Join(cmds, "\n"))
	}
	if indexOf(cmds, "sudo nala") >= 0 {
		t.Error("nala must not be used when it is not installed")
	}
}

func TestInstallerBatchFailureFallsBackToIndividualInstalls(t *testing.T) {
	env := newEngineEnv(t, true)
	plugins := []ToolPlugin{
		env.plugin("goodpkg", "apt", 0, 0, "post"),
		env.plugin("badpkg", "apt", 0, 0),
	}

	sender, result := env.run(plugins, nil)
	cmds := env.commands()

	if got := strings.Join(result.Successful, ","); got != "GOODPKG" {
		t.Errorf("successful = %q, want GOODPKG\n%s", got, strings.Join(cmds, "\n"))
	}
	if got := strings.Join(result.Failed, ","); got != "BADPKG" {
		t.Errorf("failed = %q, want BADPKG", got)
	}
	if indexOf(cmds, sudoEnv+"nala install -y goodpkg badpkg") < 0 {
		t.Error("batch install should have been attempted first")
	}
	if indexOf(cmds, sudoEnv+"nala install -y goodpkg") < 0 || indexOf(cmds, sudoEnv+"nala install -y badpkg") < 0 {
		t.Errorf("each package should be retried individually:\n%s", strings.Join(cmds, "\n"))
	}
	if indexOf(cmds, "goodpkg post_install") < 0 {
		t.Error("post_install must still run for packages that installed individually")
	}
	if indexOf(cmds, "badpkg check") >= 0 {
		t.Error("failed packages must not be verified")
	}
	if c := sender.completions()["badpkg"]; c.Success || c.Error == "" {
		t.Errorf("badpkg should be reported as failed with a reason, got %+v", c)
	}
}

func TestInstallerVerificationFailure(t *testing.T) {
	env := newEngineEnv(t, true)
	plugins := []ToolPlugin{
		env.plugin("ghost", "binary", 1, 0),  // installs fine, check fails
		env.plugin("broken", "binary", 0, 7), // install exits non-zero
		env.plugin("fine", "binary", 0, 0),
	}

	sender, result := env.run(plugins, nil)

	if got := strings.Join(result.Successful, ","); got != "FINE" {
		t.Errorf("successful = %q, want FINE", got)
	}
	if got := strings.Join(result.Failed, ","); got != "GHOST,BROKEN" {
		t.Errorf("failed = %q, want GHOST,BROKEN", got)
	}

	var errPhases []string
	for _, m := range sender.msgs {
		if e, ok := m.(ErrorMsg); ok {
			errPhases = append(errPhases, e.Package+":"+e.Phase)
		}
	}
	if got := strings.Join(errPhases, ","); got != "ghost:verify,broken:install" {
		t.Errorf("error phases = %q", got)
	}
	if c := sender.completions()["broken"]; !strings.Contains(c.Error, "Install failed") {
		t.Errorf("install failure should carry the reason, got %+v", c)
	}
}

func TestInstallerSkipVerify(t *testing.T) {
	env := newEngineEnv(t, true)
	plugins := []ToolPlugin{env.plugin("ghost", "binary", 1, 0)}

	sender := &recordingSender{}
	in := &Installer{Plugins: plugins, Sender: sender, SkipVerify: true, lookPath: env.look, prereqs: func() []string { return nil }}
	result := in.Run(context.Background())

	if !result.OK() || indexOf(env.commands(), "ghost check") >= 0 {
		t.Errorf("SkipVerify must skip check(): %+v\n%s", result, strings.Join(env.commands(), "\n"))
	}
}

func TestInstallerRepoSetupFailureFailsAptOnly(t *testing.T) {
	env := newEngineEnv(t, true)
	broken := env.plugin("repoonly", "apt", 0, 0, "repo")
	// Make setup_repo fail by pointing RepoSetup at a function that exits non-zero.
	broken.RepoSetup = "does_not_exist"
	plugins := []ToolPlugin{broken, env.plugin("other", "apt", 0, 0), env.plugin("binx", "binary", 0, 0)}

	_, result := env.run(plugins, nil)
	cmds := env.commands()

	if got := strings.Join(result.Failed, ","); got != "REPOONLY,OTHER" {
		t.Errorf("failed = %q, want REPOONLY,OTHER", got)
	}
	if got := strings.Join(result.Successful, ","); got != "BINX" {
		t.Errorf("binary plugins must still install, successful = %q", got)
	}
	if indexOf(cmds, sudoEnv+"nala install") >= 0 {
		t.Error("no APT install may run after repository setup fails")
	}
}

func TestInstallerInstallsMissingPrerequisites(t *testing.T) {
	env := newEngineEnv(t, true)
	plugins := []ToolPlugin{env.plugin("gitx", "apt", 0, 0, "repo")}

	sender := &recordingSender{}
	in := &Installer{Plugins: plugins, Sender: sender, lookPath: env.look, prereqs: func() []string { return []string{"curl", "gnupg"} }}
	result := in.Run(context.Background())
	cmds := env.commands()

	if !result.OK() {
		t.Fatalf("expected success: %+v\n%s", result, strings.Join(cmds, "\n"))
	}
	update := indexOf(cmds, sudoEnv+"apt-get update")
	prereq := indexOf(cmds, sudoEnv+"apt-get install -y curl gnupg")
	repo := indexOf(cmds, "gitx setup_repo")
	if update < 0 || prereq < 0 || repo < 0 || update >= prereq || prereq >= repo {
		t.Errorf("prerequisites must be installed before repository setup:\n%s", strings.Join(cmds, "\n"))
	}
}

func TestInstallerNoAptPluginsSkipsPackageManager(t *testing.T) {
	env := newEngineEnv(t, true)
	plugins := []ToolPlugin{env.plugin("binx", "binary", 0, 0)}

	_, result := env.run(plugins, nil)
	for _, c := range env.commands() {
		if strings.HasPrefix(c, "sudo") {
			t.Errorf("no package manager command expected, got %q", c)
		}
	}
	if !result.OK() {
		t.Errorf("expected success, got %+v", result)
	}
}

func TestInstallerCancellation(t *testing.T) {
	env := newEngineEnv(t, true)
	plugins := []ToolPlugin{env.plugin("a", "binary", 0, 0), env.plugin("b", "binary", 0, 0)}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	sender := &recordingSender{}
	in := &Installer{Plugins: plugins, Sender: sender, lookPath: env.look, prereqs: func() []string { return nil }}
	result := in.Run(ctx)

	if len(result.Failed) != 2 {
		t.Errorf("cancelled run must report every plugin as failed, got %+v", result)
	}
	if indexOf(env.commands(), "a install") >= 0 {
		t.Error("no script may run after cancellation")
	}
}

func TestScriptEnv(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("PATH", "/usr/bin:/bin")

	env := scriptEnv(map[string]string{"GIT_USER_NAME": "Jane"})
	get := func(key string) string {
		for _, kv := range env {
			if strings.HasPrefix(kv, key+"=") {
				return strings.TrimPrefix(kv, key+"=")
			}
		}
		return ""
	}

	path := get("PATH")
	for _, dir := range userBinDirs {
		if !strings.Contains(path, filepath.Join(home, dir)) {
			t.Errorf("PATH should contain %s: %s", dir, path)
		}
	}
	if !strings.HasSuffix(path, "/usr/bin:/bin") {
		t.Errorf("original PATH must be preserved: %s", path)
	}
	if get("GIT_USER_NAME") != "Jane" {
		t.Error("inputs must be exported to the script")
	}
	if get("DEBIAN_FRONTEND") != "noninteractive" {
		t.Error("DEBIAN_FRONTEND must be noninteractive")
	}
	if get("HOME") != home {
		t.Error("HOME must be inherited by the script")
	}

	// Calling twice must not duplicate PATH entries.
	t.Setenv("PATH", path)
	again := scriptEnv(nil)
	for _, kv := range again {
		if strings.HasPrefix(kv, "PATH=") && strings.Count(kv, filepath.Join(home, ".local/bin")) != 1 {
			t.Errorf("PATH entries duplicated: %s", kv)
		}
	}
}

func TestPlainSenderOutput(t *testing.T) {
	var sb strings.Builder
	plugins := []ToolPlugin{{ID: "git", Name: "Git"}}
	s := NewPlainSender(&sb, plugins)
	s.Send(PhaseStartMsg{Phase: "apt_batch", Count: 1})
	s.Send(PackageStartMsg{PackageID: "git", Phase: "install"})
	s.Send(LogMsg{Level: "info", Package: "git", Message: "hello"})
	s.Send(PackageCompleteMsg{PackageID: "git", Success: false, Error: "boom"})
	s.Send(PackageCompleteMsg{PackageID: "other", Skipped: true, Error: "manual"})

	out := sb.String()
	for _, want := range []string{"apt_batch", "Git", "Installing...", "hello", "boom", "other", "manual"} {
		if !strings.Contains(out, want) {
			t.Errorf("plain output missing %q:\n%s", want, out)
		}
	}

	sb.Reset()
	PrintSummary(&sb, InstallResult{Successful: []string{"Git"}, Failed: []string{"Zsh"}, Skipped: []string{"Ghostty"}})
	for _, want := range []string{"Installed (1)", "Git", "Failed (1)", "Zsh", "Manual steps required (1)", "Ghostty"} {
		if !strings.Contains(sb.String(), want) {
			t.Errorf("summary missing %q:\n%s", want, sb.String())
		}
	}
}
