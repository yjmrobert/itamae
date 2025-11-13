#!/bin/bash
#
# METADATA
# NAME: Rust
# OMAKASE: true
# DESCRIPTION: A multi-paradigm, general-purpose programming language.
# INSTALL_METHOD: binary
#

install() {
    echo "Installing Rust..."
    curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
    echo "✅ Rust installed."
}

remove() {
    echo "Removing Rust..."
    rustup self uninstall -y
    echo "✅ Rust removed."
}

check() {
    command -v rustc &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
