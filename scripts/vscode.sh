#!/bin/bash
#
# METADATA
# NAME: Visual Studio Code
# OMAKASE: false
# DESCRIPTION: The popular code editor.
#

install() {
    echo "Installing Visual Studio Code..."
    local DEB_URL="https://code.visualstudio.com/sha/download?build=stable&os=linux-deb-x64"
    local TMP_DEB="/tmp/vscode-itamae.deb"

    curl -fL "$DEB_URL" -o "$TMP_DEB"
    sudo dpkg -i "$TMP_DEB"
    sudo apt-get install -f -y # Fix dependencies
    rm "$TMP_DEB"
    echo "✅ Visual Studio Code installed."
}

remove() {
    echo "Removing Visual Studio Code..."
    sudo apt-get purge -y code
    echo "✅ Visual Studio Code removed."
}

# --- ROUTER ---
case "$1" in
    install)
        install
        ;;
    remove)
        remove
        ;;
    *)
        echo "Usage: $0 {install|remove}"
        exit 1
        ;;
esac
