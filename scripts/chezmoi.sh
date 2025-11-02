#!/bin/bash
#
# METADATA
# NAME: Chezmoi
# OMAKASE: false
# DESCRIPTION: The powerful, template-driven dotfile manager.
#

install() {
    echo "Installing Chezmoi..."
    local BINDIR="$HOME/.local/bin"
    sh -c "$(curl -fsLS get.chezmoi.io)" -- -b "$BINDIR"
    echo "✅ Chezmoi installed to $BINDIR"
}

remove() {
    echo "Removing Chezmoi..."
    rm "$HOME/.local/bin/chezmoi"
    echo "✅ Chezmoi removed."
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    *) echo "Usage: $0 {install|remove}" && exit 1 ;;
esac
