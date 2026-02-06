# Onyx Reverse Proxy/experimental firewall/waf

Onyx is a custom build of the Caddy Web Server, tailored for Rocky Linux 10 deployments as a reverse proxy with TLS termination. The beta includes pre-bundled support for **OVH DNS** and **Coraza WAF**.

Requires firewalld

Beta version not suitable for production environments.

## License & Credits

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

**Open Source Components:**
This binary includes software developed by:
* **Caddy Web Server** (Apache 2.0) - Copyright The Caddy Authors
* **Coraza WAF** (Apache 2.0) - Copyright The Coraza Authors
* **OVH DNS Module** (Apache 2.0)

## Features
- **OVH DNS Integration:** Automatic HTTPS for internal/private servers using DNS challenges.
- **Coraza WAF:** Web Application Firewall integration.
- **Environment Switching:** Seamless toggle between Dev (self-signed) and Prod (public HTTPS) modes.

## Planned Features
- More Let's Encrypt methods
- WAF management interfaces
- Automated updates

## Easy Install

```bash
# 1. Download the latest binary directly
curl -L -o onyx "https://github.com/DiscoMouse/onyx/releases/latest/download/onyx"

# 2. Set permissions and move to path
chmod +x onyx
sudo mv onyx /usr/bin/onyx

# 3. Update capabilities
sudo setcap cap_net_bind_service=+ep /usr/bin/onyx
```
## Building from Source

```bash
# 1. Clone the repo
git clone https://github.com/DiscoMouse/onyx.git
cd onyx

# 2. Build
go build ./cmd/onyx
```