#!/bin/bash
#
# METADATA
# NAME: ripgrep (rg)
# OMAKASE: true
# DESCRIPTION: A fast, modern replacement for grep that respects .gitignore.
#

install() {
    echo "Installing ripgrep..."
    if command -v nala &> /dev/null; then
        sudo nala install -y ripgrep
    else
        sudo apt-get install -y ripgrep
    fi
    echo "✅ ripgrep installed."
}

remove() {
    echo "Removing ripgrep..."
    sudo apt-get purge -y ripgrep
    echo "✅ ripgrep removed."
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    *) echo "Usage: $0 {install|remove}" && exit 1 ;;
esac
