#!/bin/bash
#
# METADATA
# NAME: Visual Studio Code
# DESCRIPTION: A popular code editor.
# INSTALL_METHOD: binary
#

vscode_arch() {
    case "$(uname -m)" in
        x86_64) echo x64 ;;
        aarch64|arm64) echo arm64 ;;
        *) echo "Unsupported architecture: $(uname -m)" >&2; return 1 ;;
    esac
}

install() {
    echo "Installing Visual Studio Code..."
    local arch
    arch=$(vscode_arch) || return 1
    curl -fL "https://code.visualstudio.com/sha/download?build=stable&os=linux-deb-${arch}" -o "/tmp/vscode-itamae.deb"
    sudo apt-get install -y "/tmp/vscode-itamae.deb"
    rm -f "/tmp/vscode-itamae.deb"
    echo "✅ Visual Studio Code installed."
}

remove() {
    echo "Removing Visual Studio Code..."
    sudo apt-get purge -y code
    echo "✅ Visual Studio Code removed."
}

check() {
    command -v code &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
