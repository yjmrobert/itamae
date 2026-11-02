#!/bin/bash
#
# METADATA
# NAME: tldr (tealdeer)
# DESCRIPTION: A fast, community-driven 'man' page replacement.
# INSTALL_METHOD: binary
#

BINDIR="$HOME/.local/bin"

tealdeer_arch() {
    case "$(uname -m)" in
        x86_64) echo x86_64 ;;
        aarch64|arm64) echo aarch64 ;;
        *) echo "Unsupported architecture: $(uname -m)" >&2; return 1 ;;
    esac
}

install() {
    echo "Installing tldr (tealdeer)..."
    local arch
    arch=$(tealdeer_arch) || return 1
    # 'tealdeer' is the fast Rust client
    local tldr_url="https://github.com/tealdeer-rs/tealdeer/releases/latest/download/tealdeer-linux-${arch}-musl"
    mkdir -p "$BINDIR"
    curl -fsSL "$tldr_url" -o "$BINDIR/tldr"
    chmod +x "$BINDIR/tldr"
    echo "✅ tldr installed to $BINDIR/tldr"
}

remove() {
    echo "Removing tldr..."
    rm -f "$BINDIR/tldr"
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
