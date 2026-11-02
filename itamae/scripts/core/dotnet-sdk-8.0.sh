#!/bin/bash
#
# METADATA
# NAME: .NET SDK 8.0
# DESCRIPTION: The Microsoft .NET 8.0 Software Development Kit.
# INSTALL_METHOD: apt
# PACKAGE_NAME: dotnet-sdk-8.0
# REPO_SETUP: setup_repo
#

# Ubuntu ships .NET in its own archive (22.04+) and Microsoft advises against mixing
# it with packages.microsoft.com, so the Microsoft repository is only added on Debian.

os_id() {
    # shellcheck disable=SC1091
    . /etc/os-release && echo "${ID:-}"
}

os_version_id() {
    # shellcheck disable=SC1091
    . /etc/os-release && echo "${VERSION_ID:-}"
}

setup_repo() {
    local id
    id=$(os_id)
    if [ "$id" = "ubuntu" ]; then
        echo "Ubuntu provides .NET packages natively; no extra repository needed."
        return 0
    fi

    local version_id
    version_id=$(os_version_id)
    echo "Setting up Microsoft repository for ${id} ${version_id}..."
    curl -fsSL "https://packages.microsoft.com/config/${id}/${version_id}/packages-microsoft-prod.deb" -o /tmp/packages-microsoft-prod.deb
    sudo dpkg -i /tmp/packages-microsoft-prod.deb > /dev/null
    rm -f /tmp/packages-microsoft-prod.deb
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
    if [ "$(os_id)" != "ubuntu" ]; then
        sudo apt-get purge -y packages-microsoft-prod 2>/dev/null || true
        sudo rm -f /etc/apt/sources.list.d/microsoft-prod.list
        sudo rm -f /etc/apt/trusted.gpg.d/microsoft.gpg
    fi
    echo "✅ .NET SDK 8.0 removed."
}

check() {
    command -v dotnet &> /dev/null && dotnet --list-sdks | grep -q "^8\.0"
}

# --- ROUTER ---
case "$1" in
    setup_repo) setup_repo ;;
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {setup_repo|install|remove|check}" && exit 1 ;;
esac
