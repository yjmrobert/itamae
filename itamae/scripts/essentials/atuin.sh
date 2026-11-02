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
    echo "✅ Atuin installed."
}

remove() {
    echo "Removing Atuin..."
    # The installer places everything under ~/.atuin and sources ~/.atuin/bin/env
    # from the shell rc files; undo both.
    rm -rf "$HOME/.atuin"
    for rc in "$HOME/.bashrc" "$HOME/.zshrc" "$HOME/.profile"; do
        if [ -f "$rc" ]; then
            sed -i '/\.atuin\/bin\/env/d' "$rc"
        fi
    done
    echo "✅ Atuin removed."
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
