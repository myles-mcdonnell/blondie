#!/bin/sh

# This install script is intended to download and install the latest available
# release of blondie.
#
# You can install using this script:
# $ curl https://raw.githubusercontent.com/myles-mcdonnell/blondie/master/install.sh | sh
#
# This script was shamelessly copied from https://raw.githubusercontent.com/golang/dep/master/install.sh

set -e

RELEASES_URL="https://github.com/myles-mcdonnell/blondie/releases"

downloadJSON() {
    url="$2"

    echo "Fetching $url.."
    if test -x "$(command -v curl)"; then
        response=$(curl -s -L -w 'HTTPSTATUS:%{http_code}' -H 'Accept: application/json' "$url")
        body=$(echo "$response" | sed -e 's/HTTPSTATUS\:.*//g')
        code=$(echo "$response" | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
    elif test -x "$(command -v wget)"; then
        temp=$(mktemp)
        body=$(wget -q --header='Accept: application/json' -O - --server-response "$url" 2> "$temp")
        code=$(awk '/^  HTTP/{print $2}' < "$temp" | tail -1)
        rm "$temp"
    else
        echo "Neither curl nor wget was available to perform http requests."
        exit 1
    fi
    if [ "$code" != 200 ]; then
        echo "Request failed with code $code"
        exit 1
    fi

    eval "$1='$body'"
}

downloadFile() {
    url="$1"
    destination="$2"

    echo "Fetching $url.."
    if test -x "$(command -v curl)"; then
        code=$(curl -s -w '%{http_code}' -L "$url" -o "$destination")
    elif test -x "$(command -v wget)"; then
        code=$(wget -q -O "$destination" --server-response "$url" 2>&1 | awk '/^  HTTP/{print $2}' | tail -1)
    else
        echo "Neither curl nor wget was available to perform http requests."
        exit 1
    fi

    if [ "$code" != 200 ]; then
        echo "Request failed with code $code"
        exit 1
    fi
}

initArch() {
    ARCH=$(uname -m)
    case $ARCH in
        amd64) ARCH="amd64";;
        x86_64) ARCH="amd64";;
        i386) ARCH="386";;
        *) echo "Architecture ${ARCH} is not supported by this installation script"; exit 1;;
    esac
    echo "ARCH = $ARCH"
}

initOS() {
    OS=$(uname | tr '[:upper:]' '[:lower:]')
    case "$OS" in
        darwin) OS='darwin';;
        linux) OS='linux';;
        freebsd) OS='freebsd';;
        mingw*) OS='windows';;
        msys*) OS='windows';;
        *) echo "OS ${OS} is not supported by this installation script"; exit 1;;
    esac
    echo "OS = $OS"
}

# identify platform based on uname output
initArch
initOS

# determine install directory if required
INSTALL_DIRECTORY="/usr/local/bin/"
echo "Will install into $INSTALL_DIRECTORY"

# assemble expected release artifact name
BINARY="blondie_${OS}_${ARCH}"

# add .exe if on windows
if [ "$OS" = "windows" ]; then
    BINARY="$BINARY.exe"
fi

# if BLONDIE_INSTALL_VERSION was not provided, assume latest
if [ -z "$BLONDIE_INSTALL_VERSION" ]; then
    downloadJSON LATEST_RELEASE "$RELEASES_URL/latest"
    BLONDIE_INSTALL_VERSION=$(echo "${LATEST_RELEASE}" | tr -s '\n' ' ' | sed 's/.*"tag_name":"//' | sed 's/".*//' )
fi
echo "Release Tag = $BLONDIE_INSTALL_VERSION"

BINARY_URL="$RELEASES_URL/download/$BLONDIE_INSTALL_VERSION/$BINARY"
DOWNLOAD_FILE=$(mktemp)

downloadFile "$BINARY_URL" "$DOWNLOAD_FILE"

echo "Setting executable permissions."
chmod +x "$DOWNLOAD_FILE"

INSTALL_NAME="blondie"

if [ "$OS" = "windows" ]; then
    INSTALL_NAME="$INSTALL_NAME.exe"
fi

echo "Moving executable to $INSTALL_DIRECTORY/$INSTALL_NAME"
mv "$DOWNLOAD_FILE" "$INSTALL_DIRECTORY/$INSTALL_NAME"