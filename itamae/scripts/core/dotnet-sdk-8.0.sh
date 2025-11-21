#!/bin/bash
#
# METADATA
# NAME: .NET SDK 8.0
# DESCRIPTION: The Microsoft .NET 8.0 Software Development Kit.
# INSTALL_METHOD: apt
# PACKAGE_NAME: dotnet-sdk-8.0
# REPO_SETUP: setup_repo
#

setup_repo() {
    echo "Setting up Microsoft repository..."
    wget -q https://packages.microsoft.com/config/ubuntu/$(lsb_release -rs)/packages-microsoft-prod.deb -O /tmp/packages-microsoft-prod.deb 2>/dev/null
    sudo dpkg -i /tmp/packages-microsoft-prod.deb > /dev/null 2>&1
    rm /tmp/packages-microsoft-prod.deb
    echo "✅ Microsoft repository configured."
}

install() {
    echo "Installing .NET SDK 8.0..."
    if command -v nala &> /dev/null; then
        sudo nala install -y dotnet-sdk-8.0
    else
        sudo apt-get install -y dotnet-sdk-8.0
    fi
    echo "✅ .NET SDK 8.0 installed."
}

remove() {
    echo "Removing .NET SDK 8.0..."
    sudo apt-get purge -y dotnet-sdk-8.0
    sudo rm -f /etc/apt/sources.list.d/microsoft-prod.list
    sudo rm -f /etc/apt/trusted.gpg.d/microsoft.gpg
    echo "✅ .NET SDK 8.0 removed."
}

check() {
    command -v dotnet &> /dev/null && dotnet --list-sdks | grep -q "8.0"
}

# --- ROUTER ---
case "$1" in
    setup_repo) setup_repo ;;
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {setup_repo|install|remove|check}" && exit 1 ;;
esac
