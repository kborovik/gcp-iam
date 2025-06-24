package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/kborovik/gcp-iam/config"
	"github.com/kborovik/gcp-iam/db"
	"github.com/kborovik/gcp-iam/update"
	"github.com/urfave/cli/v3"
)

// normalizeRoleName strips the "roles/" prefix if present
func normalizeRoleName(roleName string) string {
	if strings.HasPrefix(roleName, "roles/") {
		return strings.TrimPrefix(roleName, "roles/")
	}
	return roleName
}

// completeRoleNames provides completion for role names
func completeRoleNames(ctx context.Context, cmd *cli.Command) {
	// Check if this is being called for completion
	if os.Getenv("COMP_LINE") == "" {
		return
	}

	cfg, err := config.Load()
	if err != nil {
		return
	}

	database, err := db.New(cfg.DatabasePath)
	if err != nil {
		return
	}
	defer database.Close()

	roleNames, err := database.GetRoleNames()
	if err != nil {
		return
	}

	// Filter based on current input if available
	currentArg := ""
	if args := cmd.Args(); args.Len() > 0 {
		currentArg = args.First()
	}

	for _, name := range roleNames {
		if currentArg == "" || strings.HasPrefix(name, currentArg) {
			fmt.Println(name)
		}
	}
}

var cmd = &cli.Command{
	Name:                  "gcp-iam",
	Usage:                 "Query Google Cloud IAM Roles and Permissions",
	Suggest:               true,
	EnableShellCompletion: true,
	HideHelpCommand:       true,
	Action: func(ctx context.Context, cmd *cli.Command) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		database, err := db.New(cfg.DatabasePath)
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		defer database.Close()

		if cmd.Metadata == nil {
			cmd.Metadata = make(map[string]any)
		}
		cmd.Metadata["config"] = cfg
		cmd.Metadata["db"] = database

		return cli.ShowAppHelp(cmd)
	},
	Commands: []*cli.Command{
		{
			Name:  "role",
			Usage: "Query IAM Roles",
			CommandNotFound: func(ctx context.Context, cmd *cli.Command, command string) {
				cli.ShowAppHelp(cmd)
			},
			Commands: []*cli.Command{
				{
					Name:  "show",
					Usage: "Show IAM role permissions",
					ShellComplete: func(ctx context.Context, cmd *cli.Command) {
						completeRoleNames(ctx, cmd)
					},
					Action: func(ctx context.Context, cmd *cli.Command) error {
						roleName := cmd.Args().First()
						if roleName == "" {
							return fmt.Errorf("role name is required")
						}

						// Normalize role name (strip roles/ prefix if present)
						roleName = normalizeRoleName(roleName)

						cfg, err := config.Load()
						if err != nil {
							return fmt.Errorf("failed to load config: %w", err)
						}

						database, err := db.New(cfg.DatabasePath)
						if err != nil {
							return fmt.Errorf("failed to open database: %w", err)
						}
						defer database.Close()

						role, err := database.GetRoleByName(roleName)
						if err != nil {
							return fmt.Errorf("failed to get role: %w", err)
						}

						if role == nil {
							fmt.Printf("Role '%s' not found\n", roleName)
							return nil
						}

						fmt.Printf("Role: %s\n", role.Name)
						fmt.Printf("Title: %s\n", role.Title)
						fmt.Printf("Description: %s\n", role.Description)
						fmt.Printf("Stage: %s\n", role.Stage)

						permissions, err := database.GetRolePermissions(role.Name)
						if err != nil {
							return fmt.Errorf("failed to get permissions: %w", err)
						}

						fmt.Printf("Permissions (%d):\n", len(permissions))
						for _, perm := range permissions {
							fmt.Printf("  - %s\n", perm.Permission)
						}

						return nil
					},
				},
				{
					Name:  "search",
					Usage: "Search IAM roles",
					ShellComplete: func(ctx context.Context, cmd *cli.Command) {
						completeRoleNames(ctx, cmd)
					},
					Action: func(ctx context.Context, cmd *cli.Command) error {
						query := cmd.Args().First()
						if query == "" {
							return fmt.Errorf("search query is required")
						}

						cfg, err := config.Load()
						if err != nil {
							return fmt.Errorf("failed to load config: %w", err)
						}

						database, err := db.New(cfg.DatabasePath)
						if err != nil {
							return fmt.Errorf("failed to open database: %w", err)
						}
						defer database.Close()

						roles, err := database.SearchRoles(query)
						if err != nil {
							return fmt.Errorf("failed to search roles: %w", err)
						}

						fmt.Printf("Found %d roles matching '%s':\n", len(roles), query)
						for _, role := range roles {
							fmt.Printf("  %s - %s\n", role.Name, role.Title)
						}

						return nil
					},
				},
				{
					Name:  "compare",
					Usage: "Compare permissions of 2 IAM roles",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						fmt.Println(cmd.FullName(), cmd.Args().First())
						return nil
					},
				},
			},
		},
		{
			Name:  "permission",
			Usage: "Query IAM Permissions",
			CommandNotFound: func(ctx context.Context, cmd *cli.Command, command string) {
				cli.ShowAppHelp(cmd)
			},
			Commands: []*cli.Command{
				{
					Name:  "show",
					Usage: "Show IAM roles with permission",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						permissionName := cmd.Args().First()
						if permissionName == "" {
							return fmt.Errorf("permission name is required")
						}

						cfg, err := config.Load()
						if err != nil {
							return fmt.Errorf("failed to load config: %w", err)
						}

						database, err := db.New(cfg.DatabasePath)
						if err != nil {
							return fmt.Errorf("failed to open database: %w", err)
						}
						defer database.Close()

						permission, err := database.GetPermissionByName(permissionName)
						if err != nil {
							return fmt.Errorf("failed to get permission: %w", err)
						}

						if permission == nil {
							fmt.Printf("Permission '%s' not found\n", permissionName)
							return nil
						}

						fmt.Printf("Permission: %s\n", permission.Permission)
						fmt.Printf("Created At: %s\n", permission.CreatedAt.Format("2006-01-02 15:04:05"))

						roles, err := database.GetRolesWithPermission(permission.Permission)
						if err != nil {
							return fmt.Errorf("failed to get roles with permission: %w", err)
						}

						fmt.Printf("Roles with this permission (%d):\n", len(roles))
						for _, role := range roles {
							fmt.Printf("  %s - %s\n", role.Name, role.Title)
						}

						return nil
					},
				},
				{
					Name:  "search",
					Usage: "Search IAM permissions",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						query := cmd.Args().First()
						if query == "" {
							return fmt.Errorf("search query is required")
						}

						cfg, err := config.Load()
						if err != nil {
							return fmt.Errorf("failed to load config: %w", err)
						}

						database, err := db.New(cfg.DatabasePath)
						if err != nil {
							return fmt.Errorf("failed to open database: %w", err)
						}
						defer database.Close()

						permissions, err := database.SearchPermissions(query)
						if err != nil {
							return fmt.Errorf("failed to search permissions: %w", err)
						}

						fmt.Printf("Found %d permissions matching '%s':\n", len(permissions), query)
						for _, perm := range permissions {
							fmt.Printf("  %s\n", perm.Permission)
						}

						return nil
					},
				},
			},
		},
		{
			Name:  "update",
			Usage: "Update IAM roles and permissions",
			Action: func(ctx context.Context, cmd *cli.Command) error {
				cfg, err := config.Load()
				if err != nil {
					return fmt.Errorf("failed to load config: %w", err)
				}

				database, err := db.New(cfg.DatabasePath)
				if err != nil {
					return fmt.Errorf("failed to open database: %w", err)
				}
				defer database.Close()

				updater := update.New(database)

				// First update all roles
				err = updater.UpdateRoles(ctx)
				if err != nil {
					// Check if it's an authentication error and return it directly
					if strings.Contains(err.Error(), "Authentication failed accessing Google Cloud IAM API") {
						return err
					}
					return fmt.Errorf("failed to update roles: %w", err)
				}

				// Then update permissions only for roles that need it
				fmt.Println("Identifying roles needing permission updates...")
				rolesToUpdate, err := database.GetRolesNeedingPermissionUpdate()
				if err != nil {
					return fmt.Errorf("failed to get roles needing updates: %w", err)
				}

				if len(rolesToUpdate) == 0 {
					fmt.Println("No roles need permission updates - all roles are up to date")
				} else {
					fmt.Printf("Updating permissions for %d roles that need updates...\n", len(rolesToUpdate))
					for i, role := range rolesToUpdate {
						fmt.Printf("Updating permissions for role %d/%d: %s\n", i+1, len(rolesToUpdate), role.Name)
						err = updater.UpdatePermissions(ctx, role.Name)
						if err != nil {
							fmt.Printf("Warning: failed to update permissions for role %s: %v\n", role.Name, err)
							// Continue with other roles even if one fails
						}
					}
				}

				fmt.Println("Successfully updated IAM roles and permissions")
				return nil
			},
		},
		{
			Name:   "complete-roles",
			Usage:  "List all role names for shell completion",
			Hidden: true,
			Action: func(ctx context.Context, cmd *cli.Command) error {
				cfg, err := config.Load()
				if err != nil {
					return err
				}

				database, err := db.New(cfg.DatabasePath)
				if err != nil {
					return err
				}
				defer database.Close()

				roleNames, err := database.GetRoleNames()
				if err != nil {
					return err
				}

				for _, name := range roleNames {
					fmt.Println(name)
				}
				return nil
			},
		},
		{
			Name:  "info",
			Usage: "Show application configuration",
			Action: func(ctx context.Context, cmd *cli.Command) error {
				cfg, err := config.Load()
				if err != nil {
					return fmt.Errorf("failed to load config: %w", err)
				}

				database, err := db.New(cfg.DatabasePath)
				if err != nil {
					return fmt.Errorf("failed to open database: %w", err)
				}
				defer database.Close()

				roleCount, err := database.CountRoles()
				if err != nil {
					return fmt.Errorf("failed to count roles: %w", err)
				}

				permissionCount, err := database.CountPermissions()
				if err != nil {
					return fmt.Errorf("failed to count permissions: %w", err)
				}

				configPath, _ := config.GetDefaultConfigPath()
				fmt.Println("GCP IAM Configuration:")
				fmt.Printf("  Roles:        %d\n", roleCount)
				fmt.Printf("  Permissions:  %d\n", permissionCount)
				fmt.Printf("  ConfigFile:   %s\n", configPath)
				fmt.Printf("  DatabasePath: %s\n", cfg.DatabasePath)

				return nil
			},
		},
	},
}

func main() {
	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatalf("failed to run command: %v", err)
	}
}
