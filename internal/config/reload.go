// Package config reload provides hot reload configuration on SIGHUP without restart.
package config

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// ReloadableConfig represents a configuration that can be reloaded.
type ReloadableConfig struct {
	mu     sync.RWMutex
	config Config
	reload func() (Config, error)
}

// NewReloadableConfig creates a new reloadable config.
func NewReloadableConfig(initial Config, reloadFn func() (Config, error)) *ReloadableConfig {
	rc := &ReloadableConfig{
		config: initial,
		reload: reloadFn,
	}

	// Listen for SIGHUP
	go rc.listenForReload()

	return rc
}

// Get returns the current configuration (thread-safe).
func (rc *ReloadableConfig) Get() Config {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	return rc.config
}

// Reload reloads the configuration.
func (rc *ReloadableConfig) Reload() error {
	if rc.reload == nil {
		return nil
	}

	newConfig, err := rc.reload()
	if err != nil {
		return err
	}

	rc.mu.Lock()
	rc.config = newConfig
	rc.mu.Unlock()

	return nil
}

// listenForReload listens for SIGHUP and reloads configuration.
func (rc *ReloadableConfig) listenForReload() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP)

	for range sigChan {
		if err := rc.Reload(); err != nil {
			// Log error - in production use structured logging
			// logger.Error("failed to reload configuration", "error", err)
		} else {
			// Log success - in production use structured logging
			// logger.Info("configuration reloaded successfully")
		}
	}
}
