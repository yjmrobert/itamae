#!/bin/bash
#
# METADATA
# NAME: Flameshot
# DESCRIPTION: A powerful, scriptable screenshot tool.
# INSTALL_METHOD: apt
# PACKAGE_NAME: flameshot
#

install() {
    echo "Installing Flameshot..."
    if command -v nala &> /dev/null; then
        sudo env DEBIAN_FRONTEND=noninteractive nala install -y flameshot
    else
        sudo env DEBIAN_FRONTEND=noninteractive apt-get install -y flameshot
    fi
    echo "✅ Flameshot installed."
}

remove() {
    echo "Removing Flameshot..."
    sudo env DEBIAN_FRONTEND=noninteractive apt-get purge -y flameshot
    echo "✅ Flameshot removed."
}

check() {
    command -v flameshot &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
