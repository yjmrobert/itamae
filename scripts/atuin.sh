#!/bin/bash
#
# METADATA
# NAME: Atuin
# OMAKASE: true
# DESCRIPTION: Magical shell history, synced across machines.
#

install() {
    echo "Installing Atuin..."
    # Atuin is best installed from its own script
    bash -c "$(curl --proto '=https' --tlsv1.2 -sSf https://setup.atuin.sh)"
    echo "✅ Atuin installed."
    echo "NOTE: You must add 'eval \"$(atuin init zsh)\"' to your .zshrc"
}

remove() {
    echo "Removing Atuin..."
    # The script provides an uninstall
    bash -c "$(curl --proto '=https' --tlsv1.2 -sSf https://setup.atuin.sh)" -- --uninstall
    echo "✅ Atuin removed."
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    *) echo "Usage: $0 {install|remove}" && exit 1 ;;
esac
