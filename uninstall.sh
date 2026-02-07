#!/bin/bash
# Onyx Uninstaller - Clean & Thorough
set -e

BINARY_NAME="onyx"
INSTALL_PATH="/usr/bin/$BINARY_NAME"
CONFIG_DIR="/etc/onyx"
LOG_DIR="/var/log/onyx"
STATE_DIR="/var/lib/onyx"
SERVICE_FILE="/etc/systemd/system/onyx.service"

# Colours
info() { echo -e "\033[1;34m[INFO]\033[0m $1"; }
warn() { echo -e "\033[1;33m[WARN]\033[0m $1"; }
error() { echo -e "\033[1;31m[ERROR]\033[0m $1"; exit 1; }

# 1. Root check
if [ "$EUID" -ne 0 ]; then
    error "Please run as root (use sudo)."
fi

# 2. Stop and Disable Service
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

# 3. Remove Binary
if [ -f "$INSTALL_PATH" ]; then
    info "Removing binary from $INSTALL_PATH..."
    rm "$INSTALL_PATH"
fi

# 4. Remove System User & Home (State) Directory
if id "onyx" &>/dev/null; then
    info "Removing onyx system user and state directory ($STATE_DIR)..."
    userdel -r onyx &>/dev/null || warn "Could not fully remove user/home. May require manual cleanup."
else
    info "User 'onyx' not found, skipping..."
fi

# 5. Remove Configuration and Logs
# We ask before nuking configs just in case the user wants to keep their Caddyfile
read -p "[QUESTION] Do you want to delete all configuration and log files? (y/N): " confirm
if [[ "$confirm" == [yY] ]]; then
    info "Nuking $CONFIG_DIR and $LOG_DIR..."
    rm -rf "$CONFIG_DIR"
    rm -rf "$LOG_DIR"
else
    warn "Keeping $CONFIG_DIR and $LOG_DIR."
fi

info "Onyx has been successfully uninstalled."