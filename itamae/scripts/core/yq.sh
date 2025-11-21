#!/bin/bash
#
# METADATA
# NAME: yq (Go)
# DESCRIPTION: A 'jq' for YAML. (Installs the correct Go binary, not the python wrapper).
# INSTALL_METHOD: binary
#

install() {
    echo "Installing yq (Go binary)..."
    # This is critical: apt 'yq' is the wrong tool.
    local YQ_URL="https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64"
    sudo curl --silent -L "$YQ_URL" -o /usr/local/bin/yq
    sudo chmod +x /usr/local/bin/yq
    echo "✅ yq installed."
}

remove() {
    echo "Removing yq..."
    sudo rm /usr/local/bin/yq
    echo "✅ yq removed."
}

check() {
    command -v yq &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
