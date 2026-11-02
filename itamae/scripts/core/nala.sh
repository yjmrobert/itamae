#!/bin/bash
#
# METADATA
# NAME: nala
# DESCRIPTION: A more optimized package management experience.
# INSTALL_METHOD: apt
# PACKAGE_NAME: nala
#

install() {
    echo "Installing nala..."
    sudo env DEBIAN_FRONTEND=noninteractive apt-get update
    sudo env DEBIAN_FRONTEND=noninteractive apt-get install -y nala
    echo "✅ nala installed."
}

remove() {
    echo "Removing nala..."
    sudo env DEBIAN_FRONTEND=noninteractive apt-get purge -y nala
    echo "✅ nala removed."
}

check() {
    command -v nala &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
