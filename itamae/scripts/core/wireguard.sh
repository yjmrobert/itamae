#!/bin/bash
#
# METADATA
# NAME: wireguard
# DESCRIPTION: A fast, modern, and secure VPN tunnel.
# INSTALL_METHOD: apt
# PACKAGE_NAME: wireguard
#

install() {
    echo "Installing wireguard..."
    if command -v nala &> /dev/null; then
        sudo nala install -y wireguard
    else
        sudo apt-get install -y wireguard
    fi
    echo "✅ wireguard installed."
}

remove() {
    echo "Removing wireguard..."
    sudo apt-get purge -y wireguard
    echo "✅ wireguard removed."
}

check() {
    command -v wg &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
