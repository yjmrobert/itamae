#!/bin/bash
#
# METADATA
# NAME: jq
# DESCRIPTION: A lightweight and flexible command-line JSON processor.
# INSTALL_METHOD: apt
# PACKAGE_NAME: jq
#

install() {
    echo "Installing jq..."
    if command -v nala &> /dev/null; then
        sudo nala install -y jq
    else
        sudo apt-get install -y jq
    fi
    echo "✅ jq installed."
}

remove() {
    echo "Removing jq..."
    sudo apt-get purge -y jq
    echo "✅ jq removed."
}

check() {
    command -v jq &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
