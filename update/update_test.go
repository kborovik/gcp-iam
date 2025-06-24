package update

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kborovik/gcp-iam/db"
)

func TestNewUpdater(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	database, err := db.New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer database.Close()

	updater := New(database)
	if updater == nil {
		t.Fatal("Expected updater to be created")
	}

	if updater.db != database {
		t.Fatal("Expected updater to contain the database reference")
	}
}

func TestFetchGCPIAMRoles(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	database, err := db.New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer database.Close()

	updater := New(database)

	// This test requires valid GCP credentials and network access
	// Skip in CI/CD environments or when credentials are not available
	t.Skip("Skipping integration test - requires GCP credentials")

	ctx := context.Background()
	roles, err := updater.fetchRoles(ctx)
	if err != nil {
		t.Fatalf("Failed to fetch GCP IAM roles: %v", err)
	}

	if len(roles) == 0 {
		t.Fatal("Expected at least one role, got none")
	}

	// Verify first role has required fields
	if roles[0].Name == "" {
		t.Error("Expected role name to be non-empty")
	}
}

func TestUpdateDatabase(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	database, err := db.New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer database.Close()

	updater := New(database)

	testRoles := []db.Role{
		{
			Name:        "roles/test.role",
			Title:       "Test Role",
			Description: "A test role",
			Stage:       "GA",
		},
	}

	err = updater.updateDatabase(testRoles)
	if err != nil {
		t.Fatalf("Failed to update database: %v", err)
	}

	role, err := database.GetRoleByName("roles/test.role")
	if err != nil {
		t.Fatalf("Failed to get role: %v", err)
	}

	if role == nil {
		t.Fatal("Expected role to exist in database")
	}

	if role.Title != "Test Role" {
		t.Errorf("Expected role title to be 'Test Role', got '%s'", role.Title)
	}
}

func TestFetchRolePermissions(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	database, err := db.New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer database.Close()

	updater := New(database)

	// This test requires valid GCP credentials and network access
	// Skip in CI/CD environments or when credentials are not available
	t.Skip("Skipping integration test - requires GCP credentials")

	ctx := context.Background()
	permissions, err := updater.fetchPermissions(ctx, "roles/storage.admin")
	if err != nil {
		t.Fatalf("Failed to fetch role permissions: %v", err)
	}

	if len(permissions) == 0 {
		t.Fatal("Expected at least one permission for storage.admin role")
	}

	// Check that we get expected storage permissions
	hasStoragePermission := false
	for _, perm := range permissions {
		if strings.Contains(perm, "storage") {
			hasStoragePermission = true
			break
		}
	}

	if !hasStoragePermission {
		t.Error("Expected storage.admin role to have storage-related permissions")
	}
}

func TestUpdatePermissions(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	database, err := db.New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer database.Close()

	updater := New(database)

	// This test requires valid GCP credentials and network access
	// Skip in CI/CD environments or when credentials are not available
	t.Skip("Skipping integration test - requires GCP credentials")

	// First, insert a test role
	testRole := &db.Role{
		Name:        "roles/storage.admin",
		Title:       "Storage Admin",
		Description: "Full control of buckets and objects",
		Stage:       "GA",
	}

	err = database.InsertRole(testRole)
	if err != nil {
		t.Fatalf("Failed to insert test role: %v", err)
	}

	ctx := context.Background()
	err = updater.UpdatePermissions(ctx, "roles/storage.admin")
	if err != nil {
		t.Fatalf("Failed to update role permissions: %v", err)
	}

	// Verify permissions were stored in database
	permissions, err := database.GetRolePermissions("roles/storage.admin")
	if err != nil {
		t.Fatalf("Failed to get role permissions from database: %v", err)
	}

	if len(permissions) == 0 {
		t.Fatal("Expected role permissions to be stored in database")
	}
}
