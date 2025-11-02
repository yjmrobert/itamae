#!/bin/bash
#
# METADATA
# NAME: httpie
# OMAKASE: true
# DESCRIPTION: A human-friendly 'curl' replacement for testing APIs.
#

install() {
    echo "Installing httpie..."
    if command -v nala &> /dev/null; then
        sudo nala install -y httpie
    else
        sudo apt-get install -y httpie
    fi
    echo "✅ httpie installed."
}

remove() {
    echo "Removing httpie..."
    sudo apt-get purge -y httpie
    echo "✅ httpie removed."
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    *) echo "Usage: $0 {install|remove}" && exit 1 ;;
esac
