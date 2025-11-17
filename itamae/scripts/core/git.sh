#!/bin/bash
#
# METADATA
# NAME: Git
# DESCRIPTION: A free and open source distributed version control system.
# INSTALL_METHOD: apt
# PACKAGE_NAME: git
# REQUIRES: GIT_USER_NAME|Enter your Git user name
# REQUIRES: GIT_USER_EMAIL|Enter your Git user email

install() {
    echo "Installing Git..."
    if command -v nala &> /dev/null; then
        sudo nala install -y git
    else
        sudo apt-get install -y git
    fi

    if [ -n "$GIT_USER_NAME" ] && [ -n "$GIT_USER_EMAIL" ]; then
        git config --global user.name "$GIT_USER_NAME"
        git config --global user.email "$GIT_USER_EMAIL"
        echo "✅ Git configured."
    else
        echo "⚠️ Git user name and email not provided. Skipping configuration."
    fi

    echo "✅ Git installed."
}

remove() {
    echo "Removing Git..."
    sudo apt-get purge -y git
    echo "✅ Git removed."
}

check() {
    command -v git &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
