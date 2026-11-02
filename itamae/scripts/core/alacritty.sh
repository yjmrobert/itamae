#!/bin/bash
#
# METADATA
# NAME: Alacritty
# DESCRIPTION: A fast, cross-platform, OpenGL terminal emulator.
# INSTALL_METHOD: binary
#

# Alacritty is packaged in Debian 12+ and Ubuntu 24.04+. Older Ubuntu releases
# (22.04) get it from the maintainer's PPA.

os_id() {
    # shellcheck disable=SC1091
    . /etc/os-release && echo "${ID:-}"
}

apt_install() {
    if command -v nala &> /dev/null; then
        sudo nala install -y "$@"
    else
        sudo apt-get install -y "$@"
    fi
}

install() {
    echo "Installing Alacritty..."
    if apt-cache show alacritty &> /dev/null; then
        apt_install alacritty
    elif [ "$(os_id)" = "ubuntu" ]; then
        echo "alacritty is not in this Ubuntu release; adding ppa:aslatter/ppa..."
        if ! command -v add-apt-repository &> /dev/null; then
            apt_install software-properties-common
        fi
        sudo add-apt-repository -y ppa:aslatter/ppa
        sudo apt-get update
        apt_install alacritty
    else
        echo "❌ alacritty is not available from the package repositories on this system." >&2
        return 1
    fi
    echo "✅ Alacritty installed."
}

remove() {
    echo "Removing Alacritty..."
    sudo apt-get purge -y alacritty
    if [ -f /etc/apt/sources.list.d/aslatter-ubuntu-ppa-"$(. /etc/os-release && echo "${VERSION_CODENAME:-}")".sources ] || ls /etc/apt/sources.list.d/aslatter-ubuntu-ppa-* &> /dev/null; then
        sudo rm -f /etc/apt/sources.list.d/aslatter-ubuntu-ppa-*
    fi
    echo "✅ Alacritty removed."
}

check() {
    command -v alacritty &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
