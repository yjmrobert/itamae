#!/bin/bash
#
# METADATA
# NAME: pass
# OMAKASE: true
# DESCRIPTION: The standard unix password manager.
# INSTALL_METHOD: apt
# PACKAGE_NAME: pass
#

install() {
    echo "Installing pass..."
    if command -v nala &> /dev/null; then
        sudo nala install -y pass
    else
        sudo apt-get install -y pass
    fi
    echo "✅ pass installed."
}

remove() {
    echo "Removing pass..."
    sudo apt-get purge -y pass
    echo "✅ pass removed."
}

check() {
    command -v pass &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
