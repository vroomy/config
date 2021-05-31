package config

import (
	"fmt"

	"github.com/hatchify/errors"
	"github.com/vroomy/common"
)

const (
	// ErrInvalidPluginHandler is returned when a plugin handler is not valid
	ErrInvalidPluginHandler = errors.Error("plugin handler not valid")
	// ErrExpectedEndParen is returned when an ending parenthesis is missing
	ErrExpectedEndParen = errors.Error("expected ending parenthesis")
)

// Route represents a listening route
type Route struct {
	// Target plug-in handler
	HTTPHandlers []common.Handler `toml:"-"`

	// Route name/description
	Name string `toml:"name"`
	// Route group
	Group string `toml:"group"`
	// HTTP method
	Method string `toml:"method"`
	// HTTP path
	HTTPPath string `toml:"httpPath"`
	// Directory or file to serve
	Target string `toml:"target"`
	// Plugin handlers
	Handlers []string `toml:"handlers"`

	RequestExample *Request
}

// String will return a formatted version of the route
func (r *Route) String() string {
	return fmt.Sprintf(routeFmt, r.HTTPPath, r.Target, r.Handlers)
}
