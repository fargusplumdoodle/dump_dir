#!/bin/sh
set -e

# Define variables
GITHUB_REPO="fargusplumdoodle/dump_dir"
BINARY_NAME="dump_dir"
INSTALL_DIR="/usr/local/bin"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case $ARCH in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "❌ Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Fetch the latest release version
echo "🔍 Fetching the latest release of $BINARY_NAME..."
LATEST_RELEASE=$(curl -sL "https://api.github.com/repos/$GITHUB_REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_RELEASE" ]; then
    echo "❌ Failed to fetch the latest release. Please check your internet connection and try again."
    exit 1
fi

echo "🎉 Latest release: $LATEST_RELEASE"

# Construct the download URL
DOWNLOAD_URL="https://github.com/$GITHUB_REPO/releases/download/$LATEST_RELEASE/${BINARY_NAME}_${OS}_${ARCH}"

# Download the binary
echo "📥 Downloading $BINARY_NAME..."
curl -sL "$DOWNLOAD_URL" -o "$BINARY_NAME"

# Make the binary executable
chmod +x "$BINARY_NAME"

# Install the binary
echo "🚀 Installing $BINARY_NAME to $INSTALL_DIR..."
if [ -w "$INSTALL_DIR" ]; then
    mv "$BINARY_NAME" "$INSTALL_DIR"
else
    sudo mv "$BINARY_NAME" "$INSTALL_DIR"
fi

# Verify installation
if command -v "$BINARY_NAME" >/dev/null 2>&1; then
    echo "✅ $BINARY_NAME has been successfully installed to $INSTALL_DIR"
    echo "🎈 You can now use it by running: $BINARY_NAME"
    echo "🌟 Happy coding with dump_dir! 🚀📂✨"
else
    echo "❌ Installation failed. Please try again or install manually."
    exit 1
fi
