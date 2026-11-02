#!/bin/bash
#
# METADATA
# NAME: semgrep
# DESCRIPTION: A fast, open-source, static analysis tool for finding bugs.
# INSTALL_METHOD: binary
#

ensure_pipx() {
    if command -v pipx &> /dev/null; then
        return 0
    fi
    echo "Installing pipx (required)..."
    if command -v nala &> /dev/null; then
        sudo env DEBIAN_FRONTEND=noninteractive nala install -y pipx
    else
        sudo env DEBIAN_FRONTEND=noninteractive apt-get install -y pipx
    fi
}

install() {
    ensure_pipx
    echo "Installing semgrep..."
    pipx install semgrep
    echo "✅ semgrep installed."
}

remove() {
    echo "Removing semgrep..."
    pipx uninstall semgrep
    echo "✅ semgrep removed."
}

check() {
    command -v semgrep &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
