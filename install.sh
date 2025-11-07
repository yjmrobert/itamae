#!/bin/bash
# install.sh
# Installs the prerequisites and runs the itamae tool.

# Exit on error
set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Helper functions for colored output
info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

success() {
    echo -e "${GREEN}✓${NC} $1"
}

warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

error() {
    echo -e "${RED}✗${NC} $1"
}

# Compare two semantic versions
# Returns: 0 if equal, 1 if $1 > $2, 2 if $1 < $2
compare_versions() {
    local ver1=$1
    local ver2=$2
    
    # Remove 'v' prefix if present
    ver1=${ver1#v}
    ver2=${ver2#v}
    
    if [[ "$ver1" == "$ver2" ]]; then
        return 0
    fi
    
    # Split versions into arrays
    IFS='.' read -ra V1 <<< "$ver1"
    IFS='.' read -ra V2 <<< "$ver2"
    
    # Compare major, minor, patch
    for i in 0 1 2; do
        local num1=${V1[$i]:-0}
        local num2=${V2[$i]:-0}
        
        # Remove any non-numeric suffixes
        num1=$(echo "$num1" | grep -oE '^[0-9]+')
        num2=$(echo "$num2" | grep -oE '^[0-9]+')
        
        if ((num1 > num2)); then
            return 1  # ver1 is greater
        elif ((num1 < num2)); then
            return 2  # ver1 is less
        fi
    done
    
    return 0  # equal
}

# Get the latest release version from GitHub
get_latest_version() {
    # Try using GitHub API (rate limited but reliable)
    local latest=$(curl -s https://api.github.com/repos/yjmrobert/itamae/releases/latest 2>/dev/null | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    
    if [[ -z "$latest" ]]; then
        # Fallback: try git ls-remote
        latest=$(git ls-remote --tags --refs --sort="v:refname" https://github.com/yjmrobert/itamae.git 2>/dev/null | tail -n1 | sed 's/.*\///')
    fi
    
    if [[ -z "$latest" ]]; then
        echo ""  # Return empty string instead of "latest"
    else
        echo "$latest"
    fi
}

# Get the currently installed version of itamae
get_installed_version() {
    if command -v itamae &> /dev/null; then
        # Try to get version from itamae version command
        local version_output=$(itamae version 2>/dev/null | head -n1)
        
        if [[ -z "$version_output" ]]; then
            version_output=$(itamae --version 2>/dev/null | head -n1)
        fi
        
        # Extract version number (handles "itamae version v1.0.0" or "v1.0.0")
        local version=$(echo "$version_output" | grep -oE 'v?[0-9]+\.[0-9]+\.[0-9]+' | head -n1)
        
        if [[ -z "$version" ]]; then
            echo "unknown"
        else
            echo "$version"
        fi
    else
        echo ""
    fi
}

# Update package lists and install prerequisites
echo "========================================"
echo "Installing Prerequisites"
echo "========================================"
info "Updating package lists..."
sudo apt-get update -y
info "Installing prerequisites (git, nala)..."
sudo apt-get install -y git nala
success "Prerequisites installed"

success "Prerequisites installed"

echo ""
echo "========================================"
echo "Checking Itamae Installation"
echo "========================================"

# Check if itamae is already installed
if command -v itamae &> /dev/null; then
    success "itamae is already installed"
    
    # Get installed version
    INSTALLED_VERSION=$(get_installed_version)
    info "Installed version: $INSTALLED_VERSION"
    
    # Determine target version
    TARGET_VERSION=${ITAMAE_VERSION:-}
    
    if [[ -z "$TARGET_VERSION" ]]; then
        # Get latest version from GitHub
        TARGET_VERSION=$(get_latest_version)
        if [[ -n "$TARGET_VERSION" ]]; then
            info "Latest available version: $TARGET_VERSION"
        else
            info "Will update to latest available version"
        fi
    else
        info "Target version (from ITAMAE_VERSION): $TARGET_VERSION"
    fi
    
    # Determine if update is needed
    SHOULD_UPDATE=false
    
    # Skip update if installed version is unknown (old version without version command)
    if [[ "$INSTALLED_VERSION" == "unknown" ]]; then
        warning "Cannot determine installed version"
        info "Forcing update to ensure latest version..."
        SHOULD_UPDATE=true
    # If we couldn't get the latest version from GitHub, force update
    elif [[ -z "$TARGET_VERSION" ]]; then
        warning "Could not determine latest version number"
        info "Will attempt update to latest..."
        SHOULD_UPDATE=true
    # Skip update if already on target version
    elif [[ "$INSTALLED_VERSION" == "$TARGET_VERSION" ]]; then
        success "Already running the target version ($TARGET_VERSION)"
        SHOULD_UPDATE=false
    else
        # Compare versions
        compare_versions "$INSTALLED_VERSION" "$TARGET_VERSION"
        result=$?
        
        if [[ $result -eq 2 ]]; then
            info "Newer version available: $INSTALLED_VERSION → $TARGET_VERSION"
            SHOULD_UPDATE=true
        elif [[ $result -eq 1 ]]; then
            warning "Installed version ($INSTALLED_VERSION) is newer than target ($TARGET_VERSION)"
            info "Downgrading to $TARGET_VERSION..."
            SHOULD_UPDATE=true
        else
            success "Already running version $TARGET_VERSION"
            SHOULD_UPDATE=false
        fi
    fi
    
    # Perform update if needed
    if [[ "$SHOULD_UPDATE" == true ]]; then
        echo ""
        info "Updating itamae..."
        
        # Ensure Go is available
        export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin
        
        # Install with appropriate version specifier
        if [[ -z "$TARGET_VERSION" ]]; then
            # No specific version available, use @latest
            go install github.com/yjmrobert/itamae@latest
        else
            # Use specific version
            go install "github.com/yjmrobert/itamae@${TARGET_VERSION}"
        fi
        
        echo ""
        success "Updated itamae"
        info "New version:"
        itamae version 2>/dev/null || itamae --version 2>/dev/null || echo "  itamae updated (version command not available)"
    else
        info "Skipping update - already on target version"
    fi
else
    info "itamae not found - starting installation"
    
    echo ""
    echo "========================================"
    echo "Installing Go (if needed)"
    echo "========================================"
    
    # Check if Go is installed, and install it if it's not
    if ! command -v go &> /dev/null; then
        info "Go is not found, installing..."
        
        # Determine the latest Go version
        info "Checking for the latest Go version..."
        GO_VERSION=$(curl -s "https://go.dev/VERSION?m=text" | head -n 1 | sed 's/go//')
        success "Latest Go version is ${GO_VERSION}"

        # Determine OS and architecture
        OS=$(uname -s | tr '[:upper:]' '[:lower:]')
        ARCH=$(uname -m)

        case $ARCH in
          x86_64)
            ARCH="amd64"
            ;;
          aarch64)
            ARCH="arm64"
            ;;
          *)
            error "Unsupported architecture: $ARCH"
            exit 1
            ;;
        esac

        # Download and install Go
        case $OS in
          linux)
            GO_TARBALL="go${GO_VERSION}.linux-${ARCH}.tar.gz"
            DOWNLOAD_URL="https://go.dev/dl/${GO_TARBALL}"
            info "Downloading Go from ${DOWNLOAD_URL}..."
            curl -L -o "/tmp/${GO_TARBALL}" "${DOWNLOAD_URL}"
            info "Extracting Go..."
            sudo tar -C /usr/local -xzf "/tmp/${GO_TARBALL}"
            rm "/tmp/${GO_TARBALL}"
            success "Go installation complete"
            ;;
          darwin)
            error "macOS is not yet supported by this install script"
            exit 1
            ;;
          *)
            error "Unsupported OS: $OS"
            exit 1
            ;;
        esac
    else
        success "Go is already installed"
    fi

    echo ""
    echo "========================================"
    echo "Configuring Environment"
    echo "========================================"
    
    # Set PATH for the current session
    export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin

    # Add Go to the user's shell profile for future sessions
    info "Configuring shell environment..."
    SHELL_CONFIG=""
    if [[ "$SHELL" == *"zsh"* ]]; then
      SHELL_CONFIG="$HOME/.zshrc"
    elif [[ "$SHELL" == *"bash"* ]]; then
      SHELL_CONFIG="$HOME/.bashrc"
    else
      warning "Unsupported shell: $SHELL"
      info "Please add /usr/local/go/bin and \$HOME/go/bin to your PATH manually"
    fi

    if [ -n "$SHELL_CONFIG" ] && [ -f "$SHELL_CONFIG" ]; then
      if ! grep -q "/usr/local/go/bin" "$SHELL_CONFIG"; then
        echo -e "\n# Go and Itamae PATH" >> "$SHELL_CONFIG"
        echo 'export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin' >> "$SHELL_CONFIG"
        success "Updated PATH in $SHELL_CONFIG"
        source "$SHELL_CONFIG" 2>/dev/null || true
      else
        info "Go path is already in your shell config"
      fi
    else
      if [ -n "$SHELL_CONFIG" ]; then
        warning "Could not find shell configuration file: $SHELL_CONFIG"
        info "Please add /usr/local/go/bin and \$HOME/go/bin to your PATH manually"
      fi
    fi

    echo ""
    echo "========================================"
    echo "Installing Itamae"
    echo "========================================"

    # Determine OS and architecture
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    case $ARCH in
      x86_64)
        ARCH="amd64"
        ;;
      aarch64)
        ARCH="arm64"
        ;;
      *)
        error "Unsupported architecture: $ARCH"
        exit 1
        ;;
    esac

    # Get version to install
    ITAMAE_VERSION=${ITAMAE_VERSION:-$(get_latest_version)}
    if [[ -z "$ITAMAE_VERSION" ]]; then
        info "Could not determine latest version, using 'latest' tag"
        info "Downloading latest itamae for ${OS}-${ARCH}..."
        DOWNLOAD_URL="https://github.com/yjmrobert/itamae/releases/latest/download/itamae-${OS}-${ARCH}"
    else
        info "Downloading itamae ${ITAMAE_VERSION} for ${OS}-${ARCH}..."
        DOWNLOAD_URL="https://github.com/yjmrobert/itamae/releases/download/${ITAMAE_VERSION}/itamae-${OS}-${ARCH}"
    fi

    info "Download URL: ${DOWNLOAD_URL}"
    
    # Download with better error handling
    if ! curl -fL --retry 3 --retry-delay 2 -o "/tmp/itamae" "${DOWNLOAD_URL}"; then
        error "Failed to download itamae from ${DOWNLOAD_URL}"
        error "Please check that the release exists and try again"
        exit 1
    fi
    
    chmod +x "/tmp/itamae"
    sudo mv "/tmp/itamae" /usr/local/bin/itamae

    success "Itamae installed"
    info "Installed version:"
    itamae version

fi

# Ensure PATH is set for running itamae
export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin

echo ""
echo "========================================"
echo "Running Itamae Install"
echo "========================================"

# Run itamae
info "Running itamae install..."
itamae install

echo ""
echo "========================================"
success "Itamae installation complete!"
echo "========================================"
