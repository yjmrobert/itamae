#!/bin/bash
#
# METADATA
# NAME: bin
# DESCRIPTION: A tool for easier management of binary tools.
# INSTALL_METHOD: binary
#

BINDIR="$HOME/.local/bin"

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

bin_arch() {
    case "$(uname -m)" in
        x86_64) echo amd64 ;;
        aarch64|arm64) echo arm64 ;;
        *) echo "Unsupported architecture: $(uname -m)" >&2; return 1 ;;
    esac
}

install() {
    echo "Installing bin..."
    local arch version
    arch=$(bin_arch) || return 1

    # Release assets are named bin_<version>_linux_<arch>.
    version=$(latest_release_tag marcosnils/bin)
    version="${version#v}"
    if [ -z "$version" ]; then
        echo "❌ Could not determine the latest bin version." >&2
        return 1
    fi

    mkdir -p "$BINDIR"
    curl -fsSL "https://github.com/marcosnils/bin/releases/download/v${version}/bin_${version}_linux_${arch}" -o "$BINDIR/bin"
    chmod +x "$BINDIR/bin"
    echo "✅ bin ${version} installed to $BINDIR."
}

remove() {
    echo "Removing bin..."
    rm -f "$BINDIR/bin"
    echo "✅ bin removed."
}

check() {
    command -v bin &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
