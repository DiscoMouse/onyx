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
# 1. Download and run the installer script.
curl -sSL https://raw.githubusercontent.com/DiscoMouse/onyx/main/install.sh | sudo bash
# 2. Follow the instructions

```
## Building from Source

```bash
# 1. Clone the repo
git clone https://github.com/DiscoMouse/onyx.git
cd onyx

# 2. Build
go build ./cmd/onyx
```
## Configuration

To improve the clarity of the post-installation process for Onyx, the documentation should explicitly detail where configuration files reside, how the security model affects editing, and how to transition from development to production modes.

### Recommended "Post-Install Configuration" Guide

The following sections can be added to your `README.md` to clarify the setup process:

#### 1. Configuration File Locations

After installation, Onyx uses the following directory structure for configuration and state:

* **Main Configuration:** `/etc/onyx/Caddyfile` — Defines your site blocks, reverse proxies, and TLS settings.
* **Environment Variables:** `/etc/onyx/onyx.env` — Stores sensitive keys (like OVH API credentials).
* **WAF Rules:** `/var/lib/onyx/rules/` — The directory for Coraza/OWASP rule sets.
* **Logs:** `/var/log/onyx/` — Access and system logs.

#### 2. Editing the Caddyfile

The `/etc/onyx/` directory is owned by `root:onyx` with `750` permissions. To modify your configuration:

* **Use Sudo:** You must use `sudo` to edit files in this directory (e.g., `sudo nano /etc/onyx/Caddyfile`).
* **Site Configuration:** Update the site block in the `Caddyfile` to match your domain and backend service:
```caddy
yourdomain.com {
    import tls_{$ONYX_ENV:dev}
    reverse_proxy 127.0.0.1:8080 
}

```



#### 3. Configuring Production Mode (OVH DNS)

To enable public HTTPS via OVH DNS challenges:

1. **Add API Keys:** Create and edit `/etc/onyx/onyx.env` to include your OVH credentials:
```bash
OVH_ENDPOINT="ovh-eu"
OVH_APPLICATION_KEY="your_key"
OVH_APPLICATION_SECRET="your_secret"
OVH_CONSUMER_KEY="your_consumer_key"

```


2. **Toggle Environment:** Edit the service file at `/etc/systemd/system/onyx.service` and change `ONYX_ENV=dev` to `ONYX_ENV=prod`.
3. **Apply Changes:**
```bash
sudo systemctl daemon-reload
sudo systemctl restart onyx

```

#### 4. Validating and Managing the Proxy

Since Onyx is a custom Caddy build, you can use the `proxy` subcommand for direct engine management:

* **Check Syntax:** `onyx proxy validate --config /etc/onyx/Caddyfile`
* **View Dashboard:** Run `onyx status` to see the real-time TUI heartbeat.
* *Note:* Your user must be a member of the `onyx-admin` group to execute these commands.



### Key Concepts

* **The Admin Group:** The `onyx-admin` is required for humans to run the `onyx` binary and view and modify options in the dashboard, while the `onyx` group is used for the onyx service/process reading configuration data.
* **Ambience of Environment Variables:** The `ONYX_ENV` determines which TLS block (`tls_dev` or `tls_prod`) is imported by the `Caddyfile`. This prevents hammering the Let's Encrypt servers during testing.
* **Security Restrictions:** The `onyx` system user is intentionally blocked from administrative tools like `status` to maintain functional separation.