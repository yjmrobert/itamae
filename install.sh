#!/bin/bash
# install.sh - Install or update itamae, then (optionally) run it.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/yjmrobert/itamae/master/install.sh | bash
#
# Environment variables:
#   ITAMAE_VERSION       Version to install (e.g. v1.2.3). Default: latest release.
#   ITAMAE_INSTALL_DIR   Where to put the binary. Default: /usr/local/bin
#   ITAMAE_INSTALL_ARGS  Arguments for the "itamae install" run after installing,
#                        e.g. "--category core --yes" for a fully unattended setup.
#   ITAMAE_SKIP_RUN      Set to 1 to only install the binary and not run "itamae install".
#   ITAMAE_SKIP_PREREQS  Set to 1 to skip "apt-get update" and the git/nala install.
#   ITAMAE_RELEASE_URL   Base URL of the releases (for testing mirrors).
#   ITAMAE_API_URL       URL returning the latest release JSON (for testing mirrors).

set -euo pipefail

REPO="yjmrobert/itamae"
INSTALL_DIR="${ITAMAE_INSTALL_DIR:-/usr/local/bin}"
RELEASE_URL="${ITAMAE_RELEASE_URL:-https://github.com/${REPO}/releases}"
API_URL="${ITAMAE_API_URL:-https://api.github.com/repos/${REPO}/releases/latest}"

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

info()    { echo -e "${BLUE}ℹ${NC} $1"; }
success() { echo -e "${GREEN}✓${NC} $1"; }
warning() { echo -e "${YELLOW}⚠${NC} $1"; }
error()   { echo -e "${RED}✗${NC} $1" >&2; }

header() {
    echo ""
    echo "========================================"
    echo "$1"
    echo "========================================"
}

# sudo_cmd runs a command as root, using sudo only when needed.
sudo_cmd() {
    if [ "$(id -u)" -eq 0 ]; then
        "$@"
    else
        sudo "$@"
    fi
}

# compare_versions prints -1, 0 or 1 depending on whether $1 is older, equal or newer than $2.
compare_versions() {
    local ver1="${1#v}" ver2="${2#v}"
    if [[ "$ver1" == "$ver2" ]]; then
        echo 0
        return
    fi

    local -a v1 v2
    IFS='.' read -ra v1 <<< "$ver1"
    IFS='.' read -ra v2 <<< "$ver2"

    for i in 0 1 2; do
        local num1 num2
        num1=$(echo "${v1[$i]:-0}" | grep -oE '^[0-9]+' || echo 0)
        num2=$(echo "${v2[$i]:-0}" | grep -oE '^[0-9]+' || echo 0)
        if ((num1 > num2)); then
            echo 1
            return
        elif ((num1 < num2)); then
            echo -1
            return
        fi
    done
    echo 0
}

# get_latest_version prints the latest release tag, or nothing when it cannot be determined.
get_latest_version() {
    local latest
    latest=$(curl -fsSL "$API_URL" 2>/dev/null | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' || true)
    if [[ -z "$latest" ]] && command -v git &> /dev/null; then
        latest=$(git ls-remote --tags --refs --sort="v:refname" "https://github.com/${REPO}.git" 2>/dev/null | tail -n1 | sed 's/.*\///' || true)
    fi
    echo "$latest"
}

# get_installed_version prints the installed version, "unknown" when it cannot be parsed,
# or nothing when itamae is not installed.
get_installed_version() {
    if ! command -v itamae &> /dev/null; then
        echo ""
        return
    fi
    local version_output version
    version_output=$(itamae version 2>/dev/null | head -n1 || true)
    if [[ -z "$version_output" ]]; then
        version_output=$(itamae --version 2>/dev/null | head -n1 || true)
    fi
    version=$(echo "$version_output" | grep -oE 'v?[0-9]+\.[0-9]+\.[0-9]+' | head -n1 || true)
    echo "${version:-unknown}"
}

detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    case "$(uname -m)" in
        x86_64) ARCH="amd64" ;;
        aarch64|arm64) ARCH="arm64" ;;
        *) error "Unsupported architecture: $(uname -m)"; exit 1 ;;
    esac
}

# download_url prints the asset URL for a version ("" means the latest release).
download_url() {
    local version="$1"
    if [[ -z "$version" ]]; then
        echo "${RELEASE_URL}/latest/download/itamae-${OS}-${ARCH}"
    else
        echo "${RELEASE_URL}/download/${version}/itamae-${OS}-${ARCH}"
    fi
}

# install_binary downloads the given version (or latest) into INSTALL_DIR.
install_binary() {
    local version="$1" url tmp backup=""
    url=$(download_url "$version")
    info "Download URL: ${url}"

    tmp=$(mktemp /tmp/itamae.XXXXXX)
    if ! curl -fL --retry 3 --retry-delay 2 -o "$tmp" "$url"; then
        rm -f "$tmp"
        error "Failed to download itamae from ${url}"
        error "Please check that the release exists and try again"
        exit 1
    fi
    chmod +x "$tmp"

    if [[ -f "${INSTALL_DIR}/itamae" ]]; then
        backup="${INSTALL_DIR}/itamae.backup"
        sudo_cmd cp "${INSTALL_DIR}/itamae" "$backup"
    fi

    sudo_cmd mkdir -p "$INSTALL_DIR"
    if ! sudo_cmd mv "$tmp" "${INSTALL_DIR}/itamae"; then
        rm -f "$tmp"
        if [[ -n "$backup" ]]; then
            sudo_cmd mv "$backup" "${INSTALL_DIR}/itamae"
            info "Restored previous version from backup"
        fi
        error "Failed to install itamae to ${INSTALL_DIR}"
        exit 1
    fi
    [[ -n "$backup" ]] && sudo_cmd rm -f "$backup"

    hash -r 2>/dev/null || true
    success "Installed itamae to ${INSTALL_DIR}/itamae"
    info "Installed version:"
    "${INSTALL_DIR}/itamae" version
}

install_prerequisites() {
    header "Installing Prerequisites"
    info "Updating package lists..."
    sudo_cmd apt-get update -y
    info "Installing prerequisites (git)..."
    sudo_cmd apt-get install -y git
    if ! command -v nala &> /dev/null; then
        info "Installing nala (faster apt front end)..."
        if ! sudo_cmd apt-get install -y nala; then
            warning "nala is not available on this release; itamae will use apt-get instead"
        fi
    fi
    success "Prerequisites installed"
}

run_itamae_install() {
    header "Running Itamae Install"
    local bin="${INSTALL_DIR}/itamae"
    if [[ -n "${ITAMAE_INSTALL_ARGS:-}" ]]; then
        info "Running: itamae install ${ITAMAE_INSTALL_ARGS}"
        # shellcheck disable=SC2086 # ITAMAE_INSTALL_ARGS is intentionally word-split
        "$bin" install ${ITAMAE_INSTALL_ARGS}
    elif [[ -t 1 ]] && [[ -e /dev/tty ]]; then
        # When piped through "curl | bash", stdin is the script itself; hand the
        # terminal back to itamae so its prompts work.
        "$bin" install < /dev/tty
    else
        warning "No terminal available for the interactive installer."
        info "Run 'itamae install' yourself, or set ITAMAE_INSTALL_ARGS (e.g. \"--category core --yes\")."
        return 0
    fi
}

main() {
    detect_platform

    if [[ "${ITAMAE_SKIP_PREREQS:-0}" != "1" ]]; then
        install_prerequisites
    fi

    header "Checking Itamae Installation"
    local installed target
    installed=$(get_installed_version)
    target="${ITAMAE_VERSION:-}"

    if [[ -z "$target" ]]; then
        target=$(get_latest_version)
        if [[ -n "$target" ]]; then
            info "Latest available version: $target"
        else
            warning "Could not determine the latest version number; using the 'latest' release"
        fi
    else
        info "Target version (from ITAMAE_VERSION): $target"
    fi

    if [[ -z "$installed" ]]; then
        info "itamae not found - starting installation"
        header "Installing Itamae"
        install_binary "$target"
    else
        success "itamae is already installed"
        info "Installed version: $installed"

        local should_update=false
        if [[ "$installed" == "unknown" ]]; then
            warning "Cannot determine installed version; forcing update"
            should_update=true
        elif [[ -z "$target" ]]; then
            warning "Could not determine latest version number; updating to latest"
            should_update=true
        else
            case "$(compare_versions "$installed" "$target")" in
                -1) info "Newer version available: $installed → $target"; should_update=true ;;
                1)  warning "Installed version ($installed) is newer than target ($target); downgrading"; should_update=true ;;
                0)  success "Already running the target version ($target)" ;;
            esac
        fi

        if [[ "$should_update" == true ]]; then
            header "Updating Itamae"
            install_binary "$target"
            success "Updated itamae from ${installed} to ${target:-latest}"
        else
            info "Skipping update - already on target version"
        fi
    fi

    if [[ "${ITAMAE_SKIP_RUN:-0}" == "1" ]]; then
        info "ITAMAE_SKIP_RUN is set; not running 'itamae install'"
    else
        run_itamae_install
    fi

    header "Itamae installation complete!"
}

main "$@"
