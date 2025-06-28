---
allowed-tools: Bash(git add:*), Bash(git status:*), Bash(git commit:*), Bash(go run:*)
description: Create a git commit
---

# Your task

Based on the recent code changes:

- Run `go fmt ./...` before committing
- Determine version from committed changes, use golang version format `v1.0.0`
- Update application version in `main.go` file
- Create a single git commit.
