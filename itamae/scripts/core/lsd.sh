#!/bin/bash
#
# METADATA
# NAME: lsd
# DESCRIPTION: The next gen 'ls' command, with icons and colours.
# INSTALL_METHOD: binary
#

# lsd is packaged in Debian 12+ and Ubuntu 23.04+. On older releases the
# official .deb from GitHub is used instead.

# latest_release_tag prints the tag of a GitHub repository's latest release.
# It follows the /releases/latest redirect instead of calling the REST API,
# which is rate limited for unauthenticated clients (shared CI runners hit
# that limit constantly); the API is only used as a fallback.
latest_release_tag() {
    local tag
    tag=$(curl -fsSL -o /dev/null -w '%{url_effective}' "https://github.com/$1/releases/latest" 2>/dev/null | sed -n 's#.*/tag/##p')
    if [ -z "$tag" ]; then
        tag=$(curl -fsSL "https://api.github.com/repos/$1/releases/latest" 2>/dev/null | grep -oP '"tag_name": *"\K[^"]+')
    fi
    echo "$tag"
}

# apt_has_candidate succeeds when the package can actually be installed from
# the configured repositories (apt-cache show also succeeds for packages that
# are known but have no installation candidate).
apt_has_candidate() {
    apt-cache policy "$1" 2>/dev/null | grep -q "Candidate: [0-9]"
}

deb_arch() {
    dpkg --print-architecture 2>/dev/null || echo amd64
}

install() {
    echo "Installing lsd..."
    if apt_has_candidate lsd; then
        if command -v nala &> /dev/null; then
            sudo env DEBIAN_FRONTEND=noninteractive nala install -y lsd
        else
            sudo env DEBIAN_FRONTEND=noninteractive apt-get install -y lsd
        fi
    else
        echo "lsd is not in the distribution repositories, using the GitHub release..."
        local version arch deb
        version=$(latest_release_tag lsd-rs/lsd)
        version="${version#v}"
        if [ -z "$version" ]; then
            echo "❌ Could not determine the latest lsd version." >&2
            return 1
        fi
        arch=$(deb_arch)
        deb="/tmp/lsd_${version}_${arch}.deb"
        curl -fsSL "https://github.com/lsd-rs/lsd/releases/download/v${version}/lsd_${version}_${arch}.deb" -o "$deb"
        sudo env DEBIAN_FRONTEND=noninteractive apt-get install -y "$deb"
        rm -f "$deb"
    fi
    echo "✅ lsd installed."
}

remove() {
    echo "Removing lsd..."
    sudo env DEBIAN_FRONTEND=noninteractive apt-get purge -y lsd
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
