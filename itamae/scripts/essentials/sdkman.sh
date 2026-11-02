#!/bin/bash
#
# METADATA
# NAME: SDKMan
# DESCRIPTION: A tool for managing parallel versions of multiple Software Development Kits.
# INSTALL_METHOD: binary
#

ensure_deps() {
    # The SDKMAN installer refuses to run without zip and unzip.
    local missing=()
    command -v zip &> /dev/null || missing+=(zip)
    command -v unzip &> /dev/null || missing+=(unzip)
    if [ ${#missing[@]} -gt 0 ]; then
        echo "Installing dependencies: ${missing[*]}"
        if command -v nala &> /dev/null; then
            sudo nala install -y "${missing[@]}"
        else
            sudo apt-get install -y "${missing[@]}"
        fi
    fi
}

install() {
    echo "Installing SDKMan..."
    ensure_deps
    curl -s "https://get.sdkman.io?rcupdate=true" | bash
    echo "✅ SDKMan installed."
}

remove() {
    echo "Removing SDKMan..."
    rm -rf "$HOME/.sdkman"
    echo "✅ SDKMan removed."
}

check() {
    [[ -s "$HOME/.sdkman/bin/sdkman-init.sh" ]]
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
