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
    local zellij_url="https://github.com/zellij-project/zellij/releases/latest/download/zellij-${arch}-unknown-linux-musl.tar.gz"
    curl -fsSL "$zellij_url" | sudo tar -xz -C /usr/local/bin zellij
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
