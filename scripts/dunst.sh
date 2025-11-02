#!/bin/bash
#
# METADATA
# NAME: Dunst
# OMAKASE: false
# DESCRIPTION: A lightweight notification daemon for WMs.
#

install() {
    echo "Installing Dunst..."
    if command -v nala &> /dev/null; then
        sudo nala install -y dunst
    else
        sudo apt-get install -y dunst
    fi
    echo "✅ Dunst installed."
}

remove() {
    echo "Removing Dunst..."
    sudo apt-get purge -y dunst
    echo "✅ Dunst removed."
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    *) echo "Usage: $0 {install|remove}" && exit 1 ;;
esac
