#!/bin/bash
#
# METADATA
# NAME: lsd (ls deluxe)
# OMAKASE: true
# DESCRIPTION: A modern 'ls' with pretty colors and icons.
#

install() {
    echo "Installing lsd..."
    if command -v nala &> /dev/null; then
        sudo nala install -y lsd
    else
        sudo apt-get install -y lsd
    fi
    echo "✅ lsd installed."
}

remove() {
    echo "Removing lsd..."
    sudo apt-get purge -y lsd
    echo "✅ lsd removed."
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    *) echo "Usage: $0 {install|remove}" && exit 1 ;;
esac
