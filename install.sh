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
        info "Updating itamae from ${INSTALLED_VERSION} to ${TARGET_VERSION}..."
        
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
        
        # Determine download URL
        if [[ -z "$TARGET_VERSION" ]]; then
            info "Downloading latest itamae for ${OS}-${ARCH}..."
            DOWNLOAD_URL="https://github.com/yjmrobert/itamae/releases/latest/download/itamae-${OS}-${ARCH}"
        else
            info "Downloading itamae ${TARGET_VERSION} for ${OS}-${ARCH}..."
            DOWNLOAD_URL="https://github.com/yjmrobert/itamae/releases/download/${TARGET_VERSION}/itamae-${OS}-${ARCH}"
        fi
        
        info "Download URL: ${DOWNLOAD_URL}"
        
        # Backup current binary
        if [[ -f /usr/local/bin/itamae ]]; then
            info "Creating backup of current version..."
            sudo cp /usr/local/bin/itamae /usr/local/bin/itamae.backup
        fi
        
        # Download with error handling
        if ! curl -fL --retry 3 --retry-delay 2 -o "/tmp/itamae" "${DOWNLOAD_URL}"; then
            error "Failed to download itamae from ${DOWNLOAD_URL}"
            error "Please check that the release exists and try again"
            # Restore backup if download failed
            if [[ -f /usr/local/bin/itamae.backup ]]; then
                sudo mv /usr/local/bin/itamae.backup /usr/local/bin/itamae
                info "Restored previous version from backup"
            fi
            exit 1
        fi
        
        chmod +x "/tmp/itamae"
        sudo mv "/tmp/itamae" /usr/local/bin/itamae
        
        # Clean up backup on success
        sudo rm -f /usr/local/bin/itamae.backup
        
        echo ""
        success "Updated itamae to ${TARGET_VERSION}"
        info "New version:"
        itamae version
        
        # After successful update, run itamae install
        echo ""
        echo "========================================"
        echo "Running Itamae Install"
        echo "========================================"
        info "Running itamae install..."
        itamae install
        
        echo ""
        echo "========================================"
        success "Itamae installation complete!"
        echo "========================================"
        exit 0
    else
        info "Skipping update - already on target version"
        
        # Even if no update needed, run itamae install
        echo ""
        echo "========================================"
        echo "Running Itamae Install"
        echo "========================================"
        info "Running itamae install..."
        itamae install
        
        echo ""
        echo "========================================"
        success "Itamae installation complete!"
        echo "========================================"
        exit 0
    fi
else
    info "itamae not found - starting installation"
    
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
    
    # After successful installation, run itamae install
    echo ""
    echo "========================================"
    echo "Running Itamae Install"
    echo "========================================"
    info "Running itamae install..."
    itamae install
    
    echo ""
    echo "========================================"
    success "Itamae installation complete!"
    echo "========================================"
fi
