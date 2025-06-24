package db

import (
	"database/sql"
	"time"
)

type Role struct {
	Name        string    `json:"name"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Stage       string    `json:"stage"`
	Deleted     bool      `json:"deleted"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Permission struct {
	Permission string    `json:"permission"`
	Role       string    `json:"role"`
	CreatedAt  time.Time `json:"created_at"`
}

func (db *DB) InsertRole(role *Role) error {
	query := `
		INSERT INTO roles (name, title, description, stage, deleted)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(name) DO UPDATE SET
			title = excluded.title,
			description = excluded.description,
			stage = excluded.stage,
			deleted = excluded.deleted,
			updated_at = CURRENT_TIMESTAMP
	`
	_, err := db.conn.Exec(query, role.Name, role.Title, role.Description, role.Stage, role.Deleted)
	return err
}

func (db *DB) InsertPermission(perm *Permission) error {
	query := `
		INSERT OR IGNORE INTO permissions (permission, role)
		VALUES (?, ?)
	`
	_, err := db.conn.Exec(query, perm.Permission, perm.Role)
	return err
}

func (db *DB) LinkRolePermission(roleName, permissionName string) error {
	perm := &Permission{
		Permission: permissionName,
		Role:       roleName,
	}
	return db.InsertPermission(perm)
}

func (db *DB) GetRoleByName(name string) (*Role, error) {
	query := `
		SELECT name, title, description, stage, deleted, created_at, updated_at
		FROM roles
		WHERE name = ? AND deleted = FALSE
	`
	row := db.conn.QueryRow(query, name)

	var role Role
	err := row.Scan(&role.Name, &role.Title, &role.Description, &role.Stage, &role.Deleted, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &role, nil
}

func (db *DB) GetPermissionByName(name string) (*Permission, error) {
	query := `
		SELECT permission, role, created_at
		FROM permissions
		WHERE permission = ?
		LIMIT 1
	`
	row := db.conn.QueryRow(query, name)

	var perm Permission
	err := row.Scan(&perm.Permission, &perm.Role, &perm.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &perm, nil
}

func (db *DB) GetRolePermissions(roleName string) ([]Permission, error) {
	query := `
		SELECT permission, role, created_at
		FROM permissions
		WHERE role = ?
		ORDER BY permission
	`
	rows, err := db.conn.Query(query, roleName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []Permission
	for rows.Next() {
		var perm Permission
		err := rows.Scan(&perm.Permission, &perm.Role, &perm.CreatedAt)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}

	return permissions, rows.Err()
}

func (db *DB) GetAllRoles() ([]Role, error) {
	sqlQuery := `
		SELECT name, title, description, stage, deleted, created_at, updated_at
		FROM roles
		WHERE deleted = FALSE
		ORDER BY name
	`
	rows, err := db.conn.Query(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []Role
	for rows.Next() {
		var role Role
		err := rows.Scan(&role.Name, &role.Title, &role.Description, &role.Stage, &role.Deleted, &role.CreatedAt, &role.UpdatedAt)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, rows.Err()
}

func (db *DB) SearchRoles(query string) ([]Role, error) {
	sqlQuery := `
		SELECT name, title, description, stage, deleted, created_at, updated_at
		FROM roles
		WHERE (name LIKE ? OR title LIKE ? OR description LIKE ?) AND deleted = FALSE
		ORDER BY name
	`
	pattern := "%" + query + "%"
	rows, err := db.conn.Query(sqlQuery, pattern, pattern, pattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []Role
	for rows.Next() {
		var role Role
		err := rows.Scan(&role.Name, &role.Title, &role.Description, &role.Stage, &role.Deleted, &role.CreatedAt, &role.UpdatedAt)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, rows.Err()
}

func (db *DB) SearchPermissions(query string) ([]Permission, error) {
	sqlQuery := `
		SELECT DISTINCT permission
		FROM permissions
		WHERE permission LIKE ?
		ORDER BY permission
	`
	pattern := "%" + query + "%"
	rows, err := db.conn.Query(sqlQuery, pattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []Permission
	for rows.Next() {
		var perm Permission
		err := rows.Scan(&perm.Permission)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}

	return permissions, rows.Err()
}

func (db *DB) GetRolesWithPermission(permissionName string) ([]Role, error) {
	query := `
		SELECT r.name, r.title, r.description, r.stage, r.deleted, r.created_at, r.updated_at
		FROM roles r
		JOIN permissions p ON r.name = p.role
		WHERE p.permission = ? AND r.deleted = FALSE
		ORDER BY r.name
	`
	rows, err := db.conn.Query(query, permissionName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []Role
	for rows.Next() {
		var role Role
		err := rows.Scan(&role.Name, &role.Title, &role.Description, &role.Stage, &role.Deleted, &role.CreatedAt, &role.UpdatedAt)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, rows.Err()
}

func (db *DB) GetRolesNeedingPermissionUpdate() ([]Role, error) {
	query := `
		SELECT name, title, description, stage, deleted, created_at, updated_at
		FROM roles
		WHERE deleted = FALSE
		AND NOT EXISTS (
			SELECT 1 FROM permissions WHERE role = roles.name
		)
		ORDER BY name
	`
	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []Role
	for rows.Next() {
		var role Role
		err := rows.Scan(&role.Name, &role.Title, &role.Description, &role.Stage, &role.Deleted, &role.CreatedAt, &role.UpdatedAt)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, rows.Err()
}

func (db *DB) HasPermissions(roleName string) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM permissions
		WHERE role = ?
	`
	var count int
	err := db.conn.QueryRow(query, roleName).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (db *DB) CountRoles() (int, error) {
	query := `SELECT COUNT(*) FROM roles WHERE deleted = FALSE`
	var count int
	err := db.conn.QueryRow(query).Scan(&count)
	return count, err
}

func (db *DB) CountPermissions() (int, error) {
	query := `SELECT COUNT(DISTINCT permission) FROM permissions`
	var count int
	err := db.conn.QueryRow(query).Scan(&count)
	return count, err
}
