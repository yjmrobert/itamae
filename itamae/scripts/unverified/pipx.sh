#!/bin/bash
#
# METADATA
# NAME: pipx
# OMAKASE: true
# DESCRIPTION: A tool to install and run Python applications in isolated environments.
# INSTALL_METHOD: apt
# PACKAGE_NAME: pipx
#

install() {
    echo "Installing pipx..."
    if command -v nala &> /dev/null; then
        sudo nala install -y pipx
    else
        sudo apt-get install -y pipx
    fi
    pipx ensurepath
    echo "✅ pipx installed."
}

remove() {
    echo "Removing pipx..."
    sudo apt-get purge -y pipx
    echo "✅ pipx removed."
}

check() {
    command -v pipx &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
