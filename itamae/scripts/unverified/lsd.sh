#!/bin/bash
#
# METADATA
# NAME: lsd (ls deluxe)
# DESCRIPTION: A modern 'ls' with pretty colors and icons.
# INSTALL_METHOD: apt
# PACKAGE_NAME: lsd
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

check() {
    command -v lsd &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
