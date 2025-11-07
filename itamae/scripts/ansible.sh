#!/bin/bash
#
# METADATA
# NAME: Ansible
# OMAKASE: true
# DESCRIPTION: An open-source automation tool, including ansible-core and ansible-runner.
# INSTALL_METHOD: binary
#

install() {
    echo "Installing Ansible..."
    pipx install --include-deps ansible
    pipx inject ansible ansible-runner
    ansible-galaxy collection install ansible.posix
    ansible-galaxy collection install ansible.windows
    ansible-galaxy collection install community.windows
    echo "✅ Ansible installed."
}

remove() {
    echo "Removing Ansible..."
    pipx uninstall ansible
    echo "✅ Ansible removed."
}

check() {
    command -v ansible &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
