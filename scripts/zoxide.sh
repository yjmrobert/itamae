#!/bin/bash
#
# METADATA
# NAME: zoxide
# OMAKASE: true
# DESCRIPTION: A smarter 'cd' command that remembers your directories.
#

install() {
    echo "Installing zoxide..."
    curl -sS https://raw.githubusercontent.com/ajeetdsouza/zoxide/main/install.sh | bash
    echo "✅ zoxide installed."
    echo "NOTE: You must add 'eval \"$(zoxide init zsh)\"' to your .zshrc"
}

remove() {
    echo "Removing zoxide..."
    rm "$HOME/.local/bin/zoxide"
    echo "✅ zoxide removed."
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    *) echo "Usage: $0 {install|remove}" && exit 1 ;;
esac
