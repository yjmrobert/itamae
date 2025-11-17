#!/bin/bash
#
# METADATA
# NAME: wget
# DESCRIPTION: A utility for non-interactive download of files from the web.
# INSTALL_METHOD: apt
# PACKAGE_NAME: wget
#

install() {
    echo "Installing wget..."
    if command -v nala &> /dev/null; then
        sudo nala install -y wget
    else
        sudo apt-get install -y wget
    fi
    echo "✅ wget installed."
}

remove() {
    echo "Removing wget..."
    sudo apt-get purge -y wget
    echo "✅ wget removed."
}

check() {
    command -v wget &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
