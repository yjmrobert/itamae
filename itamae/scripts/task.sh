#!/bin/bash
#
# METADATA
# NAME: Task
# OMAKASE: true
# DESCRIPTION: A task runner / build tool that aims to be simpler than GNU Make.
# INSTALL_METHOD: binary
#

BINDIR="$HOME/.local/bin"

install() {
    echo "Installing Task..."
    sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b "$BINDIR"
    echo "✅ Task installed."
}

remove() {
    echo "Removing Task..."
    rm -f "$BINDIR/task"
    echo "✅ Task removed."
}

check() {
    command -v task &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
