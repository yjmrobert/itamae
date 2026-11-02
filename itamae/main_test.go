package itamae

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

var (
	mockDir     string
	logFilePath string
	plugins     []ToolPlugin // every plugin of every category
)

// mockCommands are replaced by logging stubs so scripts can run without touching the system.
var mockCommands = []string{
	"sudo", "nala", "apt-get", "apt-cache", "dpkg", "ln", "mkdir", "rm", "chsh", "git", "curl", "unzip",
	"stow", "chmod", "bash", "sh", "pipx", "rustup", "wget", "tar", "fc-cache", "tee", "mv", "cp",
}

// pluginAssertion defines a command fragment expected in the log for each script action.
type pluginAssertion struct {
	install string
	remove  string
}

var pluginAssertions = map[string]pluginAssertion{
	// Core
	"alacritty":           {install: "sudo nala install -y alacritty", remove: "sudo apt-get purge -y alacritty"},
	"apt-transport-https": {install: "sudo nala install -y apt-transport-https", remove: "sudo apt-get purge -y apt-transport-https"},
	"ca-certificates":     {install: "sudo nala install -y ca-certificates", remove: "sudo apt-get purge -y ca-certificates"},
	"curl":                {install: "sudo nala install -y curl", remove: "sudo apt-get purge -y curl"},
	"dotnet-sdk-8.0":      {install: "sudo nala install -y dotnet-sdk-8.0", remove: "sudo apt-get purge -y dotnet-sdk-8.0"},
	"fd":                  {install: "sudo nala install -y fd-find", remove: "sudo apt-get purge -y fd-find"},
	"fzf":                 {install: "sudo nala install -y fzf", remove: "sudo apt-get purge -y fzf"},
	"gh":                  {install: "sudo nala install -y gh", remove: "sudo apt-get purge -y gh"},
	"git":                 {install: "sudo nala install -y git", remove: "sudo apt-get purge -y git"},
	"gnupg":               {install: "sudo nala install -y gnupg", remove: "sudo apt-get purge -y gnupg"},
	"helm":                {install: "curl --silent https://raw.githubusercontent.com/helm/helm", remove: "sudo rm -f /usr/local/bin/helm"},
	"jq":                  {install: "sudo nala install -y jq", remove: "sudo apt-get purge -y jq"},
	"kubectl":             {install: "curl -fsSL https://dl.k8s.io/release/stable.txt", remove: "sudo rm -f /usr/local/bin/kubectl"},
	"kubecolor":           {install: "https://github.com/kubecolor/kubecolor/releases/latest", remove: "rm -f"},
	"lsd":                 {install: "sudo nala install -y lsd", remove: "sudo apt-get purge -y lsd"},
	"nala":                {install: "sudo apt-get install -y nala", remove: "sudo apt-get purge -y nala"},
	"nodejs":              {install: "sudo nala install -y nodejs", remove: "sudo apt-get purge -y nodejs"},
	"pipx":                {install: "sudo nala install -y pipx", remove: "sudo apt-get purge -y pipx"},
	"python3-full":        {install: "sudo nala install -y python3-full", remove: "sudo apt-get purge -y python3-full"},
	"task":                {install: "curl --silent --location https://taskfile.dev/install.sh", remove: "rm -f"},
	"wget":                {install: "sudo nala install -y wget", remove: "sudo apt-get purge -y wget"},
	"wireguard":           {install: "sudo nala install -y wireguard-tools", remove: "sudo apt-get purge -y wireguard-tools"},
	"yq":                  {install: "sudo curl --silent -fL https://github.com/mikefarah/yq/releases/latest/download/yq_linux_", remove: "sudo rm -f /usr/local/bin/yq"},

	// Essentials
	"atuin":    {install: "bash", remove: "rm -rf"},
	"bat":      {install: "sudo nala install -y bat", remove: "sudo apt-get purge -y bat"},
	"java":     {install: "sudo nala install -y temurin-21-jdk", remove: "sudo apt-get purge -y temurin-21-jdk"},
	"maven":    {install: "curl -fsSL https://dlcdn.apache.org/maven/maven-3/", remove: "sudo rm -rf /opt/maven"},
	"ripgrep":  {install: "sudo nala install -y ripgrep", remove: "sudo apt-get purge -y ripgrep"},
	"rust":     {install: "curl --proto", remove: "rustup self uninstall -y"},
	"sdkman":   {install: "curl -s https://get.sdkman.io", remove: "rm -rf"},
	"starship": {install: "curl -sS https://starship.rs/install.sh", remove: ""},
	"stow":     {install: "sudo nala install -y stow", remove: "sudo apt-get purge -y stow"},
	"zoxide":   {install: "curl -sS https://raw.githubusercontent.com/ajeetdsouza/zoxide/main/install.sh", remove: "rm"},

	// Unverified
	"ansible":       {install: "pipx install", remove: "pipx uninstall ansible"},
	"bin":           {install: "https://github.com/marcosnils/bin/releases/latest", remove: "rm -f"},
	"btop":          {install: "sudo nala install -y btop", remove: "sudo apt-get purge -y btop"},
	"btop-desktop":  {install: "sudo nala install -y btop", remove: "sudo apt-get purge -y btop"},
	"cascadia-code": {install: "curl -fsSL https://github.com/microsoft/cascadia-code/releases/download", remove: "rm -f"},
	"chezmoi":       {install: "sh -c", remove: "rm"},
	"dunst":         {install: "sudo nala install -y dunst", remove: "sudo apt-get purge -y dunst"},
	"flameshot":     {install: "sudo nala install -y flameshot", remove: "sudo apt-get purge -y flameshot"},
	"ghostty":       {install: "mkdir -p", remove: "rm -f"},
	"httpie":        {install: "sudo nala install -y httpie", remove: "sudo apt-get purge -y httpie"},
	"meld":          {install: "sudo nala install -y meld", remove: "sudo apt-get purge -y meld"},
	"ncdu":          {install: "sudo nala install -y ncdu", remove: "sudo apt-get purge -y ncdu"},
	"pass":          {install: "sudo nala install -y pass", remove: "sudo apt-get purge -y pass"},
	"polybar":       {install: "sudo nala install -y polybar", remove: "sudo apt-get purge -y polybar"},
	"rofi":          {install: "sudo nala install -y rofi", remove: "sudo apt-get purge -y rofi"},
	"ruby":          {install: "sudo nala install -y ruby-full", remove: "sudo apt-get purge -y ruby-full"},
	"semgrep":       {install: "pipx install semgrep", remove: "pipx uninstall semgrep"},
	"tldr":          {install: "curl -fsSL https://github.com/tealdeer-rs/tealdeer/releases/latest/download/tealdeer-linux-", remove: "rm -f"},
	"vscode":        {install: "sudo apt-get install -y /tmp/vscode-itamae.deb", remove: "sudo apt-get purge -y code"},
	"zellij":        {install: "curl -fsSL https://github.com/zellij-org/zellij/releases/latest/download/zellij-", remove: "sudo rm -f /usr/local/bin/zellij"},
	"zsh":           {install: "sudo nala install -y zsh", remove: "sudo apt-get purge -y zsh"},
}

// repoSetupAssertions lists the command expected when a plugin configures its APT repository.
var repoSetupAssertions = map[string]string{
	"gh":     "curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg",
	"nodejs": "curl --silent -fsSL https://deb.nodesource.com/setup_lts.x",
	"java":   "curl -fsSL https://packages.adoptium.net/artifactory/api/gpg/key/public",
}

// expectedMetadata is the authoritative list of plugins and their install method.
// Adding or removing a script must be reflected here on purpose.
var expectedMetadata = map[string]struct {
	category    string
	method      string
	packageName string
}{
	// Core
	"alacritty":           {"core", "binary", ""},
	"apt-transport-https": {"core", "apt", "apt-transport-https"},
	"ca-certificates":     {"core", "apt", "ca-certificates"},
	"curl":                {"core", "apt", "curl"},
	"dotnet-sdk-8.0":      {"core", "apt", "dotnet-sdk-8.0"},
	"fd":                  {"core", "apt", "fd-find"},
	"fzf":                 {"core", "apt", "fzf"},
	"gh":                  {"core", "apt", "gh"},
	"git":                 {"core", "apt", "git"},
	"gnupg":               {"core", "apt", "gnupg"},
	"helm":                {"core", "binary", ""},
	"jq":                  {"core", "apt", "jq"},
	"kubectl":             {"core", "binary", ""},
	"kubecolor":           {"core", "binary", ""},
	"lsd":                 {"core", "binary", ""},
	"nala":                {"core", "apt", "nala"},
	"nodejs":              {"core", "apt", "nodejs"},
	"pipx":                {"core", "apt", "pipx"},
	"python3-full":        {"core", "apt", "python3-full"},
	"task":                {"core", "binary", ""},
	"wget":                {"core", "apt", "wget"},
	"wireguard":           {"core", "apt", "wireguard-tools"},
	"yq":                  {"core", "binary", ""},

	// Essentials
	"atuin":    {"essentials", "binary", ""},
	"bat":      {"essentials", "apt", "bat"},
	"java":     {"essentials", "apt", "temurin-21-jdk"},
	"maven":    {"essentials", "binary", ""},
	"ripgrep":  {"essentials", "apt", "ripgrep"},
	"rust":     {"essentials", "binary", ""},
	"sdkman":   {"essentials", "binary", ""},
	"starship": {"essentials", "binary", ""},
	"stow":     {"essentials", "apt", "stow"},
	"zoxide":   {"essentials", "binary", ""},

	// Unverified
	"ansible":       {"unverified", "binary", ""},
	"bin":           {"unverified", "binary", ""},
	"btop":          {"unverified", "apt", "btop"},
	"btop-desktop":  {"unverified", "apt", "btop"},
	"cascadia-code": {"unverified", "binary", ""},
	"chezmoi":       {"unverified", "binary", ""},
	"dunst":         {"unverified", "apt", "dunst"},
	"flameshot":     {"unverified", "apt", "flameshot"},
	"ghostty":       {"unverified", "manual", ""},
	"httpie":        {"unverified", "apt", "httpie"},
	"meld":          {"unverified", "apt", "meld"},
	"ncdu":          {"unverified", "apt", "ncdu"},
	"pass":          {"unverified", "apt", "pass"},
	"polybar":       {"unverified", "apt", "polybar"},
	"rofi":          {"unverified", "apt", "rofi"},
	"ruby":          {"unverified", "apt", "ruby-full"},
	"semgrep":       {"unverified", "binary", ""},
	"tldr":          {"unverified", "binary", ""},
	"vscode":        {"unverified", "binary", ""},
	"zellij":        {"unverified", "binary", ""},
	"zsh":           {"unverified", "apt", "zsh"},
}

var expectedPostInstall = map[string]string{
	"bat": "post_install",
	"fd":  "post_install",
}

var expectedRepoSetup = map[string]string{
	"gh":             "setup_repo",
	"nodejs":         "setup_repo",
	"dotnet-sdk-8.0": "setup_repo",
	"java":           "setup_repo",
}

// TestMain loads every plugin once and prepares the mock command directory.
func TestMain(m *testing.M) {
	var cleanupPlugins func()
	var err error
	plugins, cleanupPlugins, err = LoadAllPlugins()
	if err != nil {
		fmt.Printf("Failed to load plugins in TestMain: %v\n", err)
		os.Exit(1)
	}

	var cleanupMocks func()
	mockDir, logFilePath, cleanupMocks = setupTestEnvironment()

	exitCode := m.Run()

	cleanupPlugins()
	cleanupMocks()
	os.Exit(exitCode)
}

// setupTestEnvironment creates a temporary directory with logging stubs for mockCommands.
func setupTestEnvironment() (string, string, func()) {
	dir, err := os.MkdirTemp("", "itamae-test-mocks-")
	if err != nil {
		panic(fmt.Sprintf("Failed to create mock dir: %v", err))
	}

	logFile, err := os.CreateTemp("", "itamae-test-log-")
	if err != nil {
		panic(fmt.Sprintf("Failed to create log file: %v", err))
	}
	path := logFile.Name()
	_ = logFile.Close()

	for _, cmdName := range mockCommands {
		scriptContent := fmt.Sprintf("#!/bin/bash\n# Mock for '%s'\necho \"%s $*\" >> %q\n", cmdName, cmdName, path)
		if cmdName == "apt-cache" {
			// Pretend every package has an installation candidate.
			scriptContent += `case "$*" in policy*) echo "  Candidate: 1.0-1" ;; esac
`
		}
		if cmdName == "curl" {
			// Scripts that resolve "the latest version" parse curl output; give them something plausible.
			scriptContent += `case "$*" in
  *stable.txt*) echo "v1.99.0" ;;
  *releases/latest*) echo "https://github.com/mock/mock/releases/tag/v9.9.9" ;;
  *dlcdn.apache.org/maven*) echo '<a href="3.9.99/">3.9.99/</a>' ;;
esac
`
		}
		if err := os.WriteFile(filepath.Join(dir, cmdName), []byte(scriptContent), 0o755); err != nil {
			panic(fmt.Sprintf("Failed to write mock script for %s: %v", cmdName, err))
		}
	}

	cleanup := func() {
		_ = os.RemoveAll(dir)
		_ = os.Remove(path)
	}
	return dir, path, cleanup
}

// runScriptWithMocks executes one script action with the mocks first on PATH and
// returns the recorded command log.
func runScriptWithMocks(t *testing.T, plugin ToolPlugin, command string) string {
	t.Helper()

	if err := os.WriteFile(logFilePath, []byte{}, 0o644); err != nil {
		t.Fatalf("Failed to clear log file: %v", err)
	}

	homeDir := t.TempDir()
	cmd := exec.Command("/bin/bash", plugin.ScriptPath, command)
	cmd.Env = append(os.Environ(),
		"PATH="+mockDir+string(os.PathListSeparator)+os.Getenv("PATH"),
		"HOME="+homeDir,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Script %s %s failed. Output:\n%s\nError: %v", plugin.ID, command, string(output), err)
	}

	logBytes, err := os.ReadFile(logFilePath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}
	return strings.TrimSpace(string(logBytes))
}

func runPluginTest(t *testing.T, plugin ToolPlugin, command, expectedLog string) {
	t.Helper()
	logContent := runScriptWithMocks(t, plugin, command)
	// An empty expectation means the action must run cleanly but issues no mocked command
	// (for example a remove() that finds nothing to delete).
	if expectedLog != "" && !strings.Contains(logContent, expectedLog) {
		t.Errorf("Expected log to contain %q, but it was:\n%s", expectedLog, logContent)
	}
}

// TestInstall runs the install action of every plugin against the mocks.
func TestInstall(t *testing.T) {
	for _, plugin := range plugins {
		t.Run("install_"+plugin.ID, func(t *testing.T) {
			assertion, ok := pluginAssertions[plugin.ID]
			if !ok {
				t.Fatalf("No install/remove assertion defined for plugin %s", plugin.ID)
			}
			runPluginTest(t, plugin, "install", assertion.install)
		})
	}
}

// TestUninstall runs the remove action of every plugin against the mocks.
func TestUninstall(t *testing.T) {
	for _, plugin := range plugins {
		t.Run("uninstall_"+plugin.ID, func(t *testing.T) {
			assertion, ok := pluginAssertions[plugin.ID]
			if !ok {
				t.Fatalf("No install/remove assertion defined for plugin %s", plugin.ID)
			}
			runPluginTest(t, plugin, "remove", assertion.remove)
		})
	}
}

// TestRepoSetup runs the repository setup of plugins that declare one.
func TestRepoSetup(t *testing.T) {
	for _, plugin := range plugins {
		if plugin.RepoSetup == "" {
			continue
		}
		t.Run("setup_repo_"+plugin.ID, func(t *testing.T) {
			logContent := runScriptWithMocks(t, plugin, plugin.RepoSetup)
			expected, ok := repoSetupAssertions[plugin.ID]
			if !ok {
				// Distribution dependent setup (e.g. .NET) only has to run cleanly.
				return
			}
			if !strings.Contains(logContent, expected) {
				t.Errorf("Expected log to contain %q, but it was:\n%s", expected, logContent)
			}
		})
	}
}

// TestPostInstall runs the post-install hook of plugins that declare one.
func TestPostInstall(t *testing.T) {
	for _, plugin := range plugins {
		if plugin.PostInstall == "" {
			continue
		}
		t.Run("post_install_"+plugin.ID, func(t *testing.T) {
			runScriptWithMocks(t, plugin, plugin.PostInstall)
		})
	}
}

// TestRouterRejectsUnknownAction verifies every script prints usage and fails for bad input.
func TestRouterRejectsUnknownAction(t *testing.T) {
	for _, plugin := range plugins {
		t.Run(plugin.ID, func(t *testing.T) {
			cmd := exec.Command("/bin/bash", plugin.ScriptPath, "definitely-not-an-action")
			cmd.Env = append(os.Environ(), "PATH="+mockDir+string(os.PathListSeparator)+os.Getenv("PATH"))
			output, err := cmd.CombinedOutput()
			if err == nil {
				t.Fatalf("script must exit non-zero for an unknown action, output:\n%s", output)
			}
			if !strings.Contains(string(output), "Usage:") {
				t.Errorf("script should print usage, got:\n%s", output)
			}
		})
	}
}

// TestMetadataParsing verifies the embedded scripts match the expected catalogue exactly.
func TestMetadataParsing(t *testing.T) {
	seen := make(map[string]bool)
	for _, plugin := range plugins {
		seen[plugin.ID] = true
		t.Run("metadata_"+plugin.ID, func(t *testing.T) {
			want, ok := expectedMetadata[plugin.ID]
			if !ok {
				t.Fatalf("Plugin %q is not listed in expectedMetadata; add it (and an install/remove assertion)", plugin.ID)
			}
			if plugin.Category != want.category {
				t.Errorf("category = %q, want %q", plugin.Category, want.category)
			}
			if plugin.InstallMethod != want.method {
				t.Errorf("InstallMethod = %q, want %q", plugin.InstallMethod, want.method)
			}
			if plugin.PackageName != want.packageName {
				t.Errorf("PackageName = %q, want %q", plugin.PackageName, want.packageName)
			}
			if plugin.PostInstall != expectedPostInstall[plugin.ID] {
				t.Errorf("PostInstall = %q, want %q", plugin.PostInstall, expectedPostInstall[plugin.ID])
			}
			if plugin.RepoSetup != expectedRepoSetup[plugin.ID] {
				t.Errorf("RepoSetup = %q, want %q", plugin.RepoSetup, expectedRepoSetup[plugin.ID])
			}
			if plugin.Description == "" {
				t.Error("every plugin needs a DESCRIPTION")
			}
		})
	}

	var missing []string
	for id := range expectedMetadata {
		if !seen[id] {
			missing = append(missing, id)
		}
	}
	sort.Strings(missing)
	if len(missing) > 0 {
		t.Errorf("expectedMetadata lists plugins that no longer exist: %v", missing)
	}
}

// TestScriptsDeclareRequiredFunctions makes sure every script implements the plugin contract.
func TestScriptsDeclareRequiredFunctions(t *testing.T) {
	for _, plugin := range plugins {
		t.Run(plugin.ID, func(t *testing.T) {
			content, err := os.ReadFile(plugin.ScriptPath)
			if err != nil {
				t.Fatal(err)
			}
			text := string(content)
			for _, fn := range []string{"install()", "remove()", "check()"} {
				if !strings.Contains(text, fn) {
					t.Errorf("script is missing the %s function", fn)
				}
			}
			if plugin.RepoSetup != "" && !strings.Contains(text, plugin.RepoSetup+"()") {
				t.Errorf("REPO_SETUP references %s, which is not defined", plugin.RepoSetup)
			}
			if plugin.PostInstall != "" && !strings.Contains(text, plugin.PostInstall+"()") {
				t.Errorf("POST_INSTALL references %s, which is not defined", plugin.PostInstall)
			}
			if !strings.HasPrefix(text, "#!/bin/bash") {
				t.Error("script must start with a #!/bin/bash shebang")
			}
		})
	}
}

// TestCategoryComposition guards the split between batched APT installs and individual installs.
func TestCategoryComposition(t *testing.T) {
	for _, category := range Categories {
		apt, other := 0, 0
		for _, p := range plugins {
			if p.Category != category {
				continue
			}
			if p.InstallMethod == "apt" {
				apt++
				if p.PackageName == "" {
					t.Errorf("APT plugin %q is missing PackageName", p.ID)
				}
			} else {
				other++
			}
		}
		if apt+other == 0 {
			t.Errorf("category %s is empty", category)
		}
		t.Logf("%s: %d apt, %d individual", category, apt, other)
	}
}
