#!/bin/bash
#
# METADATA
# NAME: semgrep
# OMAKASE: true
# DESCRIPTION: A fast, open-source, static analysis tool for finding bugs.
# INSTALL_METHOD: binary
#

install() {
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
