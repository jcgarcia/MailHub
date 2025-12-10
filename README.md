# MailHub Admin

Web-based administration tool for the CMH mail server.

## Overview

MailHub Admin provides a simple web interface to manage:
- **Domains**: Add/remove mail domains
- **Users**: Create/delete email accounts, change passwords
- **Audit Log**: Track all administrative changes

## Technology Stack

- **Backend**: Go 1.21+ with Chi router
- **Frontend**: HTMX + PicoCSS (server-side rendering)
- **Database**: SQLite (audit logs)
- **Auth**: Caddy AuthCrunch SSO (X-Auth-User header)

## Development

### Prerequisites

- Go 1.21+
- Docker (for building containers)

### Local Development

```bash
# Run with dev mode (no auth required)
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

Deployed to K8s on oracledev. See PMO docs for deployment guides.

## Project Structure

```
cmd/
  mailhub-admin/        # Main application entry point
internal/
  config/               # Configuration loading
  handlers/             # HTTP handlers
  middleware/           # Auth middleware
  services/             # SSH and database services (TODO)
web/
  static/               # Static assets
  templates/            # HTML templates (TODO)
```

## License

MIT License - See LICENSE file
