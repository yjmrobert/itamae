#!/bin/bash
#
# METADATA
# NAME: Zellij
# OMAKASE: true
# DESCRIPTION: A modern terminal multiplexer (like tmux/screen).
#

install() {
    echo "Installing Zellij..."
    # Install from binary release
    local ZELLIJ_URL="https://github.com/zellij-project/zellij/releases/latest/download/zellij-x86_64-unknown-linux-musl.tar.gz"
    curl -L "$ZELLIJ_URL" | sudo tar -xz -C /usr/local/bin
    echo "✅ Zellij installed."
}

remove() {
    echo "Removing Zellij..."
    sudo rm /usr/local/bin/zellij
    echo "✅ Zellij removed."
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    *) echo "Usage: $0 {install|remove}" && exit 1 ;;
esac
