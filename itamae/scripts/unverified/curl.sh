#!/bin/bash
#
# METADATA
# NAME: curl
# OMAKASE: true
# DESCRIPTION: A tool to transfer data from or to a server.
# INSTALL_METHOD: apt
# PACKAGE_NAME: curl
#

install() {
    echo "Installing curl..."
    if command -v nala &> /dev/null; then
        sudo nala install -y curl
    else
        sudo apt-get install -y curl
    fi
    echo "✅ curl installed."
}

remove() {
    echo "Removing curl..."
    sudo apt-get purge -y curl
    echo "✅ curl removed."
}

check() {
    command -v curl &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
