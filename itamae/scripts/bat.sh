#!/bin/bash
#
# METADATA
# NAME: bat (batcat)
# OMAKASE: true
# DESCRIPTION: A 'cat' clone with syntax highlighting and Git integration.
# INSTALL_METHOD: apt
# PACKAGE_NAME: batcat
# POST_INSTALL: post_install
#

post_install() {
    # Create the 'bat' symlink that all tools expect
    mkdir -p "$HOME/.local/bin"
    ln -sf "$(command -v batcat)" "$HOME/.local/bin/bat"
    echo "✅ Created symlink: bat -> batcat"
}

install() {
    echo "Installing bat..."
    # Debian/Ubuntu package it as 'batcat'
    if command -v nala &> /dev/null; then
        sudo nala install -y batcat
    else
        sudo apt-get install -y batcat
    fi
    post_install
}

remove() {
    echo "Removing bat..."
    sudo apt-get purge -y batcat
    rm -f "$HOME/.local/bin/bat"
    echo "✅ bat removed."
}

check() {
    command -v bat &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    post_install) post_install ;;
    *) echo "Usage: $0 {install|remove|check|post_install}" && exit 1 ;;
esac
