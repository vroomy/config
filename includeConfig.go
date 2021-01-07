package config

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
