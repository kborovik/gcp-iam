package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kborovik/gcp-iam/internal/constants"
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
		permission TEXT,
		role TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (permission, role),
		FOREIGN KEY (role) REFERENCES roles(name) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_permissions_role ON permissions(role);
	CREATE INDEX IF NOT EXISTS idx_permissions_permission ON permissions(permission);

	CREATE TABLE IF NOT EXISTS services (
		name TEXT PRIMARY KEY,
		title TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_services_name ON services(name);
	CREATE INDEX IF NOT EXISTS idx_services_title ON services(title);
	`

	_, err := db.conn.Exec(schema)
	return err
}

func ensureDir(filePath string) error {
	dir := filepath.Dir(filePath)
	return os.MkdirAll(dir, constants.DefaultDirPermissions)
}
