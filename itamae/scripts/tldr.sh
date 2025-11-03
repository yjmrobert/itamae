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
    local TLDR_URL="https://github.com/tealdeer-rs/tealdeer/releases/latest/download/tealdeer-linux-x86_64-musl"
    mkdir -p "$HOME/.local/bin"
    curl -L "$TLDR_URL" -o "$HOME/.local/bin/tldr"
    chmod +x "$HOME/.local/bin/tldr"
    echo "✅ tldr installed to ~/.local/bin/tldr"
}

remove() {
    echo "Removing tldr..."
    rm "$HOME/.local/bin/tldr"
    echo "✅ tldr removed."
}

check() {
    command -v tldr &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
