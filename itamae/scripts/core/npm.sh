#!/bin/bash
#
# METADATA
# NAME: npm
# DESCRIPTION: The package manager for Node.js.
# INSTALL_METHOD: apt
# PACKAGE_NAME: npm
#

install() {
    echo "Installing npm..."
    if command -v nala &> /dev/null; then
        sudo nala install -y npm
    else
        sudo apt-get install -y npm
    fi
    echo "✅ npm installed."
}

remove() {
    echo "Removing npm..."
    sudo apt-get purge -y npm
    echo "✅ npm removed."
}

check() {
    command -v npm &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
