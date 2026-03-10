#!/bin/sh
set -e

REPO="crisecheverria/svenska"
INSTALL_DIR="/usr/local/bin"

# Detect OS
OS="$(uname -s)"
case "$OS" in
  Darwin)  OS="darwin" ;;
  Linux)   OS="linux" ;;
  MINGW*|MSYS*|CYGWIN*) OS="windows" ;;
  *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

# Detect architecture
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64|amd64)  ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Get latest version
VERSION="$(curl -sI "https://github.com/$REPO/releases/latest" | grep -i "^location:" | sed 's|.*/v||' | tr -d '\r\n')"
if [ -z "$VERSION" ]; then
  echo "Error: could not determine latest version"
  exit 1
fi

echo "Installing svenska v${VERSION} (${OS}/${ARCH})..."

# Download
EXT="tar.gz"
if [ "$OS" = "windows" ]; then
  EXT="zip"
fi
URL="https://github.com/$REPO/releases/download/v${VERSION}/svenska_${OS}_${ARCH}.${EXT}"

TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

curl -sL "$URL" -o "$TMP/svenska.${EXT}"

# Extract
cd "$TMP"
if [ "$EXT" = "zip" ]; then
  unzip -q "svenska.${EXT}"
else
  tar xzf "svenska.${EXT}"
fi

# Install
if [ -w "$INSTALL_DIR" ]; then
  mv svenska "$INSTALL_DIR/svenska"
else
  echo "Need sudo to install to $INSTALL_DIR"
  sudo mv svenska "$INSTALL_DIR/svenska"
fi

chmod +x "$INSTALL_DIR/svenska"

echo "Done! Run 'svenska' to start learning Swedish."
