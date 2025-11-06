#!/bin/bash
#
# METADATA
# NAME: chezmoi
# OMAKASE: true
# DESCRIPTION: Manages dotfiles across multiple machines.
# INSTALL_METHOD: binary
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

check() {
    command -v chezmoi &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
