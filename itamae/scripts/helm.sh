#!/bin/bash
#
# METADATA
# NAME: Helm
# OMAKASE: true
# DESCRIPTION: The package manager for Kubernetes.
# INSTALL_METHOD: binary
#

install() {
    echo "Installing Helm..."
    curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
    echo "✅ Helm installed."
}

remove() {
    echo "Removing Helm..."
    sudo rm -f /usr/local/bin/helm
    echo "✅ Helm removed."
}

check() {
    command -v helm &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
