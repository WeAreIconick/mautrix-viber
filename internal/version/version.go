// Package version provides version information for the bridge.
package version

import (
	_ "embed"
	"strings"
)

//go:embed ../../VERSION
var versionStr string

// Version returns the current version string.
func Version() string {
	return strings.TrimSpace(versionStr)
}

// BuildInfo contains build information.
var BuildInfo struct {
	Version string
	GitCommit string
	BuildDate string
}

func init() {
	BuildInfo.Version = Version()
	// GitCommit and BuildDate would be set at build time via ldflags
}

