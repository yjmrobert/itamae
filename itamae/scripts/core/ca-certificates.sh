#!/bin/bash
#
# METADATA
# NAME: ca-certificates
# DESCRIPTION: Provides common CA certificates for SSL/TLS.
# INSTALL_METHOD: apt
# PACKAGE_NAME: ca-certificates
#

install() {
    echo "Installing ca-certificates..."
    if command -v nala &> /dev/null; then
        sudo env DEBIAN_FRONTEND=noninteractive nala install -y ca-certificates
    else
        sudo env DEBIAN_FRONTEND=noninteractive apt-get install -y ca-certificates
    fi
    echo "✅ ca-certificates installed."
}

remove() {
    echo "Removing ca-certificates..."
    sudo env DEBIAN_FRONTEND=noninteractive apt-get purge -y ca-certificates
    echo "✅ ca-certificates removed."
}

check() {
    dpkg -l | grep -q ca-certificates
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
