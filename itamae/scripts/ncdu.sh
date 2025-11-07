#!/bin/bash
#
# METADATA
# NAME: ncdu
# OMAKASE: false
# DESCRIPTION: A disk usage analyzer with an ncurses interface.
# INSTALL_METHOD: apt
# PACKAGE_NAME: ncdu
#

install() {
    echo "Installing ncdu..."
    if command -v nala &> /dev/null; then
        sudo nala install -y ncdu
    else
        sudo apt-get install -y ncdu
    fi
    echo "✅ ncdu installed."
}

remove() {
    echo "Removing ncdu..."
    sudo apt-get purge -y ncdu
    echo "✅ ncdu removed."
}

check() {
    command -v ncdu &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
