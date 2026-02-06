#!/bin/bash
# Onyx Installer - Minimalist & Functional
set -e

REPO="DiscoMouse/onyx"
BINARY_NAME="onyx"
INSTALL_PATH="/usr/bin/$BINARY_NAME"
CONFIG_DIR="/etc/onyx"
LOG_DIR="/var/log/onyx"

# Colours
info() { echo -e "\033[1;34m[INFO]\033[0m $1"; }
warn() { echo -e "\033[1;33m[WARN]\033[0m $1"; }
error() { echo -e "\033[1;31m[ERROR]\033[0m $1"; exit 1; }

# 1. Root check
if [ "$EUID" -ne 0 ]; then
    error "Please run as root (use sudo)."
fi

# 2. Grab the latest binary from GitHub
info "Fetching latest Onyx binary..."
curl -L -o "$BINARY_NAME" "https://github.com/$REPO/releases/latest/download/$BINARY_NAME"
chmod +x "$BINARY_NAME"

# 3. Setup System User
if ! id "onyx" &>/dev/null; then
    info "Creating onyx system user..."
    useradd --system --create-home --home-dir /var/lib/onyx --shell /usr/sbin/nologin onyx
else
    info "User 'onyx' already exists, skipping..."
fi

# 4. Directory permissions
info "Preparing directories..."
mkdir -p "$CONFIG_DIR" "$LOG_DIR"
chown -R root:onyx "$CONFIG_DIR"
chown onyx:onyx "$LOG_DIR"

# 5. Move binary and set capabilities
info "Installing binary to $INSTALL_PATH..."
mv "$BINARY_NAME" "$INSTALL_PATH"
setcap cap_net_bind_service=+ep "$INSTALL_PATH"

# 6. Service setup
SERVICE_URL="https://raw.githubusercontent.com/$REPO/main/onyx.service"

if [ ! -f "onyx.service" ]; then
    info "Service file not found locally. Downloading from GitHub..."
    curl -sSL -o /etc/systemd/system/onyx.service "$SERVICE_URL"
else
    info "Using local onyx.service..."
    cp onyx.service /etc/systemd/system/
fi

systemctl daemon-reload
info "Installation complete! Run 'onyx version' to verify."
info "Service registered. Run 'systemctl enable --now onyx' to start."
