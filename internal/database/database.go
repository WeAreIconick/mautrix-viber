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

// CacheInterface defines the caching interface for database methods.
// This allows for optional caching without requiring Redis to be configured.
type CacheInterface interface {
	GetJSON(ctx context.Context, key string, v interface{}) error
	SetJSON(ctx context.Context, key string, v interface{}) error
	Delete(ctx context.Context, key string) error
}

// DB wraps sql.DB with bridge-specific methods.
// Includes optional caching layer for frequently accessed data.
type DB struct {
	db    *sql.DB
	cache CacheInterface // Optional cache for frequently accessed queries
}

// Open opens or creates a SQLite database at path and runs migrations.
// Configures connection pool settings for production use:
// - SetMaxOpenConns: Maximum number of open connections (25 for SQLite)
// - SetMaxIdleConns: Maximum idle connections (5)
// - SetConnMaxLifetime: Maximum connection lifetime (5 minutes)
// - SetConnMaxIdleTime: Maximum idle time before closing (10 minutes)
// Optionally accepts a cache interface for frequently accessed data.
func Open(path string) (*DB, error) {
	return OpenWithCache(path, nil)
}

// OpenWithCache opens a database with an optional cache layer.
// The cache is used for frequently accessed queries like room mappings and user lookups.
func OpenWithCache(path string, cache CacheInterface) (*DB, error) {
	db, err := sql.Open("sqlite3", path+"?_foreign_keys=1&_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	// Configure connection pool for production
	db.SetMaxOpenConns(25)                  // SQLite default is usually unlimited, but we limit for safety
	db.SetMaxIdleConns(5)                   // Keep some connections ready
	db.SetConnMaxLifetime(5 * time.Minute)  // Refresh connections periodically
	db.SetConnMaxIdleTime(10 * time.Minute) // Close idle connections

	// Retry database ping with exponential backoff (useful for network databases or startup delays)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var pingErr error
	for attempt := 0; attempt < 3; attempt++ {
		if attempt > 0 {
			// Exponential backoff: 100ms, 200ms, 400ms
			delay := time.Duration(100*(1<<uint(attempt-1))) * time.Millisecond
			select {
			case <-ctx.Done():
				db.Close()
				return nil, fmt.Errorf("database ping timeout: %w", pingErr)
			case <-time.After(delay):
			}
		}

		pingCtx, pingCancel := context.WithTimeout(ctx, 5*time.Second)
		pingErr = db.PingContext(pingCtx)
		pingCancel()

		if pingErr == nil {
			break
		}
	}

	if pingErr != nil {
		db.Close()
		return nil, fmt.Errorf("ping database after retries: %w", pingErr)
	}
	d := &DB{db: db, cache: cache}
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
	ViberID      string
	ViberName    string
	MatrixUserID *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
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
// The context controls cancellation and timeout for the operation.
// Invalidates cache if configured.
func (d *DB) UpsertViberUser(ctx context.Context, viberID, viberName string) error {
	if viberID == "" {
		return fmt.Errorf("%w: viber_id cannot be empty", ErrInvalidInput)
	}
	if viberName == "" {
		return fmt.Errorf("%w: viber_name cannot be empty", ErrInvalidInput)
	}
	_, err := d.db.ExecContext(ctx, `
		INSERT INTO viber_users (viber_id, viber_name, updated_at)
		VALUES (?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(viber_id) DO UPDATE SET
			viber_name = excluded.viber_name,
			updated_at = CURRENT_TIMESTAMP
	`, viberID, viberName)
	if err != nil {
		return fmt.Errorf("upsert viber user %s: %w", viberID, err)
	}

	// Invalidate cache if configured
	if d.cache != nil {
		key := "user:viber:" + viberID
		_ = d.cache.Delete(ctx, key) // Best-effort cache invalidation, ignore errors
	}

	return nil
}

// GetViberUser retrieves a Viber user by ID.
// Returns ErrNotFound if the user does not exist.
// The context controls cancellation and timeout for the operation.
// Uses cache if configured for improved performance.
func (d *DB) GetViberUser(ctx context.Context, viberID string) (*ViberUser, error) {
	if viberID == "" {
		return nil, fmt.Errorf("%w: viber_id cannot be empty", ErrInvalidInput)
	}

	// Try cache first if configured
	if d.cache != nil {
		var user ViberUser
		key := "user:viber:" + viberID
		if err := d.cache.GetJSON(ctx, key, &user); err == nil {
			return &user, nil
		}
	}

	var user ViberUser
	var matrixUserID sql.NullString
	err := d.db.QueryRowContext(ctx, `
		SELECT viber_id, viber_name, matrix_user_id, created_at, updated_at
		FROM viber_users
		WHERE viber_id = ?
	`, viberID).Scan(&user.ViberID, &user.ViberName, &matrixUserID, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("%w: viber user %s", ErrNotFound, viberID)
	}
	if err != nil {
		return nil, fmt.Errorf("query viber user %s: %w", viberID, err)
	}

	if matrixUserID.Valid {
		user.MatrixUserID = &matrixUserID.String
	}

	// Cache the result if cache is configured
	if d.cache != nil {
		key := "user:viber:" + viberID
		_ = d.cache.SetJSON(ctx, key, &user) // Best-effort cache, ignore errors
	}

	return &user, nil
}

// LinkViberUser links a Viber user to a Matrix user ID.
// The context controls cancellation and timeout for the operation.
// Invalidates cache if configured.
func (d *DB) LinkViberUser(ctx context.Context, viberID, matrixUserID string) error {
	if viberID == "" {
		return fmt.Errorf("%w: viber_id cannot be empty", ErrInvalidInput)
	}
	if matrixUserID == "" {
		return fmt.Errorf("%w: matrix_user_id cannot be empty", ErrInvalidInput)
	}
	_, err := d.db.ExecContext(ctx, `
		UPDATE viber_users
		SET matrix_user_id = ?, updated_at = CURRENT_TIMESTAMP
		WHERE viber_id = ?
	`, matrixUserID, viberID)
	if err != nil {
		return fmt.Errorf("link viber user %s to matrix user %s: %w", viberID, matrixUserID, err)
	}

	// Invalidate cache if configured
	if d.cache != nil {
		key := "user:viber:" + viberID
		_ = d.cache.Delete(ctx, key) // Best-effort cache invalidation
	}

	return nil
}

// UnlinkMatrixUser removes the Matrix user link from a Viber user.
// The context controls cancellation and timeout for the operation.
func (d *DB) UnlinkMatrixUser(ctx context.Context, matrixUserID string) error {
	if matrixUserID == "" {
		return fmt.Errorf("%w: matrix_user_id cannot be empty", ErrInvalidInput)
	}
	_, err := d.db.ExecContext(ctx, `
		UPDATE viber_users
		SET matrix_user_id = NULL, updated_at = CURRENT_TIMESTAMP
		WHERE matrix_user_id = ?
	`, matrixUserID)
	if err != nil {
		return fmt.Errorf("unlink matrix user %s: %w", matrixUserID, err)
	}
	return nil
}

// GetViberUserByMatrixID retrieves a Viber user by their linked Matrix user ID.
// The context controls cancellation and timeout for the operation.
// Returns ErrNotFound if the user does not exist.
func (d *DB) GetViberUserByMatrixID(ctx context.Context, matrixUserID string) (*ViberUser, error) {
	if matrixUserID == "" {
		return nil, fmt.Errorf("%w: matrix_user_id cannot be empty", ErrInvalidInput)
	}
	var user ViberUser
	var matrixUserIDNullable sql.NullString
	err := d.db.QueryRowContext(ctx, `
		SELECT viber_id, viber_name, matrix_user_id, created_at, updated_at
		FROM viber_users
		WHERE matrix_user_id = ?
	`, matrixUserID).Scan(&user.ViberID, &user.ViberName, &matrixUserIDNullable, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("%w: no viber user linked to matrix user %s", ErrNotFound, matrixUserID)
	}
	if err != nil {
		return nil, fmt.Errorf("query viber user by matrix id %s: %w", matrixUserID, err)
	}

	if matrixUserIDNullable.Valid {
		user.MatrixUserID = &matrixUserIDNullable.String
	}
	return &user, nil
}

// ListLinkedUsers returns all users with linked Matrix accounts.
// The context controls cancellation and timeout for the operation.
func (d *DB) ListLinkedUsers(ctx context.Context) ([]*ViberUser, error) {
	rows, err := d.db.QueryContext(ctx, `
		SELECT viber_id, viber_name, matrix_user_id, created_at, updated_at
		FROM viber_users
		WHERE matrix_user_id IS NOT NULL
		ORDER BY updated_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("query linked users: %w", err)
	}
	defer rows.Close()

	var users []*ViberUser
	for rows.Next() {
		var user ViberUser
		var matrixUserID sql.NullString
		if err := rows.Scan(&user.ViberID, &user.ViberName, &matrixUserID, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan linked user: %w", err)
		}
		if matrixUserID.Valid {
			user.MatrixUserID = &matrixUserID.String
		}
		users = append(users, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate linked users: %w", err)
	}
	return users, nil
}

// CreateRoomMapping creates a mapping between a Viber chat and Matrix room.
// The context controls cancellation and timeout for the operation.
// Invalidates cache if configured.
func (d *DB) CreateRoomMapping(ctx context.Context, viberChatID, matrixRoomID string) error {
	if viberChatID == "" {
		return fmt.Errorf("%w: viber_chat_id cannot be empty", ErrInvalidInput)
	}
	if matrixRoomID == "" {
		return fmt.Errorf("%w: matrix_room_id cannot be empty", ErrInvalidInput)
	}
	_, err := d.db.ExecContext(ctx, `
		INSERT INTO room_mappings (viber_chat_id, matrix_room_id)
		VALUES (?, ?)
		ON CONFLICT(viber_chat_id) DO UPDATE SET
			matrix_room_id = excluded.matrix_room_id
	`, viberChatID, matrixRoomID)
	if err != nil {
		return fmt.Errorf("create room mapping for chat %s -> room %s: %w", viberChatID, matrixRoomID, err)
	}

	// Invalidate cache if configured
	if d.cache != nil {
		key1 := "room:viber:" + viberChatID
		key2 := "room:matrix:" + matrixRoomID
		_ = d.cache.Delete(ctx, key1) // Best-effort cache invalidation
		_ = d.cache.Delete(ctx, key2)
	}

	return nil
}

// GetMatrixRoomID retrieves the Matrix room ID for a Viber chat.
// Returns empty string and nil error if mapping does not exist (not an error condition).
// The context controls cancellation and timeout for the operation.
// Uses cache if configured for improved performance.
func (d *DB) GetMatrixRoomID(ctx context.Context, viberChatID string) (string, error) {
	if viberChatID == "" {
		return "", fmt.Errorf("%w: viber_chat_id cannot be empty", ErrInvalidInput)
	}

	// Try cache first if configured
	if d.cache != nil {
		var cachedRoomID string
		key := "room:viber:" + viberChatID
		if err := d.cache.GetJSON(ctx, key, &cachedRoomID); err == nil {
			return cachedRoomID, nil
		}
	}

	var matrixRoomID string
	err := d.db.QueryRowContext(ctx, `
		SELECT matrix_room_id
		FROM room_mappings
		WHERE viber_chat_id = ?
	`, viberChatID).Scan(&matrixRoomID)
	if err == sql.ErrNoRows {
		return "", nil // Not found is not an error - mapping may not exist yet
	}
	if err != nil {
		return "", fmt.Errorf("query matrix room id for chat %s: %w", viberChatID, err)
	}

	// Cache the result if cache is configured
	if d.cache != nil && matrixRoomID != "" {
		key := "room:viber:" + viberChatID
		_ = d.cache.SetJSON(ctx, key, matrixRoomID) // Best-effort cache, ignore errors
	}

	return matrixRoomID, nil
}

// GetViberChatID retrieves the Viber chat ID for a Matrix room.
// Returns empty string and nil error if mapping does not exist (not an error condition).
// The context controls cancellation and timeout for the operation.
// Uses cache if configured for improved performance.
func (d *DB) GetViberChatID(ctx context.Context, matrixRoomID string) (string, error) {
	if matrixRoomID == "" {
		return "", fmt.Errorf("%w: matrix_room_id cannot be empty", ErrInvalidInput)
	}

	// Try cache first if configured
	if d.cache != nil {
		var cachedChatID string
		key := "room:matrix:" + matrixRoomID
		if err := d.cache.GetJSON(ctx, key, &cachedChatID); err == nil {
			return cachedChatID, nil
		}
	}

	var viberChatID string
	err := d.db.QueryRowContext(ctx, `
		SELECT viber_chat_id
		FROM room_mappings
		WHERE matrix_room_id = ?
	`, matrixRoomID).Scan(&viberChatID)
	if err == sql.ErrNoRows {
		return "", nil // Not found is not an error - mapping may not exist yet
	}
	if err != nil {
		return "", fmt.Errorf("query viber chat id for room %s: %w", matrixRoomID, err)
	}

	// Cache the result if cache is configured
	if d.cache != nil && viberChatID != "" {
		key := "room:matrix:" + matrixRoomID
		_ = d.cache.SetJSON(ctx, key, viberChatID) // Best-effort cache, ignore errors
	}

	return viberChatID, nil
}

// ListRoomMappings returns all room mappings.
// The context controls cancellation and timeout for the operation.
func (d *DB) ListRoomMappings(ctx context.Context) ([]RoomMapping, error) {
	rows, err := d.db.QueryContext(ctx, `
		SELECT viber_chat_id, matrix_room_id, created_at
		FROM room_mappings
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("query room mappings: %w", err)
	}
	defer rows.Close()

	var mappings []RoomMapping
	for rows.Next() {
		var m RoomMapping
		if err := rows.Scan(&m.ViberChatID, &m.MatrixRoomID, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan room mapping: %w", err)
		}
		mappings = append(mappings, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate room mappings: %w", err)
	}
	return mappings, nil
}

// RoomMapping represents a room mapping in the database.
type RoomMapping struct {
	ViberChatID  string
	MatrixRoomID string
	CreatedAt    time.Time
}

// StoreMessageMapping stores a mapping between Viber message ID and Matrix event ID.
// The context controls cancellation and timeout for the operation.
func (d *DB) StoreMessageMapping(ctx context.Context, viberMessageID, matrixEventID, viberChatID string) error {
	if viberMessageID == "" {
		return fmt.Errorf("%w: viber_message_id cannot be empty", ErrInvalidInput)
	}
	if matrixEventID == "" {
		return fmt.Errorf("%w: matrix_event_id cannot be empty", ErrInvalidInput)
	}
	if viberChatID == "" {
		return fmt.Errorf("%w: viber_chat_id cannot be empty", ErrInvalidInput)
	}
	_, err := d.db.ExecContext(ctx, `
		INSERT INTO message_mappings (viber_message_id, matrix_event_id, viber_chat_id)
		VALUES (?, ?, ?)
		ON CONFLICT(viber_message_id) DO UPDATE SET
			matrix_event_id = excluded.matrix_event_id
	`, viberMessageID, matrixEventID, viberChatID)
	if err != nil {
		return fmt.Errorf("store message mapping %s -> %s: %w", viberMessageID, matrixEventID, err)
	}
	return nil
}

// GetMatrixEventID retrieves the Matrix event ID for a Viber message.
// Returns empty string and nil error if mapping does not exist (not an error condition).
// The context controls cancellation and timeout for the operation.
func (d *DB) GetMatrixEventID(ctx context.Context, viberMessageID string) (string, error) {
	if viberMessageID == "" {
		return "", fmt.Errorf("%w: viber_message_id cannot be empty", ErrInvalidInput)
	}
	var matrixEventID string
	err := d.db.QueryRowContext(ctx, `
		SELECT matrix_event_id
		FROM message_mappings
		WHERE viber_message_id = ?
	`, viberMessageID).Scan(&matrixEventID)
	if err == sql.ErrNoRows {
		return "", nil // Not found is not an error - mapping may not exist yet
	}
	if err != nil {
		return "", fmt.Errorf("query matrix event id for message %s: %w", viberMessageID, err)
	}
	return matrixEventID, nil
}

// UpsertGroupMember adds or updates a group member in a Viber chat.
// The context controls cancellation and timeout for the operation.
func (d *DB) UpsertGroupMember(ctx context.Context, viberChatID, viberUserID string, viberUserName ...string) error {
	if viberChatID == "" {
		return fmt.Errorf("%w: viber_chat_id cannot be empty", ErrInvalidInput)
	}
	if viberUserID == "" {
		return fmt.Errorf("%w: viber_user_id cannot be empty", ErrInvalidInput)
	}
	name := viberUserID
	if len(viberUserName) > 0 && viberUserName[0] != "" {
		name = viberUserName[0]
	}
	_, err := d.db.ExecContext(ctx, `
		INSERT INTO group_members (viber_chat_id, viber_user_id, viber_user_name, joined_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(viber_chat_id, viber_user_id) DO UPDATE SET
			viber_user_name = excluded.viber_user_name
	`, viberChatID, viberUserID, name)
	if err != nil {
		return fmt.Errorf("upsert group member %s in chat %s: %w", viberUserID, viberChatID, err)
	}
	return nil
}

// ListGroupMembers lists all members of a Viber group chat.
// The context controls cancellation and timeout for the operation.
func (d *DB) ListGroupMembers(ctx context.Context, viberChatID string) ([]string, error) {
	if viberChatID == "" {
		return nil, fmt.Errorf("%w: viber_chat_id cannot be empty", ErrInvalidInput)
	}
	rows, err := d.db.QueryContext(ctx, `
		SELECT viber_user_id
		FROM group_members
		WHERE viber_chat_id = ?
	`, viberChatID)
	if err != nil {
		return nil, fmt.Errorf("query group members for chat %s: %w", viberChatID, err)
	}
	defer rows.Close()

	var members []string
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, fmt.Errorf("scan group member: %w", err)
		}
		members = append(members, userID)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate group members: %w", err)
	}
	return members, nil
}
