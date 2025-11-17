#!/bin/bash
#
# METADATA
# NAME: Atuin
# DESCRIPTION: A better shell history.
# INSTALL_METHOD: binary
#

install() {
    echo "Installing Atuin..."
    curl --proto '=https' --tlsv1.2 -sSf https://setup.atuin.sh | bash
}

remove() {
    echo "Removing Atuin..."
    # As of 2024-05-21, the Atuin installer supports an uninstall command.
    curl --proto '=https' --tlsv1.2 -sSf https://setup.atuin.sh | bash -s -- --uninstall
}

check() {
    command -v atuin &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
