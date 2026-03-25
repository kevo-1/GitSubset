#!/usr/bin/env bash
set -e

REPO="kevo-1/GitSubset"
echo "Fetching latest release from $REPO..."

# Detect OS and Arch
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

if [ "$ARCH" = "x86_64" ]; then
    ARCH="amd64"
elif [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
    ARCH="arm64"
fi

echo "Detected OS: $OS, Architecture: $ARCH"

# Extract the browser_download_url from the GitHub API response
URLS=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep -Eo '"browser_download_url": "[^"]+"' | cut -d'"' -f4)

DOWNLOAD_URL=""
for url in $URLS; do
    lower_url=$(echo "$url" | tr '[:upper:]' '[:lower:]')
    if [[ "$lower_url" == *"$OS"* ]] && [[ "$lower_url" == *"$ARCH"* ]]; then
        DOWNLOAD_URL=$url
        break
    fi
done

if [ -z "$DOWNLOAD_URL" ]; then
    echo "Error: Could not find a release binary for $OS ($ARCH)."
    exit 1
fi

echo "Downloading $DOWNLOAD_URL..."

INSTALL_DIR="/usr/local/bin"
TMP_DIR=$(mktemp -d)

FILENAME=$(basename "$DOWNLOAD_URL")
curl -sL "$DOWNLOAD_URL" -o "$TMP_DIR/$FILENAME"

if [[ "$FILENAME" == *.tar.gz ]]; then
    tar -xzf "$TMP_DIR/$FILENAME" -C "$TMP_DIR"
    EXE_PATH=$(find "$TMP_DIR" -type f -name "gitsubset" | head -n 1)
elif [[ "$FILENAME" == *.zip ]]; then
    unzip -q "$TMP_DIR/$FILENAME" -d "$TMP_DIR"
    EXE_PATH=$(find "$TMP_DIR" -type f -name "gitsubset" | head -n 1)
else
    EXE_PATH="$TMP_DIR/$FILENAME"
fi

if [ -z "$EXE_PATH" ] || [ ! -f "$EXE_PATH" ]; then
    echo "Error: Target executable 'gitsubset' not found in downloaded file."
    rm -rf "$TMP_DIR"
    exit 1
fi

# Ensure it is executable
chmod +x "$EXE_PATH"

echo "Installing gitsubset to $INSTALL_DIR..."
if [ -w "$INSTALL_DIR" ]; then
    mv "$EXE_PATH" "$INSTALL_DIR/gitsubset"
else
    echo "Requesting sudo privileges to install to $INSTALL_DIR (you may be prompted for your password)..."
    sudo mv "$EXE_PATH" "$INSTALL_DIR/gitsubset"
fi

rm -rf "$TMP_DIR"
echo "Installation complete! Run 'gitsubset' to get started."
