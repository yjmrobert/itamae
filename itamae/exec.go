package itamae

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// userBinDirs are the user-level directories plugin scripts install binaries into.
// They are added to PATH for every script so check() works even when the user's
// shell profile has not picked them up yet.
var userBinDirs = []string{".local/bin", ".cargo/bin", ".atuin/bin"}

// scriptEnv builds the environment for a plugin script: the current environment,
// the user supplied inputs, and the user-level bin directories on PATH.
func scriptEnv(inputs map[string]string) []string {
	env := os.Environ()

	home := os.Getenv("HOME")
	if home != "" {
		path := os.Getenv("PATH")
		for _, dir := range userBinDirs {
			full := filepath.Join(home, dir)
			if !containsPathEntry(path, full) {
				path = full + string(os.PathListSeparator) + path
			}
		}
		env = setEnv(env, "PATH", path)
	}

	// Never let an interactive frontend block the installation.
	env = setEnv(env, "DEBIAN_FRONTEND", "noninteractive")

	for key, value := range inputs {
		env = setEnv(env, key, value)
	}
	return env
}

func containsPathEntry(path, entry string) bool {
	for _, p := range strings.Split(path, string(os.PathListSeparator)) {
		if p == entry {
			return true
		}
	}
	return false
}

func setEnv(env []string, key, value string) []string {
	prefix := key + "="
	for i, kv := range env {
		if strings.HasPrefix(kv, prefix) {
			env[i] = prefix + value
			return env
		}
	}
	return append(env, prefix+value)
}

// debugWriter forwards everything written to it to the debug log with a prefix.
type debugWriter struct {
	prefix string
}

func (w debugWriter) Write(p []byte) (int, error) {
	for _, line := range strings.Split(strings.TrimRight(string(p), "\n"), "\n") {
		if strings.TrimSpace(line) != "" {
			DebugLog("%s%s", w.prefix, line)
		}
	}
	return len(p), nil
}

// runScript executes one function of a plugin script ("install", "remove", ...).
// Output is always captured into the debug log and, when out is not nil, copied there too.
// The returned error includes the tail of the combined output so failures are self-explanatory.
func runScript(ctx context.Context, plugin ToolPlugin, command string, inputs map[string]string, out io.Writer) error {
	cmd := exec.CommandContext(ctx, "bash", plugin.ScriptPath, command)
	cmd.Env = scriptEnv(inputs)
	cmd.Stdin = nil

	var tail bytes.Buffer
	writers := []io.Writer{&tail, debugWriter{prefix: fmt.Sprintf("[%s:%s] ", plugin.ID, command)}}
	if out != nil {
		writers = append(writers, out)
	}
	w := io.MultiWriter(writers...)
	cmd.Stdout = w
	cmd.Stderr = w

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s %s failed: %w%s", plugin.ID, command, err, outputTail(tail.String(), 5))
	}
	return nil
}

// checkInstalled runs the plugin's check() function and reports whether it passed.
func checkInstalled(ctx context.Context, plugin ToolPlugin, inputs map[string]string) bool {
	cmd := exec.CommandContext(ctx, "bash", plugin.ScriptPath, "check")
	cmd.Env = scriptEnv(inputs)
	cmd.Stdin = nil
	output, err := cmd.CombinedOutput()
	if strings.TrimSpace(string(output)) != "" {
		DebugLog("[%s:check] %s", plugin.ID, strings.TrimSpace(string(output)))
	}
	return err == nil
}

// sudoPrefix is prepended to every privileged package manager call. sudo resets
// the environment, so DEBIAN_FRONTEND has to be passed explicitly or debconf
// blocks on prompts (tzdata, keyboard-configuration) in non-interactive runs.
var sudoPrefix = []string{"env", "DEBIAN_FRONTEND=noninteractive"}

// runSudo runs a privileged system command (package manager calls) through sudo, with
// output captured to the debug log and optionally copied to out.
func runSudo(ctx context.Context, out io.Writer, args ...string) error {
	args = append(append([]string{}, sudoPrefix...), args...)
	DebugLog("$ sudo %s", strings.Join(args, " "))
	cmd := exec.CommandContext(ctx, "sudo", args...)
	cmd.Env = scriptEnv(nil)
	cmd.Stdin = nil

	var tail bytes.Buffer
	writers := []io.Writer{&tail, debugWriter{prefix: "  | "}}
	if out != nil {
		writers = append(writers, out)
	}
	w := io.MultiWriter(writers...)
	cmd.Stdout = w
	cmd.Stderr = w

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("sudo %s failed: %w%s", strings.Join(args, " "), err, outputTail(tail.String(), 8))
	}
	return nil
}

// outputTail returns the last n non-empty lines of output, formatted for an error message.
func outputTail(output string, n int) string {
	var lines []string
	for _, line := range strings.Split(output, "\n") {
		if strings.TrimSpace(line) != "" {
			lines = append(lines, strings.TrimRight(line, "\r"))
		}
	}
	if len(lines) == 0 {
		return ""
	}
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	return "\n  " + strings.Join(lines, "\n  ")
}
