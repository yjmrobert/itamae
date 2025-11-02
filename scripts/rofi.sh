#!/bin/bash
#
# METADATA
# NAME: Rofi
# OMAKASE: false
# DESCRIPTION: A fast, keyboard-driven application launcher for WMs.
#

install() {
    echo "Installing Rofi..."
    if command -v nala &> /dev/null; then
        sudo nala install -y rofi
    else
        sudo apt-get install -y rofi
    fi
    echo "✅ Rofi installed."
}

remove() {
    echo "Removing Rofi..."
    sudo apt-get purge -y rofi
    echo "✅ Rofi removed."
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    *) echo "Usage: $0 {install|remove}" && exit 1 ;;
esac
