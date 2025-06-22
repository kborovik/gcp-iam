# Google Cloud IAM Roles and Permissions

A command-line interface (CLI) tool written in Go that displays Google Cloud IAM predefined roles and their associated permissions. This tool helps developers, system administrators, and security professionals understand and explore Google Cloud's IAM role hierarchy.

# Features

## Core Functionality

The application consist of 4 main commands

### Role

Command:

- `role` - Query IAM Roles

Sub-commands:

- `search` - Search pre-defined roles, simple search, display role name i.e. `compute.imageUser`,
- `show` - Show pre-defined role permissions formatted as table
- `compare` - Compare permissions between 2 roles

### Permission

Command:

- `permission` - Query IAM Permissions

Sub-commands:

- `show` - Show permission details
- `search` - Search permissions
- `compare` - Compare 2 permissions roles

### Update

Command:

- `update` - Update IAM roles and permissions

  - Creates DB schema
  - Stores Google pre-defined IAM Roles and Permissions in sqlite3 database

### Info

Command:

- `info` - displays application configuration details

## Software Components

- Command-line interface built using github.com/urfave/cli/v3
- Visual interface elements (table) using github.com/charmbracelet/lipgloss/table
- Data storage utilizes sqlite3 database
