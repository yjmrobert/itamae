#!/bin/bash
#
# METADATA
# NAME: kubectl
# DESCRIPTION: The command-line tool for controlling Kubernetes clusters.
# INSTALL_METHOD: binary
#

k8s_arch() {
    case "$(uname -m)" in
        x86_64) echo amd64 ;;
        aarch64|arm64) echo arm64 ;;
        *) echo "Unsupported architecture: $(uname -m)" >&2; return 1 ;;
    esac
}

install() {
    echo "Installing kubectl..."
    local arch version
    arch=$(k8s_arch) || return 1
    version=$(curl -fsSL https://dl.k8s.io/release/stable.txt)
    if [ -z "$version" ]; then
        echo "❌ Could not determine the latest kubectl version." >&2
        return 1
    fi

    curl -fsSL -o /tmp/kubectl "https://dl.k8s.io/release/${version}/bin/linux/${arch}/kubectl"
    chmod +x /tmp/kubectl
    sudo mv /tmp/kubectl /usr/local/bin/kubectl

    echo "✅ kubectl ${version} installed."
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
