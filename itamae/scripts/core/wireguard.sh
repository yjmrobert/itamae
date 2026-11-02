#!/bin/bash
#
# METADATA
# NAME: WireGuard
# DESCRIPTION: Fast, modern VPN tooling (wg, wg-quick); the kernel module is built in.
# INSTALL_METHOD: apt
# PACKAGE_NAME: wireguard-tools
#

install() {
    echo "Installing WireGuard tools..."
    if command -v nala &> /dev/null; then
        sudo nala install -y wireguard-tools
    else
        sudo apt-get install -y wireguard-tools
    fi
    echo "✅ WireGuard tools installed."
}

remove() {
    echo "Removing WireGuard tools..."
    sudo apt-get purge -y wireguard-tools
    echo "✅ WireGuard tools removed."
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
