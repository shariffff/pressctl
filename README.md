# pressctl

Automated WordPress hosting on Ubuntu servers. One command to provision, one command to deploy.

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/shariffff/pressctl/main/install.sh | bash
```

## Quick Start

```bash
# 1. Add your server
press server add

# 2. Provision it (installs Nginx, PHP, MariaDB, Redis, SSL)
press server provision myserver

# 3. Create a WordPress site
press site create

# 4. Issue SSL certificate
press domain ssl
```

## What It Does

**Server provisioning:**

- Nginx from official repo
- PHP 8.3 with optimized FPM pools
- MariaDB with secure defaults
- Redis for object caching
- Let's Encrypt SSL via Certbot
- UFW firewall + Fail2ban

**Site isolation:**

- Each site runs as its own Linux user
- Dedicated PHP-FPM pool per site
- Isolated file permissions

## Commands

```bash
press server add          # Add a server
press server provision    # Provision server with LEMP stack
press server list         # List servers

press site create         # Create WordPress site
press site list           # List sites
press site delete         # Delete site

press domain add          # Add domain to site
press domain ssl          # Issue SSL certificate

press config show         # Show configuration
```

All commands support `--help` for details.

## Requirements

- Ansible 2.14+ on your local machine
- Ubuntu 24.04 target server with root SSH access

## Documentation

- [CLI Reference](cli/README.md)
- [Ansible Playbooks](ansible/README.md)

## Development

```bash
git clone https://github.com/shariffff/pressctl.git
cd pressctl
make build    # Build CLI
make test     # Run tests
```

## License

MIT
