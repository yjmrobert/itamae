#!/bin/bash
#
# METADATA
# NAME: btop
# OMAKASE: true
# DESCRIPTION: A beautiful, modern resource monitor.
#

install() {
    echo "Installing btop..."
    # btop is in modern repos
    if command -v nala &> /dev/null; then
        sudo nala install -y btop
    else
        sudo apt-get install -y btop
    fi
    echo "✅ btop installed."
}

remove() {
    echo "Removing btop..."
    sudo apt-get purge -y btop
    echo "✅ btop removed."
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    *) echo "Usage: $0 {install|remove}" && exit 1 ;;
esac
