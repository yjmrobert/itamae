#!/bin/bash
#
# METADATA
# NAME: Ruby
# DESCRIPTION: A dynamic, open source programming language.
# INSTALL_METHOD: apt
# PACKAGE_NAME: ruby-full
#

install() {
    echo "Installing Ruby..."
    if command -v nala &> /dev/null; then
        sudo nala install -y ruby-full
    else
        sudo apt-get install -y ruby-full
    fi
    echo "✅ Ruby installed."
}

remove() {
    echo "Removing Ruby..."
    sudo apt-get purge -y ruby-full
    echo "✅ Ruby removed."
}

check() {
    command -v ruby &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
