#!/bin/bash
#
# METADATA
# NAME: Rust
# DESCRIPTION: A multi-paradigm, general-purpose programming language.
# INSTALL_METHOD: binary
#

ensure_deps() {
    # rustc needs a C linker; without build-essential every cargo build fails.
    if ! command -v cc &> /dev/null; then
        echo "Installing build-essential (C toolchain for the Rust linker)..."
        if command -v nala &> /dev/null; then
            sudo nala install -y build-essential
        else
            sudo apt-get install -y build-essential
        fi
    fi
}

install() {
    echo "Installing Rust..."
    ensure_deps
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
