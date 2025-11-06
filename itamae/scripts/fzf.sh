#!/bin/bash
#
# METADATA
# NAME: fzf
# OMAKASE: true
# DESCRIPTION: A command-line fuzzy finder.
# INSTALL_METHOD: apt
# PACKAGE_NAME: fzf
#

install() {
    echo "Installing fzf..."
    if command -v nala &> /dev/null; then
        sudo nala install -y fzf
    else
        sudo apt-get install -y fzf
    fi
    echo "✅ fzf installed."
}

remove() {
    echo "Removing fzf..."
    sudo apt-get purge -y fzf
    echo "✅ fzf removed."
}

check() {
    command -v fzf &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
