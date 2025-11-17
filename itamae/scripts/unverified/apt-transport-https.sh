#!/bin/bash
#
# METADATA
# NAME: apt-transport-https
# DESCRIPTION: Allows using repositories over HTTPS.
# INSTALL_METHOD: apt
# PACKAGE_NAME: apt-transport-https
#

install() {
    echo "Installing apt-transport-https..."
    if command -v nala &> /dev/null; then
        sudo nala install -y apt-transport-https
    else
        sudo apt-get install -y apt-transport-https
    fi
    echo "✅ apt-transport-https installed."
}

remove() {
    echo "Removing apt-transport-https..."
    sudo apt-get purge -y apt-transport-https
    echo "✅ apt-transport-https removed."
}

check() {
    dpkg -l | grep -q apt-transport-https
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
