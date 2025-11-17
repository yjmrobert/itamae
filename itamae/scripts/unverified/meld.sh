#!/bin/bash
#
# METADATA
# NAME: Meld
# DESCRIPTION: A visual diff and merge tool for developers.
# INSTALL_METHOD: apt
# PACKAGE_NAME: meld
#

install() {
    echo "Installing Meld..."
    if command -v nala &> /dev/null; then
        sudo nala install -y meld
    else
        sudo apt-get install -y meld
    fi
    echo "✅ Meld installed."
}

remove() {
    echo "Removing Meld..."
    sudo apt-get purge -y meld
    echo "✅ Meld removed."
}

check() {
    command -v meld &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
