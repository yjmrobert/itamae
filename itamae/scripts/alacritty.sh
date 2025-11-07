#!/bin/bash
#
# METADATA
# NAME: Alacritty
# OMAKASE: true
# DESCRIPTION: A fast, cross-platform, OpenGL terminal emulator.
# INSTALL_METHOD: apt
# PACKAGE_NAME: alacritty
#

install() {
    echo "Installing Alacritty..."
    if command -v nala &> /dev/null; then
        sudo nala install -y alacritty
    else
        sudo apt-get install -y alacritty
    fi
    echo "✅ Alacritty installed."
}

remove() {
    echo "Removing Alacritty..."
    sudo apt-get purge -y alacritty
    echo "✅ Alacritty removed."
}

check() {
    command -v alacritty &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
