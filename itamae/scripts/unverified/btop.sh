#!/bin/bash
#
# METADATA
# NAME: btop
# DESCRIPTION: A beautiful, modern resource monitor.
# INSTALL_METHOD: apt
# PACKAGE_NAME: btop
#

install() {
    echo "Installing btop..."
    # btop is in modern repos
    if command -v nala &> /dev/null; then
        sudo env DEBIAN_FRONTEND=noninteractive nala install -y btop
    else
        sudo env DEBIAN_FRONTEND=noninteractive apt-get install -y btop
    fi
    echo "✅ btop installed."
}

remove() {
    echo "Removing btop..."
    sudo env DEBIAN_FRONTEND=noninteractive apt-get purge -y btop
    echo "✅ btop removed."
}

check() {
    command -v btop &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
