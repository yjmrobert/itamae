package itamae

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// ExecuteScriptAsyncCmd executes a plugin script asynchronously and returns a tea.Cmd
// that sends log messages and completion status back to the TUI
func ExecuteScriptAsyncCmd(plugin ToolPlugin, command string, env map[string]string) tea.Cmd {
	return func() tea.Msg {
		// Build the command
		cmd := exec.Command("bash", "-c", fmt.Sprintf("source %s && %s", plugin.ScriptPath, command))
		
		// Set environment variables
		for k, v := range env {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
		}
		
		// Create pipes for stdout and stderr
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return ErrorMsg{
				Package: plugin.ID,
				Phase:   command,
				Message: fmt.Sprintf("Failed to create stdout pipe: %v", err),
			}
		}
		
		stderr, err := cmd.StderrPipe()
		if err != nil {
			return ErrorMsg{
				Package: plugin.ID,
				Phase:   command,
				Message: fmt.Sprintf("Failed to create stderr pipe: %v", err),
			}
		}
		
		// Start the command
		if err := cmd.Start(); err != nil {
			return ErrorMsg{
				Package: plugin.ID,
				Phase:   command,
				Message: fmt.Sprintf("Failed to start command: %v", err),
			}
		}
		
		// Read stdout in a goroutine and collect logs
		var stdoutBuf bytes.Buffer
		stdoutDone := make(chan bool)
		go func() {
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				line := scanner.Text()
				stdoutBuf.WriteString(line + "\n")
				// Note: We're collecting output but not sending individual line messages
				// to avoid overwhelming the TUI with too many updates
			}
			stdoutDone <- true
		}()
		
		// Read stderr in a goroutine
		var stderrBuf bytes.Buffer
		stderrDone := make(chan bool)
		go func() {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				line := scanner.Text()
				stderrBuf.WriteString(line + "\n")
			}
			stderrDone <- true
		}()
		
		// Wait for output collection to complete
		<-stdoutDone
		<-stderrDone
		
		// Wait for command to complete
		err = cmd.Wait()
		
		// Determine success/failure
		if err != nil {
			// Command failed
			errorMsg := stderrBuf.String()
			if errorMsg == "" {
				errorMsg = fmt.Sprintf("Command exited with error: %v", err)
			}
			
			return PackageCompleteMsg{
				PackageID: plugin.ID,
				Success:   false,
				Error:     strings.TrimSpace(errorMsg),
			}
		}
		
		// Command succeeded
		return PackageCompleteMsg{
			PackageID: plugin.ID,
			Success:   true,
			Error:     "",
		}
	}
}

// ExecuteBatchAPTCmd executes a batch APT install command asynchronously
func ExecuteBatchAPTCmd(packages []string, useNala bool) tea.Cmd {
	return func() tea.Msg {
		var cmd *exec.Cmd
		
		if useNala {
			args := append([]string{"install", "-y"}, packages...)
			cmd = exec.Command("sudo", append([]string{"nala"}, args...)...)
		} else {
			args := append([]string{"install", "-y"}, packages...)
			cmd = exec.Command("sudo", append([]string{"apt-get"}, args...)...)
		}
		
		// Create output buffer
		var outputBuf bytes.Buffer
		cmd.Stdout = &outputBuf
		cmd.Stderr = &outputBuf
		
		// Run the command
		err := cmd.Run()
		
		if err != nil {
			return LogMsg{
				Level:   "error",
				Package: "",
				Message: fmt.Sprintf("Batch APT install failed: %v\nOutput: %s", err, outputBuf.String()),
			}
		}
		
		return LogMsg{
			Level:   "success",
			Package: "",
			Message: fmt.Sprintf("Successfully installed %d APT packages", len(packages)),
		}
	}
}

// ExecuteSystemCommandCmd executes a system command (like apt-get update) asynchronously
func ExecuteSystemCommandCmd(description string, cmdName string, args ...string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command(cmdName, args...)
		
		var outputBuf bytes.Buffer
		cmd.Stdout = &outputBuf
		cmd.Stderr = &outputBuf
		
		err := cmd.Run()
		
		if err != nil {
			return ErrorMsg{
				Package: "",
				Phase:   description,
				Message: fmt.Sprintf("Command '%s %s' failed: %v\nOutput: %s",
					cmdName, strings.Join(args, " "), err, outputBuf.String()),
			}
		}
		
		return LogMsg{
			Level:   "success",
			Package: "",
			Message: fmt.Sprintf("%s completed successfully", description),
		}
	}
}

// StreamedExecuteScriptCmd executes a script and streams output line-by-line to the TUI
// This is useful for long-running commands where we want real-time feedback
func StreamedExecuteScriptCmd(plugin ToolPlugin, command string, env map[string]string, program *tea.Program) tea.Cmd {
	return func() tea.Msg {
		// Build the command
		cmd := exec.Command("bash", "-c", fmt.Sprintf("source %s && %s", plugin.ScriptPath, command))
		
		// Set environment variables
		for k, v := range env {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
		}
		
		// Create pipes
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return ErrorMsg{
				Package: plugin.ID,
				Phase:   command,
				Message: fmt.Sprintf("Failed to create stdout pipe: %v", err),
			}
		}
		
		stderr, err := cmd.StderrPipe()
		if err != nil {
			return ErrorMsg{
				Package: plugin.ID,
				Phase:   command,
				Message: fmt.Sprintf("Failed to create stderr pipe: %v", err),
			}
		}
		
		// Start the command
		if err := cmd.Start(); err != nil {
			return ErrorMsg{
				Package: plugin.ID,
				Phase:   command,
				Message: fmt.Sprintf("Failed to start command: %v", err),
			}
		}
		
		// Stream stdout
		go func() {
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.TrimSpace(line) != "" {
					program.Send(LogMsg{
						Level:   "info",
						Package: plugin.ID,
						Message: line,
					})
				}
			}
		}()
		
		// Stream stderr
		go func() {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.TrimSpace(line) != "" {
					program.Send(LogMsg{
						Level:   "warning",
						Package: plugin.ID,
						Message: line,
					})
				}
			}
		}()
		
		// Wait for command to complete
		err = cmd.Wait()
		
		// Send completion message
		if err != nil {
			return PackageCompleteMsg{
				PackageID: plugin.ID,
				Success:   false,
				Error:     fmt.Sprintf("Command failed: %v", err),
			}
		}
		
		return PackageCompleteMsg{
			PackageID: plugin.ID,
			Success:   true,
			Error:     "",
		}
	}
}

// CaptureCommandOutput executes a command and returns its output as a string
// This is a synchronous helper function (not a tea.Cmd)
func CaptureCommandOutput(cmdName string, args ...string) (string, error) {
	cmd := exec.Command(cmdName, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}
