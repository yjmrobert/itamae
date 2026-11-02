#!/bin/bash
#
# METADATA
# NAME: fd (fd-find)
# DESCRIPTION: A fast and user-friendly alternative to 'find'.
# INSTALL_METHOD: apt
# PACKAGE_NAME: fd-find
# POST_INSTALL: post_install
#

post_install() {
    # Create the 'fd' symlink that all tools expect
    mkdir -p "$HOME/.local/bin"
    ln -sf "$(command -v fdfind)" "$HOME/.local/bin/fd"
    echo "✅ Created symlink: fd -> fdfind"
}

install() {
    echo "Installing fd..."
    # Debian/Ubuntu package it as 'fd-find'
    if command -v nala &> /dev/null; then
        sudo env DEBIAN_FRONTEND=noninteractive nala install -y fd-find
    else
        sudo env DEBIAN_FRONTEND=noninteractive apt-get install -y fd-find
    fi
    post_install
}

remove() {
    echo "Removing fd..."
    sudo env DEBIAN_FRONTEND=noninteractive apt-get purge -y fd-find
    rm -f "$HOME/.local/bin/fd"
    echo "✅ fd removed."
}

check() {
    command -v fd &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    post_install) post_install ;;
    *) echo "Usage: $0 {install|remove|check|post_install}" && exit 1 ;;
esac
