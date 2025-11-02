#!/bin/bash
#
# METADATA
# NAME: chezmoi
# OMAKASE: true
# DESCRIPTION: Manages dotfiles across multiple machines.
#

BINDIR="$HOME/.local/bin"

install() {
    echo "Installing chezmoi..."
    # Install from get.chezmoi.io
    sh -c "$(curl -fsLS get.chezmoi.io)" -- -b "$BINDIR"
}

remove() {
    echo "Removing chezmoi..."
    # The installation script places the binary in ~/.local/bin
    rm "$HOME/.local/bin/chezmoi"
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    *) echo "Usage: $0 {install|remove}" && exit 1 ;;
esac
