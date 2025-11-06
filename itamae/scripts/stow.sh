#!/bin/bash
#
# METADATA
# NAME: GNU Stow
# OMAKASE: true
# DESCRIPTION: A simple symlink manager for dotfiles.
# INSTALL_METHOD: apt
# PACKAGE_NAME: stow
#

install() {
    echo "Installing GNU Stow..."
    if command -v nala &> /dev/null; then
        sudo nala install -y stow
    else
        sudo apt-get install -y stow
    fi
    echo "✅ GNU Stow installed."
}

remove() {
    echo "Removing GNU Stow..."
    sudo apt-get purge -y stow
    echo "✅ GNU Stow removed."
}

check() {
    command -v stow &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
