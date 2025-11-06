#!/bin/bash
#
# METADATA
# NAME: Rofi
# OMAKASE: false
# DESCRIPTION: A fast, keyboard-driven application launcher for WMs.
# INSTALL_METHOD: apt
# PACKAGE_NAME: rofi
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

check() {
    command -v rofi &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
