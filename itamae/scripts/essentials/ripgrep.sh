#!/bin/bash
#
# METADATA
# NAME: ripgrep (rg)
# DESCRIPTION: A fast, modern replacement for grep that respects .gitignore.
# INSTALL_METHOD: apt
# PACKAGE_NAME: ripgrep
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

check() {
    command -v rg &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
