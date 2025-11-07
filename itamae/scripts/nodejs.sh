#!/bin/bash
#
# METADATA
# NAME: Node.js
# OMAKASE: true
# DESCRIPTION: A JavaScript runtime environment.
# INSTALL_METHOD: binary
#

install() {
    echo "Installing Node.js..."
    # Add NodeSource repository
    curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -
    if command -v nala &> /dev/null; then
        sudo nala install -y nodejs
    else
        sudo apt-get install -y nodejs
    fi
    echo "✅ Node.js installed."
}

remove() {
    echo "Removing Node.js..."
    sudo apt-get purge -y nodejs
    sudo rm -f /etc/apt/sources.list.d/nodesource.list*
    sudo rm -f /etc/apt/keyrings/nodesource.gpg
    echo "✅ Node.js removed."
}

check() {
    command -v node &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
