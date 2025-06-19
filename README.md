# Google Cloud IAM Roles and Permissions

A command-line interface (CLI) tool written in Go that displays Google Cloud IAM predefined roles and their associated permissions. This tool helps developers, system administrators, and security professionals understand and explore Google Cloud's IAM role hierarchy.

## Features

### Core Functionality

1. Roles:

- Search pre-defined roles, simple search, display role name i.e. `compute.imageUser`,
- List pre-defined role permissions, display role permissions i.e `compute.imageUser` has permissions `compute.images.create`, `compute.images.delete`, `compute.images.get`, `compute.images.list`, `compute.images.update`
- Compare permissions between 2 roles

1. Permissions:

- Search permissions and display roles that have those permissions

1. Sync:

- Synchronize permissions and roles

## Architecture

Storage database is `sqlite3`
