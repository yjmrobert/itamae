#!/bin/bash
#
# METADATA
# NAME: Maven
# DESCRIPTION: A build automation tool used primarily for Java projects.
# INSTALL_METHOD: binary
#

# Fallback when the Apache mirror listing cannot be read. Old releases disappear
# from dlcdn.apache.org, so the fallback is fetched from the permanent archive.
FALLBACK_VERSION="3.9.9"
INSTALL_DIR="/opt/maven"

latest_version() {
    curl -fsSL https://dlcdn.apache.org/maven/maven-3/ 2>/dev/null \
        | grep -oE 'href="3\.[0-9]+\.[0-9]+/"' \
        | grep -oE '3\.[0-9]+\.[0-9]+' \
        | sort -V | tail -n1
}

install() {
    echo "Installing Maven..."

    local version url tarball
    version=$(latest_version)
    if [ -n "$version" ]; then
        url="https://dlcdn.apache.org/maven/maven-3/${version}/binaries/apache-maven-${version}-bin.tar.gz"
    else
        version="$FALLBACK_VERSION"
        url="https://archive.apache.org/dist/maven/maven-3/${version}/binaries/apache-maven-${version}-bin.tar.gz"
    fi
    echo "Downloading Maven ${version}..."

    tarball="/tmp/apache-maven-${version}-bin.tar.gz"
    curl -fsSL "$url" -o "$tarball"
    sudo mkdir -p "$INSTALL_DIR"
    sudo tar -xzf "$tarball" -C "$INSTALL_DIR" --strip-components=1
    rm -f "$tarball"

    sudo ln -sf "$INSTALL_DIR/bin/mvn" /usr/local/bin/mvn

    echo "✅ Maven ${version} installed."
}

remove() {
    echo "Removing Maven..."
    sudo rm -rf "$INSTALL_DIR"
    sudo rm -f /usr/local/bin/mvn
    echo "✅ Maven removed."
}

check() {
    command -v mvn &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
