// Package database provides SQLite-based persistence for user/room mappings,
// message IDs, and bridge state. Thread-safe and transaction-aware.
package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// DB wraps sql.DB with bridge-specific methods.
type DB struct {
	db *sql.DB
}

// Open opens or creates a SQLite database at path and runs migrations.
func Open(path string) (*DB, error) {
	db, err := sql.Open("sqlite3", path+"?_foreign_keys=1&_journal_mode=WAL")
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}
	d := &DB{db: db}
	if err := d.migrate(); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return d, nil
}

// Close closes the database connection.
func (d *DB) Close() error {
	return d.db.Close()
}

// migrate creates tables if they don't exist.
func (d *DB) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS viber_users (
		viber_id TEXT PRIMARY KEY,
		viber_name TEXT NOT NULL,
		matrix_user_id TEXT,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS room_mappings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		viber_chat_id TEXT UNIQUE NOT NULL,
		matrix_room_id TEXT UNIQUE NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS message_mappings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		viber_message_id TEXT UNIQUE NOT NULL,
		matrix_event_id TEXT UNIQUE NOT NULL,
		viber_chat_id TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (viber_chat_id) REFERENCES room_mappings(viber_chat_id)
	);

	CREATE INDEX IF NOT EXISTS idx_viber_users_matrix ON viber_users(matrix_user_id);
	CREATE INDEX IF NOT EXISTS idx_message_mappings_viber ON message_mappings(viber_message_id);
	CREATE INDEX IF NOT EXISTS idx_message_mappings_matrix ON message_mappings(matrix_event_id);
	`
	if _, err := d.db.Exec(schema); err != nil {
		return fmt.Errorf("create tables: %w", err)
	}
	return nil
}

// ViberUser represents a Viber user with optional Matrix mapping.
type ViberUser struct {
	ViberID     string
	ViberName   string
	MatrixUserID *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// UpsertViberUser creates or updates a Viber user record.
func (d *DB) UpsertViberUser(viberID, viberName string) error {
	_, err := d.db.Exec(
		`INSERT INTO viber_users (viber_id, viber_name, updated_at)
		 VALUES (?, ?, CURRENT_TIMESTAMP)
		 ON CONFLICT(viber_id) DO UPDATE SET
		 viber_name = excluded.viber_name,
		 updated_at = CURRENT_TIMESTAMP`,
		viberID, viberName,
	)
	return err
}

// LinkViberUser links a Viber user to a Matrix user ID.
func (d *DB) LinkViberUser(viberID, matrixUserID string) error {
	_, err := d.db.Exec(
		`UPDATE viber_users SET matrix_user_id = ?, updated_at = CURRENT_TIMESTAMP WHERE viber_id = ?`,
		matrixUserID, viberID,
	)
	return err
}

// GetViberUser returns a Viber user by ID.
func (d *DB) GetViberUser(viberID string) (*ViberUser, error) {
	var u ViberUser
	var mxID sql.NullString
	err := d.db.QueryRow(
		`SELECT viber_id, viber_name, matrix_user_id, created_at, updated_at
		 FROM viber_users WHERE viber_id = ?`,
		viberID,
	).Scan(&u.ViberID, &u.ViberName, &mxID, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if mxID.Valid {
		u.MatrixUserID = &mxID.String
	}
	return &u, nil
}

// CreateRoomMapping creates a mapping between a Viber chat and Matrix room.
func (d *DB) CreateRoomMapping(viberChatID, matrixRoomID string) error {
	_, err := d.db.Exec(
		`INSERT INTO room_mappings (viber_chat_id, matrix_room_id) VALUES (?, ?)`,
		viberChatID, matrixRoomID,
	)
	return err
}

// GetMatrixRoomID returns the Matrix room ID for a Viber chat, or empty if not found.
func (d *DB) GetMatrixRoomID(viberChatID string) (string, error) {
	var roomID string
	err := d.db.QueryRow(
		`SELECT matrix_room_id FROM room_mappings WHERE viber_chat_id = ?`,
		viberChatID,
	).Scan(&roomID)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return roomID, err
}

// GetViberChatID returns the Viber chat ID for a Matrix room, or empty if not found.
func (d *DB) GetViberChatID(matrixRoomID string) (string, error) {
	var chatID string
	err := d.db.QueryRow(
		`SELECT viber_chat_id FROM room_mappings WHERE matrix_room_id = ?`,
		matrixRoomID,
	).Scan(&chatID)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return chatID, err
}

// StoreMessageMapping stores a mapping between Viber and Matrix message IDs.
func (d *DB) StoreMessageMapping(viberMsgID, matrixEventID, viberChatID string) error {
	_, err := d.db.Exec(
		`INSERT INTO message_mappings (viber_message_id, matrix_event_id, viber_chat_id)
		 VALUES (?, ?, ?)`,
		viberMsgID, matrixEventID, viberChatID,
	)
	return err
}

// GetMatrixEventID returns the Matrix event ID for a Viber message, or empty if not found.
func (d *DB) GetMatrixEventID(viberMsgID string) (string, error) {
	var eventID string
	err := d.db.QueryRow(
		`SELECT matrix_event_id FROM message_mappings WHERE viber_message_id = ?`,
		viberMsgID,
	).Scan(&eventID)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return eventID, err
}

