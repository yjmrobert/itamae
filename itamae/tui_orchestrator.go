package itamae

import (
	"fmt"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

// RunInstallTUI runs the installation with the new TUI interface
func RunInstallTUI(plugins []ToolPlugin, category string) {
	// Request sudo access upfront
	if err := ensureSudoAccess(); err != nil {
		fmt.Println("\n❌ Failed to obtain sudo access. Installation cancelled.")
		return
	}

	var selectedPlugins []ToolPlugin
	if category == "core" || category == "essentials" {
		// For core and essentials, install everything without prompting
		fmt.Printf("Installing %d %s packages\n", len(plugins), category)
		selectedPlugins = plugins
	} else {
		// For unverified, show multiselect
		selectedPlugins = selectPlugins(plugins)
		if len(selectedPlugins) == 0 {
			fmt.Println("No plugins selected. Exiting.")
			return
		}
	}

	// Gather all required inputs upfront
	requiredInputs := make(map[string]string)
	for _, p := range selectedPlugins {
		for _, input := range p.RequiredInputs {
			if _, ok := requiredInputs[input.Name]; !ok {
				defaultValue := getDefaultValue(input.DefaultCmd)
				value, err := RunTextInput(input.Prompt, defaultValue)
				if err != nil {
					Logger.Errorf("Error getting input for %s: %v", input.Name, err)
					return
				}
				requiredInputs[input.Name] = value
			}
		}
	}

	// Confirm before proceeding
	if !confirmInstallation() {
		fmt.Println("\nInstallation cancelled.")
		return
	}

	// Initialize TUI model
	model := NewInstallModel(selectedPlugins)
	
	// Create Bubbletea program
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	// Start installation in the background
	go processInstallTUI(p, selectedPlugins, requiredInputs)

	// Run the TUI
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
		return
	}
}

// processInstallTUI orchestrates the installation and sends messages to the TUI
func processInstallTUI(p *tea.Program, selectedPlugins []ToolPlugin, requiredInputs map[string]string) {
	// Separate plugins by install method
	aptPlugins := []ToolPlugin{}
	otherPlugins := []ToolPlugin{}
	for _, plugin := range selectedPlugins {
		if plugin.InstallMethod == "apt" {
			aptPlugins = append(aptPlugins, plugin)
		} else {
			otherPlugins = append(otherPlugins, plugin)
		}
	}

	// Track success/failure
	successful := []string{}
	failed := []string{}

	// Phase 0: Repository Setup
	repoPlugins := []ToolPlugin{}
	for _, plugin := range aptPlugins {
		if plugin.RepoSetup != "" {
			repoPlugins = append(repoPlugins, plugin)
		}
	}

	if len(repoPlugins) > 0 {
		p.Send(PhaseStartMsg{Phase: "repo_setup", Count: len(repoPlugins)})
		
		for _, plugin := range repoPlugins {
			p.Send(PackageStartMsg{PackageID: plugin.ID, Phase: "repo_setup"})
			p.Send(LogMsg{Level: "info", Package: plugin.ID, Message: "Setting up custom repository..."})
			
			// Execute repo setup synchronously (must be sequential)
			if err := executeScript(plugin, "setup_repo", requiredInputs); err != nil {
				p.Send(ErrorMsg{
					Package: plugin.ID,
					Phase:   "repo_setup",
					Message: fmt.Sprintf("Repository setup failed: %v", err),
				})
				p.Send(PackageCompleteMsg{PackageID: plugin.ID, Success: false, Error: err.Error()})
				p.Send(LogMsg{Level: "error", Package: "", Message: "Repository setup failed. Cannot proceed with installation."})
				p.Send(SummaryMsg{Successful: successful, Failed: append(failed, plugin.Name)})
				return
			}
			
			p.Send(PackageCompleteMsg{PackageID: plugin.ID, Success: true})
		}
		
		p.Send(PhaseCompleteMsg{Phase: "repo_setup"})
		
		// Run single apt-get update
		p.Send(LogMsg{Level: "info", Package: "", Message: "Updating package lists..."})
		
		var updateCmd *exec.Cmd
		if _, err := exec.LookPath("nala"); err == nil {
			updateCmd = exec.Command("sudo", "nala", "update")
		} else {
			updateCmd = exec.Command("sudo", "apt-get", "update")
		}
		
		output, err := updateCmd.CombinedOutput()
		if err != nil {
			p.Send(ErrorMsg{
				Package: "",
				Phase:   "update",
				Message: fmt.Sprintf("Package list update failed: %v\nOutput: %s", err, output),
			})
			p.Send(LogMsg{Level: "error", Package: "", Message: "Package list update failed. Cannot proceed with installation."})
			p.Send(SummaryMsg{Successful: successful, Failed: failed})
			return
		}
		
		p.Send(LogMsg{Level: "success", Package: "", Message: "Package lists updated successfully"})
	}

	// Phase 1: Batch install APT packages
	if len(aptPlugins) > 0 {
		p.Send(PhaseStartMsg{Phase: "apt_batch", Count: len(aptPlugins)})
		p.Send(LogMsg{Level: "info", Package: "", Message: fmt.Sprintf("Installing %d APT packages in parallel...", len(aptPlugins))})
		
		// Mark all APT packages as running
		for _, plugin := range aptPlugins {
			p.Send(PackageStartMsg{PackageID: plugin.ID, Phase: "install"})
		}
		
		// Check if nala is available
		_, err := exec.LookPath("nala")
		useNala := err == nil
		
		// Collect package names
		packages := []string{}
		for _, plugin := range aptPlugins {
			if plugin.PackageName != "" {
				packages = append(packages, plugin.PackageName)
			}
		}
		
		// Build install command
		var cmd *exec.Cmd
		if useNala {
			args := append([]string{"nala", "install", "-y"}, packages...)
			cmd = exec.Command("sudo", args...)
		} else {
			args := append([]string{"apt-get", "install", "-y"}, packages...)
			cmd = exec.Command("sudo", args...)
		}
		
		// Execute with output capture
		output, err := cmd.CombinedOutput()
		
		if err != nil {
			p.Send(LogMsg{Level: "error", Package: "", Message: fmt.Sprintf("Batch APT installation failed: %v", err)})
			p.Send(ErrorMsg{Package: "", Phase: "apt_batch", Message: string(output)})
			
			// Mark all as failed
			for _, plugin := range aptPlugins {
				p.Send(PackageCompleteMsg{PackageID: plugin.ID, Success: false, Error: "Batch installation failed"})
				failed = append(failed, plugin.Name)
			}
		} else {
			p.Send(LogMsg{Level: "success", Package: "", Message: fmt.Sprintf("Successfully installed %d APT packages", len(aptPlugins))})
			
			// Mark all as successful
			for _, plugin := range aptPlugins {
				p.Send(PackageCompleteMsg{PackageID: plugin.ID, Success: true})
				successful = append(successful, plugin.Name)
			}
			
			// Run post-install tasks
			for _, plugin := range aptPlugins {
				if plugin.PostInstall != "" {
					p.Send(PackageStartMsg{PackageID: plugin.ID, Phase: "post_install"})
					p.Send(LogMsg{Level: "info", Package: plugin.ID, Message: "Running post-installation tasks..."})
					
					if err := executeScript(plugin, "post_install", requiredInputs); err != nil {
						p.Send(LogMsg{Level: "warning", Package: plugin.ID, Message: fmt.Sprintf("Post-install failed: %v", err)})
					} else {
						p.Send(LogMsg{Level: "success", Package: plugin.ID, Message: "Post-installation complete"})
					}
					
					p.Send(PackageCompleteMsg{PackageID: plugin.ID, Success: true})
				}
			}
		}
		
		p.Send(PhaseCompleteMsg{Phase: "apt_batch"})
	}

	// Phase 2: Install other plugins individually
	if len(otherPlugins) > 0 {
		p.Send(PhaseStartMsg{Phase: "individual", Count: len(otherPlugins)})
		
		for _, plugin := range otherPlugins {
			p.Send(PackageStartMsg{PackageID: plugin.ID, Phase: "install"})
			p.Send(LogMsg{Level: "info", Package: plugin.ID, Message: "Installing..."})
			
			if err := executeScript(plugin, "install", requiredInputs); err != nil {
				p.Send(ErrorMsg{
					Package: plugin.ID,
					Phase:   "install",
					Message: fmt.Sprintf("Installation failed: %v", err),
				})
				p.Send(PackageCompleteMsg{PackageID: plugin.ID, Success: false, Error: err.Error()})
				failed = append(failed, plugin.Name)
			} else {
				p.Send(PackageCompleteMsg{PackageID: plugin.ID, Success: true})
				successful = append(successful, plugin.Name)
			}
		}
		
		p.Send(PhaseCompleteMsg{Phase: "individual"})
	}

	// Send summary
	p.Send(SummaryMsg{Successful: successful, Failed: failed})
}
