#!/bin/bash
#
# METADATA
# NAME: Cascadia Code
# DESCRIPTION: A monospaced font from Microsoft that includes programming ligatures.
# INSTALL_METHOD: binary
#

FONT_VERSION="2404.23"
FONT_DIR="$HOME/.local/share/fonts"

install() {
    echo "Installing Cascadia Code..."
    
    mkdir -p "$FONT_DIR"
    wget https://github.com/microsoft/cascadia-code/releases/download/v${FONT_VERSION}/CascadiaCode-${FONT_VERSION}.zip -P /tmp
    unzip -o /tmp/CascadiaCode-${FONT_VERSION}.zip -d /tmp/cascadia-code
    cp /tmp/cascadia-code/ttf/*.ttf "$FONT_DIR/"
    rm -rf /tmp/CascadiaCode-${FONT_VERSION}.zip /tmp/cascadia-code
    
    # Refresh font cache
    fc-cache -f -v
    
    echo "✅ Cascadia Code installed."
}

remove() {
    echo "Removing Cascadia Code..."
    rm -f "$FONT_DIR"/CascadiaCode*.ttf
    rm -f "$FONT_DIR"/CascadiaMono*.ttf
    fc-cache -f -v
    echo "✅ Cascadia Code removed."
}

check() {
    fc-list | grep -q "Cascadia"
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
