package update

import (
	"context"
	"fmt"
	"log"

	"github.com/kborovik/gcp-iam/db"
	"google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
)

type Updater struct {
	db *db.DB
}

// =============================================================================
// PUBLIC API - Constructor and Main Functions
// =============================================================================

// New creates a new Updater instance with the provided database connection
func New(database *db.DB) *Updater {
	return &Updater{
		db: database,
	}
}

// UpdateRoles fetches all IAM roles from Google Cloud and stores them in the database
func (u *Updater) UpdateRoles(ctx context.Context) error {
	fmt.Println("Updating GCP IAM pre-defined roles and permissions...")

	roles, err := u.fetchRoles(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch GCP IAM roles: %w", err)
	}

	fmt.Printf("Fetched %d roles from GCP\n", len(roles))

	err = u.updateDatabase(roles)
	if err != nil {
		return fmt.Errorf("failed to update database: %w", err)
	}

	fmt.Println("Successfully updated IAM roles and permissions")
	return nil
}

// UpdatePermissions fetches and stores permissions for a specific role in the database
func (u *Updater) UpdatePermissions(ctx context.Context, roleName string) error {
	permissions, err := u.fetchPermissions(ctx, roleName)
	if err != nil {
		return fmt.Errorf("failed to fetch permissions for role %s: %w", roleName, err)
	}

	// Insert permissions into database and link them to the role
	for _, permName := range permissions {
		// Create permission record
		perm := &db.Permission{
			Name:        permName,
			Title:       permName, // API doesn't provide separate title for permissions
			Description: "",       // API doesn't provide description for individual permissions
			Stage:       "GA",     // Default stage
		}

		err = u.db.InsertPermission(perm)
		if err != nil {
			log.Printf("Warning: failed to insert permission %s: %v", permName, err)
			// Continue even if permission insert fails
		}

		// Link permission to role
		err = u.db.LinkRolePermission(roleName, permName)
		if err != nil {
			log.Printf("Warning: failed to link permission %s to role %s: %v", permName, roleName, err)
		}
	}

	log.Printf("Updated %d permissions for role %s", len(permissions), roleName)
	return nil
}

// =============================================================================
// PRIVATE IMPLEMENTATION - Helper Functions
// =============================================================================

// fetchRoles fetches all IAM roles from Google Cloud API
func (u *Updater) fetchRoles(ctx context.Context) ([]db.Role, error) {
	service, err := iam.NewService(ctx, option.WithScopes(iam.CloudPlatformScope))
	if err != nil {
		return nil, fmt.Errorf("failed to create IAM service: %w", err)
	}

	var roles []db.Role

	// List predefined roles
	call := service.Roles.List().ShowDeleted(false).View("FULL").PageSize(1000)

	err = call.Pages(ctx, func(page *iam.ListRolesResponse) error {
		for _, role := range page.Roles {
			dbRole := db.Role{
				Name:        role.Name,
				Title:       role.Title,
				Description: role.Description,
				Stage:       role.Stage,
			}
			roles = append(roles, dbRole)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}

	return roles, nil
}

// fetchPermissions fetches the detailed permissions for a specific role
func (u *Updater) fetchPermissions(ctx context.Context, roleName string) ([]string, error) {
	service, err := iam.NewService(ctx, option.WithScopes(iam.CloudPlatformScope))
	if err != nil {
		return nil, fmt.Errorf("failed to create IAM service: %w", err)
	}

	role, err := service.Roles.Get(roleName).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get role %s: %w", roleName, err)
	}

	return role.IncludedPermissions, nil
}

// updateDatabase stores roles in the database
func (u *Updater) updateDatabase(roles []db.Role) error {
	for _, role := range roles {
		err := u.db.InsertRole(&role)
		if err != nil {
			log.Printf("Warning: failed to create role %s: %v", role.Name, err)
			continue
		}
	}

	return nil
}
