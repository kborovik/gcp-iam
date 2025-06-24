# Fish shell completion for gcp-iam

# Function to get role names for completion
function __gcp_iam_role_names
    # Try to find gcp-iam in PATH, fallback to common locations
    if command -q gcp-iam
        gcp-iam complete-roles 2>/dev/null
    else if test -x ./gcp-iam
        ./gcp-iam complete-roles 2>/dev/null
    else if test -x /usr/local/bin/gcp-iam
        /usr/local/bin/gcp-iam complete-roles 2>/dev/null
    end
end

# Complete role names for 'gcp-iam role show' and 'gcp-iam role search'
complete -c gcp-iam -n '__fish_seen_subcommand_from role; and __fish_seen_subcommand_from show; and not __fish_seen_subcommand_from help' -f -a '(__gcp_iam_role_names)'
complete -c gcp-iam -n '__fish_seen_subcommand_from role; and __fish_seen_subcommand_from search; and not __fish_seen_subcommand_from help' -f -a '(__gcp_iam_role_names)'

# Basic command completion
complete -c gcp-iam -n '__fish_use_subcommand' -f -a 'role' -d 'Query IAM Roles'
complete -c gcp-iam -n '__fish_use_subcommand' -f -a 'permission' -d 'Query IAM Permissions'
complete -c gcp-iam -n '__fish_use_subcommand' -f -a 'update' -d 'Update IAM roles and permissions'
complete -c gcp-iam -n '__fish_use_subcommand' -f -a 'info' -d 'Show application configuration'

# Role subcommands
complete -c gcp-iam -n '__fish_seen_subcommand_from role; and not __fish_seen_subcommand_from show search compare' -f -a 'show' -d 'Show IAM role permissions'
complete -c gcp-iam -n '__fish_seen_subcommand_from role; and not __fish_seen_subcommand_from show search compare' -f -a 'search' -d 'Search IAM roles'
complete -c gcp-iam -n '__fish_seen_subcommand_from role; and not __fish_seen_subcommand_from show search compare' -f -a 'compare' -d 'Compare permissions of 2 IAM roles'

# Permission subcommands
complete -c gcp-iam -n '__fish_seen_subcommand_from permission; and not __fish_seen_subcommand_from show search' -f -a 'show' -d 'Show IAM roles with permission'
complete -c gcp-iam -n '__fish_seen_subcommand_from permission; and not __fish_seen_subcommand_from show search' -f -a 'search' -d 'Search IAM permissions'