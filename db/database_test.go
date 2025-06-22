package db

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDatabaseCreation(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Error("Database file was not created")
	}
}

func TestRoleOperations(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	role := &Role{
		Name:        "roles/test.role",
		Title:       "Test Role",
		Description: "A test role for testing",
		Stage:       "GA",
		Deleted:     false,
	}

	err = db.InsertRole(role)
	if err != nil {
		t.Fatalf("Failed to insert role: %v", err)
	}

	retrieved, err := db.GetRoleByName("roles/test.role")
	if err != nil {
		t.Fatalf("Failed to get role: %v", err)
	}

	if retrieved == nil {
		t.Fatal("Role not found")
	}

	if retrieved.Name != role.Name {
		t.Errorf("Expected role name %s, got %s", role.Name, retrieved.Name)
	}

	if retrieved.Title != role.Title {
		t.Errorf("Expected role title %s, got %s", role.Title, retrieved.Title)
	}
}

func TestPermissionOperations(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	perm := &Permission{
		Name:        "compute.instances.get",
		Title:       "Get Instance",
		Description: "Get compute instance details",
		Stage:       "GA",
		APIDisabled: false,
	}

	err = db.InsertPermission(perm)
	if err != nil {
		t.Fatalf("Failed to insert permission: %v", err)
	}

	retrieved, err := db.GetPermissionByName("compute.instances.get")
	if err != nil {
		t.Fatalf("Failed to get permission: %v", err)
	}

	if retrieved == nil {
		t.Fatal("Permission not found")
	}

	if retrieved.Name != perm.Name {
		t.Errorf("Expected permission name %s, got %s", perm.Name, retrieved.Name)
	}
}

func TestRolePermissionLink(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	role := &Role{
		Name:        "roles/test.role",
		Title:       "Test Role",
		Description: "A test role",
		Stage:       "GA",
	}

	perm := &Permission{
		Name:        "compute.instances.get",
		Title:       "Get Instance",
		Description: "Get compute instance details",
		Stage:       "GA",
	}

	err = db.InsertRole(role)
	if err != nil {
		t.Fatalf("Failed to insert role: %v", err)
	}

	err = db.InsertPermission(perm)
	if err != nil {
		t.Fatalf("Failed to insert permission: %v", err)
	}

	err = db.LinkRolePermission(role.Name, perm.Name)
	if err != nil {
		t.Fatalf("Failed to link role and permission: %v", err)
	}

	permissions, err := db.GetRolePermissions(role.Name)
	if err != nil {
		t.Fatalf("Failed to get role permissions: %v", err)
	}

	if len(permissions) != 1 {
		t.Fatalf("Expected 1 permission, got %d", len(permissions))
	}

	if permissions[0].Name != perm.Name {
		t.Errorf("Expected permission name %s, got %s", perm.Name, permissions[0].Name)
	}
}

func TestSearchRoles(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	roles := []*Role{
		{Name: "roles/compute.admin", Title: "Compute Admin", Description: "Full compute access"},
		{Name: "roles/storage.viewer", Title: "Storage Viewer", Description: "View storage resources"},
		{Name: "roles/compute.viewer", Title: "Compute Viewer", Description: "View compute resources"},
	}

	for _, role := range roles {
		err = db.InsertRole(role)
		if err != nil {
			t.Fatalf("Failed to insert role: %v", err)
		}
	}

	results, err := db.SearchRoles("compute")
	if err != nil {
		t.Fatalf("Failed to search roles: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}
}

func TestSearchPermissions(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	permissions := []*Permission{
		{Name: "compute.instances.get", Title: "Get Instance", Description: "Get compute instance"},
		{Name: "compute.instances.list", Title: "List Instances", Description: "List compute instances"},
		{Name: "storage.buckets.get", Title: "Get Bucket", Description: "Get storage bucket"},
	}

	for _, perm := range permissions {
		err = db.InsertPermission(perm)
		if err != nil {
			t.Fatalf("Failed to insert permission: %v", err)
		}
	}

	results, err := db.SearchPermissions("compute")
	if err != nil {
		t.Fatalf("Failed to search permissions: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}
}
