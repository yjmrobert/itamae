#!/bin/bash
# install.sh
# Installs the prerequisites and runs the itamae tool.

# Exit on error
set -e

# Update package lists and install prerequisites
echo "Updating package lists..."
sudo apt-get update -y
echo "Installing prerequisites (git, nala)..."
sudo apt-get install -y git nala

# Check if itamae is already installed
if command -v itamae &> /dev/null; then
    echo "itamae is already installed."
    exit 0
fi

# Check if Go is installed, and install it if it's not
if ! command -v go &> /dev/null; then
    echo "Go is not found, installing..."
    # Determine the latest Go version
    echo "Checking for the latest Go version..."
    GO_VERSION=$(curl -s "https://go.dev/VERSION?m=text" | head -n 1 | sed 's/go//')
    echo "Latest Go version is ${GO_VERSION}"

    # Determine OS and architecture
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    case $ARCH in
      x86_64)
        ARCH="amd64"
        ;;
      aarch64)
        ARCH="arm64"
        ;;
      *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
    esac

    # Download and install Go
    case $OS in
      linux)
        GO_TARBALL="go${GO_VERSION}.linux-${ARCH}.tar.gz"
        DOWNLOAD_URL="https://go.dev/dl/${GO_TARBALL}"
        echo "Downloading Go from ${DOWNLOAD_URL}..."
        curl -L -o "/tmp/${GO_TARBALL}" "${DOWNLOAD_URL}"
        echo "Extracting Go..."
        sudo tar -C /usr/local -xzf "/tmp/${GO_TARBALL}"
        rm "/tmp/${GO_TARBALL}"
        ;;
      darwin)
        echo "macOS is not yet supported by this bootstrap script."
        exit 1
        ;;
      *)
        echo "Unsupported OS: $OS"
        exit 1
        ;;
    esac

    echo "Go installation complete."
else
    echo "Go is already installed."
fi

# Set PATH for the current session
export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin

# Add Go to the user's shell profile for future sessions
echo "Configuring shell environment..."
SHELL_CONFIG=""
if [[ "$SHELL" == *"zsh"* ]]; then
  SHELL_CONFIG="$HOME/.zshrc"
elif [[ "$SHELL" == *"bash"* ]]; then
  SHELL_CONFIG="$HOME/.bashrc"
else
  echo "Unsupported shell: $SHELL. Please add /usr/local/go/bin and \$HOME/go/bin to your PATH manually."
fi

if [ -n "$SHELL_CONFIG" ] && [ -f "$SHELL_CONFIG" ]; then
  if ! grep -q "/usr/local/go/bin" "$SHELL_CONFIG"; then
    echo -e "\n# Go and Itamae PATH" >> "$SHELL_CONFIG"
    echo 'export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin' >> "$SHELL_CONFIG"
    echo "Updated PATH in $SHELL_CONFIG. Sourcing the file now to apply changes."
    source "$SHELL_CONFIG"
  else
    echo "Go path is already in your shell config."
  fi
else
  if [ -z "$SHELL_CONFIG" ]; then
    # already printed the "unsupported shell" message
    :
  else
    echo "Could not find shell configuration file: $SHELL_CONFIG"
    echo "Please add /usr/local/go/bin and \$HOME/go/bin to your PATH manually."
  fi
fi

# Install and run itamae
echo "Installing itamae..."
go install github.com/yjmrobert/itamae@latest

echo "Running itamae install..."
itamae install

echo "itamae installation is complete."
