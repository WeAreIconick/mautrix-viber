// Package version provides version information for the bridge.
package version

import (
	"strings"
)

// versionStr is set at build time via ldflags, or defaults to "dev"
var versionStr = "dev"

// Version returns the current version string.
func Version() string {
	return strings.TrimSpace(versionStr)
}

// BuildInfo contains build information.
var BuildInfo struct {
	Version   string
	GitCommit string
	BuildDate string
}

func init() {
	BuildInfo.Version = Version()
	// GitCommit and BuildDate would be set at build time via ldflags
}
