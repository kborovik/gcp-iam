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

func New(database *db.DB) *Updater {
	return &Updater{
		db: database,
	}
}

func (u *Updater) UpdateRolesAndPermissions(ctx context.Context) error {
	fmt.Println("Updating GCP IAM pre-defined roles and permissions...")

	roles, err := u.fetchGCPIAMRoles(ctx)
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

func (u *Updater) fetchGCPIAMRoles(ctx context.Context) ([]db.Role, error) {
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
