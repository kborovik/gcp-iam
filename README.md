# Google Cloud IAM Roles and Permissions

A command-line interface (CLI) tool written in Go that displays Google Cloud IAM predefined roles and their associated permissions. This tool helps developers, system administrators, and security professionals understand and explore Google Cloud's IAM role hierarchy.

## Features

### Local SQLite Database

- Stores Google Cloud IAM roles and permissions locally
- Fast querying and searching capabilities
- Automatic database creation and schema management
- Uses name-based primary keys for efficient storage

### Configuration Management

- YAML-based configuration file at `~/.gcp-iam/config.yaml`
- Automatic directory and file creation
- Default database location: `~/.gcp-iam/database.sqlite`
- Configurable log levels and cache directory

## Commands

### Role Commands

- `role show <role-name>` - Display role details and permissions
- `role search <query>` - Search roles by name, title, or description
- `role compare <role1> <role2>` - Compare permissions between two roles

### Permission Commands

- `permission show <permission-name>` - Show permission details
- `permission search <query>` - Search permissions by name or description

### Management Commands

- `update` - Update local database with latest Google Cloud IAM data
- `info` - Display current application configuration

## Quick Start

```bash
# View application configuration
gcp-iam info

# Search for compute-related roles
gcp-iam role search compute

# Show details of a specific role
gcp-iam role show roles/compute.admin
```

## Architecture

### Software Components

- **CLI Framework**: Built with `github.com/urfave/cli/v3`
- **Database**: SQLite with pure Go driver (`modernc.org/sqlite`)
- **Configuration**: YAML-based config with `gopkg.in/yaml.v3`
- **Testing**: Comprehensive test coverage for all components

### Database Schema

```sql
-- Roles table (name as primary key)
CREATE TABLE roles (
    name TEXT PRIMARY KEY,
    title TEXT,
    description TEXT,
    stage TEXT,
    deleted BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Permissions table (name as primary key)
CREATE TABLE permissions (
    name TEXT PRIMARY KEY,
    title TEXT,
    description TEXT,
    stage TEXT,
    api_disabled BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Many-to-many relationship
CREATE TABLE role_permissions (
    role_name TEXT,
    permission_name TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (role_name, permission_name),
    FOREIGN KEY (role_name) REFERENCES roles(name) ON DELETE CASCADE,
    FOREIGN KEY (permission_name) REFERENCES permissions(name) ON DELETE CASCADE
);
```

### Project Structure

```
├── config/          # Configuration management
│   ├── config.go    # Config loading and defaults
│   └── config_test.go
├── db/              # Database layer
│   ├── database.go  # Connection and schema
│   ├── models.go    # Data models and CRUD operations
│   └── database_test.go
├── main.go          # CLI application entry point
├── main_test.go     # CLI functionality tests
└── go.mod           # Go module dependencies
```

## Configuration

The application uses a YAML configuration file located at `~/.gcp-iam/config.yaml`:

```yaml
log_level: info
database_path: /home/user/.gcp-iam/database.sqlite
cache_dir: /home/user/.gcp-iam/cache
```

All paths and directories are created automatically on first run.
