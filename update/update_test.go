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

	if roles[0].Title == "" {
		t.Error("Expected role title to be non-empty")
	}

	// Check that we have common GCP roles
	hasCommonRole := false
	for _, role := range roles {
		if role.Name == "viewer" || role.Name == "editor" || role.Name == "owner" {
			hasCommonRole = true
			break
		}
	}

	if !hasCommonRole {
		t.Error("Expected to find at least one common GCP role (viewer, editor, or owner)")
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

	ctx := context.Background()
	permissions, err := updater.fetchPermissions(ctx, "storage.admin")
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

	// First, insert a test role
	testRole := &db.Role{
		Name:        "storage.admin",
		Title:       "Storage Admin",
		Description: "Full control of buckets and objects",
		Stage:       "GA",
	}

	err = database.InsertRole(testRole)
	if err != nil {
		t.Fatalf("Failed to insert test role: %v", err)
	}

	ctx := context.Background()
	err = updater.UpdatePermissions(ctx, "storage.admin")
	if err != nil {
		t.Fatalf("Failed to update role permissions: %v", err)
	}

	// Verify permissions were stored in database
	permissions, err := database.GetRolePermissions("storage.admin")
	if err != nil {
		t.Fatalf("Failed to get role permissions from database: %v", err)
	}

	if len(permissions) == 0 {
		t.Fatal("Expected role permissions to be stored in database")
	}

	// Check that we have storage-related permissions
	hasStoragePermission := false
	for _, perm := range permissions {
		if strings.Contains(perm.Permission, "storage") {
			hasStoragePermission = true
			break
		}
	}

	if !hasStoragePermission {
		t.Error("Expected storage.admin role to have storage-related permissions in database")
	}
}

func TestFetchServices(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	database, err := db.New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer database.Close()

	updater := New(database)

	ctx := context.Background()
	services, err := updater.fetchServices(ctx)
	if err != nil {
		t.Fatalf("Failed to fetch GCP services: %v", err)
	}

	if len(services) == 0 {
		t.Fatal("Expected at least one service, got none")
	}

	// Verify first service has required fields
	if services[0].Name == "" {
		t.Error("Expected service name to be non-empty")
	}

	if services[0].Title == "" {
		t.Error("Expected service title to be non-empty")
	}

	// Check that we get googleapis.com services
	hasGoogleAPIsService := false
	for _, service := range services {
		if strings.Contains(service.Name, "googleapis.com") {
			hasGoogleAPIsService = true
			break
		}
	}

	if !hasGoogleAPIsService {
		t.Error("Expected to find at least one googleapis.com service")
	}
}

func TestUpdateServices(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	database, err := db.New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer database.Close()

	updater := New(database)

	ctx := context.Background()
	err = updater.UpdateServices(ctx)
	if err != nil {
		t.Fatalf("Failed to update services: %v", err)
	}

	// Verify services were stored in database
	services, err := database.GetAllServices()
	if err != nil {
		t.Fatalf("Failed to get services from database: %v", err)
	}

	if len(services) == 0 {
		t.Fatal("Expected services to be stored in database")
	}

	// Verify first service has required fields
	if services[0].Name == "" {
		t.Error("Expected service name to be non-empty")
	}

	if services[0].Title == "" {
		t.Error("Expected service title to be non-empty")
	}

	// Check that we have common Google Cloud services
	hasCommonService := false
	commonServices := []string{
		"compute.googleapis.com",
		"storage.googleapis.com",
		"iam.googleapis.com",
		"cloudresourcemanager.googleapis.com",
	}

	for _, service := range services {
		for _, common := range commonServices {
			if service.Name == common {
				hasCommonService = true
				break
			}
		}
		if hasCommonService {
			break
		}
	}

	if !hasCommonService {
		t.Error("Expected to find at least one common Google Cloud service")
	}

	// Verify that all services contain googleapis.com
	for _, service := range services {
		if !strings.Contains(service.Name, "googleapis.com") {
			t.Errorf("Expected service name %s to contain googleapis.com", service.Name)
		}
	}
}
