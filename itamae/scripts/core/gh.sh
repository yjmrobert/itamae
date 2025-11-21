#!/bin/bash
#
# METADATA
# NAME: GitHub CLI
# DESCRIPTION: The official GitHub command-line tool.
# INSTALL_METHOD: apt
# PACKAGE_NAME: gh
# REPO_SETUP: setup_repo
#

setup_repo() {
    echo "Setting up GitHub CLI repository..."
    sudo mkdir -p /etc/apt/keyrings
    wget -qO- https://cli.github.com/packages/githubcli-archive-keyring.gpg 2>/dev/null | sudo tee /etc/apt/keyrings/githubcli-archive-keyring.gpg > /dev/null
    sudo chmod go+r /etc/apt/keyrings/githubcli-archive-keyring.gpg
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null
    echo "✅ GitHub CLI repository configured."
}

install() {
    echo "Installing GitHub CLI..."
    if command -v nala &> /dev/null; then
        sudo nala install -y gh
    else
        sudo apt-get install -y gh
    fi
    echo "✅ GitHub CLI installed."
}

remove() {
    echo "Removing GitHub CLI..."
    sudo apt-get purge -y gh
    sudo rm -f /etc/apt/sources.list.d/github-cli.list
    sudo rm -f /etc/apt/keyrings/githubcli-archive-keyring.gpg
    echo "✅ GitHub CLI removed."
}

check() {
    command -v gh &> /dev/null
}

# --- ROUTER ---
case "$1" in
    setup_repo) setup_repo ;;
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {setup_repo|install|remove|check}" && exit 1 ;;
esac
