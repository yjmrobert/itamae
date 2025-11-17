#!/bin/bash
#
# METADATA
# NAME: kubectl
# DESCRIPTION: The command-line tool for controlling Kubernetes clusters.
# INSTALL_METHOD: binary
#

install() {
    echo "Installing kubectl..."
    
    # Download the latest stable version
    curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
    chmod +x kubectl
    sudo mv kubectl /usr/local/bin/
    
    echo "✅ kubectl installed."
}

remove() {
    echo "Removing kubectl..."
    sudo rm -f /usr/local/bin/kubectl
    echo "✅ kubectl removed."
}

check() {
    command -v kubectl &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
