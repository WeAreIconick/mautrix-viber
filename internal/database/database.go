// Package database provides SQLite-based persistence for user/room mappings,
// message IDs, and bridge state. Thread-safe and transaction-aware.
package database

import (
	"context"
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
// Configures connection pool settings for production use:
// - SetMaxOpenConns: Maximum number of open connections (25 for SQLite)
// - SetMaxIdleConns: Maximum idle connections (5)
// - SetConnMaxLifetime: Maximum connection lifetime (5 minutes)
// - SetConnMaxIdleTime: Maximum idle time before closing (10 minutes)
func Open(path string) (*DB, error) {
	db, err := sql.Open("sqlite3", path+"?_foreign_keys=1&_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	
	// Configure connection pool for production
	db.SetMaxOpenConns(25)           // SQLite default is usually unlimited, but we limit for safety
	db.SetMaxIdleConns(5)            // Keep some connections ready
	db.SetConnMaxLifetime(5 * time.Minute)  // Refresh connections periodically
	db.SetConnMaxIdleTime(10 * time.Minute) // Close idle connections
	
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

// Ping checks database connectivity with a timeout.
func (d *DB) Ping(ctx context.Context) error {
	return d.db.PingContext(ctx)
}

// ViberUser represents a Viber user in the database.
type ViberUser struct {
	ViberID     string
	ViberName   string
	MatrixUserID *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// migrate creates all necessary database tables if they don't exist.
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

	CREATE TABLE IF NOT EXISTS group_members (
		viber_chat_id TEXT NOT NULL,
		viber_user_id TEXT NOT NULL,
		viber_user_name TEXT NOT NULL,
		joined_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (viber_chat_id, viber_user_id),
		FOREIGN KEY (viber_chat_id) REFERENCES room_mappings(viber_chat_id)
	);

	CREATE INDEX IF NOT EXISTS idx_room_mappings_viber ON room_mappings(viber_chat_id);
	CREATE INDEX IF NOT EXISTS idx_room_mappings_matrix ON room_mappings(matrix_room_id);
	CREATE INDEX IF NOT EXISTS idx_message_mappings_viber ON message_mappings(viber_message_id);
	CREATE INDEX IF NOT EXISTS idx_message_mappings_matrix ON message_mappings(matrix_event_id);
	`
	
	if _, err := d.db.Exec(schema); err != nil {
		return fmt.Errorf("execute migration: %w", err)
	}
	return nil
}

// UpsertViberUser creates or updates a Viber user in the database.
func (d *DB) UpsertViberUser(viberID, viberName string) error {
	_, err := d.db.Exec(`
		INSERT INTO viber_users (viber_id, viber_name, updated_at)
		VALUES (?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(viber_id) DO UPDATE SET
			viber_name = excluded.viber_name,
			updated_at = CURRENT_TIMESTAMP
	`, viberID, viberName)
	return err
}

// GetViberUser retrieves a Viber user by ID.
func (d *DB) GetViberUser(viberID string) (*ViberUser, error) {
	var user ViberUser
	var matrixUserID sql.NullString
	err := d.db.QueryRow(`
		SELECT viber_id, viber_name, matrix_user_id, created_at, updated_at
		FROM viber_users
		WHERE viber_id = ?
	`, viberID).Scan(&user.ViberID, &user.ViberName, &matrixUserID, &user.CreatedAt, &user.UpdatedAt)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	if matrixUserID.Valid {
		user.MatrixUserID = &matrixUserID.String
	}
	return &user, nil
}

// LinkViberUser links a Viber user to a Matrix user ID.
func (d *DB) LinkViberUser(viberID, matrixUserID string) error {
	if viberID == "" {
		return fmt.Errorf("viber_id cannot be empty")
	}
	_, err := d.db.Exec(`
		UPDATE viber_users
		SET matrix_user_id = ?, updated_at = CURRENT_TIMESTAMP
		WHERE viber_id = ?
	`, matrixUserID, viberID)
	return err
}

// UnlinkMatrixUser removes the Matrix user link from a Viber user.
func (d *DB) UnlinkMatrixUser(matrixUserID string) error {
	_, err := d.db.Exec(`
		UPDATE viber_users
		SET matrix_user_id = NULL, updated_at = CURRENT_TIMESTAMP
		WHERE matrix_user_id = ?
	`, matrixUserID)
	return err
}

// GetViberUserByMatrixID retrieves a Viber user by their linked Matrix user ID.
func (d *DB) GetViberUserByMatrixID(matrixUserID string) (*ViberUser, error) {
	var user ViberUser
	var matrixUserIDNullable sql.NullString
	err := d.db.QueryRow(`
		SELECT viber_id, viber_name, matrix_user_id, created_at, updated_at
		FROM viber_users
		WHERE matrix_user_id = ?
	`, matrixUserID).Scan(&user.ViberID, &user.ViberName, &matrixUserIDNullable, &user.CreatedAt, &user.UpdatedAt)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	if matrixUserIDNullable.Valid {
		user.MatrixUserID = &matrixUserIDNullable.String
	}
	return &user, nil
}

// ListLinkedUsers returns all users with linked Matrix accounts.
func (d *DB) ListLinkedUsers() ([]*ViberUser, error) {
	rows, err := d.db.Query(`
		SELECT viber_id, viber_name, matrix_user_id, created_at, updated_at
		FROM viber_users
		WHERE matrix_user_id IS NOT NULL
		ORDER BY updated_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var users []*ViberUser
	for rows.Next() {
		var user ViberUser
		var matrixUserID sql.NullString
		if err := rows.Scan(&user.ViberID, &user.ViberName, &matrixUserID, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		if matrixUserID.Valid {
			user.MatrixUserID = &matrixUserID.String
		}
		users = append(users, &user)
	}
	return users, rows.Err()
}

// CreateRoomMapping creates a mapping between a Viber chat and Matrix room.
func (d *DB) CreateRoomMapping(viberChatID, matrixRoomID string) error {
	_, err := d.db.Exec(`
		INSERT INTO room_mappings (viber_chat_id, matrix_room_id)
		VALUES (?, ?)
		ON CONFLICT(viber_chat_id) DO UPDATE SET
			matrix_room_id = excluded.matrix_room_id
	`, viberChatID, matrixRoomID)
	return err
}

// GetMatrixRoomID retrieves the Matrix room ID for a Viber chat.
func (d *DB) GetMatrixRoomID(viberChatID string) (string, error) {
	var matrixRoomID string
	err := d.db.QueryRow(`
		SELECT matrix_room_id
		FROM room_mappings
		WHERE viber_chat_id = ?
	`, viberChatID).Scan(&matrixRoomID)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return matrixRoomID, err
}

// GetViberChatID retrieves the Viber chat ID for a Matrix room.
func (d *DB) GetViberChatID(matrixRoomID string) (string, error) {
	var viberChatID string
	err := d.db.QueryRow(`
		SELECT viber_chat_id
		FROM room_mappings
		WHERE matrix_room_id = ?
	`, matrixRoomID).Scan(&viberChatID)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return viberChatID, err
}

// ListRoomMappings returns all room mappings.
func (d *DB) ListRoomMappings() ([]RoomMapping, error) {
	rows, err := d.db.Query(`
		SELECT viber_chat_id, matrix_room_id, created_at
		FROM room_mappings
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var mappings []RoomMapping
	for rows.Next() {
		var m RoomMapping
		if err := rows.Scan(&m.ViberChatID, &m.MatrixRoomID, &m.CreatedAt); err != nil {
			return nil, err
		}
		mappings = append(mappings, m)
	}
	return mappings, rows.Err()
}

// RoomMapping represents a room mapping in the database.
type RoomMapping struct {
	ViberChatID  string
	MatrixRoomID string
	CreatedAt    time.Time
}

// StoreMessageMapping stores a mapping between Viber message ID and Matrix event ID.
func (d *DB) StoreMessageMapping(viberMessageID, matrixEventID, viberChatID string) error {
	_, err := d.db.Exec(`
		INSERT INTO message_mappings (viber_message_id, matrix_event_id, viber_chat_id)
		VALUES (?, ?, ?)
		ON CONFLICT(viber_message_id) DO UPDATE SET
			matrix_event_id = excluded.matrix_event_id
	`, viberMessageID, matrixEventID, viberChatID)
	return err
}

// GetMatrixEventID retrieves the Matrix event ID for a Viber message.
func (d *DB) GetMatrixEventID(viberMessageID string) (string, error) {
	var matrixEventID string
	err := d.db.QueryRow(`
		SELECT matrix_event_id
		FROM message_mappings
		WHERE viber_message_id = ?
	`, viberMessageID).Scan(&matrixEventID)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return matrixEventID, err
}

// UpsertGroupMember adds or updates a group member in a Viber chat.
func (d *DB) UpsertGroupMember(viberChatID, viberUserID string, viberUserName ...string) error {
	name := viberUserID
	if len(viberUserName) > 0 && viberUserName[0] != "" {
		name = viberUserName[0]
	}
	_, err := d.db.Exec(`
		INSERT INTO group_members (viber_chat_id, viber_user_id, viber_user_name, joined_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(viber_chat_id, viber_user_id) DO UPDATE SET
			viber_user_name = excluded.viber_user_name
	`, viberChatID, viberUserID, name)
	return err
}

// ListGroupMembers lists all members of a Viber group chat.
func (d *DB) ListGroupMembers(viberChatID string) ([]string, error) {
	rows, err := d.db.Query(`
		SELECT viber_user_id
		FROM group_members
		WHERE viber_chat_id = ?
	`, viberChatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var members []string
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		members = append(members, userID)
	}
	return members, rows.Err()
}
