# Onyx Reverse Proxy/experimental firewall/waf

Onyx is a custom build of the Caddy Web Server, tailored for Rocky Linux 10 deployments as a reverse proxy with TLS termination. The beta includes pre-bundled support for **OVH DNS** and **Coraza WAF**.

Requires firewalld

Beta version not suitable for production environments.

## Features
- **OVH DNS Integration:** Automatic HTTPS for internal/private servers using DNS challenges.
- **Coraza WAF:** Web Application Firewall integration.
- **Environment Switching:** Seamless toggle between Dev (self-signed) and Prod (public HTTPS) modes.

## Planned Features
- More Let's Encrypt methods
- WAF management interfaces
- Automated updates

## Building from Source

```bash
# 1. Clone the repo
git clone [https://github.com/DiscoMouse/onyx.git](https://github.com/DiscoMouse/onyx.git)
cd onyx

# 2. Build
go build ./cmd/onyx
