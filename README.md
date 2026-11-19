# qcli

Official Go command-line interface for [QEMU Interface](https://github.com/nycon/qemu-wi) panels.

Manage VMs, snapshots, backups, tenants, networks, public IPs, and security from the terminal.

## Installation

### Prebuilt binaries (recommended)

Download the latest release for your OS from:

**https://github.com/nycon/qcli/releases**

| Platform | Archive |
|---|---|
| macOS Apple Silicon | `qcli_*_darwin_arm64.tar.gz` |
| macOS Intel | `qcli_*_darwin_amd64.tar.gz` |
| Linux x86_64 | `qcli_*_linux_amd64.tar.gz` |
| Linux ARM64 | `qcli_*_linux_arm64.tar.gz` |
| Windows x86_64 | `qcli_*_windows_amd64.zip` |

Example (macOS Apple Silicon):

```bash
curl -fsSL -o qcli.tar.gz \
  https://github.com/nycon/qcli/releases/latest/download/qcli_1.2.0_darwin_arm64.tar.gz
tar -xzf qcli.tar.gz
chmod +x qcli_1.2.0_darwin_arm64
sudo mv qcli_1.2.0_darwin_arm64 /usr/local/bin/qcli
qcli version
```

### Build from source

```bash
git clone https://github.com/nycon/qcli.git
cd qcli
make install
```

Or:

```bash
go install github.com/nycon/qcli@latest
```

Ensure `$(go env GOPATH)/bin` (and optionally `~/.local/bin`) is on your `PATH`.

## Getting Started

1. Create an API token in the panel: **User Profile → API Keys → New API Key**
2. Create a CLI context:

```bash
qcli context create my-server \
  --url https://panel.example.com \
  --token qi_your_api_key
```

For self-signed TLS certificates add `--insecure`.

3. Start managing VMs:

```bash
qcli vm list
qcli vm start my-webserver
```

## Commands

```
qcli context create <name>     Create a new server context
qcli context list              List all contexts
qcli context use <name>        Switch active context
qcli context delete <name>     Remove a context
qcli context active            Show active context

qcli vm list                   List all VMs
qcli vm create <name>          Create a new VM
qcli vm describe <id|name>     Show VM details
qcli vm start <id|name>        Start a VM
qcli vm stop <id|name>         Gracefully stop a VM
qcli vm reboot <id|name>       Reboot a VM
qcli vm poweroff <id|name>     Force power-off a VM
qcli vm delete <id|name>       Delete a VM
qcli vm resize <id|name>       Resize CPU/RAM/Disk

qcli snapshot list <vm>        List snapshots
qcli snapshot create <vm>      Create a snapshot
qcli snapshot restore <vm> <n> Restore a snapshot
qcli snapshot delete <vm> <n>  Delete a snapshot

qcli backup list <vm>          List backups
qcli backup create <vm>        Create a backup
qcli backup restore <vm> <id>  Restore from backup
qcli backup delete <vm> <id>   Delete a backup

qcli storage-pool list         List storage pools
qcli image list                List cloud images

qcli firewall list             List banned IPs
qcli firewall ban <ip>         Ban an IP
qcli firewall unban <ip>       Unban an IP
qcli firewall check <ip>       Check if IP is banned
qcli firewall jails            Show Fail2Ban jail config
qcli firewall set-jail <jail>  Update jail settings
qcli firewall log              View Fail2Ban log

qcli tenant list               List tenants
qcli tenant use <id|name>      Select active tenant
qcli tenant public-ips         List assigned public IPs
qcli tenant set-email <email>  Set notification email
qcli tenant assign-public-ips  Assign free public IPs (admin)

qcli network list              List NAT/Internal networks
qcli network create <name>     Create NAT or Internal network
qcli network update <id>       Update network settings
qcli network delete <id>       Delete a network

qcli ip-pool list              List IP pools
qcli ip-pool available         Free public IPs (admin)

qcli user list                 List users

qcli server-info               Show host CPU/RAM/Disk
qcli version                   Show CLI version
qcli completion zsh            Shell completion script
```

## Output Formats

```bash
qcli vm list
qcli -o json vm list
qcli -o yaml vm describe my-webserver
```

## Multiple Servers

```bash
qcli context create production --url https://panel.example.com --token qiqi_...
qcli context create staging --url https://staging.example.com --token qiqi_...

qcli context use production
qcli vm list
```

## Partial IDs and Names

```bash
qcli vm start 8b859dca-b3d5-4f8a-8a19-ae403426073e
qcli vm start 8b859
qcli vm start my-webserver
```

## Creating VMs

```bash
# Public IP assigned to the tenant
qcli vm create my-webserver \
  --image "Ubuntu 24.04" \
  --project <project-id> \
  --public-ip 159.69.172.28 \
  --random-password

# Or pick a reserved allocation UUID
qcli vm create my-webserver \
  --image "Ubuntu 24.04" \
  --project <project-id> \
  --ip-allocation <uuid> \
  --random-password

# NAT / Internal network via IP pool
qcli vm create my-webserver \
  --image "Ubuntu 24.04" \
  --project <project-id> \
  --ip-pool <ip-pool-id> \
  --cpu 2 --ram 2048 --disk 30 \
  --password "MySecurePassword123!" \
  --ssh-key "ssh-ed25519 AAAA... user@host"
```

### VM Create Options

| Option | Default | Description |
|---|---|---|
| `--image` | (required) | Cloud image ID or name |
| `--project` | (required) | Project UUID |
| `--ip-pool` | - | IP pool UUID (NAT/Internal or public pool) |
| `--public-ip` | - | Assigned public IP (resolves pool + allocation) |
| `--ip-allocation` | - | Reserved public IP allocation UUID |
| `--cpu` | 1 | Number of vCPUs |
| `--ram` | 1024 | RAM in MB |
| `--disk` | 20 | Disk size in GB |
| `--password` | - | Guest password |
| `--random-password` | - | Generate a 20-char random password |
| `--ssh-key` | - | SSH public key |
| `--hostname` | - | Custom hostname |
| `--user` | cloud | Default cloud-init user |
| `--storage-pool` | - | Storage pool UUID |

## Networks & public IPs

```bash
qcli network create web-dmz --type nat
qcli network create db --type internal --network 10.42.0.0/24 --gateway 10.42.0.1
qcli network list
qcli network update <id> --dhcp=false

qcli tenant public-ips
qcli ip-pool available
qcli tenant assign-public-ips Default --allocation-id <uuid> --allocation-id <uuid>
qcli tenant set-email ops@example.com
```

## Tenant isolation

API keys authenticate as a panel user. Non-superadmin users only see their tenants; VM actions are scoped to `/tenants/{id}/...` membership checks.

## Configuration

Stored in `~/.qcli/config.json` (mode `0600`). Existing `~/.qemu-cli/config.json` is migrated automatically on first run:

```json
{
  "active_context": "production",
  "contexts": {
    "production": {
      "server_url": "https://panel.example.com",
      "token": "qi_abc123...",
      "tenant_id": "...",
      "insecure": false
    }
  }
}
```

## Shell Completion

```bash
qcli completion zsh > "${fpath[1]}/_qcli"
qcli completion bash > /etc/bash_completion.d/qcli
qcli completion fish > ~/.config/fish/completions/qcli.fish
```

## Related

- Panel: https://github.com/nycon/qemu-wi

## License

MIT
