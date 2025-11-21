#!/bin/bash
#
# METADATA
# NAME: Node.js
# DESCRIPTION: A JavaScript runtime environment.
# INSTALL_METHOD: apt
# PACKAGE_NAME: nodejs
# REPO_SETUP: setup_repo
#

setup_repo() {
    echo "Setting up NodeSource repository..."
    curl --silent -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash - > /dev/null 2>&1
    echo "✅ NodeSource repository configured."
}

install() {
    echo "Installing Node.js..."
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
    setup_repo) setup_repo ;;
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {setup_repo|install|remove|check}" && exit 1 ;;
esac
