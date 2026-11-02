#!/bin/bash
#
# METADATA
# NAME: httpie
# DESCRIPTION: A human-friendly 'curl' replacement for testing APIs.
# INSTALL_METHOD: apt
# PACKAGE_NAME: httpie
#

install() {
    echo "Installing httpie..."
    if command -v nala &> /dev/null; then
        sudo env DEBIAN_FRONTEND=noninteractive nala install -y httpie
    else
        sudo env DEBIAN_FRONTEND=noninteractive apt-get install -y httpie
    fi
    echo "✅ httpie installed."
}

remove() {
    echo "Removing httpie..."
    sudo env DEBIAN_FRONTEND=noninteractive apt-get purge -y httpie
    echo "✅ httpie removed."
}

check() {
    command -v http &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
