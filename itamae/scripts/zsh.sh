#!/bin/bash
#
# METADATA
# NAME: Zsh
# OMAKASE: true
# DESCRIPTION: The Z Shell, a powerful foundation for the terminal.
# INSTALL_METHOD: apt
# PACKAGE_NAME: zsh
#

install() {
    echo "Installing Zsh..."
    if command -v nala &> /dev/null; then
        sudo nala install -y zsh
    else
        sudo apt-get install -y zsh
    fi
    echo "✅ Zsh installed."
    echo "Run 'chsh -s $(which zsh)' to make it your default."
}

remove() {
    echo "Removing Zsh..."
    sudo apt-get purge -y zsh
    echo "✅ Zsh removed."
}

check() {
    command -v zsh &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
