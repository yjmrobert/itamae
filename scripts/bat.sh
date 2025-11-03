#!/bin/bash
#
# METADATA
# NAME: bat (batcat)
# OMAKASE: true
# DESCRIPTION: A 'cat' clone with syntax highlighting and Git integration.
#

install() {
    echo "Installing bat..."
    # Debian/Ubuntu package it as 'batcat'
    if command -v nala &> /dev/null; then
        sudo nala install -y batcat
    else
        sudo apt-get install -y batcat
    fi
    # Create the 'bat' symlink that all tools expect
    mkdir -p "$HOME/.local/bin"
    ln -s "$(command -v batcat)" "$HOME/.local/bin/bat"
    echo "✅ bat installed and linked to ~/.local/bin/bat"
}

remove() {
    echo "Removing bat..."
    sudo apt-get purge -y batcat
    rm "$HOME/.local/bin/bat"
    echo "✅ bat removed."
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    *) echo "Usage: $0 {install|remove}" && exit 1 ;;
esac
