#!/bin/bash
#
# METADATA
# NAME: yq (Go)
# DESCRIPTION: A 'jq' for YAML. (Installs the correct Go binary, not the python wrapper).
# INSTALL_METHOD: binary
#

yq_arch() {
    case "$(uname -m)" in
        x86_64) echo amd64 ;;
        aarch64|arm64) echo arm64 ;;
        *) echo "Unsupported architecture: $(uname -m)" >&2; return 1 ;;
    esac
}

install() {
    echo "Installing yq (Go binary)..."
    local arch
    arch=$(yq_arch) || return 1
    # This is critical: apt 'yq' is the wrong tool.
    local yq_url="https://github.com/mikefarah/yq/releases/latest/download/yq_linux_${arch}"
    sudo curl --silent -fL "$yq_url" -o /usr/local/bin/yq
    sudo chmod +x /usr/local/bin/yq
    echo "✅ yq installed."
}

remove() {
    echo "Removing yq..."
    sudo rm -f /usr/local/bin/yq
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
