#!/bin/bash
#
# METADATA
# NAME: bin
# OMAKASE: true
# DESCRIPTION: A tool for easier management of binary tools.
# INSTALL_METHOD: binary
#

BINDIR="$HOME/.local/bin"

install() {
    echo "Installing bin..."
    curl -sL https://raw.githubusercontent.com/marcosnils/bin/master/install.sh | bash
    echo "✅ bin installed."
}

remove() {
    echo "Removing bin..."
    rm -f "$BINDIR/bin"
    echo "✅ bin removed."
}

check() {
    command -v bin &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
