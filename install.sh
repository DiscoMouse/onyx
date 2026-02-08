#!/bin/bash
# Onyx Installer - v0.1.5 Secure Functional Separation
set -e

REPO="DiscoMouse/onyx"
BINARY_NAME="onyx"
INSTALL_PATH="/usr/bin/$BINARY_NAME"
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

# 2. Grab the latest binary from GitHub
info "Fetching latest Onyx binary..."
curl -L -o "$BINARY_NAME" "https://github.com/$REPO/releases/latest/download/$BINARY_NAME"
chmod +x "$BINARY_NAME"

# 3. Setup System Groups and Users
if ! id "onyx" &>/dev/null; then
    info "Creating onyx system user..."
    useradd --system --create-home --home-dir /var/lib/onyx --shell /usr/sbin/nologin onyx
else
    info "User 'onyx' already exists, skipping..."
fi

if ! getent group "$ADMIN_GROUP" &>/dev/null; then
    info "Creating $ADMIN_GROUP group..."
    groupadd "$ADMIN_GROUP"
fi

# EXECUTION ACCESS: Add system user to admin group so the service can run the binary
usermod -aG "$ADMIN_GROUP" onyx

# ADMIN ACCESS: Human gets both execution (admin) and data (onyx) access
if [ -n "$SUDO_USER" ]; then
    info "Provisioning human admin '$SUDO_USER'..."
    usermod -aG "$ADMIN_GROUP" "$SUDO_USER"
    usermod -aG "onyx" "$SUDO_USER"
fi

# 4. Directory permissions
info "Preparing directories..."
mkdir -p "$CONFIG_DIR" "$LOG_DIR" "/var/lib/onyx/rules"

if [ ! -f "$CONFIG_DIR/Caddyfile" ]; then
    info "Installing Caddyfile template..."
    curl -sSL -o "$CONFIG_DIR/Caddyfile" "https://raw.githubusercontent.com/$REPO/main/exampleCaddyfile"
fi

chown -R root:onyx "$CONFIG_DIR"
chmod 750 "$CONFIG_DIR"
chown -R onyx:onyx "/var/lib/onyx"
chown onyx:onyx "$LOG_DIR"

# 5. Binary lockdown
info "Installing binary to $INSTALL_PATH..."
mv "$BINARY_NAME" "$INSTALL_PATH"

chown root:"$ADMIN_GROUP" "$INSTALL_PATH"
chmod 750 "$INSTALL_PATH"
setcap cap_net_bind_service=+ep "$INSTALL_PATH"

# 6. Service setup
SERVICE_URL="https://raw.githubusercontent.com/$REPO/main/onyx.service"
curl -sSL -o /etc/systemd/system/onyx.service "$SERVICE_URL"

systemctl daemon-reload

# 7. Finalize
info "Installation complete!"
if [ -n "$SUDO_USER" ]; then
    warn "USER ACTION REQUIRED: Run 'exec su -l $SUDO_USER' to refresh your group permissions."
fi