package update

import (
	"context"
	"path/filepath"
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
	roles, err := updater.fetchGCPIAMRoles(ctx)
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
