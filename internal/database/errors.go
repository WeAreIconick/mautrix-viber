package database

import "errors"

// Sentinel errors for database operations.
// These can be checked using errors.Is() for specific error handling.
var (
	// ErrNotFound indicates the requested record was not found.
	ErrNotFound = errors.New("database: record not found")

	// ErrAlreadyExists indicates a record with the same unique constraint already exists.
	ErrAlreadyExists = errors.New("database: record already exists")

	// ErrInvalidInput indicates the input data is invalid (e.g., empty required field).
	ErrInvalidInput = errors.New("database: invalid input")
)
