package config

import (
	"github.com/Hatch1fy/errors"
	"github.com/Hatch1fy/httpserve"
	"github.com/vroomy/plugins"
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

	HttpHandlers []httpserve.Handler `toml:"-"`

	G httpserve.Group `toml:"-"`
}

// Init will init a group
func (g *Group) Init(p *plugins.Plugins) (err error) {
	for _, handlerKey := range g.Handlers {
		var h httpserve.Handler
		if h, err = newPluginHandler(p, handlerKey); err != nil {
			return
		}

		g.HttpHandlers = append(g.HttpHandlers, h)
	}

	return
}
