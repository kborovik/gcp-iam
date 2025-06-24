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
		Permission: "compute.instances.get",
		Role:       "roles/compute.admin",
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

	if retrieved.Permission != perm.Permission {
		t.Errorf("Expected permission name %s, got %s", perm.Permission, retrieved.Permission)
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
		Permission: "compute.instances.get",
		Role:       role.Name,
	}

	err = db.InsertRole(role)
	if err != nil {
		t.Fatalf("Failed to insert role: %v", err)
	}

	err = db.InsertPermission(perm)
	if err != nil {
		t.Fatalf("Failed to insert permission: %v", err)
	}

	permissions, err := db.GetRolePermissions(role.Name)
	if err != nil {
		t.Fatalf("Failed to get role permissions: %v", err)
	}

	if len(permissions) != 1 {
		t.Fatalf("Expected 1 permission, got %d", len(permissions))
	}

	if permissions[0].Permission != perm.Permission {
		t.Errorf("Expected permission name %s, got %s", perm.Permission, permissions[0].Permission)
	}
}

func TestGetAllRoles(t *testing.T) {
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

	results, err := db.GetAllRoles()
	if err != nil {
		t.Fatalf("Failed to get all roles: %v", err)
	}

	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	// Verify roles are sorted by name
	expectedOrder := []string{"roles/compute.admin", "roles/compute.viewer", "roles/storage.viewer"}
	for i, role := range results {
		if role.Name != expectedOrder[i] {
			t.Errorf("Expected role at index %d to be %s, got %s", i, expectedOrder[i], role.Name)
		}
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

func TestGetRolesWithPermission(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Insert test roles
	roles := []*Role{
		{Name: "roles/compute.admin", Title: "Compute Admin", Description: "Full compute access"},
		{Name: "roles/storage.admin", Title: "Storage Admin", Description: "Full storage access"},
	}

	for _, role := range roles {
		err = db.InsertRole(role)
		if err != nil {
			t.Fatalf("Failed to insert role: %v", err)
		}
	}

	// Insert test permission
	perm := &Permission{
		Permission: "compute.instances.get",
		Role:       "roles/compute.admin",
	}

	err = db.InsertPermission(perm)
	if err != nil {
		t.Fatalf("Failed to insert permission: %v", err)
	}

	// Test getting roles with permission
	rolesWithPermission, err := db.GetRolesWithPermission("compute.instances.get")
	if err != nil {
		t.Fatalf("Failed to get roles with permission: %v", err)
	}

	if len(rolesWithPermission) != 1 {
		t.Errorf("Expected 1 role with permission, got %d", len(rolesWithPermission))
	}

	if rolesWithPermission[0].Name != "roles/compute.admin" {
		t.Errorf("Expected role name 'roles/compute.admin', got %s", rolesWithPermission[0].Name)
	}
}

func TestHasPermissions(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Insert test role
	role := &Role{
		Name:        "roles/compute.admin",
		Title:       "Compute Admin",
		Description: "Full compute access",
	}

	err = db.InsertRole(role)
	if err != nil {
		t.Fatalf("Failed to insert role: %v", err)
	}

	// Test role without permissions
	hasPerms, err := db.HasPermissions("roles/compute.admin")
	if err != nil {
		t.Fatalf("Failed to check permissions: %v", err)
	}

	if hasPerms {
		t.Error("Expected role to have no permissions initially")
	}

	// Add a permission
	perm := &Permission{
		Permission: "compute.instances.get",
		Role:       "roles/compute.admin",
	}

	err = db.InsertPermission(perm)
	if err != nil {
		t.Fatalf("Failed to insert permission: %v", err)
	}

	// Test role with permissions
	hasPerms, err = db.HasPermissions("roles/compute.admin")
	if err != nil {
		t.Fatalf("Failed to check permissions: %v", err)
	}

	if !hasPerms {
		t.Error("Expected role to have permissions after linking")
	}
}

func TestGetRolesNeedingPermissionUpdate(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Insert roles - one with permissions, one without
	roleWithPerms := &Role{
		Name:        "roles/compute.admin",
		Title:       "Compute Admin",
		Description: "Full compute access",
	}

	roleWithoutPerms := &Role{
		Name:        "roles/storage.admin",
		Title:       "Storage Admin",
		Description: "Full storage access",
	}

	err = db.InsertRole(roleWithPerms)
	if err != nil {
		t.Fatalf("Failed to insert role: %v", err)
	}

	err = db.InsertRole(roleWithoutPerms)
	if err != nil {
		t.Fatalf("Failed to insert role: %v", err)
	}

	// Add permission to first role
	perm := &Permission{
		Permission: "compute.instances.get",
		Role:       "roles/compute.admin",
	}

	err = db.InsertPermission(perm)
	if err != nil {
		t.Fatalf("Failed to insert permission: %v", err)
	}

	// Get roles needing updates
	rolesToUpdate, err := db.GetRolesNeedingPermissionUpdate()
	if err != nil {
		t.Fatalf("Failed to get roles needing updates: %v", err)
	}

	// Should return at least the role without permissions and possibly recently updated roles
	if len(rolesToUpdate) == 0 {
		t.Error("Expected at least one role needing updates")
	}

	// Check that the role without permissions is included
	foundRoleWithoutPerms := false
	for _, role := range rolesToUpdate {
		if role.Name == "roles/storage.admin" {
			foundRoleWithoutPerms = true
			break
		}
	}

	if !foundRoleWithoutPerms {
		t.Error("Expected role without permissions to be in roles needing updates")
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
		{Permission: "compute.instances.get", Role: "roles/compute.admin"},
		{Permission: "compute.instances.list", Role: "roles/compute.admin"},
		{Permission: "storage.buckets.get", Role: "roles/storage.admin"},
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
