#!/bin/bash
#
# METADATA
# NAME: Flameshot
# OMAKASE: false
# DESCRIPTION: A powerful, scriptable screenshot tool.
#

install() {
    echo "Installing Flameshot..."
    if command -v nala &> /dev/null; then
        sudo nala install -y flameshot
    else
        sudo apt-get install -y flameshot
    fi
    echo "✅ Flameshot installed."
}

remove() {
    echo "Removing Flameshot..."
    sudo apt-get purge -y flameshot
    echo "✅ Flameshot removed."
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    *) echo "Usage: $0 {install|remove}" && exit 1 ;;
esac
