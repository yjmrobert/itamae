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

# apt_has_candidate succeeds when the package can actually be installed from
# the configured repositories (apt-cache show also succeeds for packages that
# are known but have no installation candidate).
apt_has_candidate() {
    apt-cache policy "$1" 2>/dev/null | grep -q "Candidate: [0-9]"
}

apt_install() {
    if command -v nala &> /dev/null; then
        sudo env DEBIAN_FRONTEND=noninteractive nala install -y "$@"
    else
        sudo env DEBIAN_FRONTEND=noninteractive apt-get install -y "$@"
    fi
}

install() {
    echo "Installing Alacritty..."
    if apt_has_candidate alacritty; then
        apt_install alacritty || return 1
    elif [ "$(os_id)" = "ubuntu" ]; then
        echo "alacritty is not in this Ubuntu release; adding ppa:aslatter/ppa..."
        if ! command -v add-apt-repository &> /dev/null; then
            apt_install software-properties-common || return 1
        fi
        sudo env DEBIAN_FRONTEND=noninteractive add-apt-repository -y ppa:aslatter/ppa || return 1
        sudo env DEBIAN_FRONTEND=noninteractive apt-get update
        apt_install alacritty || return 1
    else
        echo "❌ alacritty is not available from the package repositories on this system." >&2
        return 1
    fi
    echo "✅ Alacritty installed."
}

remove() {
    echo "Removing Alacritty..."
    sudo env DEBIAN_FRONTEND=noninteractive apt-get purge -y alacritty
    if ls /etc/apt/sources.list.d/aslatter-ubuntu-ppa-* &> /dev/null; then
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
