package config

import (
	"github.com/hatchify/errors"
	"github.com/vroomy/common"
)

const (
	// ErrGroupNotFound is returned when a group cannot be found by name
	ErrGroupNotFound = errors.Error("group not found")
)

// Group represents a route group
type Group struct {
	Name string `toml:"name"`
	// Route group
	Group string `toml:"group"`
	// HTTP method
	Method string `toml:"method"`
	// HTTP path
	HTTPPath string `toml:"httpPath"`
	// Plugin handlers
	Handlers []string `toml:"handlers"`

	HTTPHandlers []common.Handler `toml:"-"`

	G common.Group `toml:"-"`

	// Requests are keys to the request map which includes example request/response data for docs and tests
	Requests map[string]*Request
}
