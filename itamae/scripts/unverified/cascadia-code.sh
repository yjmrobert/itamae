#!/bin/bash
#
# METADATA
# NAME: Cascadia Code
# DESCRIPTION: A monospaced font from Microsoft that includes programming ligatures.
# INSTALL_METHOD: binary
#

FONT_VERSION="2404.23"
FONT_DIR="$HOME/.local/share/fonts"

ensure_deps() {
    local missing=()
    command -v unzip &> /dev/null || missing+=(unzip)
    command -v fc-cache &> /dev/null || missing+=(fontconfig)
    if [ ${#missing[@]} -gt 0 ]; then
        echo "Installing dependencies: ${missing[*]}"
        if command -v nala &> /dev/null; then
            sudo nala install -y "${missing[@]}"
        else
            sudo apt-get install -y "${missing[@]}"
        fi
    fi
}

install() {
    echo "Installing Cascadia Code..."
    ensure_deps

    mkdir -p "$FONT_DIR"
    curl -fsSL "https://github.com/microsoft/cascadia-code/releases/download/v${FONT_VERSION}/CascadiaCode-${FONT_VERSION}.zip" -o "/tmp/CascadiaCode-${FONT_VERSION}.zip"
    unzip -oq "/tmp/CascadiaCode-${FONT_VERSION}.zip" -d /tmp/cascadia-code
    cp /tmp/cascadia-code/ttf/*.ttf "$FONT_DIR/"
    rm -rf "/tmp/CascadiaCode-${FONT_VERSION}.zip" /tmp/cascadia-code

    # Refresh font cache
    fc-cache -f

    echo "✅ Cascadia Code installed."
}

remove() {
    echo "Removing Cascadia Code..."
    rm -f "$FONT_DIR"/CascadiaCode*.ttf
    rm -f "$FONT_DIR"/CascadiaMono*.ttf
    fc-cache -f
    echo "✅ Cascadia Code removed."
}

check() {
    fc-list 2>/dev/null | grep -q "Cascadia"
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
