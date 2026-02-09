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

### 1. Installation

#### A. Server Installation (The Engine)
Run this on your **Rocky Linux 10 VPS**. This automated script installs the `onyx` engine, sets up the systemd service, and provisions the necessary users and security groups.

```bash
# Download and run the automated installer
curl -sSL https://raw.githubusercontent.com/DiscoMouse/onyx/main/install.sh | sudo bash
```