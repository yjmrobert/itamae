#!/bin/bash
#
# METADATA
# NAME: tldr (tealdeer)
# OMAKASE: true
# DESCRIPTION: A fast, community-driven 'man' page replacement.
#

install() {
    echo "Installing tldr (tealdeer)..."
    # 'tealdeer' is the fast Rust client
    if command -v nala &> /dev/null; then
        sudo nala install -y tealdeer
    else
        sudo apt-get install -y tealdeer
    fi
    # Link 'tldr' to 'tldr'
    mkdir -p "$HOME/.local/bin"
    ln -s "$(command -v tldr)" "$HOME/.local/bin/tldr"
    echo "✅ tldr installed and linked to ~/.local/bin/tldr"
}

remove() {
    echo "Removing tldr..."
    sudo apt-get purge -y tealdeer
    rm "$HOME/.local/bin/tldr"
    echo "✅ tldr removed."
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    *) echo "Usage: $0 {install|remove}" && exit 1 ;;
esac
