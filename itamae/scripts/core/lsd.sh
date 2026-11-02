#!/bin/bash
#
# METADATA
# NAME: lsd
# DESCRIPTION: The next gen 'ls' command, with icons and colours.
# INSTALL_METHOD: binary
#

# lsd is packaged in Debian 12+ and Ubuntu 23.04+. On older releases the
# official .deb from GitHub is used instead.

deb_arch() {
    dpkg --print-architecture 2>/dev/null || echo amd64
}

install() {
    echo "Installing lsd..."
    if apt-cache show lsd &> /dev/null; then
        if command -v nala &> /dev/null; then
            sudo nala install -y lsd
        else
            sudo apt-get install -y lsd
        fi
    else
        echo "lsd is not in the distribution repositories, using the GitHub release..."
        local version
        version=$(curl -fsSL https://api.github.com/repos/lsd-rs/lsd/releases/latest | grep -oP '"tag_name": *"v?\K[^"]+')
        if [ -z "$version" ]; then
            echo "❌ Could not determine the latest lsd version." >&2
            return 1
        fi
        local arch deb
        arch=$(deb_arch)
        deb="/tmp/lsd_${version}_${arch}.deb"
        curl -fsSL "https://github.com/lsd-rs/lsd/releases/download/v${version}/lsd_${version}_${arch}.deb" -o "$deb"
        sudo apt-get install -y "$deb"
        rm -f "$deb"
    fi
    echo "✅ lsd installed."
}

remove() {
    echo "Removing lsd..."
    sudo apt-get purge -y lsd
    echo "✅ lsd removed."
}

check() {
    command -v lsd &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
