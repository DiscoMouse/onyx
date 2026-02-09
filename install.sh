#!/bin/bash
# Onyx Installer - v0.1.6 Secure Functional Separation
# Strictly implements the Permissions & Ownership Manifest
set -e

REPO="DiscoMouse/onyx"
BINARY_NAME="onyx"
ADMIN_BINARY_NAME="onyx-admin"
INSTALL_PATH="/usr/bin/$BINARY_NAME"
ADMIN_INSTALL_PATH="/usr/bin/$ADMIN_BINARY_NAME"
CONFIG_DIR="/etc/onyx"
LOG_DIR="/var/log/onyx"
ADMIN_GROUP="onyx-admin"

# Colours
info() { echo -e "\033[1;34m[INFO]\033[0m $1"; }
warn() { echo -e "\033[1;33m[WARN]\033[0m $1"; }
error() { echo -e "\033[1;31m[ERROR]\033[0m $1"; exit 1; }

# 1. Root check
if [ "$EUID" -ne 0 ]; then
    error "Please run as root (use sudo)."
fi

# 2. Create Groups First
info "Configuring security groups..."
# Create the Admin group (for humans)
if ! getent group "$ADMIN_GROUP" &>/dev/null; then
    groupadd "$ADMIN_GROUP"
fi
# Create the Engine group (system)
if ! getent group "onyx" &>/dev/null; then
    groupadd "onyx"
fi

# 3. Create System User
if ! id "onyx" &>/dev/null; then
    info "Creating onyx system user..."
    # Note: We do NOT add onyx to the onyx-admin group.
    useradd --system --create-home --home-dir /var/lib/onyx --shell /usr/sbin/nologin -g onyx onyx
else
    info "User 'onyx' already exists, skipping..."
fi

# 4. Fetch Binaries
info "Fetching latest Onyx binaries..."
curl -L -o "$BINARY_NAME" "https://github.com/$REPO/releases/latest/download/$BINARY_NAME"
curl -L -o "$ADMIN_BINARY_NAME" "https://github.com/$REPO/releases/latest/download/$ADMIN_BINARY_NAME"

# 5. Install & Apply Strict Permissions (The Manifest Logic)
info "Installing binaries with Functional Separation..."

# A. The Proxy Engine (onyx) -> root:onyx
mv "$BINARY_NAME" "$INSTALL_PATH"
chown root:onyx "$INSTALL_PATH"
chmod 750 "$INSTALL_PATH"
# Allow binding to ports 80/443 without root
setcap cap_net_bind_service=+ep "$INSTALL_PATH"

# B. The Admin Console (onyx-admin) -> root:onyx-admin
mv "$ADMIN_BINARY_NAME" "$ADMIN_INSTALL_PATH"
chown root:"$ADMIN_GROUP" "$ADMIN_INSTALL_PATH"
chmod 750 "$ADMIN_INSTALL_PATH"

# 6. Configure Directories
info "Setting up directories..."
mkdir -p "$CONFIG_DIR" "$LOG_DIR" "/var/lib/onyx/rules"

# Config: root owns it, onyx group can read it
chown -R root:onyx "$CONFIG_DIR"
chmod 750 "$CONFIG_DIR"

if [ ! -f "$CONFIG_DIR/Caddyfile" ]; then
    curl -sSL -o "$CONFIG_DIR/Caddyfile" "https://raw.githubusercontent.com/$REPO/main/exampleCaddyfile"
fi

# Logs/State: onyx owns it
chown -R onyx:onyx "/var/lib/onyx"
chown onyx:onyx "$LOG_DIR"
chmod 755 "$LOG_DIR"

# 7. Provision Human Admin
if [ -n "$SUDO_USER" ]; then
    info "Provisioning human admin '$SUDO_USER'..."
    # Human gets access to admin tools
    usermod -aG "$ADMIN_GROUP" "$SUDO_USER"
    # Optional: Human gets read access to logs via onyx group
    usermod -aG "onyx" "$SUDO_USER"
fi

# 8. Service Setup
SERVICE_URL="https://raw.githubusercontent.com/$REPO/main/onyx.service"
curl -sSL -o /etc/systemd/system/onyx.service "$SERVICE_URL"
systemctl daemon-reload

info "Installation complete! (v0.1.6)"
if [ -n "$SUDO_USER" ]; then
    warn "USER ACTION REQUIRED: Run 'exec su -l $SUDO_USER' to refresh your group permissions."
fi