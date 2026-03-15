#!/bin/sh
set -e

# Writebook CLI installer
# Usage: curl -fsSL https://raw.githubusercontent.com/gezibash/writebook/main/install.sh | sh

REPO="gezibash/writebook"
BINARY="writebook"
INSTALL_DIR="/usr/local/bin"

# Colors (if terminal supports it)
if [ -t 1 ]; then
    BOLD="\033[1m"
    GREEN="\033[32m"
    RED="\033[31m"
    RESET="\033[0m"
else
    BOLD=""
    GREEN=""
    RED=""
    RESET=""
fi

info() { printf "${BOLD}${GREEN}==>${RESET} ${BOLD}%s${RESET}\n" "$1"; }
error() { printf "${BOLD}${RED}error:${RESET} %s\n" "$1" >&2; exit 1; }

# Detect OS
detect_os() {
    case "$(uname -s)" in
        Linux*)  echo "linux" ;;
        Darwin*) echo "darwin" ;;
        CYGWIN*|MINGW*|MSYS*) echo "windows" ;;
        *) error "Unsupported operating system: $(uname -s)" ;;
    esac
}

# Detect architecture
detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64)  echo "amd64" ;;
        arm64|aarch64) echo "arm64" ;;
        *) error "Unsupported architecture: $(uname -m)" ;;
    esac
}

# Detect latest version from GitHub
detect_version() {
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/'
    elif command -v wget >/dev/null 2>&1; then
        wget -qO- "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/'
    else
        error "curl or wget is required"
    fi
}

# Download file
download() {
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL -o "$2" "$1"
    elif command -v wget >/dev/null 2>&1; then
        wget -qO "$2" "$1"
    fi
}

main() {
    OS=$(detect_os)
    ARCH=$(detect_arch)
    VERSION=${WRITEBOOK_VERSION:-$(detect_version)}

    if [ -z "$VERSION" ]; then
        error "Could not determine latest version. Set WRITEBOOK_VERSION manually."
    fi

    info "Installing ${BINARY} ${VERSION} (${OS}/${ARCH})"

    # Construct download URL
    if [ "$OS" = "windows" ]; then
        ARCHIVE="${BINARY}_${OS}_${ARCH}.zip"
    else
        ARCHIVE="${BINARY}_${OS}_${ARCH}.tar.gz"
    fi
    URL="https://github.com/${REPO}/releases/download/${VERSION}/${ARCHIVE}"

    # Download to temp directory
    TMP_DIR=$(mktemp -d)
    trap 'rm -rf "$TMP_DIR"' EXIT

    info "Downloading ${URL}"
    download "$URL" "${TMP_DIR}/${ARCHIVE}"

    # Extract
    info "Extracting"
    if [ "$OS" = "windows" ]; then
        unzip -q "${TMP_DIR}/${ARCHIVE}" -d "$TMP_DIR"
    else
        tar -xzf "${TMP_DIR}/${ARCHIVE}" -C "$TMP_DIR"
    fi

    # Verify checksum if sha256sum is available
    CHECKSUMS_URL="https://github.com/${REPO}/releases/download/${VERSION}/checksums.txt"
    if command -v sha256sum >/dev/null 2>&1; then
        download "$CHECKSUMS_URL" "${TMP_DIR}/checksums.txt"
        EXPECTED=$(grep "${ARCHIVE}" "${TMP_DIR}/checksums.txt" | awk '{print $1}')
        ACTUAL=$(sha256sum "${TMP_DIR}/${ARCHIVE}" | awk '{print $1}')
        if [ "$EXPECTED" != "$ACTUAL" ]; then
            error "Checksum verification failed"
        fi
        info "Checksum verified"
    elif command -v shasum >/dev/null 2>&1; then
        download "$CHECKSUMS_URL" "${TMP_DIR}/checksums.txt"
        EXPECTED=$(grep "${ARCHIVE}" "${TMP_DIR}/checksums.txt" | awk '{print $1}')
        ACTUAL=$(shasum -a 256 "${TMP_DIR}/${ARCHIVE}" | awk '{print $1}')
        if [ "$EXPECTED" != "$ACTUAL" ]; then
            error "Checksum verification failed"
        fi
        info "Checksum verified"
    fi

    # Install
    INSTALL_PATH="${INSTALL_DIR}/${BINARY}"
    if [ -w "$INSTALL_DIR" ]; then
        mv "${TMP_DIR}/${BINARY}" "$INSTALL_PATH"
        chmod +x "$INSTALL_PATH"
    else
        info "Elevated permissions required to install to ${INSTALL_DIR}"
        sudo mv "${TMP_DIR}/${BINARY}" "$INSTALL_PATH"
        sudo chmod +x "$INSTALL_PATH"
    fi

    info "Installed ${BINARY} to ${INSTALL_PATH}"
    printf "\n"
    "${INSTALL_PATH}" --version
    printf "\n"
    info "Run 'writebook login' to get started"
}

main
