#!/bin/bash
#
# METADATA
# NAME: Dunst
# DESCRIPTION: A lightweight notification daemon for WMs.
# INSTALL_METHOD: apt
# PACKAGE_NAME: dunst
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

check() {
    command -v dunst &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
