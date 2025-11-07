#!/bin/bash
#
# METADATA
# NAME: Java
# OMAKASE: true
# DESCRIPTION: The Java Development Kit (Temurin/Adoptium).
# INSTALL_METHOD: binary
#

install() {
    echo "Installing Java (Temurin)..."
    
    # Add Adoptium repository
    sudo mkdir -p /etc/apt/keyrings
    wget -O - https://packages.adoptium.net/artifactory/api/gpg/key/public | sudo tee /etc/apt/keyrings/adoptium.asc > /dev/null
    echo "deb [signed-by=/etc/apt/keyrings/adoptium.asc] https://packages.adoptium.net/artifactory/deb $(awk -F= '/^VERSION_CODENAME/{print$2}' /etc/os-release) main" | sudo tee /etc/apt/sources.list.d/adoptium.list
    
    sudo apt-get update
    if command -v nala &> /dev/null; then
        sudo nala install -y temurin-21-jdk
    else
        sudo apt-get install -y temurin-21-jdk
    fi
    echo "✅ Java (Temurin) installed."
}

remove() {
    echo "Removing Java (Temurin)..."
    sudo apt-get purge -y temurin-21-jdk
    sudo rm -f /etc/apt/sources.list.d/adoptium.list
    sudo rm -f /etc/apt/keyrings/adoptium.asc
    echo "✅ Java (Temurin) removed."
}

check() {
    command -v java &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
