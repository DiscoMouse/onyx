# Onyx Reverse Proxy / Security Engine

Onyx is a security-hardened reverse proxy and administrative control plane tailored for Rocky Linux 10. Built on the Caddy engine, it provides seamless mTLS-protected management, an interactive dashboard, and integrated Coraza WAF.

> [!WARNING]
> **Beta version:** Currently in active development and not suitable for production environments.

## New in v0.1.6: Administrative Control Plane
Onyx now features a zero-trust administrative layer. To manage a remote Onyx engine, you must perform a secure **Pairing Handshake** which exchanges Ed25519-backed certificates.

## Features
- **mTLS Pairing:** Securely link admin consoles to servers using one-time tokens and Ed25519 certificates.
- **Interactive TUI:** Real-time health monitoring via the `onyx-admin status` dashboard.
- **OVH DNS Integration:** Automatic HTTPS for internal/private servers using DNS challenges.
- **Coraza WAF:** Web Application Firewall integration with pre-bundled OWASP Core Rule Sets.
- **Environment Toggling:** Switch between `dev` (self-signed) and `prod` (Let's Encrypt) modes instantly.

---

## Getting Started

### 1. Installation
```bash
# Download and run the automated installer
curl -sSL [https://raw.githubusercontent.com/DiscoMouse/onyx/main/install.sh](https://raw.githubusercontent.com/DiscoMouse/onyx/main/install.sh) | sudo bash