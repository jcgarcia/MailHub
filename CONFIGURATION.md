# MailHub Admin Configuration Guide

MailHub Admin is configured entirely through **environment variables**. No configuration files are needed—all settings are read from the environment at startup.

## Required Configuration

Before running MailHub Admin, you must set these environment variables with values specific to your environment:

### Mail Server Connection (SSH)

| Environment Variable | Default Value | Description | Example |
|---|---|---|---|
| `CMH_SSH_HOST` | `localhost` | Direct SSH host to your mail server | `192.168.1.10` or `mail.yourdomain.com` |
| `CMH_SSH_PORT` | `22` | SSH port on the mail server | `22` or `2222` |
| `CMH_SSH_USER` | `postman` | SSH username for mail server access | `postman` or `mailadmin` |
| `CMH_SSH_KEY_PATH` | `/secrets/mailhub_key` | Path to SSH private key file | `/home/user/.ssh/id_rsa` or `/var/secrets/mailhub_key` |
| `CMH_SSH_JUMP_HOST` | `jump.yourdomain.com` | **[IMPORTANT]** SSH jump host (bastion) | `jump.yourdomain.com` or `bastion.yourdomain.com` |
| `CMH_SSH_JUMP_USER` | `ubuntu` | SSH username on the jump host | `ubuntu` or `ec2-user` |
| `CMH_SSH_JUMP_KEY_PATH` | `/secrets/jump_key` | Path to jump host SSH private key | `/var/secrets/jump_key` |

### Application Settings

| Environment Variable | Default Value | Description |
|---|---|---|
| `PORT` | `8080` | Port where MailHub Admin listens |
| `DATABASE_PATH` | `/data/mailhub.db` | SQLite database file location (for audit logs) |
| `LOG_LEVEL` | `info` | Logging level: `debug`, `info`, `warn`, `error` |
| `DEV_MODE` | `false` | Enable development mode (disables auth, uses mock SSH) |
| `DEV_AUTH_EMAIL` | `dev@example.com` | Dev mode authentication email |

## Configuration Examples

### Example 1: Local Development

```bash
# Dev mode with mock SSH (no real connection needed)
export DEV_MODE=true
export DEV_AUTH_EMAIL=admin@example.com
export PORT=8080
export DATABASE_PATH=/tmp/mailhub.db

go run ./cmd/mailhub-admin
# Access at http://localhost:8080
```

### Example 2: Docker Container (On-Premises)

Create an `.env` file:

```bash
# Mail server SSH access (direct, no jump host)
CMH_SSH_HOST=192.168.1.100
CMH_SSH_PORT=22
CMH_SSH_USER=postman
CMH_SSH_KEY_PATH=/secrets/mailhub_key

# Application settings
PORT=8080
DATABASE_PATH=/data/mailhub.db
LOG_LEVEL=info
```

Run the Docker container:

```bash
docker run -d \
  --name mailhub \
  --env-file .env \
  -v /path/to/ssh/key:/secrets/mailhub_key:ro \
  -v mailhub-data:/data \
  -p 8080:8080 \
  ghcr.io/jcgarcia/mailhub-admin:latest
```

### Example 3: Kubernetes with Jump Host

Create a Kubernetes Secret for SSH keys:

```bash
kubectl create secret generic mailhub-ssh-keys \
  --from-file=mailhub_key=/path/to/mailhub_key \
  --from-file=jump_key=/path/to/jump_key \
  -n mailhub
```

Update the Kubernetes deployment with environment variables:

```yaml
env:
- name: CMH_SSH_HOST
  value: "192.168.1.100"  # Internal mail server IP
- name: CMH_SSH_PORT
  value: "22"
- name: CMH_SSH_USER
  value: "postman"
- name: CMH_SSH_KEY_PATH
  value: "/secrets/mailhub_key"

# Jump host configuration (critical for remote access)
- name: CMH_SSH_JUMP_HOST
  value: "jump.yourdomain.com"
- name: CMH_SSH_JUMP_USER
  value: "ubuntu"
- name: CMH_SSH_JUMP_KEY_PATH
  value: "/secrets/jump_key"

- name: PORT
  value: "8080"
- name: DATABASE_PATH
  value: "/data/mailhub.db"
- name: LOG_LEVEL
  value: "info"
```

### Example 4: Cloud Deployment (AWS/Azure/OracleCloud)

```bash
# Environment file for cloud deployment
CMH_SSH_HOST=internal-mail-server.aws.internal
CMH_SSH_PORT=22
CMH_SSH_USER=postman
CMH_SSH_KEY_PATH=/var/secrets/mailhub_key

# Jump host is critical for cloud deployments where mail server is in private subnet
CMH_SSH_JUMP_HOST=bastion.yourdomain.com
CMH_SSH_JUMP_USER=ec2-user
CMH_SSH_JUMP_KEY_PATH=/var/secrets/jump_key

PORT=8080
DATABASE_PATH=/data/mailhub.db
LOG_LEVEL=info
```

## Understanding SSH Architecture

MailHub Admin connects to your mail server through an SSH tunnel. Two connection models are supported:

### Direct SSH (Single Server)

```
MailHub Admin → SSH → Mail Server
                (22)
```

**Use when:**
- Mail server is directly reachable from MailHub Admin
- Both are on the same network
- Both are accessible from the same public IP

**Configuration:**
```bash
CMH_SSH_HOST=192.168.1.100      # Direct IP or hostname
CMH_SSH_PORT=22                  # Standard SSH port
CMH_SSH_USER=postman
CMH_SSH_KEY_PATH=/secrets/mailhub_key
```

### SSH with Jump Host (Secure - Recommended)

```
MailHub Admin → SSH → Jump Host (Bastion) → SSH → Mail Server
                (22)                        (22 or custom)
```

**Use when:**
- Mail server is in a private subnet
- You want to restrict direct internet access to mail server
- You have a bastion/jump host for secure access (AWS, Azure, etc.)
- Multi-cloud or hybrid cloud setup

**Configuration:**
```bash
CMH_SSH_HOST=192.168.1.100                    # Private IP of mail server
CMH_SSH_PORT=22
CMH_SSH_USER=postman
CMH_SSH_KEY_PATH=/secrets/mailhub_key

CMH_SSH_JUMP_HOST=jump.yourdomain.com         # Public bastion host
CMH_SSH_JUMP_USER=ubuntu
CMH_SSH_JUMP_KEY_PATH=/secrets/jump_key
```

## Verifying Configuration

### Test SSH Connection

Before starting MailHub Admin, verify your SSH configuration works:

```bash
# Test direct connection
ssh -i /path/to/mailhub_key -p 22 postman@192.168.1.100

# Test jump host connection (if using bastion)
ssh -i /path/to/jump_key -J ubuntu@jump.yourdomain.com \
    -i /path/to/mailhub_key -p 22 postman@192.168.1.100
```

### Check MailHub Logs

When MailHub Admin starts, check the logs for connection status:

```bash
# Docker
docker logs mailhub

# Kubernetes
kubectl logs -f deployment/mailhub-admin -n mailhub

# Local development
go run ./cmd/mailhub-admin 2>&1 | grep -i "ssh\|connection"
```

Look for these log messages:

✅ **Success:** `SSH connection to mail server established`

❌ **Error:** `WARNING: SSH connection test failed`

## Troubleshooting

### "Connection refused" error

- Verify `CMH_SSH_HOST` and `CMH_SSH_PORT` are correct
- Ensure firewall allows SSH access
- Check SSH key permissions: `chmod 600 /path/to/key`

### "Permission denied" error

- Verify `CMH_SSH_USER` exists on the mail server
- Ensure the SSH key is authorized on mail server (`~/.ssh/authorized_keys`)
- Check SSH key is readable by MailHub process

### "Could not resolve hostname" error

- Verify `CMH_SSH_JUMP_HOST` or `CMH_SSH_HOST` DNS name is correct
- Check network DNS resolution: `nslookup jump.yourdomain.com`

### Jump host connection fails

- Test jump host manually: `ssh -i /path/to/jump_key ubuntu@jump.yourdomain.com`
- Verify jump host can reach mail server: `ssh from-jump-host-to-mail-server`
- Ensure both SSH keys are in correct locations and readable

## Security Best Practices

1. **Restrict SSH Key Permissions**
   ```bash
   chmod 600 /path/to/mailhub_key
   chmod 600 /path/to/jump_key
   ```

2. **Use Secret Management**
   - Kubernetes: Store keys in Secrets, not in Deployments
   - Docker: Use secrets volumes, not environment variables for sensitive data
   - Cloud: Use cloud provider secret management (AWS Secrets Manager, Azure Key Vault, etc.)

3. **Limit SSH User Privileges**
   - Create a dedicated SSH user with minimal permissions
   - Only allow SSH key authentication (disable password auth)
   - Restrict user to specific mail server management commands

4. **Use Jump Host for Security**
   - Don't expose mail server directly to internet
   - Use bastion/jump host for centralized access control
   - Enable audit logging on jump host

5. **Encrypt SSH Keys in Transit**
   - Use TLS for MailHub Admin web interface
   - Mount SSH keys from secure locations (Secret volumes, encrypted storage)

## Environment Variable Quick Reference

```bash
# Copy and customize this template
export PORT=8080
export DATABASE_PATH=/data/mailhub.db
export LOG_LEVEL=info
export DEV_MODE=false

export CMH_SSH_HOST=192.168.1.100
export CMH_SSH_PORT=22
export CMH_SSH_USER=postman
export CMH_SSH_KEY_PATH=/secrets/mailhub_key

export CMH_SSH_JUMP_HOST=jump.yourdomain.com
export CMH_SSH_JUMP_USER=ubuntu
export CMH_SSH_JUMP_KEY_PATH=/secrets/jump_key
```

## Support

For issues or questions about configuration, please check:
1. Verify all SSH settings are correct (see "Verifying Configuration")
2. Test SSH connection manually (see "Troubleshooting")
3. Check application logs for detailed error messages
4. Review the example configurations above for your deployment model
