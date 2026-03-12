#!/bin/sh
set -e

REPO="chinmaymk/acli"
BINARY_NAME="acli"
INSTALL_DIR="/usr/local/bin"

# Allow overriding the install directory
if [ -n "$ACLI_INSTALL_DIR" ]; then
    INSTALL_DIR="$ACLI_INSTALL_DIR"
fi

detect_os() {
    case "$(uname -s)" in
        Linux*)  echo "linux" ;;
        Darwin*) echo "darwin" ;;
        MINGW*|MSYS*|CYGWIN*) echo "windows" ;;
        *)
            echo "Unsupported operating system: $(uname -s)" >&2
            exit 1
            ;;
    esac
}

detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64)  echo "amd64" ;;
        aarch64|arm64) echo "arm64" ;;
        *)
            echo "Unsupported architecture: $(uname -m)" >&2
            exit 1
            ;;
    esac
}

get_latest_version() {
    if command -v curl >/dev/null 2>&1; then
        curl -sI "https://github.com/${REPO}/releases/latest" | grep -i "^location:" | sed 's#.*/tag/##' | tr -d '\r\n'
    elif command -v wget >/dev/null 2>&1; then
        wget --spider --max-redirect=0 "https://github.com/${REPO}/releases/latest" 2>&1 | grep "Location:" | sed 's#.*/tag/##' | tr -d '\r\n'
    else
        echo "Error: curl or wget is required" >&2
        exit 1
    fi
}

download() {
    url="$1"
    output="$2"
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL -o "$output" "$url"
    elif command -v wget >/dev/null 2>&1; then
        wget -qO "$output" "$url"
    fi
}

main() {
    os="$(detect_os)"
    arch="$(detect_arch)"

    echo "Detected platform: ${os}/${arch}"

    # Determine version
    if [ -n "$ACLI_VERSION" ]; then
        version="$ACLI_VERSION"
    else
        echo "Fetching latest version..."
        version="$(get_latest_version)"
        if [ -z "$version" ]; then
            echo "Error: could not determine latest version. Set ACLI_VERSION to install a specific version." >&2
            exit 1
        fi
    fi

    echo "Installing acli ${version}..."

    # Strip leading 'v' from version for archive name (goreleaser uses version without 'v' prefix)
    version_num="${version#v}"

    # Build download URL — goreleaser produces tar.gz (linux/mac) and zip (windows)
    if [ "$os" = "windows" ]; then
        filename="${BINARY_NAME}_${version_num}_${os}_${arch}.zip"
    else
        filename="${BINARY_NAME}_${version_num}_${os}_${arch}.tar.gz"
    fi

    url="https://github.com/${REPO}/releases/download/${version}/${filename}"
    echo "Downloading ${url}..."

    tmpdir="$(mktemp -d)"
    trap 'rm -rf "$tmpdir"' EXIT

    download "$url" "${tmpdir}/${filename}"

    # Extract the binary from the archive
    if [ "$os" = "windows" ]; then
        if command -v unzip >/dev/null 2>&1; then
            unzip -q "${tmpdir}/${filename}" -d "${tmpdir}"
        else
            echo "Error: unzip is required to extract the archive" >&2
            exit 1
        fi
        binary="${BINARY_NAME}.exe"
        target="${INSTALL_DIR}/${BINARY_NAME}.exe"
    else
        tar -xzf "${tmpdir}/${filename}" -C "${tmpdir}"
        binary="${BINARY_NAME}"
        target="${INSTALL_DIR}/${BINARY_NAME}"
        chmod +x "${tmpdir}/${binary}"
    fi

    # Install the binary
    if [ -w "$INSTALL_DIR" ]; then
        mv "${tmpdir}/${binary}" "$target"
    else
        echo "Installing to ${INSTALL_DIR} (requires sudo)..."
        sudo mv "${tmpdir}/${binary}" "$target"
    fi

    echo "acli ${version} installed to ${target}"
}

main
