# GCP IAM Explorer

A fast command-line tool for exploring Google Cloud IAM roles and permissions. Discover what permissions are included in roles, find roles that have specific permissions, and understand GCP's IAM structure with ease.

## ✨ Features

- 🔍 **Search & Explore** - Find roles and permissions with instant search
- ⚡ **Fast Performance** - Local database for lightning-fast queries
- 🐚 **TAB Completion** - Fish shell completion for role and permission names
- 🔄 **Always Current** - Update from live GCP IAM API
- 💡 **User-Friendly** - Clean output and helpful error messages

## 🚀 Quick Start

```bash
# Update local database with latest GCP IAM data
gcp-iam update --roles --services

# See what's in your database
gcp-iam info

# Explore roles and permissions
gcp-iam role search storage
gcp-iam role show editor
gcp-iam permission show storage.objects.get
```

## 📋 Commands

### 🔍 Explore Roles

```bash
# Search for roles by name or description
gcp-iam role search compute
gcp-iam role search "storage admin"

# Show detailed information about a role
gcp-iam role show editor
gcp-iam role show compute.admin
gcp-iam role show roles/storage.admin  # also works with full name

# Compare two roles to see permission differences
gcp-iam role compare editor viewer
```

### 🔐 Explore Permissions

```bash
# Search for permissions
gcp-iam permission search storage
gcp-iam permission search "compute.instances"

# See which roles include a specific permission
gcp-iam permission show storage.objects.get
gcp-iam permission show compute.instances.create
```

### 🔄 Data Management

```bash
# Update your local database with latest GCP IAM data
gcp-iam update --roles --services   # Update both roles and services
gcp-iam update --roles              # Update only roles and permissions
gcp-iam update --services           # Update only services

# View database statistics and configuration
gcp-iam info
```

## 💡 Example Workflows

### Find the right role for storage access

```bash
# Find storage-related roles
$ gcp-iam role search storage

Found 15 roles matching 'storage':
  storage.admin - Storage Admin
  storage.objectAdmin - Storage Object Admin
  storage.objectCreator - Storage Object Creator
  storage.objectViewer - Storage Object Viewer
  ...

# Check what the Storage Object Admin can do
$ gcp-iam role show storage.objectAdmin

Role: storage.objectAdmin
Title: Storage Object Admin
Description: Full control of GCS objects
Stage: GA
Permissions (12):
  - storage.objects.create
  - storage.objects.delete
  - storage.objects.get
  - storage.objects.list
  - storage.objects.update
  ...
```

### Find which roles can create compute instances

```bash
# Search for the permission
$ gcp-iam permission show compute.instances.create

Permission: compute.instances.create
Roles with this permission (8):
  compute.admin - Compute Admin
  compute.instanceAdmin - Compute Instance Admin
  editor - Editor
  owner - Owner
  ...
```

### Compare similar roles

```bash
# See the difference between editor and viewer
$ gcp-iam role compare editor viewer

Permissions only in 'editor': 2,847 permissions
Permissions only in 'viewer': 0 permissions
Common permissions: 7,189 permissions

# Editor has all viewer permissions plus 2,847 additional permissions
```

## 🛠️ Setup TAB Completion (Fish Shell)

Enable instant TAB completion for role and permission names:

```bash
# Copy completion script
cp tools/gcp-iam.fish ~/.config/fish/completions/

# Now you can TAB complete!
gcp-iam role show ed<TAB>           # completes to 'editor'
gcp-iam permission show storage.<TAB>  # shows all storage.* permissions
```

## 🔐 Authentication

To update data from GCP, you need to authenticate:

```bash
# Authenticate with Google Cloud
gcloud auth login --update-adc

# Then update your local database
gcp-iam update --roles --services
```

If you see authentication errors, the tool will guide you with the exact command to run.

## 📊 What's Included

- **1,892 IAM roles** - All predefined GCP roles
- **11,530 permissions** - Every GCP permission available
- **Fast search** - Find roles/permissions in milliseconds
- **Local storage** - No API calls needed for browsing
- **Auto-updates** - Keep data current with GCP changes

## 🆘 Common Use Cases

**Security Auditing**: Understand what permissions a role actually grants

```bash
gcp-iam role show iam.serviceAccountUser
```

**Least Privilege**: Find the minimal role for specific permissions

```bash
gcp-iam permission show pubsub.topics.publish
```

**Role Discovery**: Find roles for specific GCP services

```bash
gcp-iam role search kubernetes
gcp-iam role search bigquery
```

**Permission Research**: Understand GCP permission structure

```bash
gcp-iam permission search "instances.create"
gcp-iam permission search "buckets"
```

---

**Note**: This tool reads GCP IAM data but never modifies your actual GCP resources or permissions. It's safe for exploration and auditing.
