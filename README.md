# MailHub Admin

Web-based administration tool for the Central Mail Hub (CMH) mail server.

## Overview

MailHub Admin provides a simple web interface to manage:
- **Domains**: Add/remove mail domains hosted on CMH
- **Users**: Create/delete email accounts, change passwords
- **Audit Log**: Track all administrative changes with timestamps
- **Mail Client Config**: Instructions for setting up email clients

## Features

- Modern responsive UI with mobile support
- Real-time domain and mailbox management via SSH
- Secure authentication via Caddy AuthCrunch SSO
- Complete audit trail of all operations

## Technology Stack

- **Backend**: Go 1.24 with Chi router
- **Frontend**: HTMX (server-side rendering, no JavaScript frameworks)
- **Database**: SQLite (audit logs)
- **Auth**: Caddy AuthCrunch SSO (X-Auth-User/X-Auth-Name headers)
- **Deployment**: Kubernetes on OracleCloud

## Architecture

The application connects to the CMH mail server via SSH through a jump host:
```
K8s Pod → jump.ingasti.com:22 → localhost:2223 → CMH:22
```

Mail server configuration is managed via:
- `/etc/postfix/virtual_domains` - Domain list
- `/etc/postfix/virtual_mailbox` - Mailbox mappings
- `/etc/dovecot/users` - User authentication

## Development

### Prerequisites

- Go 1.24+
- Docker (for building containers)
- SSH keys for CMH access (for production)

### Local Development

```bash
# Run with dev mode (no auth required, mock SSH)
DEV_MODE=true DEV_AUTH_EMAIL=test@example.com go run ./cmd/mailhub-admin

# Access at http://localhost:8080
```

### Build

```bash
# Build binary
go build -o mailhub-admin ./cmd/mailhub-admin

# Build Docker image
docker build -t mailhub-admin .
```

## Deployment

Deployed via Jenkins CI/CD to K8s on oracledev:
- Registry: `ghcr.io/jcgarcia/mailhub-admin`
- Namespace: `mailhub`
- URL: https://mailhub.ingasti.com

## Project Structure

```
cmd/
  mailhub-admin/        # Main application entry point
internal/
  config/               # Configuration loading
  handlers/             # HTTP handlers (domains, users, audit)
  middleware/           # Auth middleware
  services/             # SSH client, mail operations, audit store
  templates/            # Go template functions and CSS
web/
  static/               # Static assets (logo)
k8s/
  deployment.yaml       # Kubernetes manifests
```

## License

MIT License - See LICENSE file
