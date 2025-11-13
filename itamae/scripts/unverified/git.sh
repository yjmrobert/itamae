#!/bin/bash
#
# METADATA
# NAME: Git
# OMAKASE: true
# DESCRIPTION: A free and open source distributed version control system.
# INSTALL_METHOD: apt
# PACKAGE_NAME: git
#

install() {
    echo "Installing Git..."
    if command -v nala &> /dev/null; then
        sudo nala install -y git
    else
        sudo apt-get install -y git
    fi
    echo "✅ Git installed."
}

remove() {
    echo "Removing Git..."
    sudo apt-get purge -y git
    echo "✅ Git removed."
}

check() {
    command -v git &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
