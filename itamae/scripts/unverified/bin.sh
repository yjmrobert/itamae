#!/bin/bash
#
# METADATA
# NAME: bin
# DESCRIPTION: A tool for easier management of binary tools.
# INSTALL_METHOD: binary
#

BINDIR="$HOME/.local/bin"

bin_arch() {
    case "$(uname -m)" in
        x86_64) echo amd64 ;;
        aarch64|arm64) echo arm64 ;;
        *) echo "Unsupported architecture: $(uname -m)" >&2; return 1 ;;
    esac
}

install() {
    echo "Installing bin..."
    local arch version
    arch=$(bin_arch) || return 1

    # Release assets are named bin_<version>_linux_<arch>.
    version=$(curl -fsSL https://api.github.com/repos/marcosnils/bin/releases/latest | grep -oP '"tag_name": *"v?\K[^"]+')
    if [ -z "$version" ]; then
        echo "❌ Could not determine the latest bin version." >&2
        return 1
    fi

    mkdir -p "$BINDIR"
    curl -fsSL "https://github.com/marcosnils/bin/releases/download/v${version}/bin_${version}_linux_${arch}" -o "$BINDIR/bin"
    chmod +x "$BINDIR/bin"
    echo "✅ bin ${version} installed to $BINDIR."
}

remove() {
    echo "Removing bin..."
    rm -f "$BINDIR/bin"
    echo "✅ bin removed."
}

check() {
    command -v bin &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
