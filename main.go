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

const Version = "1.0.2"

// normalizeRoleName strips the "roles/" prefix if present
func normalizeRoleName(roleName string) string {
	if after, ok := strings.CutPrefix(roleName, "roles/"); ok {
		return after
	}
	return roleName
}

// completeRoleNames provides completion for role names
func completeRoleNames(cmd *cli.Command) {
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

// completePermissionNames provides completion for permission names
func completePermissionNames(cmd *cli.Command) {
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

	permissionNames, err := database.GetPermissionNames()
	if err != nil {
		return
	}

	// Filter based on current input if available
	currentArg := ""
	if args := cmd.Args(); args.Len() > 0 {
		currentArg = args.First()
	}

	for _, name := range permissionNames {
		if currentArg == "" || strings.HasPrefix(name, currentArg) {
			fmt.Println(name)
		}
	}
}

// completeServiceNames provides completion for service names
func completeServiceNames(cmd *cli.Command) {
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

	serviceNames, err := database.GetServiceNames()
	if err != nil {
		return
	}

	// Filter based on current input if available
	currentArg := ""
	if args := cmd.Args(); args.Len() > 0 {
		currentArg = args.First()
	}

	for _, name := range serviceNames {
		if currentArg == "" || strings.HasPrefix(name, currentArg) {
			fmt.Println(name)
		}
	}
}

var cmd = &cli.Command{
	Name:                  "gcp-iam",
	Usage:                 "Query Google Cloud IAM Roles and Permissions",
	Version:               Version,
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
					Name:      "show",
					Usage:     "Show IAM role permissions",
					ArgsUsage: "<role-name>",
					Description: "Display detailed information about a specific IAM role including its permissions.\n\n" +
						"Examples:\n" +
						"  gcp-iam role show viewer\n" +
						"  gcp-iam role show compute.instanceAdmin.v1",
					ShellComplete: func(ctx context.Context, cmd *cli.Command) {
						completeRoleNames(cmd)
					},
					Action: func(ctx context.Context, cmd *cli.Command) error {
						roleName := cmd.Args().First()
						if roleName == "" {
							return cli.ShowSubcommandHelp(cmd)
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
					Name:      "search",
					Usage:     "Search IAM roles",
					ArgsUsage: "<search-query>",
					Description: "Search for IAM roles by name or title using a query string.\n\n" +
						"Examples:\n" +
						"  gcp-iam role search storage\n" +
						"  gcp-iam role search admin\n" +
						"  gcp-iam role search compute",
					ShellComplete: func(ctx context.Context, cmd *cli.Command) {
						completeRoleNames(cmd)
					},
					Action: func(ctx context.Context, cmd *cli.Command) error {
						query := cmd.Args().First()
						if query == "" {
							return cli.ShowSubcommandHelp(cmd)
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
					Name:      "compare",
					Usage:     "Compare permissions of 2 IAM roles",
					ArgsUsage: "<role1> <role2>",
					Description: "Compare the permissions between two IAM roles, showing common permissions and differences.\n\n" +
						"Examples:\n" +
						"  gcp-iam role compare viewer editor\n" +
						"  gcp-iam role compare storage.admin storage.objectAdmin\n" +
						"  gcp-iam role compare roles/compute.admin compute.instanceAdmin",
					ShellComplete: func(ctx context.Context, cmd *cli.Command) {
						completeRoleNames(cmd)
					},
					Action: func(ctx context.Context, cmd *cli.Command) error {
						args := cmd.Args().Slice()
						if len(args) < 2 {
							return cli.ShowSubcommandHelp(cmd)
						}

						role1Name := normalizeRoleName(args[0])
						role2Name := normalizeRoleName(args[1])

						cfg, err := config.Load()
						if err != nil {
							return fmt.Errorf("failed to load config: %w", err)
						}

						database, err := db.New(cfg.DatabasePath)
						if err != nil {
							return fmt.Errorf("failed to open database: %w", err)
						}
						defer database.Close()

						// Get both roles
						role1, err := database.GetRoleByName(role1Name)
						if err != nil {
							return fmt.Errorf("failed to get role '%s': %w", role1Name, err)
						}
						if role1 == nil {
							return fmt.Errorf("role '%s' not found", role1Name)
						}

						role2, err := database.GetRoleByName(role2Name)
						if err != nil {
							return fmt.Errorf("failed to get role '%s': %w", role2Name, err)
						}
						if role2 == nil {
							return fmt.Errorf("role '%s' not found", role2Name)
						}

						// Get permissions for both roles
						perms1, err := database.GetRolePermissions(role1.Name)
						if err != nil {
							return fmt.Errorf("failed to get permissions for role '%s': %w", role1.Name, err)
						}

						perms2, err := database.GetRolePermissions(role2.Name)
						if err != nil {
							return fmt.Errorf("failed to get permissions for role '%s': %w", role2.Name, err)
						}

						// Create maps for easier comparison
						perms1Map := make(map[string]bool)
						for _, perm := range perms1 {
							perms1Map[perm.Permission] = true
						}

						perms2Map := make(map[string]bool)
						for _, perm := range perms2 {
							perms2Map[perm.Permission] = true
						}

						// Find permissions unique to each role and common permissions
						var onlyInRole1, onlyInRole2, common []string

						for perm := range perms1Map {
							if perms2Map[perm] {
								common = append(common, perm)
							} else {
								onlyInRole1 = append(onlyInRole1, perm)
							}
						}

						for perm := range perms2Map {
							if !perms1Map[perm] {
								onlyInRole2 = append(onlyInRole2, perm)
							}
						}

						// Display comparison results
						fmt.Printf("Comparing roles:\n")
						fmt.Printf("  Role 1: %s (%s)\n", role1.Name, role1.Title)
						fmt.Printf("  Role 2: %s (%s)\n\n", role2.Name, role2.Title)

						fmt.Printf("Common permissions (%d):\n", len(common))
						for _, perm := range common {
							fmt.Printf("  ✓ %s\n", perm)
						}

						fmt.Printf("\nPermissions only in '%s' (%d):\n", role1.Name, len(onlyInRole1))
						for _, perm := range onlyInRole1 {
							fmt.Printf("  - %s\n", perm)
						}

						fmt.Printf("\nPermissions only in '%s' (%d):\n", role2.Name, len(onlyInRole2))
						for _, perm := range onlyInRole2 {
							fmt.Printf("  + %s\n", perm)
						}

						fmt.Printf("\nSummary:\n")
						fmt.Printf("  Total permissions in '%s': %d\n", role1.Name, len(perms1))
						fmt.Printf("  Total permissions in '%s': %d\n", role2.Name, len(perms2))
						fmt.Printf("  Common permissions: %d\n", len(common))
						fmt.Printf("  Unique to '%s': %d\n", role1.Name, len(onlyInRole1))
						fmt.Printf("  Unique to '%s': %d\n", role2.Name, len(onlyInRole2))

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
					Name:      "show",
					Usage:     "Show IAM roles with permission",
					ArgsUsage: "<permission-name>",
					Description: "Display all IAM roles that include a specific permission.\n\n" +
						"Examples:\n" +
						"  gcp-iam permission show storage.objects.get\n" +
						"  gcp-iam permission show compute.instances.create\n" +
						"  gcp-iam permission show iam.serviceAccounts.actAs",
					ShellComplete: func(ctx context.Context, cmd *cli.Command) {
						completePermissionNames(cmd)
					},
					Action: func(ctx context.Context, cmd *cli.Command) error {
						permissionName := cmd.Args().First()
						if permissionName == "" {
							return cli.ShowSubcommandHelp(cmd)
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
					Name:      "search",
					Usage:     "Search IAM permissions",
					ArgsUsage: "<search-query>",
					Description: "Search for IAM permissions by name using a query string.\n\n" +
						"Examples:\n" +
						"  gcp-iam permission search storage\n" +
						"  gcp-iam permission search create\n" +
						"  gcp-iam permission search compute.instances",
					ShellComplete: func(ctx context.Context, cmd *cli.Command) {
						completePermissionNames(cmd)
					},
					Action: func(ctx context.Context, cmd *cli.Command) error {
						query := cmd.Args().First()
						if query == "" {
							return cli.ShowSubcommandHelp(cmd)
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
			Name:  "service",
			Usage: "Query Google Cloud Services",
			CommandNotFound: func(ctx context.Context, cmd *cli.Command, command string) {
				cli.ShowAppHelp(cmd)
			},
			Commands: []*cli.Command{
				{
					Name:      "show",
					Usage:     "Show service details",
					ArgsUsage: "<service-name>",
					Description: "Display detailed information about a specific Google Cloud service.\n\n" +
						"Examples:\n" +
						"  gcp-iam service show storage.googleapis.com\n" +
						"  gcp-iam service show compute.googleapis.com",
					ShellComplete: func(ctx context.Context, cmd *cli.Command) {
						completeServiceNames(cmd)
					},
					Action: func(ctx context.Context, cmd *cli.Command) error {
						serviceName := cmd.Args().First()
						if serviceName == "" {
							return cli.ShowSubcommandHelp(cmd)
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

						service, err := database.GetServiceByName(serviceName)
						if err != nil {
							return fmt.Errorf("failed to get service: %w", err)
						}

						if service == nil {
							fmt.Printf("Service '%s' not found\n", serviceName)
							return nil
						}

						fmt.Printf("Service: %s\n", service.Name)
						fmt.Printf("Title: %s\n", service.Title)
						if service.Description != "" && service.Description != service.Title {
							fmt.Printf("Description: %s\n", service.Description)
						}
						fmt.Printf("Created: %s\n", service.CreatedAt.Format("2006-01-02 15:04:05"))
						fmt.Printf("Updated: %s\n", service.UpdatedAt.Format("2006-01-02 15:04:05"))

						return nil
					},
				},
				{
					Name:      "search",
					Usage:     "Search for services",
					ArgsUsage: "<query>",
					Description: "Search for Google Cloud services by name or title.\n\n" +
						"The search is case-insensitive and matches partial strings in:\n" +
						"  • Service name (e.g., 'storage.googleapis.com')\n" +
						"  • Service title (e.g., 'Cloud Storage')\n\n" +
						"Examples:\n" +
						"  gcp-iam service search storage\n" +
						"  gcp-iam service search compute\n" +
						"  gcp-iam service search 'cloud sql'",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						query := strings.Join(cmd.Args().Slice(), " ")
						if query == "" {
							return cli.ShowSubcommandHelp(cmd)
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

						services, err := database.SearchServices(query)
						if err != nil {
							return fmt.Errorf("failed to search services: %w", err)
						}

						if len(services) == 0 {
							fmt.Printf("No services found matching '%s'\n", query)
							return nil
						}

						fmt.Printf("Found %d services matching '%s':\n\n", len(services), query)
						for _, service := range services {
							fmt.Printf("%-50s %s\n", service.Name, service.Title)
						}

						return nil
					},
				},
			},
		},
		{
			Name:  "update",
			Usage: "Update IAM roles, permissions, and services",
			Description: "Fetch the latest data from Google Cloud Platform and update the local database.\n\n" +
				"By default, updates both roles and permissions. Use flags to update specific resources:\n" +
				"  --roles    Only update IAM roles and permissions\n" +
				"  --services Only update Google Cloud services\n\n" +
				"Examples:\n" +
				"  gcp-iam update           # Update both roles and services\n" +
				"  gcp-iam update --roles   # Update only roles and permissions\n" +
				"  gcp-iam update --services # Update only services",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "roles",
					Usage: "Only update IAM roles and permissions",
				},
				&cli.BoolFlag{
					Name:  "services",
					Usage: "Only update Google Cloud services",
				},
			},
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

				// Determine what to update based on flags
				updateRoles := cmd.Bool("roles")
				updateServices := cmd.Bool("services")

				// If no flags specified, update both
				if !updateRoles && !updateServices {
					updateRoles = true
					updateServices = true
				}

				// Update roles and permissions if requested
				if updateRoles {
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
				}

				// Update services if requested
				if updateServices {
					err = updater.UpdateServices(ctx)
					if err != nil {
						return fmt.Errorf("failed to update services: %w", err)
					}
				}

				fmt.Println("Update completed successfully")
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
			Name:   "complete-permissions",
			Usage:  "List all permission names for shell completion",
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

				permissionNames, err := database.GetPermissionNames()
				if err != nil {
					return err
				}

				for _, name := range permissionNames {
					fmt.Println(name)
				}
				return nil
			},
		},
		{
			Name:   "complete-services",
			Usage:  "List all service names for shell completion",
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

				serviceNames, err := database.GetServiceNames()
				if err != nil {
					return err
				}

				for _, name := range serviceNames {
					fmt.Println(name)
				}
				return nil
			},
		},
		{
			Name:  "info",
			Usage: "Show application configuration",
			Description: "Display current application configuration including database statistics and file paths.\n\n" +
				"Shows:\n" +
				"  • Number of roles, permissions, and services in database\n" +
				"  • Configuration file location\n" +
				"  • Database file location\n\n" +
				"Examples:\n" +
				"  gcp-iam info",
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

				serviceCount, err := database.CountServices()
				if err != nil {
					return fmt.Errorf("failed to count services: %w", err)
				}

				configPath, _ := config.GetDefaultConfigPath()
				fmt.Println("GCP IAM Configuration:")
				fmt.Printf("  Roles:        %d\n", roleCount)
				fmt.Printf("  Permissions:  %d\n", permissionCount)
				fmt.Printf("  Services:     %d\n", serviceCount)
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
