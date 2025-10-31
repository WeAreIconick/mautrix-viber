// Package matrix provides Matrix client functionality for the bridge.
package matrix

// Bridge is an interface for Matrix bridge functionality.
type Bridge interface {
	Start() error
	Stop() error
}
