#!/bin/bash
# Onyx Uninstaller - v0.1.6
set -e

BINARY_NAME="onyx"
ADMIN_BINARY_NAME="onyx-admin"
INSTALL_PATH="/usr/bin/$BINARY_NAME"
ADMIN_INSTALL_PATH="/usr/bin/$ADMIN_BINARY_NAME"
CONFIG_DIR="/etc/onyx"
LOG_DIR="/var/log/onyx"
STATE_DIR="/var/lib/onyx"
SERVICE_FILE="/etc/systemd/system/onyx.service"
ADMIN_GROUP="onyx-admin"

# Colours
info() { echo -e "\033[1;34m[INFO]\033[0m $1"; }
warn() { echo -e "\033[1;33m[WARN]\033[0m $1"; }
error() { echo -e "\033[1;31m[ERROR]\033[0m $1"; exit 1; }

# 1. Root check
if [ "$EUID" -ne 0 ]; then
    error "Please run as root (use sudo)."
fi

# 2. Stop Service
if systemctl is-active --quiet "$BINARY_NAME"; then
    info "Stopping Onyx service..."
    systemctl stop "$BINARY_NAME"
fi

if [ -f "$SERVICE_FILE" ]; then
    info "Removing systemd service..."
    systemctl disable "$BINARY_NAME" &>/dev/null
    rm "$SERVICE_FILE"
    systemctl daemon-reload
    systemctl reset-failed
fi

# 3. Remove Binaries
if [ -f "$INSTALL_PATH" ]; then
    rm "$INSTALL_PATH"
    info "Removed engine binary."
fi
if [ -f "$ADMIN_INSTALL_PATH" ]; then
    rm "$ADMIN_INSTALL_PATH"
    info "Removed admin binary."
fi

# 4. Remove User and Groups
if id "onyx" &>/dev/null; then
    info "Removing onyx system user..."
    userdel -r onyx &>/dev/null || warn "Manual cleanup needed for user onyx."
fi

if getent group "$ADMIN_GROUP" &>/dev/null; then
    info "Removing $ADMIN_GROUP group..."
    groupdel "$ADMIN_GROUP"
fi

# 5. Remove Data (Nuclear Option)
echo -n "[QUESTION] Do you want to delete all configuration and log files? (y/N): "
read -r confirm < /dev/tty || confirm="n"

case "$confirm" in
    [yY][eE][sS]|[yY])
        info "Removing config and logs..."
        rm -rf "$CONFIG_DIR" "$LOG_DIR"
        ;;
    *)
        warn "Keeping $CONFIG_DIR and $LOG_DIR."
        ;;
esac

info "Onyx uninstalled."