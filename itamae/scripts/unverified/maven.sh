#!/bin/bash
#
# METADATA
# NAME: Maven
# DESCRIPTION: A build automation tool used primarily for Java projects.
# INSTALL_METHOD: binary
#

MAVEN_VERSION="3.9.9"
INSTALL_DIR="/opt/maven"

install() {
    echo "Installing Maven..."
    
    # Download and extract Maven
    wget https://dlcdn.apache.org/maven/maven-3/${MAVEN_VERSION}/binaries/apache-maven-${MAVEN_VERSION}-bin.tar.gz -P /tmp
    sudo mkdir -p "$INSTALL_DIR"
    sudo tar -xzf /tmp/apache-maven-${MAVEN_VERSION}-bin.tar.gz -C "$INSTALL_DIR" --strip-components=1
    rm /tmp/apache-maven-${MAVEN_VERSION}-bin.tar.gz
    
    # Create symlink
    sudo ln -sf "$INSTALL_DIR/bin/mvn" /usr/local/bin/mvn
    
    echo "✅ Maven installed."
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
