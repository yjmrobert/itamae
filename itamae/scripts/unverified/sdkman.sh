#!/bin/bash
#
# METADATA
# NAME: SDKMan
# DESCRIPTION: A tool for managing parallel versions of multiple Software Development Kits.
# INSTALL_METHOD: binary
#

install() {
    echo "Installing SDKMan..."
    curl -s "https://get.sdkman.io" | bash
    source "$HOME/.sdkman/bin/sdkman-init.sh"
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
