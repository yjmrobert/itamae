#!/bin/bash
#
# METADATA
# NAME: kubecolor
# DESCRIPTION: A tool to colorize kubectl output.
# INSTALL_METHOD: binary
#

BINDIR="$HOME/.local/bin"

kubecolor_arch() {
    case "$(uname -m)" in
        x86_64) echo amd64 ;;
        aarch64|arm64) echo arm64 ;;
        *) echo "Unsupported architecture: $(uname -m)" >&2; return 1 ;;
    esac
}

install() {
    echo "Installing kubecolor..."
    local arch version tarball
    arch=$(kubecolor_arch) || return 1

    version=$(curl -fsSL https://api.github.com/repos/kubecolor/kubecolor/releases/latest | grep -oP '"tag_name": *"\K[^"]+')
    if [ -z "$version" ]; then
        echo "❌ Could not determine the latest kubecolor version." >&2
        return 1
    fi

    tarball="/tmp/kubecolor_${version#v}_linux_${arch}.tar.gz"
    curl -fsSL "https://github.com/kubecolor/kubecolor/releases/download/${version}/kubecolor_${version#v}_linux_${arch}.tar.gz" -o "$tarball"
    mkdir -p "$BINDIR"
    tar -xzf "$tarball" -C "$BINDIR" kubecolor
    chmod +x "$BINDIR/kubecolor"
    rm -f "$tarball"

    echo "✅ kubecolor ${version} installed to $BINDIR."
}

remove() {
    echo "Removing kubecolor..."
    rm -f "$BINDIR/kubecolor"
    echo "✅ kubecolor removed."
}

check() {
    command -v kubecolor &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
