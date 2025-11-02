#!/bin/bash
#
# METADATA
# NAME: fd (fd-find)
# OMAKASE: true
# DESCRIPTION: A fast and user-friendly alternative to 'find'.
#

install() {
    echo "Installing fd..."
    # Debian/Ubuntu package it as 'fd-find'
    if command -v nala &> /dev/null; then
        sudo nala install -y fd-find
    else
        sudo apt-get install -y fd-find
    fi
    # Create the 'fd' symlink that all tools expect
    ln -s "$(command -v fdfind)" "$HOME/.local/bin/fd"
    echo "✅ fd installed and linked to ~/.local/bin/fd"
}

remove() {
    echo "Removing fd..."
    sudo apt-get purge -y fd-find
    rm "$HOME/.local/bin/fd"
    echo "✅ fd removed."
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    *) echo "Usage: $0 {install|remove}" && exit 1 ;;
esac
