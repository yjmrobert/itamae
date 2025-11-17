#!/bin/bash
#
# METADATA
# NAME: bat
# DESCRIPTION: A 'cat' clone with syntax highlighting and Git integration.
# INSTALL_METHOD: apt
# PACKAGE_NAME: bat
# POST_INSTALL: post_install
#

post_install() {
    echo "✅ bat installed successfully"
}

install() {
    echo "Installing bat..."
    if command -v nala &> /dev/null; then
        sudo nala install -y bat
    else
        sudo apt-get install -y bat
    fi
    post_install
}

remove() {
    echo "Removing bat..."
    sudo apt-get purge -y bat
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
