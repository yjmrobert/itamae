#!/bin/bash
#
# METADATA
# NAME: Starship
# OMAKASE: true
# DESCRIPTION: The minimal, fast, and customizable prompt.
#

install() {
    echo "Installing Starship..."
    # Use -y to bypass the prompt
    curl -sS https://starship.rs/install.sh | sh -s -- -y
    echo "✅ Starship installed."
    echo "NOTE: You must add 'eval \"$(starship init zsh)\"' to your .zshrc"
}

remove() {
    echo "Removing Starship..."
    sh -c 'rm "$(command -v starship)"'
    echo "✅ Starship removed."
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    *) echo "Usage: $0 {install|remove}" && exit 1 ;;
esac
