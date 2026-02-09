# Onyx Reverse Proxy / Security Engine

Onyx is a security-hardened reverse proxy and administrative control plane tailored for Rocky Linux 10. Built on the Caddy engine, it provides seamless mTLS-protected management, an interactive dashboard, and integrated Coraza WAF.

> [!WARNING]
> **Beta version:** Currently in active development and not suitable for production environments.

## New in v0.1.7: Distributed Architecture
Onyx now operates as a dual-binary system. The **Engine** runs on your server, while the **Admin Console** runs locally on your workstation. Communication is secured via a strict mTLS control plane.

## Features
- **mTLS Pairing:** Securely link admin consoles to servers using one-time tokens and Ed25519 certificates.
- **Interactive TUI:** Real-time health monitoring via the `onyx-admin status` dashboard.
- **OVH DNS Integration:** Automatic HTTPS for internal/private servers using DNS challenges.
- **Coraza WAF:** Web Application Firewall integration with pre-bundled OWASP Core Rule Sets.
- **Environment Toggling:** Switch between `dev` (self-signed) and `prod` (Let's Encrypt) modes instantly.

---

## Getting Started

### 1. Server Installation (The Engine)
Run this on your **Rocky Linux 10 VPS**. This automated script installs the `onyx` engine, sets up the systemd service, and provisions the necessary users and security groups.

```bash
# Download and run the automated installer
curl -sSL https://raw.githubusercontent.com/DiscoMouse/onyx/main/install.sh | sudo bash
```

### 2. Client Installation (The Admin Tool)
Run this on your local Linux machine (Laptop/Workstation/VM). This installs only the onyx-admin CLI tool used to pair with and manage your remote servers.

```bash
# 1. Download the admin binary
curl -L -o onyx-admin "[https://github.com/DiscoMouse/onyx/releases/latest/download/onyx-admin](https://github.com/DiscoMouse/onyx/releases/latest/download/onyx-admin)"
chmod +x onyx-admin

# 2. Install to your path
sudo mv onyx-admin /usr/local/bin/

# 3. Initialize local config directory
mkdir -p ~/.config/onyx/certs
```

### 3. Usage
Step 1: Initialize the Server
SSH into your VPS and start the pairing window. This will generate a one-time token valid for 5 minutes.

```bash
# On the VPS
sudo onyx --pair
```

Take note of the token (e.g., ABCD-1234) and ensure Port 2305 is open/acessable so you can connect to it. (VPN/Wireguard/SSH Tunnel highly reccomended for security to protect the service)

Step 2: Pair the Client
On your local machine, use the token from the previous step to establish the secure trust relationship.

```bash
# On your local machine
onyx-admin pair <VPS_IP_ADDRESS> --token <TOKEN>
```
Step 3: Launch Dashboard
Once paired, you can monitor your remote node in real-time.

```bash
# On your local machine
onyx-admin status
```

Security Architecture
Onyx enforces Security by Isolation.

- The Engine (onyx): Runs as a restricted system user. It cannot access administrative keys or execute administrative logic.
- The Admin Tool (onyx-admin): Restricted to the onyx-admin group. Only authorized human users can initiate pairing or view system status.
- Control Plane: All remote communication is encrypted via mTLS using Ed25519 keys generated locally on your devices.

License
This project is licensed under the Apache 2.0 License - see the LICENSE file for details.

