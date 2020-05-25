package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/Hatch1fy/errors"
)

const (
	// RouteFmt specifies expected route definition syntax
	routeFmt = "{ HTTPPath: \"%s\", Target: \"%s\" Plugin Handler: \"%v\" }"
	// ErrProtectedFlag is returned when a protected flag is used
	ErrProtectedFlag = errors.Error("cannot use protected flag")
)

// NewConfig will return a new configuration
func NewConfig(loc string) (cfg *Config, err error) {
	var c Config
	if _, err = toml.DecodeFile(loc, &c); err != nil {
		return
	}

	if err = c.loadIncludes(); err != nil {
		return
	}

	if c.Dir == "" {
		c.Dir = "./"
	}

	cfg = &c
	return
}

// PluginConfig is abbreviated import for vpm use case
type PluginConfig struct {
	// Specify which plugins are in scope
	Plugins []string `toml:"plugins"`
}

// Config is the configuration needed to initialize a new instance of Service
type Config struct {
	Name string `toml:"name"`

	Dir     string `toml:"dir"`
	Port    uint16 `toml:"port"`
	TLSPort uint16 `toml:"tlsPort"`
	TLSDir  string `toml:"tlsDir"`

	IncludeConfig

	PerformUpdate bool `toml:"-"`

	Flags map[string]string `toml:"-"`

	// Plugin keys as they are referenced by the plugins store
	PluginKeys []string
}

func (c *Config) loadIncludes() (err error) {
	for _, include := range c.Include {
		// Include each file or directory
		if err = c.loadInclude(include); err != nil {
			// Include failed
			return
		}
	}

	return
}

func (c *Config) loadInclude(include string) (err error) {
	if path.Ext(include) == ".toml" {
		// Attempt to decode toml
		var icfg IncludeConfig
		if _, err = toml.DecodeFile(include, &icfg); err != nil {
			return
		}

		c.IncludeConfig.merge(&icfg)
	} else {
		// Attempt to parse directory
		var files []os.FileInfo
		if files, err = ioutil.ReadDir(include); err != nil {
			return fmt.Errorf("%s is not a .toml file or directory", include)
		}

		// Call recursively
		for _, file := range files {
			if err = c.loadInclude(path.Join(include, file.Name())); err != nil {
				return
			}
		}
	}

	return
}

func (c *Config) GetGroup(name string) (g *Group, err error) {
	if len(name) == 0 {
		return
	}

	// TODO: Make this a map for faster lookups?
	for _, group := range c.Groups {
		if group.Name != name {
			continue
		}

		g = group
		return
	}

	err = ErrGroupNotFound
	return
}

// IncludeConfig will include routes
type IncludeConfig struct {
	// Application environment
	Environment map[string]string `toml:"env"`

	// Allow included files to add includes
	Include []string `toml:"include"`

	// Specify which plugins are in scope
	Plugins []string `toml:"plugins"`

	// Groups are the route groups
	Groups []*Group `toml:"group"`
	// Routes are the routes to listen for and serve
	Routes []*Route `toml:"route"`

	// Commands are the dynamic commands specified in config
	CommandEntries []*Command `toml:"command"`
	// Flags are the dynamic flags specified in config
	FlagEntries []*Flag `toml:"flag"`
}

func (i *IncludeConfig) merge(merge *IncludeConfig) {
	if i.Environment == nil {
		i.Environment = make(map[string]string)
	}

	for key, val := range merge.Environment {
		i.Environment[key] = val
	}

	i.Include = append(i.Include, merge.Include...)

	i.Plugins = append(i.Plugins, merge.Plugins...)

	i.Groups = append(i.Groups, merge.Groups...)
	i.Routes = append(i.Routes, merge.Routes...)

	i.CommandEntries = append(i.CommandEntries, merge.CommandEntries...)
	i.FlagEntries = append(i.FlagEntries, merge.FlagEntries...)
}
