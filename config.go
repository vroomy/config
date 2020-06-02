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

	ExampleRequests  map[string]*Request
	ExampleResponses map[string]*Response
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

// GetGroup will return group with name
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

	// Commands are the dynamic commands specified in config
	CommandEntries []*Command `toml:"command"`
	// Flags are the dynamic flags specified in config
	FlagEntries []*Flag `toml:"flag"`

	// Groups are the route groups
	Groups []*Group `toml:"group"`
	// Routes are the routes to listen for and serve
	Routes []*Route `toml:"route"`

	// Requests are example requests for docs/tests
	Requests []*Request `toml:"request"`
	// Responses are example responses for docs/tests
	Responses []*Response `toml:"response"`
}

// Request is an example requests for docs/tests
type Request struct {
	// ID
	Name  string `toml:"name"`
	Group string `toml:"group"`

	Query map[string]string `toml:"query"`
	Body  map[string]string `toml:"body"`

	Responses        []string    `toml:"responses"`
	ResponseExamples []*Response `toml:"-"`

	// Links
	Parent        string `toml:"parent"`
	parentExample *Request
}

// InheritFrom ensures family tree is populated and overridden
func (r *Request) InheritFrom(parents map[string]*Request) {
	if parent, ok := parents[r.Parent]; ok {
		if parent.Parent != "" && parent.parentExample == nil {
			// Recursive call to populate family tree from the bottom up
			parent.InheritFrom(parents)
		}

		if r.Query == nil {
			r.Query = parent.Query
		}

		if r.Body == nil {
			r.Body = parent.Body
		}

		if len(parent.Responses) > 0 {
			r.Responses = append(r.Responses, parent.Responses...)

			if len(parent.ResponseExamples) == 0 {
				r.ResponseExamples = append(r.ResponseExamples, parent.ResponseExamples...)
			}
		}
	}
}

// Response is an example response for docs/tests
type Response struct {
	// ID
	Name string `toml:"name"`

	// Links
	Parent        string `toml:"parent"`
	parentExample *Response

	// Inheritable vals
	Code int               `toml:"code"`
	Data map[string]string `toml:"data"`
}

// InheritFrom ensures family tree is populated and overridden
func (r *Response) InheritFrom(parents map[string]*Response) {
	if parent, ok := parents[r.Parent]; ok {
		if parent.Parent != "" && parent.parentExample == nil {
			// Recursive call to populate family tree from the bottom up
			parent.InheritFrom(parents)
		}

		if r.Code == 0 {
			r.Code = parent.Code
		}

		if r.Data == nil {
			r.Data = parent.Data
		}

		r.parentExample = parent
	}
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

	i.CommandEntries = append(i.CommandEntries, merge.CommandEntries...)
	i.FlagEntries = append(i.FlagEntries, merge.FlagEntries...)

	i.Groups = append(i.Groups, merge.Groups...)
	i.Routes = append(i.Routes, merge.Routes...)

	i.Requests = append(i.Requests, merge.Requests...)
	i.Responses = append(i.Responses, merge.Responses...)
}
