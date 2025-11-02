#!/bin/bash
#
# METADATA
# NAME: Atuin
# OMAKASE: true
# DESCRIPTION: A better shell history.
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

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    *) echo "Usage: $0 {install|remove}" && exit 1 ;;
esac
