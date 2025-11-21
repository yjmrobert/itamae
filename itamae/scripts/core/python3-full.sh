#!/bin/bash
#
# METADATA
# NAME: python3-full
# DESCRIPTION: The complete Python 3 programming language environment.
# INSTALL_METHOD: apt
# PACKAGE_NAME: python3-full
#

install() {
    echo "Installing python3-full..."
    if command -v nala &> /dev/null; then
        sudo nala install -y python3-full
    else
        sudo apt-get install -y python3-full
    fi
    echo "✅ python3-full installed."
}

remove() {
    echo "Removing python3-full..."
    sudo apt-get purge -y python3-full
    echo "✅ python3-full removed."
}

check() {
    command -v python3 &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
