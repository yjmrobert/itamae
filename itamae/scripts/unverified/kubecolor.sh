#!/bin/bash
#
# METADATA
# NAME: kubecolor
# DESCRIPTION: A tool to colorize kubectl output.
# INSTALL_METHOD: binary
#

BINDIR="$HOME/.local/bin"

install() {
    echo "Installing kubecolor..."
    
    # Get the latest version
    KUBECOLOR_VERSION=$(curl -s https://api.github.com/repos/kubecolor/kubecolor/releases/latest | grep -oP '"tag_name": "\K(.*)(?=")')
    
    # Download and install
    wget https://github.com/kubecolor/kubecolor/releases/download/${KUBECOLOR_VERSION}/kubecolor_${KUBECOLOR_VERSION#v}_linux_amd64.tar.gz -P /tmp
    tar -xzf /tmp/kubecolor_${KUBECOLOR_VERSION#v}_linux_amd64.tar.gz -C /tmp
    mkdir -p "$BINDIR"
    mv /tmp/kubecolor "$BINDIR/"
    chmod +x "$BINDIR/kubecolor"
    rm /tmp/kubecolor_${KUBECOLOR_VERSION#v}_linux_amd64.tar.gz
    
    echo "✅ kubecolor installed."
}

remove() {
    echo "Removing kubecolor..."
    rm -f "$BINDIR/kubecolor"
    echo "✅ kubecolor removed."
}

check() {
    command -v kubecolor &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
