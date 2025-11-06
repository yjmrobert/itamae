#!/bin/bash
#
# METADATA
# NAME: Ghostty
# OMAKASE: false
# DESCRIPTION: A GPU-accelerated terminal emulator.
# INSTALL_METHOD: manual
#

install() {
    echo "Installing Ghostty..."
    echo "Ghostty installation is manual."
    echo "Please download the binary from: https://github.com/ghostty-org/ghostty"
    echo "This is a placeholder for a future automated script."
    # A real script would curl the .deb
    mkdir -p "$HOME/.local/bin"
}

remove() {
    echo "Removing Ghostty..."
    echo "Please remove the binary manually."
    rm -f "$HOME/.local/bin/ghostty"
}

check() {
    command -v ghostty &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
