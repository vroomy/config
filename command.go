package config

// Command represents a command entry
type Command struct {
	Name    string `toml:"name"`
	Usage   string `toml:"usage"`
	Require string `toml:"require"`

	Prehook  string `toml:"prehook"`
	Handler  string `toml:"handler"`
	Posthook string `toml:"posthook"`
}
