#!/bin/bash
#
# METADATA
# NAME: gnupg
# DESCRIPTION: The GNU Privacy Guard, for encryption and signing.
# INSTALL_METHOD: apt
# PACKAGE_NAME: gnupg
#

install() {
    echo "Installing gnupg..."
    if command -v nala &> /dev/null; then
        sudo nala install -y gnupg
    else
        sudo apt-get install -y gnupg
    fi
    echo "✅ gnupg installed."
}

remove() {
    echo "Removing gnupg..."
    sudo apt-get purge -y gnupg
    echo "✅ gnupg removed."
}

check() {
    command -v gpg &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
