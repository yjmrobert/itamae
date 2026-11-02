#!/bin/bash
#
# METADATA
# NAME: Zellij
# DESCRIPTION: A modern terminal multiplexer (like tmux/screen).
# INSTALL_METHOD: binary
#

zellij_arch() {
    case "$(uname -m)" in
        x86_64) echo x86_64 ;;
        aarch64|arm64) echo aarch64 ;;
        *) echo "Unsupported architecture: $(uname -m)" >&2; return 1 ;;
    esac
}

install() {
    echo "Installing Zellij..."
    local arch
    arch=$(zellij_arch) || return 1
    local zellij_url="https://github.com/zellij-org/zellij/releases/latest/download/zellij-${arch}-unknown-linux-musl.tar.gz"
    local tarball="/tmp/zellij-${arch}.tar.gz"
    curl -fsSL "$zellij_url" -o "$tarball"
    sudo tar -xzf "$tarball" -C /usr/local/bin zellij
    rm -f "$tarball"
    echo "✅ Zellij installed."
}

remove() {
    echo "Removing Zellij..."
    sudo rm -f /usr/local/bin/zellij
    echo "✅ Zellij removed."
}

check() {
    command -v zellij &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
