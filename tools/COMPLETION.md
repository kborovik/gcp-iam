# Fish Shell Completion Setup for gcp-iam

This document explains how to set up TAB completion for role names in the `gcp-iam` command.

## Setup Instructions

### Option 1: Install globally (recommended)

1. **Copy the completion script to fish completions directory:**
   ```bash
   cp gcp-iam.fish ~/.config/fish/completions/
   ```

2. **Make sure the `gcp-iam` binary is in your PATH:**
   ```bash
   # Add to your PATH or copy to a directory in PATH
   sudo cp gcp-iam /usr/local/bin/
   ```

3. **Reload fish shell or restart your terminal**

### Option 2: Source manually

1. **Source the completion script in your fish config:**
   ```bash
   echo "source $(pwd)/gcp-iam.fish" >> ~/.config/fish/config.fish
   ```

2. **Reload fish configuration:**
   ```bash
   source ~/.config/fish/config.fish
   ```

## Usage

Once set up, you can use TAB completion for role names:

```bash
# TAB completion for role names
gcp-iam role show ed<TAB>
# Completes to: gcp-iam role show editor

gcp-iam role show compute.<TAB>
# Shows all roles starting with "compute."

gcp-iam role search storage<TAB>
# Shows all roles containing "storage"
```

## Features

- **Role name completion** for `gcp-iam role show <TAB>`
- **Role name completion** for `gcp-iam role search <TAB>`
- **Command and subcommand completion** for all gcp-iam commands
- **Fast completion** using dedicated `complete-roles` command

## Troubleshooting

### Completion not working
1. Verify fish can find the completion script:
   ```bash
   ls ~/.config/fish/completions/gcp-iam.fish
   ```

2. Test the completion function manually:
   ```bash
   gcp-iam complete-roles
   ```

3. Reload fish completions:
   ```bash
   fish_reload_completions
   ```

### Slow completion
The completion uses a dedicated `gcp-iam complete-roles` command that quickly retrieves role names from the local database. If completion is slow, ensure:

1. Database exists and is populated:
   ```bash
   gcp-iam info
   ```

2. Run update if needed:
   ```bash
   gcp-iam update
   ```

## Technical Details

The completion system works by:

1. **Fish completion script** (`gcp-iam.fish`) defines completion rules
2. **Dedicated completion command** (`gcp-iam complete-roles`) provides role names
3. **Fast database query** retrieves role names from local SQLite database
4. **Context-aware completion** only triggers for appropriate commands

The completion is context-aware and only suggests role names when appropriate (e.g., after `gcp-iam role show` or `gcp-iam role search`).