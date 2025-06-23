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
	Name        string    `json:"name"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Stage       string    `json:"stage"`
	APIDisabled bool      `json:"api_disabled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RolePermission struct {
	RoleName       string    `json:"role_name"`
	PermissionName string    `json:"permission_name"`
	CreatedAt      time.Time `json:"created_at"`
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
		INSERT INTO permissions (name, title, description, stage, api_disabled)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(name) DO UPDATE SET
			title = excluded.title,
			description = excluded.description,
			stage = excluded.stage,
			api_disabled = excluded.api_disabled,
			updated_at = CURRENT_TIMESTAMP
	`
	_, err := db.conn.Exec(query, perm.Name, perm.Title, perm.Description, perm.Stage, perm.APIDisabled)
	return err
}

func (db *DB) LinkRolePermission(roleName, permissionName string) error {
	query := `
		INSERT OR IGNORE INTO role_permissions (role_name, permission_name)
		VALUES (?, ?)
	`
	_, err := db.conn.Exec(query, roleName, permissionName)
	return err
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
		SELECT name, title, description, stage, api_disabled, created_at, updated_at
		FROM permissions
		WHERE name = ? AND api_disabled = FALSE
	`
	row := db.conn.QueryRow(query, name)

	var perm Permission
	err := row.Scan(&perm.Name, &perm.Title, &perm.Description, &perm.Stage, &perm.APIDisabled, &perm.CreatedAt, &perm.UpdatedAt)
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
		SELECT p.name, p.title, p.description, p.stage, p.api_disabled, p.created_at, p.updated_at
		FROM permissions p
		JOIN role_permissions rp ON p.name = rp.permission_name
		WHERE rp.role_name = ? AND p.api_disabled = FALSE
		ORDER BY p.name
	`
	rows, err := db.conn.Query(query, roleName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []Permission
	for rows.Next() {
		var perm Permission
		err := rows.Scan(&perm.Name, &perm.Title, &perm.Description, &perm.Stage, &perm.APIDisabled, &perm.CreatedAt, &perm.UpdatedAt)
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
		SELECT name, title, description, stage, api_disabled, created_at, updated_at
		FROM permissions
		WHERE (name LIKE ? OR title LIKE ? OR description LIKE ?) AND api_disabled = FALSE
		ORDER BY name
	`
	pattern := "%" + query + "%"
	rows, err := db.conn.Query(sqlQuery, pattern, pattern, pattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []Permission
	for rows.Next() {
		var perm Permission
		err := rows.Scan(&perm.Name, &perm.Title, &perm.Description, &perm.Stage, &perm.APIDisabled, &perm.CreatedAt, &perm.UpdatedAt)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}

	return permissions, rows.Err()
}
