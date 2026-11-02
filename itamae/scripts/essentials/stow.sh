#!/bin/bash
#
# METADATA
# NAME: GNU Stow
# DESCRIPTION: A simple symlink manager for dotfiles.
# INSTALL_METHOD: apt
# PACKAGE_NAME: stow
#

install() {
    echo "Installing GNU Stow..."
    if command -v nala &> /dev/null; then
        sudo env DEBIAN_FRONTEND=noninteractive nala install -y stow
    else
        sudo env DEBIAN_FRONTEND=noninteractive apt-get install -y stow
    fi
    echo "✅ GNU Stow installed."
}

remove() {
    echo "Removing GNU Stow..."
    sudo env DEBIAN_FRONTEND=noninteractive apt-get purge -y stow
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
