#!/bin/bash
#
# METADATA
# NAME: Polybar
# DESCRIPTION: A fast and easy-to-use status bar for X11.
# INSTALL_METHOD: apt
# PACKAGE_NAME: polybar
#

install() {
    echo "Installing Polybar..."
    if command -v nala &> /dev/null; then
        sudo nala install -y polybar
    else
        sudo apt-get install -y polybar
    fi
    echo "✅ Polybar installed."
}

remove() {
    echo "Removing Polybar..."
    sudo apt-get purge -y polybar
    echo "✅ Polybar removed."
}

check() {
    command -v polybar &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
