#!/bin/bash
#
# METADATA
# NAME: bat
# DESCRIPTION: A 'cat' clone with syntax highlighting and Git integration.
# INSTALL_METHOD: apt
# PACKAGE_NAME: bat
# POST_INSTALL: post_install
#

BINDIR="$HOME/.local/bin"

post_install() {
    # Debian and Ubuntu install the binary as 'batcat' to avoid a name clash;
    # provide the 'bat' name everyone expects.
    if ! command -v bat &> /dev/null && command -v batcat &> /dev/null; then
        mkdir -p "$BINDIR"
        ln -sf "$(command -v batcat)" "$BINDIR/bat"
        echo "✅ Created symlink: bat -> batcat"
    fi
    echo "✅ bat installed successfully"
}

install() {
    echo "Installing bat..."
    if command -v nala &> /dev/null; then
        sudo env DEBIAN_FRONTEND=noninteractive nala install -y bat
    else
        sudo env DEBIAN_FRONTEND=noninteractive apt-get install -y bat
    fi
    post_install
}

remove() {
    echo "Removing bat..."
    sudo env DEBIAN_FRONTEND=noninteractive apt-get purge -y bat
    rm -f "$BINDIR/bat"
    echo "✅ bat removed."
}

check() {
    command -v bat &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    post_install) post_install ;;
    *) echo "Usage: $0 {install|remove|check|post_install}" && exit 1 ;;
esac
