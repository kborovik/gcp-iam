package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type DB struct {
	conn *sql.DB
}

func New(dbPath string) (*DB, error) {
	if err := ensureDir(dbPath); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db := &DB{conn: conn}
	if err := db.createTables(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) createTables() error {
	schema := `
	CREATE TABLE IF NOT EXISTS roles (
		name TEXT PRIMARY KEY,
		title TEXT,
		description TEXT,
		stage TEXT,
		deleted BOOLEAN DEFAULT FALSE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS permissions (
		name TEXT PRIMARY KEY,
		title TEXT,
		description TEXT,
		stage TEXT,
		api_disabled BOOLEAN DEFAULT FALSE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS role_permissions (
		role_name TEXT,
		permission_name TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (role_name, permission_name),
		FOREIGN KEY (role_name) REFERENCES roles(name) ON DELETE CASCADE,
		FOREIGN KEY (permission_name) REFERENCES permissions(name) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_role_permissions_role ON role_permissions(role_name);
	CREATE INDEX IF NOT EXISTS idx_role_permissions_permission ON role_permissions(permission_name);
	`

	_, err := db.conn.Exec(schema)
	return err
}

func ensureDir(filePath string) error {
	dir := filepath.Dir(filePath)
	return os.MkdirAll(dir, 0755)
}
