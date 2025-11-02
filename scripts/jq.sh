#!/bin/bash
#
# METADATA
# NAME: jq
# OMAKASE: true
# DESCRIPTION: A lightweight and flexible command-line JSON processor.
#

install() {
    echo "Installing jq..."
    if command -v nala &> /dev/null; then
        sudo nala install -y jq
    else
        sudo apt-get install -y jq
    fi
    echo "✅ jq installed."
}

remove() {
    echo "Removing jq..."
    sudo apt-get purge -y jq
    echo "✅ jq removed."
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    *) echo "Usage: $0 {install|remove}" && exit 1 ;;
esac
