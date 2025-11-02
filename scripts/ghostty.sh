#!/bin/bash
#
# METADATA
# NAME: Ghostty
# OMAKASE: false
# DESCRIPTION: A GPU-accelerated terminal emulator.
#

install() {
    echo "Installing Ghostty..."
    echo "Ghostty installation is manual."
    echo "Please download the binary from: https://github.com/ghostty-org/ghostty"
    echo "This is a placeholder for a future automated script."
    # A real script would curl the .deb
}

remove() {
    echo "Removing Ghostty..."
    echo "Please remove the binary manually."
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    *) echo "Usage: $0 {install|remove}" && exit 1 ;;
esac
