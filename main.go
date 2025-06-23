package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/kborovik/gcp-iam/config"
	"github.com/kborovik/gcp-iam/db"
	"github.com/kborovik/gcp-iam/update"
	"github.com/urfave/cli/v3"
)

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
					Action: func(ctx context.Context, cmd *cli.Command) error {
						roleName := cmd.Args().First()
						if roleName == "" {
							return fmt.Errorf("role name is required")
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
							fmt.Printf("  - %s\n", perm.Name)
						}

						return nil
					},
				},
				{
					Name:  "search",
					Usage: "Search IAM roles",
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

						fmt.Printf("Permission: %s\n", permission.Name)
						fmt.Printf("Title: %s\n", permission.Title)
						fmt.Printf("Description: %s\n", permission.Description)
						fmt.Printf("Stage: %s\n", permission.Stage)

						roles, err := database.GetRolesWithPermission(permission.Name)
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
							fmt.Printf("  %s - %s\n", perm.Name, perm.Title)
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
					return fmt.Errorf("failed to update roles: %w", err)
				}

				// Then update permissions for all roles
				fmt.Println("Fetching permissions for all roles...")
				roles, err := database.GetAllRoles()
				if err != nil {
					return fmt.Errorf("failed to get roles from database: %w", err)
				}

				for i, role := range roles {
					fmt.Printf("Updating permissions for role %d/%d: %s\n", i+1, len(roles), role.Name)
					err = updater.UpdatePermissions(ctx, role.Name)
					if err != nil {
						fmt.Printf("Warning: failed to update permissions for role %s: %v\n", role.Name, err)
						// Continue with other roles even if one fails
					}
				}

				fmt.Println("Successfully updated IAM roles and permissions")
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

				configPath, _ := config.GetDefaultConfigPath()
				fmt.Println("GCP IAM Configuration:")
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
