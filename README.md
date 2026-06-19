<div align="center">

<img src="assets/banner.svg" alt="Claude Sync" width="100%">

<br>

*Encrypted with [age](https://github.com/FiloSottile/age) • R2 / S3 / GCS / Azure / WebDAV supported*

[![Release](https://img.shields.io/github/v/release/tawanorg/claude-sync)](https://github.com/eschweighofer/claude-sync/releases)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![npm](https://img.shields.io/npm/v/@tawandotorg/claude-sync)](https://www.npmjs.com/package/@tawandotorg/claude-sync)
[![Socket Badge](https://badge.socket.dev/npm/package/@tawandotorg/claude-sync/1.11.1)](https://badge.socket.dev/npm/package/@tawandotorg/claude-sync/1.11.1)

[Quick Start](#quick-start) • [Setup Guide](#setup-guide) • [Commands](#commands) • [Shell Integration](#shell-integration) • [Security](#security)

</div>

---

## Features

- **Cross-device sync**: Continue Claude Code conversations on any laptop
- **Multi-provider storage**: Cloudflare R2, AWS S3, Google Cloud Storage, Azure Blob Storage, S3-compatible (Backblaze B2, MinIO, Wasabi), or WebDAV (Nextcloud, ownCloud)
- **End-to-end encryption**: All files encrypted with age before upload
- **Passphrase-based keys**: Same passphrase = same key on any device (no file copying)
- **Selective sync**: Choose `--scope sessions` to sync only conversation data (skip plugins/node_modules)
- **Interactive wizard**: Arrow-key driven setup with validation
- **Secure self-updating**: `claude-sync update` downloads and verifies SHA256 checksums
- **Simple CLI**: `push`, `pull`, `status`, `diff`, `conflicts` commands
- **Compression**: Gzip compression before encryption for faster syncs
- **Shell integration**: Optional shell hooks for automatic push/pull

<div align="center">
<img src="assets/claude-sync.gif" alt="Claude Sync Demo" width="100%">
</div>

## Quick Start

### First Device

```bash
# Install
npm install -g @tawandotorg/claude-sync

# Set up (interactive wizard)
claude-sync init

# Push your sessions
claude-sync push
```

### Second Device

```bash
# Install
npm install -g @tawandotorg/claude-sync

# Set up with SAME storage credentials
claude-sync init
# Select same provider (R2/S3/GCS/WebDAV)
# Enter same bucket name and credentials
# Choose "Passphrase" for encryption
# Enter the SAME passphrase as first device
# ✓ Encryption key verified  <-- confirms passphrase matches!

# Preview what would be synced
claude-sync pull --dry-run

# Pull sessions (creates backup if you have existing files)
claude-sync pull
```

**Same passphrase = same encryption key.** The init verifies your passphrase can decrypt remote files before completing.

## Setup Guide

### Step 1: Choose a Storage Provider

| Provider | Free Tier | Best For |
|----------|-----------|----------|
| **Cloudflare R2** | 10GB storage | Personal use (recommended) |
| **AWS S3** | 5GB (12 months) | AWS users |
| **Google Cloud Storage** | 5GB | GCP users |
| **S3-compatible** | varies | Backblaze B2, MinIO, Wasabi, DigitalOcean Spaces, self-hosted |
| **Azure Blob Storage** | 5GB (12 months) | Azure / Microsoft ecosystem users |
| **WebDAV** | Self-hosted (unlimited) | Nextcloud/ownCloud users |
| **Local WSL** | N/A (local only) | WSL-only Windows -> WSL VS Code extension sync (no cloud) |

### Step 2: Create a Bucket

If you choose **Local WSL** (visible only when running inside WSL), skip bucket creation: `init` will prompt for a local source directory and accepts either:

- a Windows path (for example `D:\PortableApps\VSCodeInsiders`) — automatically converted to `/mnt/d/...`
- a direct WSL-mounted path (for example `/mnt/c/Users/<you>/AppData/Roaming/Code`)

Common Windows VS Code locations you can enter directly:

- Installed VS Code (Stable): `C:\Users\<you>\AppData\Roaming\Code`
- Installed VS Code (Insiders): `C:\Users\<you>\AppData\Roaming\Code - Insiders`
- Portable VS Code: `<VSCodeFolder>\data\user-data`

During `init`, Local WSL now provides a source-type menu:

- Installed VS Code (Stable)
- Installed VS Code (Insiders)
- Portable / Other (manual `user-data` path input)

Manual path input is required for **Portable / Other** and should point to the VS Code `data\user-data` directory.

For cloud providers (R2/S3/GCS/S3-compatible/WebDAV), `init` now also asks:

- **Also sync VS Code extension data?** (yes/no)
- If **yes**, select source type: Stable, Insiders, or Portable/Other and provide the VS Code user-data location.

This VS Code data sync uses the same selected external provider and encryption key as the regular Claude sync.

<details>
<summary><b>Cloudflare R2</b> (recommended)</summary>

1. Go to [Cloudflare Dashboard](https://dash.cloudflare.com/) → R2 Object Storage
2. Click "Create bucket" → name it `claude-sync`
3. Go to "Manage R2 API Tokens" → "Create API Token"
4. Select **Object Read & Write** permission → Create

You'll need: Account ID, Access Key ID, Secret Access Key
</details>

<details>
<summary><b>AWS S3</b></summary>

1. Go to [S3 Console](https://s3.console.aws.amazon.com/s3/bucket/create) → Create bucket
2. Go to [IAM Security Credentials](https://console.aws.amazon.com/iam/home#/security_credentials)
3. Create Access Keys

You'll need: Access Key ID, Secret Access Key, Region
</details>

<details>
<summary><b>Google Cloud Storage</b></summary>

1. Go to [Cloud Storage](https://console.cloud.google.com/storage/create-bucket) → Create bucket
2. Go to [Service Accounts](https://console.cloud.google.com/iam-admin/serviceaccounts) → Create service account
3. Grant "Storage Object Admin" role → Create JSON key

You'll need: Project ID, Service Account JSON file (or use `gcloud auth application-default login`)
</details>

<details>
<summary><b>S3-compatible</b> (Backblaze B2, MinIO, Wasabi, DigitalOcean Spaces, ...)</summary>

Any provider exposing an S3-compatible API works through the **S3-compatible (custom endpoint)** option. Create a bucket and an application key with your provider, then supply its S3 endpoint URL.

Example (Backblaze B2):

```bash
claude-sync init --provider s3-compatible --endpoint https://s3.us-west-004.backblazeb2.com
```

You'll need: Endpoint URL, Access Key ID, Secret Access Key, Bucket. The signing region is auto-detected from the endpoint (e.g. `us-west-004`); for providers that ignore it, `auto` is used.

> Custom endpoints automatically relax the AWS SDK's default integrity-checksum headers, which some S3-compatible providers reject. AWS S3 behavior is unchanged.
</details>

<details>
<summary><b>WebDAV (Nextcloud, ownCloud, etc.)</b></summary>

No bucket to create — just point at your existing WebDAV server.

1. **Nextcloud**: Go to Settings → Security → Devices & sessions → Create app password
2. Note your WebDAV URL: `https://your-server/remote.php/dav/files/USERNAME/`

You'll need: WebDAV URL, Username, App password

The wizard will create a `claude-sync` subdirectory automatically.
</details>

<details>
<summary><b>Azure Blob Storage</b></summary>

**Option A — Azure Portal (no CLI required):**

1. Go to [Azure Portal](https://portal.azure.com) → **Storage accounts** → open your account (or create one)
2. Go to **Containers** → **+ Container** → name it `claude-sync` → Create
3. Open the `claude-sync` container → **Shared access tokens** in the left menu
4. Check all permissions: **Read, Add, Create, Write, Delete, List** → set Expiry (e.g. 2099-01-01) → click **Generate SAS token and URL**
5. Copy the **Blob SAS URL** — this is your full SAS URL

**Option B — Azure CLI:**

```bash
# Get your account key, then generate a long-lived SAS token
KEY=$(az storage account keys list --account-name <account> --query "[0].value" -o tsv)
TOKEN=$(az storage container generate-sas \
  --account-name <account> --name claude-sync \
  --permissions racwdl --expiry 2099-01-01T00:00Z \
  --account-key "$KEY" -o tsv)
echo "https://<account>.blob.core.windows.net/claude-sync?$TOKEN"
```

You'll need: the full SAS URL — `https://<account>.blob.core.windows.net/claude-sync?sv=...`

> The SAS URL is stored in `~/.claude-sync/config.yaml`. Generate it once and use the same URL on every device. You can also set `CLAUDE_SYNC_AZURE_URL` in your environment to avoid passing it on the command line.
</details>

### Step 3: Run Init

```bash
claude-sync init
```

The interactive wizard will guide you through:

1. **Select storage provider** (R2, S3, GCS, or WebDAV)
2. **Enter credentials** (provider-specific)
3. **Choose encryption method**:
   - **Passphrase** (recommended) - same passphrase on all devices = same key
   - **Random key** - must copy `~/.claude-sync/age-key.txt` to other devices
4. **Test the connection** to verify everything works

### Step 4: Push and Pull

```bash
# Upload local changes
claude-sync push

# Download remote changes
claude-sync pull
```

## What Gets Synced

| Path | Content |
|------|---------|
| `~/.claude/projects/` | Session files, auto-memory |
| `~/.claude/plans/` | Implementation plans from plan mode |
| `~/.claude/tasks/` | Task tracking state |
| `~/.claude/history.jsonl` | Command history |
| `~/.claude/agents/` | Custom agents |
| `~/.claude/skills/` | Custom skills |
| `~/.claude/plugins/` | Plugins |
| `~/.claude/rules/` | Custom rules |
| `~/.claude/settings.json` | Settings |
| `~/.claude/settings.local.json` | Local settings |
| `~/.claude/CLAUDE.md` | Global instructions |
| `~/.claude/cowork/` | Cowork project outputs and state |

**Platform-specific locations:**
- **macOS & Linux**: `~/.claude/` (home directory)
- **Windows (Desktop installer)**: `~/.claude/` (home directory)
- **Windows (Microsoft Store app)**: `AppData\Local\Packages\Claude_*\LocalCache\Roaming\Claude\` (automatically detected)

**Automatic cache exclusions:** When using the Windows Store app, claude-sync automatically excludes cache and runtime directories (~9.3 GB of VM bundles, browser caches, logs, etc.) that are regenerated on demand. This reduces sync size from 10 GB to ~20-30 MB of actual portable data.

### Sync scope

`init` asks whether to sync everything or just conversation data; you can also set it with `--scope`:

| Scope | Syncs | Use when |
|-------|-------|----------|
| `full` (default) | everything in the table above | you want settings, skills, agents, plugins, and Cowork data mirrored too |
| `sessions` | `projects/`, `history.jsonl`, `tasks/`, `plans/`, `cowork/` | you just want `claude --resume` and Cowork tasks to work across machines |

```bash
claude-sync init --scope sessions
```

**Why `sessions` exists:** `full` includes `plugins/`, whose plugin caches bundle `node_modules` and Python `.venv` trees — thousands of large, machine-/arch-specific files that are regenerated on demand and should not be synced. `sessions` skips them, keeping syncs small, fast, and portable. The scope is saved in `~/.claude-sync/config.yaml` and applies to every `push`/`pull`.

### Automatically Excluded Directories

When using the Windows Store app, the following are automatically excluded (these are regenerated on demand and shouldn't be synced):

- **vm_bundles/** — Cowork VM runtime (~9.3 GB)
- **claude-code/** — VS Code runtime binaries
- **Cache/**, **Code Cache/** — Browser and compilation caches
- **GPUCache/**, **DawnGraphiteCache/** — Graphics rendering caches
- **dxt-install-*/** — Temporary installers
- **logs/** — Application logs
- **Preferences**, **Session Storage/** — UI state (regenerated per session)
- **sentry/** — Error reporting data

If you want to manually override these exclusions, add an `exclude:` list to your config file (will disable automatic exclusions).

## Cross-Device Path Mapping

Claude Code indexes project sessions by **absolute filesystem path**:

```
/Users/alice/my-app → ~/.claude/projects/-Users-alice-my-app/
/Users/bob/my-app   → ~/.claude/projects/-Users-bob-my-app/
```

Synced verbatim, those would be **different projects** and `claude --resume` on the second machine would never find the first machine's sessions. claude-sync solves this by translating paths during sync:

- **Home directories are mapped automatically.** Sessions are stored remotely under a portable `${HOME}` token (in both remote keys and transcript content), then rewritten to each device's real home on pull. Different usernames across machines just work.
- **Other layout differences are configurable.** If one machine keeps projects in `~/work` and another in `~/Projects`, point both at the same token in `~/.claude-sync/config.yaml`:

  ```yaml
  # machine 1
  path_map:
    ~/work: WORK
  ```

  ```yaml
  # machine 2
  path_map:
    ~/Projects: WORK
  ```

  Sessions under either directory sync to the shared `${WORK}` namespace and resume correctly on both machines.

**Upgrading from an older version?** Run `claude-sync migrate` once on each device to convert existing remote data to portable keys. Paths the current device doesn't own are left for the other device's migrate run.

## Commands

```bash
claude-sync init        # Set up configuration (interactive wizard)
claude-sync push        # Upload local changes to cloud storage
claude-sync pull        # Download remote changes from cloud storage
claude-sync status      # Show pending local changes
claude-sync diff        # Show differences between local and remote
claude-sync conflicts   # List and resolve conflicts
claude-sync reset       # Reset configuration (forgot passphrase)
claude-sync migrate     # Convert legacy remote keys to portable path-mapped keys
claude-sync update      # Update to latest version (verifies release checksums)
claude-sync changelog   # Show release history
claude-sync --help      # Show all commands
```

### Pull Options

```bash
claude-sync pull              # Normal pull (prompts if existing files)
claude-sync pull --dry-run    # Preview what would change
claude-sync pull --force      # Skip confirmation prompts
```

### Init Options

```bash
claude-sync init              # Full setup wizard
claude-sync init --passphrase # Re-enter passphrase only (keeps storage config)
claude-sync init --force      # Reset everything, start fresh
```

### Quiet Mode

```bash
claude-sync push -q     # No output (for scripts)
claude-sync pull -q
```

### Check for Updates

```bash
claude-sync update --check   # Check without installing
claude-sync update           # Download and install latest version
```

### Changelog

```bash
claude-sync changelog            # Show recent releases
claude-sync changelog --limit 5  # Show last 5 releases
```

## Exclude Patterns

Skip specific files or directories during sync by adding exclude patterns to your config (`~/.claude-sync/config.yaml`):

```yaml
exclude:
  - "*.tmp"
  - "projects/*/node_modules/*"
  - "projects/*/.git/*"
```

Patterns use glob syntax and are matched against paths relative to `~/.claude`.

## Shell Integration

Add to `~/.zshrc` or `~/.bashrc`:

```bash
# Auto-pull on shell start
if command -v claude-sync &> /dev/null; then
  # Run in a subshell so the job is detached from the parent shell's
  # job table — avoids interactive `[1] 12345` / `[1] + done` noise.
  (claude-sync pull -q &) >/dev/null 2>&1
fi

# Auto-push on shell exit
trap 'claude-sync push -q' EXIT
```

> **Note:** The subshell wrapper `(cmd &)` prevents zsh/bash from printing job control
> messages (`[1] 12345` on start and `[1] + done cmd` on completion) every time you open
> a terminal. A plain `claude-sync pull -q &` works but produces noisy shell prompts.

## Pulling with Existing Files

When you pull on a device that already has `~/.claude` files, claude-sync will:

1. **Show what would change** - files that would be overwritten, kept, or downloaded
2. **Ask for confirmation** - choose to backup, overwrite, or abort
3. **Create a backup** - saves existing files to `~/.claude.backup.{timestamp}`

```bash
# Preview first
claude-sync pull --dry-run

# Pull with prompts
claude-sync pull

# Skip prompts (for scripts)
claude-sync pull --force
```

## Conflict Resolution

When both local and remote files change, the remote version is saved as `.conflict`:

```bash
claude-sync conflicts            # Interactive resolution
claude-sync conflicts --list     # Just list conflicts
claude-sync conflicts --keep local   # Keep all local versions
claude-sync conflicts --keep remote  # Keep all remote versions
```

Interactive options:
- **[l]** Keep local (delete conflict file)
- **[r]** Keep remote (replace local)
- **[d]** Show diff
- **[s]** Skip
- **[q]** Quit

## Wrong Passphrase?

If you entered the wrong passphrase on a new device:

```bash
# Re-enter passphrase (keeps your storage config)
claude-sync init --passphrase
```

The init will verify your passphrase can decrypt remote files before completing.

## Forgot Passphrase?

The passphrase is **never stored**. If you forget it:

1. Your encrypted files cannot be recovered
2. Reset and start fresh:

```bash
claude-sync reset --remote   # Delete remote files and local config
claude-sync init             # Set up again with new passphrase
claude-sync push             # Re-upload from this device
```

## Security

- Files compressed with gzip, then encrypted with [age](https://github.com/FiloSottile/age) before upload
- Passphrase-derived keys use Argon2 (memory-hard KDF)
- Passphrase is never stored - only the derived key at `~/.claude-sync/age-key.txt`
- Cloud storage is private (API key/IAM auth)
- Config files and downloads stored with 0600/0700 permissions (user-only)
- Self-update verifies SHA256 checksums before installing new binaries
- Backward compatible: can read both compressed and uncompressed remote files

## Cost

Claude sessions typically use < 50MB. Syncing is effectively **free** on any provider:

| Provider | Free Tier |
|----------|-----------|
| **Cloudflare R2** | 10GB storage, 1M writes, 10M reads/month |
| **AWS S3** | 5GB for 12 months (then ~$0.023/GB) |
| **Google Cloud Storage** | 5GB, 5K writes, 50K reads/month |
| **WebDAV** | Self-hosted — no limits, no cost beyond your own server |

## Installation Options

### npm (recommended)

**Prerequisite:** Node.js 14+ (no Go required - downloads pre-compiled binary)

```bash
# Global install
npm install -g @tawandotorg/claude-sync

# Or one-time use
npx @tawandotorg/claude-sync init
```

### GitHub Packages

**Prerequisite:** Node.js 14+

```bash
# Add to ~/.npmrc
echo "@tawanorg:registry=https://npm.pkg.github.com" >> ~/.npmrc

# Install
npm install -g @tawanorg/claude-sync
```

### Download Binary

**Prerequisite:** None

```bash
# macOS ARM (M1/M2/M3)
curl -L https://github.com/eschweighofer/claude-sync/releases/latest/download/claude-sync-darwin-arm64 -o claude-sync
chmod +x claude-sync
sudo mv claude-sync /usr/local/bin/
```

See [GitHub Releases](https://github.com/eschweighofer/claude-sync/releases) for all platforms.

### Go Install

**Prerequisite:** Go 1.21+ (for developers)

```bash
go install github.com/eschweighofer/claude-sync/cmd/claude-sync@latest
```

### Build from Source

**Prerequisite:** Go 1.21+

```bash
git clone https://github.com/eschweighofer/claude-sync
cd claude-sync
make build
./bin/claude-sync --version
```

## Windows Store App & Cowork Support

### Windows Microsoft Store App

Claude Desktop is now distributed via the Microsoft Store as an MSIX package on Windows. Claude-sync automatically detects and syncs from the Store app location:

```
AppData\Local\Packages\Claude_<hash>\LocalCache\Roaming\Claude\
```

**No configuration needed** — claude-sync checks for the Store app first, then falls back to the traditional `~/.claude/` location. Both work seamlessly.

### Cowork Data

Cowork project outputs and state are now included in sync scopes:

- **`full` scope**: All Cowork data in `~/.claude/cowork/`
- **`sessions` scope**: Cowork task files and project outputs

This allows Cowork tasks and project outputs to sync across devices without manual file management. Configure during `init`:

```bash
claude-sync init --scope full    # Sync Claude Code + Cowork
claude-sync init --scope sessions # Sync just conversations + Cowork tasks
```

## Development

```bash
make test          # Run tests
make fmt           # Format code
make check         # Run all pre-commit checks
make build-all     # Build for all platforms
make setup-hooks   # Enable git pre-commit hooks
```

## License

MIT
