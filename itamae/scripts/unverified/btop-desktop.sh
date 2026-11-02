#!/bin/bash
#
# METADATA
# NAME: btop-desktop
# DESCRIPTION: A resource monitor that shows usage and stats for processor, memory, disks, network, and processes.
# INSTALL_METHOD: apt
# PACKAGE_NAME: btop
#

install() {
    echo "Installing btop-desktop..."
    # btop is in modern repos
    if command -v nala &> /dev/null; then
        sudo env DEBIAN_FRONTEND=noninteractive nala install -y btop
    else
        sudo env DEBIAN_FRONTEND=noninteractive apt-get install -y btop
    fi
    echo "✅ btop-desktop installed."
}

remove() {
    echo "Removing btop-desktop..."
    sudo env DEBIAN_FRONTEND=noninteractive apt-get purge -y btop
    echo "✅ btop-desktop removed."
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
