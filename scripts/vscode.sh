#!/bin/bash
#
# METADATA
# NAME: Visual Studio Code
# OMAKASE: true
# DESCRIPTION: A popular code editor.
#

install() {
    echo "Installing Visual Studio Code..."
    # Download the .deb package
    curl -fL "https://code.visualstudio.com/sha/download?build=stable&os=linux-deb-x64" -o "/tmp/vscode-itamae.deb"
    # Install the package
    sudo apt-get install -y "/tmp/vscode-itamae.deb"
    # Clean up
    rm "/tmp/vscode-itamae.deb"
}

remove() {
    echo "Removing Visual Studio Code..."
    sudo apt-get purge -y code
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    *) echo "Usage: $0 {install|remove}" && exit 1 ;;
esac
