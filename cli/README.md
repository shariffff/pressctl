# pressctl CLI

A command-line interface tool for managing WordPress hosting infrastructure using Ansible. pressctl provides both an intuitive, interactive mode for manual operations and a script mode with flags for automation and CI/CD pipelines.

## Features

- **Dual Operating Modes**: Interactive prompts for manual use, script mode with flags for automation
- **YAML-based State Management**: All configuration stored in `~/.pressctl/pressctl.yaml`
- **Ansible Integration**: Seamlessly executes existing Ansible playbooks
- **Server Management**: Add, list, remove, and provision servers
- **Site Management**: Create, list, and delete WordPress sites
- **Domain Management**: Add domains and manage SSL certificates

## Installation

### Prerequisites

- Go 1.24 or higher (or Docker for building without Go)
- Ansible installed and configured
- SSH access to target servers

### Build from Source

```bash
# Clone or navigate to the repository
cd /path/to/ansible/cli

# Build the binary (requires Go 1.24+)
make build

# Or build using Docker (no Go required)
make docker-build

# Install to /usr/local/bin (requires sudo)
make install

# Or install to ~/bin (no sudo required)
make install-user
```

### Verify Installation

```bash
press version
press --help
```

## Quick Start

### 1. Initialize Configuration

```bash
press init
```

This creates `~/.pressctl/pressctl.yaml` with default settings.

### 2. Configure Ansible Path

Edit `~/.pressctl/pressctl.yaml` and set the correct Ansible project path:

```yaml
ansible:
  path: '/Users/yourname/Projects/ansible' # Update this path
  roles_path: './roles'
  inventory_path: '/tmp/pressctl-inventory-{timestamp}.ini'
  python_interpreter: '/usr/bin/python3'
```

### 3. Add a Server

```bash
press server add
```

Follow the interactive prompts to add server details:

- Server name (e.g., production-1)
- Hostname or IP address
- SSH user and port
- SSH private key file

### 4. List Servers

```bash
press server list
```

### 5. Validate Configuration

```bash
press config validate
```

## Operating Modes

pressctl CLI supports two modes of operation:

### Interactive Mode (Default)

When you run commands without flags, the CLI guides you through the process with interactive prompts.

```bash
press site create
# Prompts you for: server, domain, admin credentials
# (site ID is auto-generated from domain)
```

**Use interactive mode when:**

- Learning the tool
- Performing manual operations
- You want validation and helpful hints
- Exploring available options

### Script Mode (Non-Interactive)

Provide all parameters as command-line flags for fully automated operations.

```bash
press site create --non-interactive \
  --server production-1 \
  --domain example.com \
  --admin-user admin \
  --admin-email admin@example.com \
  --admin-password SecurePass123!
# --site-id is optional (auto-generated from domain if not provided)
```

**Use script mode when:**

- Automating deployments
- Running in CI/CD pipelines
- Scripting repetitive tasks
- No user interaction is possible

**Common flags for script mode:**

- `--non-interactive`: Required flag to enable script mode
- `--force`: Skip confirmation prompts
- `--skip-ssh-check`: Skip SSH connectivity validation

## Commands

### Configuration Management

```bash
# Show current configuration
press config show

# Validate configuration
press config validate

# Edit configuration in your preferred editor
press config edit
```

### Server Management

```bash
# Add a new server
press server add

# List all servers
press server list

# Remove a server
press server remove <name>

# Provision a server
press server provision <name>

# Provision with options
press server provision <name> --force              # Skip confirmation
press server provision <name> --skip-ssh-check     # Skip SSH connectivity test
```

### Site Management

```bash
# Create a new WordPress site (interactive)
press site create

# Create a site non-interactively (site-id auto-generated)
press site create --non-interactive \
  --server production-1 \
  --domain example.com \
  --admin-user admin \
  --admin-email admin@example.com \
  --admin-password SecurePass123!

# Create with explicit site-id
press site create --non-interactive \
  --server production-1 \
  --domain example.com \
  --site-id mysite \
  --admin-user admin \
  --admin-email admin@example.com \
  --admin-password SecurePass123!

# List all sites
press site list

# List sites on a specific server
press site list --server production-1

# Delete a site (interactive selection)
press site delete

# Delete a specific site (by site ID)
press site delete --server production-1 --site mysiteid

# Force delete without confirmation
press site delete --server production-1 --site mysiteid --force
```

### Domain Management

```bash
# Add a domain to a site (interactive)
press domain add

# Add domain with automatic SSL
# (prompts will ask if you want to issue SSL)

# Remove a domain (interactive selection)
press domain remove

# Force remove without confirmation
press domain remove --force

# Issue SSL certificate for a domain (interactive)
press domain ssl

# The CLI will:
# - Show only domains without SSL
# - Prompt for Let's Encrypt email
# - Obtain and configure SSL certificate
# - Update Nginx to use HTTPS
# - Track SSL expiration in configuration
```

## Configuration File

The configuration file is located at `~/.pressctl/pressctl.yaml`. Here's an example structure:

```yaml
version: '1.0'

ansible:
  path: '/Users/sharif/Projects/ansible'
  roles_path: './roles'
  inventory_path: '/tmp/pressctl-inventory-{timestamp}.ini'
  python_interpreter: '/usr/bin/python3'

global_vars:
  certbot_email: 'admin@example.com'
  mysql_pressctlbot_password: '${MYSQL_WORDMONBOT_PASSWORD}'
  pressctl_ssh_key: '~/.ssh/pressctl_rsa.pub'

servers:
  - name: 'production-1'
    hostname: 'prod1.example.com'
    ip: '203.0.113.10'
    ssh:
      user: 'pressctl'
      port: 22
      key_file: '~/.ssh/pressctl_rsa'
    status: 'unprovisioned'
    sites: []
```

## Development

### Build

```bash
make build
```

### Test

```bash
make test
```

### Format Code

```bash
make fmt
```

### Clean Build Artifacts

```bash
make clean
```

## Project Structure

```
cli/
├── cmd/                  # Command definitions
│   ├── root.go          # Root command
│   ├── config.go        # Config commands
│   ├── server.go        # Server commands
│   └── version.go       # Version command
├── internal/
│   ├── config/          # Configuration management
│   ├── ansible/         # Ansible integration (coming soon)
│   ├── state/           # State management (coming soon)
│   ├── prompt/          # Interactive prompts
│   └── utils/           # Utilities (coming soon)
├── pkg/
│   └── models/          # Data models
├── templates/           # Templates (inventory, etc.)
├── main.go             # Entry point
├── Makefile            # Build automation
└── README.md           # This file
```

## Roadmap

- [ ] Shell completion scripts
- [ ] Comprehensive error handling
- [ ] Installation script
- [ ] Release automation
